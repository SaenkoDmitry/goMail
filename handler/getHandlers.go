package handler

import (
	"net/http"
	"github.com/gorilla/mux"
	"goMail/dbs/tarantool"
	"goMail/utils"
	"strconv"
	"goMail/dbs/mysql"
	"errors"
	"go.uber.org/zap"
)

func getAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := mysql.GetAllUsers()
	if err != nil {
		utils.WriteResponse(w, err, nil, http.StatusInternalServerError)
		utils.Logger.Error("can't perform operation",
			zap.Error(err),
		)
		return
	}
	utils.WriteResponse(w, nil, users, http.StatusOK)
}

func getAllSpaces(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	token, _ := utils.HandleToken(r.Header.Get("Authorization"))
	utils.Mu.Lock()
	a, _ := utils.Cookies[token]
	utils.Mu.Unlock()
	user, _ := mysql.GetUser(a)
	spaces, err := mysql.GetAllSpaces(user.Id)
	if err != nil {
		utils.WriteResponse(w, err, nil, http.StatusInternalServerError)
		utils.Logger.Error("can't perform operation",
			zap.Error(err),
		)
		return
	}
	utils.WriteResponse(w, nil, spaces, http.StatusOK)
}

func getUserHistory(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	token, _ := utils.HandleToken(r.Header.Get("Authorization"))
	utils.Mu.Lock()
	a, _ := utils.Cookies[token]
	utils.Mu.Unlock()
	user, err := mysql.GetUser(a)
	history, err := mysql.GetUserHistory(user.Id)
	if err != nil {
		utils.WriteResponse(w, err, nil, http.StatusInternalServerError)
		utils.Logger.Error("can't perform operation",
			zap.Error(err),
		)
		return
	}
	utils.WriteResponse(w, nil, history, http.StatusOK)
	utils.Logger.Info("success executing function",
		zap.String("url", r.URL.Path),
	)
}

func getSpaceHistory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	r.ParseForm()
	token, _ := utils.HandleToken(r.Header.Get("Authorization"))
	utils.Mu.Lock()
	a, _ := utils.Cookies[token]
	utils.Mu.Unlock()
	user, err := mysql.GetUser(a)
	spaceName, exists := vars["name_space"]
	if !exists {
		utils.WriteResponse(w, errors.New("name_space param undefined"), nil, http.StatusBadRequest)
		utils.Logger.Error("can't perform operation",
			zap.Error(errors.New("name_space param undefined")),
		)
		return
	}
	space, err := mysql.GetSpace(spaceName, user.Id)
	if err != nil {
		utils.WriteResponse(w, err, nil, http.StatusInternalServerError)
		utils.Logger.Error("can't perform operation",
			zap.Error(err),
		)
		return
	}
	if b, err := mysql.CheckPermissionsOnSpace(user.Id, space.Id); !b {
		if err != nil {
			utils.WriteResponse(w, err, nil, http.StatusInternalServerError)
		} else {
			utils.WriteResponse(w, errors.New("check rights on space"), nil, http.StatusMethodNotAllowed)
		}
		utils.Logger.Error("can't perform operation",
			zap.Error(err),
		)
		return
	}
	history, err := mysql.GetSpaceHistory(user.Id)
	if err != nil {
		utils.WriteResponse(w, err, nil, http.StatusInternalServerError)
		utils.Logger.Error("can't perform operation",
			zap.Error(err),
		)
		return
	}
	utils.WriteResponse(w, nil, history, http.StatusOK)
	utils.Logger.Info("success executing function",
		zap.String("url", r.URL.Path),
	)
}

