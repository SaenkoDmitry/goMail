package handler

import (
	"net/http"
	"github.com/gorilla/mux"
	"goMail/dbs/mysql"
	"goMail/dbs/tarantool"
	"goMail/utils"
	"strconv"
	"encoding/json"
	"errors"
	"go.uber.org/zap"
)

func loginUser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	a := r.Form.Get("user")
	b := r.Form.Get("password")
	hash := utils.HashPassword(b)

	var token string

	c, err := mysql.GetUser(a)
	if err != nil {
		mysql.AddUser(a, hash)
		token, err = utils.CreateToken(a, b)
		utils.Mu.Lock()
		utils.Cookies[token] = a
		utils.Mu.Unlock()
		if err != nil {
			utils.WriteResponse(w, err, nil, http.StatusInternalServerError)
			utils.Logger.Error("cannot create token",
				zap.Error(err),
			)
			return
		}
	} else {
		if hash != c.HashPassword {
			utils.WriteResponse(w, errors.New("invalid password"), nil, http.StatusUnauthorized)
			utils.Logger.Error("can't perform operation",
				zap.Error(errors.New("invalid password")),
			)
			return
		}
		var exists bool
		utils.Mu.Lock()
		for k, v := range utils.Cookies {
			if v == a {
				token = k
				exists = true

			}
		}
		utils.Mu.Unlock()
		if !exists {
			var err error
			token, err = utils.CreateToken(a, b)
			utils.Mu.Lock()
			utils.Cookies[token] = a
			utils.Mu.Unlock()
			if err != nil {
				utils.WriteResponse(w, err, nil, http.StatusInternalServerError)
				utils.Logger.Error("can't perform operation",
					zap.Error(err),
				)
				return
			}
		}
	}

	resp := make(map[string]interface{})
	resp["token"] = token
	utils.WriteResponse(w, nil, resp, http.StatusOK)
	utils.Logger.Info("success executing function",
		zap.String("url", r.URL.Path),
	)
}

func addSpace(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	r.ParseForm()
	token, _ := utils.HandleToken(r.Header.Get("Authorization"))
	utils.Mu.Lock()
	username, _ := utils.Cookies[token]
	utils.Mu.Unlock()
	user, err := mysql.GetUser(username)
	if err != nil {
		utils.WriteResponse(w, err, nil, http.StatusInternalServerError)
		utils.Logger.Error("can't perform operation",
			zap.Error(err),
		)
		return
	}

	b, exists := vars["name_space"]
	if !exists {
		utils.WriteResponse(w, errors.New("name_space param is undefined"), nil, http.StatusInternalServerError)
		utils.Logger.Error("can't perform operation",
			zap.Error(errors.New("name_space param is undefined")),
		)
		return
	}

	_, err = mysql.AddSpace(b, user.Id)
	if err != nil {
		utils.WriteResponse(w, err, nil, http.StatusInternalServerError)
		utils.Logger.Error("can't perform operation",
			zap.Error(err),
		)
		return
	}

	space, err := mysql.GetSpace(b, user.Id)
	if err != nil {
		utils.WriteResponse(w, err, nil, http.StatusBadRequest)
		utils.Logger.Error("can't perform operation",
			zap.Error(err),
		)
		return
	}
	if _, err := mysql.AddPermission(user.Id, space.Id); err != nil {
		utils.WriteResponse(w, err, nil, http.StatusInternalServerError)
		utils.Logger.Error("can't perform operation",
			zap.Error(err),
		)
		return
	}

	if _, err := mysql.AddHistory(user.Id, space.Id, "added space : "+b, "OK"); err != nil {
		utils.WriteResponse(w, err, nil, http.StatusInternalServerError)
		utils.Logger.Error("can't perform operation",
			zap.Error(err),
		)
		return
	}

	resp, err := tarantool.CreateSpace(b, user.Id)
	if err != nil {
		utils.WriteResponse(w, err, nil, http.StatusInternalServerError)
		utils.Logger.Error("can't perform operation",
			zap.Error(err),
		)
		return
	}
	utils.WriteResponse(w, nil, resp, http.StatusOK)
	utils.Logger.Info("success executing function",
		zap.String("url", r.URL.Path),
	)
}

