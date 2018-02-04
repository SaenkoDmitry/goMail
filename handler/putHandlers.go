package handler

import (
	"net/http"
	"github.com/gorilla/mux"
	"goMail/utils"
	"goMail/dbs/mysql"
	"goMail/dbs/tarantool"
	"strconv"
	"encoding/json"
	"errors"
	"go.uber.org/zap"
)

func updateTuple(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token, _ := utils.HandleToken(r.Header.Get("Authorization"))
	decoder := json.NewDecoder(r.Body)
	var data []interface{}
	err := decoder.Decode(&data)
	if err != nil {
		utils.WriteResponse(w, err, nil, http.StatusBadRequest)
		utils.Logger.Error("can't perform operation",
			zap.Error(err),
		)
		return
	}
	utils.Mu.Lock()
	a, _ := utils.Cookies[token]
	utils.Mu.Unlock()
	user, err := mysql.GetUser(a)
	if err != nil {
		utils.WriteResponse(w, err, nil, http.StatusBadRequest)
		utils.Logger.Error("can't perform operation",
			zap.Error(err),
		)
		return
	}
	b, exists := vars["name_space"]
	if !exists {
		utils.WriteResponse(w, errors.New("param name_space is undefined"), nil, http.StatusBadRequest)
		utils.Logger.Error("can't perform operation",
			zap.Error(errors.New("param name_space is undefined")),
		)
		return
	}
	c, exists := vars["id_tuple"]
	if !exists {
		utils.WriteResponse(w, errors.New("param id_tuple is undefined"), nil, http.StatusBadRequest)
		utils.Logger.Error("can't perform operation",
			zap.Error(errors.New("param id_tuple is undefined")),
		)
		return
	}
	s, err := strconv.ParseUint(c, 10, 64)
	if err != nil {
		utils.WriteResponse(w, err, nil, http.StatusBadRequest)
		utils.Logger.Error("can't perform operation",
			zap.Error(err),
		)
		return
	}
	space, err := mysql.GetSpace(b, user.Id)
	if err != nil {
		utils.WriteResponse(w, err, nil, http.StatusInternalServerError)
		utils.Logger.Error("can't perform operation",
			zap.Error(err),
		)
		return
	}
	if _, err := mysql.CheckPermissionsOnSpace(user.Id, space.Id); err != nil {
		utils.WriteResponse(w, err, nil, http.StatusMethodNotAllowed)
		utils.Logger.Error("can't perform operation",
			zap.Error(err),
		)
		return
	}
	resp, ok, err := tarantool.UpdateTuple(s, b, user.Id, data)
	if err != nil {
		utils.WriteResponse(w, err, nil, http.StatusInternalServerError)
		utils.Logger.Error("can't perform operation",
			zap.Error(err),
		)
		return
	}
	if !ok {
		utils.WriteResponse(w, errors.New("cannot update tuple in tarantool"), nil, http.StatusInternalServerError)
		utils.Logger.Error("can't perform operation",
			zap.Error(errors.New("cannot update tuple in tarantool")),
		)
		return
	}
	utils.WriteResponse(w, nil, resp, http.StatusOK)
	utils.Logger.Info("success executing function",
		zap.String("url", r.URL.Path),
	)
}