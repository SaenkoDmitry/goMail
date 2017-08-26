package handler

import (
	"net/http"
	"github.com/gorilla/mux"
	"main/utils"
	"fmt"
	"main/dbs/mysql"
	"unicode/utf8"
	"strconv"
	"main/dbs/tarantool"
	"encoding/json"
	"io/ioutil"
)

func updateTuple(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	r.ParseForm()
	token := r.Header.Get("Authorization")
	if utf8.RuneCountInString(token) > 8 {
		token = token[7:]
	}
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	var data []interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println(err)
	}
	a, exists := utils.Cookies[token]
	user, exs1 := mysql.GetUser(a)
	if !exists || !exs1 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	b := vars["name_space"]
	c := vars["id_tuple"]
	s, err := strconv.ParseUint(c, 10, 64)
	checkErr(err)
	space, exists := mysql.GetSpace(b, user.Id)
	if !exists || !mysql.CheckPermissionsOnSpace(user.Id, space.Id) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	tarantool.UpdateTuple(s, b, user.Id, data)
	w.WriteHeader(http.StatusOK)
}