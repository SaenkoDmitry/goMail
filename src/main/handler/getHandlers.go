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
)

func getAllUsers(w http.ResponseWriter, r *http.Request) {
	db := mysql.GetDb()
	rows, err := db.Query("SELECT * FROM users")
	var a[]model.User
	var temp model.User
	for rows.Next() {
		err = rows.Scan(&temp.Id, &temp.Name, &temp.HashPassword)
		a = append(a, temp)
		checkErr(err)
	}
	js, err := json.Marshal(a)
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

func getAllSpaces(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	token := r.Header.Get("Authorization")[7:]
	a, exists := utils.Cookies[token]
	user, exs := mysql.GetUser(a)
	if !exists || !exs {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	b, _ := mysql.GetAllSpaces(user.Id)
	js, _ := json.Marshal(b)
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

func getUserHistory(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	token := r.Header.Get("Authorization")[7:]
	a, exists := utils.Cookies[token]
	user, exs := mysql.GetUser(a)
	if !exists || !exs {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	b, _ := mysql.GetUserHistory(user.Id)
	js, _ := json.Marshal(b)
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

func getSpaceHistory(w http.ResponseWriter, r *http.Request) {
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
	c, _ := mysql.GetSpaceHistory(user.Id)
	js, _ := json.Marshal(c)
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

func getAllSpacePermissions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	r.ParseForm()
	token := r.Header.Get("Authorization")[7:]
	a, exists := utils.Cookies[token]
	user, exs := mysql.GetUser(a)
	b := vars["name_sgit pace"]
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
	c, _ := mysql.GetSpacePermissions(s)
	js, _ := json.Marshal(c)
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

func getAllTuples(w http.ResponseWriter, r *http.Request) {
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
	//opts := tarantool.Opts{User: a[0], Pass: b[0]}
	//conn, err := tarantool.Connect("127.0.0.1:3302", opts)
	//if err != nil {
	//	fmt.Println("Connection refused: %s", err.Error())
	//}
	//resp, err := conn.Insert(10, []interface{}{99999, "BB"})
	//if err != nil {
	//	fmt.Println("Error", err)
	//	fmt.Println("Code", resp.Code)
	// -----------------------------------------------------------------------
	//js, _ := json.Marshal(c)
	w.WriteHeader(http.StatusOK)
	//w.Write(js)
}

func getTuple(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	r.ParseForm()
	token := r.Header.Get("Authorization")[7:]
	a, exists := utils.Cookies[token]
	user, exs := mysql.GetUser(a)
	b := vars["name_space"]
	//c := vars["name_tuple"]
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