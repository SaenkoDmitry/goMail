package handler

import (
	"encoding/json"
	"fmt"
	"main/dbs/mysql"
	"main/dbs/tarantool"
	"main/utils"
	"net/http"
	"github.com/gorilla/mux"
	"strconv"
	//"main/workerpool"
	"unicode/utf8"
	"io/ioutil"
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
		for k, v := range utils.Cookies {
			if v == a {
				token = k
				exists = true

			}
		}
		if exists != true {
			var err error
			token, err = utils.CreateToken(a, b)
			utils.Cookies[token] = a
			checkErr(err)
		}
	}
	w.WriteHeader(http.StatusOK)
	js, _ := json.MarshalIndent(token, "", " ")
	w.Write(js)
}

func addSpace(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	r.ParseForm()
	token := r.Header.Get("Authorization")
	if utf8.RuneCountInString(token) > 8 {
		token = token[7:]
	}
	a, exists := utils.Cookies[token]
	user, exs := mysql.GetUser(a)
	b := vars["name_space"]
	if exists && exs {
		mysql.AddSpace(b, user.Id)
		space, _ := mysql.GetSpace(b, user.Id)
		mysql.AddPermission(user.Id, space.Id)
		mysql.AddHistory(user.Id, space.Id, "added space : " + b, "OK")
		tarantool.CreateSpace(b, user.Id)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func addTuple(w http.ResponseWriter, r *http.Request) {
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
	err = json.Unmarshal([]byte(body), &data)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("data : ", data)

	a, exists := utils.Cookies[token]
	user, exs1 := mysql.GetUser(a)
	b := vars["name_space"]
	space, exs2 := mysql.GetSpace(b, user.Id)
	c := vars["id_tuple"]
	id, _ := strconv.ParseUint(c, 10, 64)
	if !exs1 || !exists {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if !exs2 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if mysql.CheckPermissionsOnSpace(user.Id, space.Id) {
		fmt.Println(c)
		mysql.AddHistory(user.Id, space.Id, "", "")
		//// execute task of pool for access to tarantool -----------------------------------------------------
		//t := workerpool.TarantoolTask{"InsertTuple", id, b, user.Id, data}
		//workerpool.MainPool.Exec(workerpool.TarantoolTask(t))
		////---------------------------------------------------------------------------------------------------
		tarantool.InsertTuple(id, b, user.Id, data)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func addPermission(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	r.ParseForm()
	token := r.Header.Get("Authorization")
	if utf8.RuneCountInString(token) > 8 {
		token = token[7:]
	}
	a, exists := utils.Cookies[token]
	user, exs1 := mysql.GetUser(a)
	b := vars["name"]
	c := vars["name_space"]
	user2, exs2 := mysql.GetUser(b)
	space, exs3 := mysql.GetSpace(c, user.Id)
	if !exs1 || !exists {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if !exs2 || !exs3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if mysql.CheckPermissionsOnSpace(user.Id, space.Id) {
		mysql.AddPermission(user2.Id, space.Id)
		mysql.AddHistory(user.Id, space.Id, "added permission for " + user2.Name + " on space " + c, "OK")
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(http.StatusOK)
}
