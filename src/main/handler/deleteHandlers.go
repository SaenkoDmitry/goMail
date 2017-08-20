package handler

import (
	"net/http"
	"github.com/gorilla/mux"
	"main/utils"
	"main/dbs/mysql"
	"strconv"
)

func deleteUser(w http.ResponseWriter, r *http.Request) {
	// -----------------------------------------------
}

func deletePermission(w http.ResponseWriter, r *http.Request) {
	// -----------------------------------------------
}

func deleteSpace(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	r.ParseForm()
	token := r.Header.Get("Authorization")[7:]
	a, exists := utils.Cookies[token]
	user, exs := mysql.GetUser(a)
	b := vars["name_space"]
	if !exists || !exs {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	s, err := strconv.ParseUint(b, 10, 64)
	checkErr(err)
	if !mysql.CheckPermissionsOnSpace(user.Id, s) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	mysql.DeleteSpace(b)
	// tarantool -------------------------------------------------------------
	//js, _ := json.Marshal(c)
	w.WriteHeader(http.StatusOK)
	//w.Write(js)
}

func deleteTuple(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	r.ParseForm()
	token := r.Header.Get("Authorization")[7:]
	a, exists := utils.Cookies[token]
	user, exs := mysql.GetUser(a)
	b := vars["name_space"]
	if !exists || !exs {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	s, err := strconv.ParseUint(b, 10, 64)
	checkErr(err)
	if !mysql.CheckPermissionsOnSpace(user.Id, s) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	// tarantool -------------------------------------------------------------
	//js, _ := json.Marshal(c)
	w.WriteHeader(http.StatusOK)
	//w.Write(js)
}