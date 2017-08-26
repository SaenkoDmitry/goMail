package handler

import (
	"net/http"
	"github.com/gorilla/mux"
	"main/utils"
	"main/dbs/mysql"
	"unicode/utf8"
	"main/dbs/tarantool"
	"strconv"
	"encoding/json"
	"fmt"
)

func deleteSpace(w http.ResponseWriter, r *http.Request) {
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
	mysql.DeleteSpace(b, user.Id)
	//mysql.DeletePermissionsOnSpace(b, user.Id)
	//tarantool.DeleteSpace(b, user.Id) // -------------------------------------------------
	//js, _ := json.Marshal(c)
	w.WriteHeader(http.StatusOK)
	//w.Write(js)
}

func deleteTuple(w http.ResponseWriter, r *http.Request) {
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
	s, err := strconv.ParseUint(c, 10, 64)
	checkErr(err)
	if !exists || !exs {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	space, exists := mysql.GetSpace(b, user.Id)
	if !exists || !mysql.CheckPermissionsOnSpace(user.Id, space.Id) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	tuple, _ := tarantool.SelectTuple(s, b, user.Id)
	_, res := tarantool.DeleteTuple(s, b, user.Id)
	if res == true {
		js, _ := json.Marshal(tuple)
		mysql.AddHistory(user.Id, space.Id, "deleted tuple : " + fmt.Sprintf("%s", js), "OK")
	}
	//js, _ := json.Marshal(c)
	w.WriteHeader(http.StatusOK)
	//w.Write(js)
}