package handler

import (
	"net/http"
	"main/utils"
	"encoding/json"
	"main/dbs/mysql"
	"github.com/gorilla/mux"
	"fmt"
)

func loginUser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	credentials := r.Form
	a := credentials["user"][0]
	b := credentials["password"][0]
	hash := utils.HashPassword(b)
	var token string
	c, exists := mysql.GetUser(a)
	if exists != true {
		var err error
		mysql.AddUser(a, hash)
		token, err = utils.CreateToken(a, b)
		utils.Cookies[token] = a
		checkErr(err)
	} else {
		if hash != c.HashPassword {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var exists bool
		token, exists = utils.Cookies[hash]
		if exists != true {
			var err error
			token, err = utils.CreateToken(a, b)
			utils.Cookies[token] = a
			checkErr(err)
		}
	}
	w.WriteHeader(http.StatusOK)
	js, _ := json.Marshal(token)
	w.Write(js)
}

func addSpace(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	r.ParseForm()
	token := r.Header.Get("Authorization")[7:]
	a, exists := utils.Cookies[token]
	user, exs := mysql.GetUser(a)
	b := vars["name_space"]
	if exists && exs {
		mysql.AddSpace(b, user.Id)
		//add tarantool space --------------------------------------------------------------------------
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func addTuple(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	r.ParseForm()
	token := r.Header.Get("Authorization")[7:]
	a, exists := utils.Cookies[token]
	user, exs1 := mysql.GetUser(a)
	b := vars["name_space"]
	space, exs2 := mysql.GetSpace(b, user.Id)
	c := vars["name_tuple"]
	if exists && exs1 && exs2 && mysql.CheckPermissionsOnSpace(user.Id, space.Id) {
		fmt.Println(c)
		//add tarantool tuple ----------------------------------------------------------------------------
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func addPermission(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	r.ParseForm()
	token := r.Header.Get("Authorization")[7:]
	a, exists := utils.Cookies[token]
	user, exs1 := mysql.GetUser(a)
	b := vars["name"]
	space, exs2 := mysql.GetSpace(b, user.Id)
	c := vars["name_space"]
	if exists && exs1 && exs2 && mysql.CheckPermissionsOnSpace(user.Id, space.Id) {
		fmt.Println(c)
		//add tarantool tuple ----------------------------------------------------------------------------
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
}