func addTuple(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var data []interface{}
	decoder.Decode(&data)

	vars := mux.Vars(r)
	r.ParseForm()
	token, _ := utils.HandleToken(r.Header.Get("Authorization"))

	utils.Mu.Lock()
	a, _ := utils.Cookies[token]
	utils.Mu.Unlock()

	user, err := mysql.GetUser(a)
	if err != nil {
		utils.WriteResponse(w, err, nil, http.StatusInternalServerError)
		utils.Logger.Error("can't perform operation",
			zap.Error(err),
		)
		return
	}
	b, exists := vars["name_space"]
	if !exists {
		utils.WriteResponse(w, errors.New("name_space param is undefined"), nil, http.StatusBadRequest)
		utils.Logger.Error("can't perform operation",
		)
		return
	}
	space, err := mysql.GetSpace(b, user.Id)
	if err != nil {
		utils.WriteResponse(w, err, nil, http.StatusBadRequest)
		utils.Logger.Error("can't perform operation",
			zap.Error(err),
		)
		return
	}
	c, exists := vars["id_tuple"]
	if !exists {
		utils.WriteResponse(w, errors.New("id_tuple param is undefined"), nil, http.StatusBadRequest)
		utils.Logger.Error("can't perform operation",
			zap.Error(err),
		)
		return
	}
	id, err := strconv.ParseUint(c, 10, 64)
	if err != nil {
		utils.WriteResponse(w, errors.New("id tuple is not valid"), nil, http.StatusBadRequest)
		utils.Logger.Error("can't perform operation",
		)
		return
	}
	if _, err := mysql.CheckPermissionsOnSpace(user.Id, space.Id); err == nil {
		mysql.AddHistory(user.Id, space.Id, "", "")
		//// execute task of pool for access to tarantool -----------------------------------------------------
		//t := workerpool.TarantoolTask{"InsertTuple", id, b, user.Id, data}
		//workerpool.MainPool.Exec(workerpool.TarantoolTask(t))
		////---------------------------------------------------------------------------------------------------
		resp, err := tarantool.InsertTuple(id, b, user.Id, data)
		if err != nil {
			utils.WriteResponse(w, err, nil, http.StatusInternalServerError)
			utils.Logger.Error("can't perform operation",
				zap.Error(err),
			)
			return
		}
		utils.WriteResponse(w, nil, resp, http.StatusOK)
		utils.Logger.Info("success executing function",
			zap.String("url", r.URL.Path),
		)
		return
	} else {
		utils.WriteResponse(w, err, nil, http.StatusMethodNotAllowed)
		utils.Logger.Error("can't perform operation",
			zap.Error(err),
		)
		return
	}
}

func addPermission(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	r.ParseForm()
	token, _ := utils.HandleToken(r.Header.Get("Authorization"))
	utils.Mu.Lock()
	a, _ := utils.Cookies[token]
	utils.Mu.Unlock()
	user, err := mysql.GetUser(a)
	if err != nil {
		utils.WriteResponse(w, err, nil, http.StatusInternalServerError)
		utils.Logger.Error("can't perform operation",
			zap.Error(err),
		)
		return
	}
	b, exists := vars["name"]
	if !exists {
		utils.WriteResponse(w, errors.New("user name param is undefined"), nil, http.StatusBadRequest)
		utils.Logger.Error("can't perform operation",
			zap.Error(errors.New("user name param is undefined")),
		)
		return
	}
	c, exists := vars["name_space"]
	if !exists {
		utils.WriteResponse(w, errors.New("user name param is undefined"), nil, http.StatusBadRequest)
		utils.Logger.Error("can't perform operation",
			zap.Error(errors.New("user name param is undefined")),
		)
		return
	}

	user2, err := mysql.GetUser(b)
	if err != nil {
		utils.WriteResponse(w, err, nil, http.StatusBadRequest)
		utils.Logger.Error("can't perform operation",
			zap.Error(err),
		)
		return
	}
	space, err := mysql.GetSpace(c, user.Id)
	if err != nil {
		utils.WriteResponse(w, err, nil, http.StatusBadRequest)
		utils.Logger.Error("can't perform operation",
			zap.Error(err),
		)
		return
	}

	if _, err := mysql.CheckPermissionsOnSpace(user.Id, space.Id); err == nil {
		resp, err := mysql.AddPermission(user2.Id, space.Id)
		if err != nil {
			utils.WriteResponse(w, err, nil, http.StatusInternalServerError)
			utils.Logger.Error("can't perform operation",
				zap.Error(err),
			)
			return
		}
		_, err = mysql.AddHistory(user.Id, space.Id, "added permission for "+user2.Name+" on space "+c, "OK")
		if err != nil {
			utils.WriteResponse(w, err, nil, http.StatusInternalServerError)
			utils.Logger.Error("can't perform operation",
				zap.Error(err),
			)
			return
		}
		utils.WriteResponse(w, nil, resp, http.StatusInternalServerError)
		utils.Logger.Info("success executing function",
			zap.String("url", r.URL.Path),
		)
		return
	} else {
		utils.WriteResponse(w, err, nil, http.StatusMethodNotAllowed)
		utils.Logger.Error("can't perform operation",
			zap.Error(err),
		)
		return
	}
}