package handler

import (
	"net/http"
	//"github.com/gorilla/mux"
	"main/dbs/mysql"
	"main/model"
	"encoding/json"
	"main/utils"
	"github.com/gorilla/mux"
	"strconv"
	"unicode/utf8"
	"main/dbs/tarantool"
)

func getAllUsers(w http.ResponseWriter, r *http.Request) {
	var a[]model.User
	a = mysql.SelectAllUser()
	js, err := json.MarshalIndent(a, "", "  ")
	checkErr(err)
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

func getAllSpaces(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	token := r.Header.Get("Authorization")
	if utf8.RuneCountInString(token) > 7 {
		token = token[7:]
	}
	a, exists := utils.Cookies[token]
	user, exs := mysql.GetUser(a)
	if !exists || !exs {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	b, _ := mysql.GetAllSpaces(user.Id)
	js, _ := json.MarshalIndent(b, "", " ")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

func getUserHistory(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	token := r.Header.Get("Authorization")
	if utf8.RuneCountInString(token) > 8 {
		token = token[7:]
	}
	a, exists := utils.Cookies[token]
	user, exs := mysql.GetUser(a)
	if !exists || !exs {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	b, _ := mysql.GetUserHistory(user.Id)
	js, _ := json.MarshalIndent(b, "", " ")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

func getSpaceHistory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	r.ParseForm()
	token := r.Header.Get("Authorization")
	if utf8.RuneCountInString(token) > 8 {
		token = token[7:]
	}
	a, exists := utils.Cookies[token]
	user, exs := mysql.GetUser(a)
	b := vars["name_space"]
	if !exists || !exs {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	space, exists := mysql.GetSpace(b, user.Id)
	if !exists || !mysql.CheckPermissionsOnSpace(user.Id, space.Id) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	c, _ := mysql.GetSpaceHistory(user.Id)
	js, _ := json.MarshalIndent(c, "", " ")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

func getAllSpacePermissions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	r.ParseForm()
	token := r.Header.Get("Authorization")
	if utf8.RuneCountInString(token) > 8 {
		token = token[7:]
	}
	a, exists := utils.Cookies[token]
	user, exs := mysql.GetUser(a)
	b := vars["name_space"]
	if !exists || !exs {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	space, exists := mysql.GetSpace(b, user.Id)
	if !exists || !mysql.CheckPermissionsOnSpace(user.Id, space.Id) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	c, _ := mysql.GetSpacePermissions(space.Id)
	js, _ := json.MarshalIndent(c, "", " ")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

func getAllTuples(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	r.ParseForm()
	token := r.Header.Get("Authorization")
	if utf8.RuneCountInString(token) > 8 {
		token = token[7:]
	}
	a, exists := utils.Cookies[token]
	user, exs := mysql.GetUser(a)
	b := vars["name_space"]
	if !exists || !exs {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	space, exists := mysql.GetSpace(b, user.Id)
	if !exists || !mysql.CheckPermissionsOnSpace(user.Id, space.Id) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	c, _ := tarantool.SelectAllTuples(b, user.Id)
	js, _ := json.MarshalIndent(c, "", " ")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

func getTuple(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	r.ParseForm()
	token := r.Header.Get("Authorization")
	if utf8.RuneCountInString(token) > 8 {
		token = token[7:]
	}
	a, exists := utils.Cookies[token]
	user, exs := mysql.GetUser(a)
	b := vars["name_space"]
	c := vars["id_tuple"]
	if !exists || !exs {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	id, err := strconv.ParseUint(c, 10, 64)
	checkErr(err)
	space, exists := mysql.GetSpace(b, user.Id)
	if !exists || !mysql.CheckPermissionsOnSpace(user.Id, space.Id) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	//// execute task of pool for access to tarantool -----------------------------------------------------
	//t := workerpool.TarantoolTask{"SelectTuple", id, b, user.Id, []interface{}{}}
	//workerpool.MainPool.Exec(workerpool.TarantoolTask(t))
	////---------------------------------------------------------------------------------------------------
	tuple, _ := tarantool.SelectTuple(id, b, user.Id)
	js, _ := json.MarshalIndent(tuple, "", " ")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}