func getAllSpacePermissions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	r.ParseForm()
	token, _ := utils.HandleToken(r.Header.Get("Authorization"))
	utils.Mu.Lock()
	a, _ := utils.Cookies[token]
	utils.Mu.Unlock()
	user, err := mysql.GetUser(a)
	b, exists := vars["name_space"]
	if !exists {
		utils.WriteResponse(w, errors.New("name_space param is undefined"), nil, http.StatusBadRequest)
		utils.Logger.Error("can't perform operation",
			zap.Error(errors.New("name_space param is undefined")),
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
	if b, err := mysql.CheckPermissionsOnSpace(user.Id, space.Id); !b {
		if err != nil {
			utils.WriteResponse(w, err, nil, http.StatusInternalServerError)
		} else {
			utils.WriteResponse(w, errors.New("check rights on space"), nil, http.StatusMethodNotAllowed)
		}
		utils.Logger.Error("can't perform operation",
			zap.Error(err),
		)
		return
	}
	c, err := mysql.GetSpacePermissions(space.Id)
	if err != nil {
		utils.WriteResponse(w, err, nil, http.StatusInternalServerError)
		utils.Logger.Error("can't perform operation",
			zap.Error(err),
		)
		return
	}
	utils.WriteResponse(w, nil, c, http.StatusOK)
	utils.Logger.Info("success executing function",
		zap.String("url", r.URL.Path),
	)
}

func getAllTuples(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	r.ParseForm()
	token, _ := utils.HandleToken(r.Header.Get("Authorization"))
	utils.Mu.Lock()
	a, _ := utils.Cookies[token]
	utils.Mu.Unlock()
	user, _ := mysql.GetUser(a)
	b, exists := vars["name_space"]
	if !exists {
		utils.WriteResponse(w, errors.New("name_space param is undefined"), nil, http.StatusBadRequest)
		utils.Logger.Error("can't perform operation",
			zap.Error(errors.New("name_space param is undefined")),
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
	if !exists {
		utils.WriteResponse(w, errors.New("space with such name doesn't exist"), nil, http.StatusBadRequest)
		utils.Logger.Error("can't perform operation",
			zap.Error(errors.New("space with such name doesn't exist")),
		)
		return
	}
	if b, err := mysql.CheckPermissionsOnSpace(user.Id, space.Id); !b {
		if err != nil {
			utils.WriteResponse(w, err, nil, http.StatusInternalServerError)
		} else {
			utils.WriteResponse(w, errors.New("check rights on space"), nil, http.StatusMethodNotAllowed)
		}
		utils.Logger.Error("can't perform operation",
			zap.Error(err),
		)
		return
	}
	tuples, err := tarantool.SelectAllTuples(b, user.Id)
	if err != nil {
		utils.WriteResponse(w, err, nil, http.StatusInternalServerError)
		utils.Logger.Error("can't perform operation",
			zap.Error(err),
		)
		return
	}
	utils.WriteResponse(w, nil, tuples, http.StatusOK)
	utils.Logger.Info("success executing function",
		zap.String("url", r.URL.Path),
	)
}

func getTuple(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	r.ParseForm()
	token, _ := utils.HandleToken(r.Header.Get("Authorization"))
	utils.Mu.Lock()
	a, _ := utils.Cookies[token]
	utils.Mu.Unlock()
	user, _ := mysql.GetUser(a)
	b, exists := vars["name_space"]
	if !exists {
		utils.WriteResponse(w, errors.New("name_space param is undefined"), nil, http.StatusBadRequest)
		utils.Logger.Error("can't perform operation",
			zap.Error(errors.New("name_space param is undefined")),
		)
		return
	}
	c, exists := vars["id_tuple"]
	if !exists {
		utils.WriteResponse(w, errors.New("id_tuple param is undefined"), nil, http.StatusBadRequest)
		utils.Logger.Error("can't perform operation",
			zap.Error(errors.New("id_tuple param is undefined")),
		)
		return
	}
	id, err := strconv.ParseUint(c, 10, 64)
	if err != nil {
		utils.WriteResponse(w, err, nil, http.StatusBadRequest)
		utils.Logger.Error("can't perform operation",
		)
		return
	}
	space, err := mysql.GetSpace(b, user.Id)
	if err != nil {
		utils.WriteResponse(w, err, nil, http.StatusInternalServerError)
		utils.Logger.Error("can't perform operation",
		)
		return
	}
	if !exists {
		utils.WriteResponse(w, errors.New("space with such name doesn't exist"), nil, http.StatusInternalServerError)
		utils.Logger.Error("can't perform operation",
			zap.Error(errors.New("space with such name doesn't exist")),
		)
		return
	}
	if b, err := mysql.CheckPermissionsOnSpace(user.Id, space.Id); !b {
		if err != nil {
			utils.WriteResponse(w, err, nil, http.StatusInternalServerError)
		} else {
			utils.WriteResponse(w, errors.New("check rights on space"), nil, http.StatusMethodNotAllowed)
		}
		utils.Logger.Error("can't perform operation",
			zap.Error(err),
		)
		return
	}
	//// execute task of pool for access to tarantool -----------------------------------------------------
	//t := workerpool.TarantoolTask{"SelectTuple", id, b, user.Id, []interface{}{}}
	//workerpool.MainPool.Exec(workerpool.TarantoolTask(t))
	////---------------------------------------------------------------------------------------------------
	tuple, err := tarantool.SelectTuple(id, b, user.Id)
	if err != nil {
		utils.WriteResponse(w, err, nil, http.StatusInternalServerError)
		utils.Logger.Error("can't perform operation",
			zap.Error(err),
		)
		return
	}
	utils.WriteResponse(w, nil, tuple, http.StatusOK)
}