package handler

import (
	"net/http"
	"github.com/gorilla/mux"
	"goMail/utils"
	"goMail/dbs/mysql"
	"goMail/dbs/tarantool"
	"strconv"
	"errors"
	"go.uber.org/zap"
)

func deleteTuple(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token, _ := utils.HandleToken(r.Header.Get("Authorization"))
	utils.Mu.Lock()
	a, _ := utils.Cookies[token]
	utils.Mu.Unlock()
	user, _ := mysql.GetUser(a)
	b, exists := vars["name_space"]
	if !exists {
		utils.WriteResponse(w, errors.New("name_space param is undefined"), nil, http.StatusBadRequest)
		utils.Logger.Error("cannot get name_space",
			zap.Error(errors.New("name_space param is undefined")),
		)
		return
	}
	c, exists := vars["id_tuple"]
	if !exists {
		utils.WriteResponse(w, errors.New("id_tuple param is undefined"), nil, http.StatusBadRequest)
		utils.Logger.Error("cannot get id_tuple",
			zap.Error(errors.New("name_space param is undefined")),
		)
		return
	}
	s, err := strconv.ParseUint(c, 10, 64)
	if err != nil {
		utils.WriteResponse(w, errors.New("id is not valid"), nil, http.StatusBadRequest)
		utils.Logger.Error("cannot parse id",
			zap.Error(errors.New("name_space param is undefined")),
		)
		return
	}

	space, err := mysql.GetSpace(b, user.Id)
	if err != nil {
		utils.WriteResponse(w, err, nil, http.StatusInternalServerError)
		utils.Logger.Error("cannot get space",
			zap.Error(err),
		)
		return
	}
	if !exists {
		utils.WriteResponse(w, errors.New("space with such name doesn't exist"), nil, http.StatusInternalServerError)
		utils.Logger.Error("cannot get space",
			zap.Error(errors.New("space with such name doesn't exist")),
		)
		return
	}
	if ok, _ := mysql.CheckPermissionsOnSpace(user.Id, space.Id); !ok {
		utils.WriteResponse(w, errors.New("get rights on this space"), nil, http.StatusMethodNotAllowed)
		utils.Logger.Error("cannot get rights",
			zap.Error(err),
				zap.Error(errors.New("get rights on this space")),
		)
		return
	}
	_, err = tarantool.SelectTuple(s, b, user.Id)
	if !exists {
		utils.WriteResponse(w, errors.New("tuple with such id doesn't exist"), nil, http.StatusBadRequest)
		utils.Logger.Error("cannot select tuple",
			zap.Error(errors.New("tuple with such id doesn't exist")),
		)
		return
	}
	resp, err := tarantool.DeleteTuple(s, b, user.Id)
	if err != nil {
		utils.WriteResponse(w, err, nil, http.StatusInternalServerError)
		utils.Logger.Error("cannot delete tuple",
			zap.Error(err),
		)
		return
	}
	utils.WriteResponse(w, nil, resp, http.StatusOK)
	utils.Logger.Info("success executing function",
		zap.String("url", r.URL.Path),
	)
}

func deleteSpace(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token, _ := utils.HandleToken(r.Header.Get("Authorization"))
	utils.Mu.Lock()
	a, _ := utils.Cookies[token]
	utils.Mu.Unlock()
	user, _ := mysql.GetUser(a)
	b, exists := vars["name_space"]
	if !exists {
		utils.WriteResponse(w, errors.New("name_space param is undefined"), nil, http.StatusBadRequest)
		utils.Logger.Error("cannot get name_space",
			zap.Error(errors.New("name_space param is undefined")),
		)
		return
	}

	space, err := mysql.GetSpace(b, user.Id)
	if err != nil {
		utils.WriteResponse(w, err, nil, http.StatusInternalServerError)
		utils.Logger.Error("cannot get space",
			zap.Error(err),
		)
		return
	}

	mysql.DeleteSpace(space.Name, user.Id)
	if err != nil {
		utils.WriteResponse(w, err, nil, http.StatusInternalServerError)
		utils.Logger.Error("cannot delete tuple",
			zap.Error(err),
		)
		return
	}

	resp, err := tarantool.DeleteSpace(b, user.Id)
	if err != nil {
		utils.WriteResponse(w, err, nil, http.StatusInternalServerError)
		utils.Logger.Error("cannot delete tuple",
			zap.Error(err),
		)
		return
	}
	utils.WriteResponse(w, nil, resp, http.StatusOK)
	utils.Logger.Info("success executing function",
		zap.String("url", r.URL.Path),
	)
}