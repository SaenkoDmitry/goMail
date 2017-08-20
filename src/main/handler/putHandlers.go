package handler

import (
	"net/http"
	"github.com/gorilla/mux"
	"main/utils"
	"fmt"
	"main/dbs/mysql"
)

func updateTuple(w http.ResponseWriter, r *http.Request) {
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
		//update tarantool tuple ----------------------------------------------------------------------------
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
}

//func updatePermissions(w http.ResponseWriter, r *http.Request) {
//	vars := mux.Vars(r)
//	r.ParseForm()
//	token := r.Header.Get("Authorization")[7:]
//	a, exists := utils.Cookies[token]
//	if exists && a == vars["name"] {
//		db := mysql.GetDb()
//		rows, err := db.Query("SELECT * FROM users where name=?", a)
//		if rows.Next() {
//			var c model.User
//			err = rows.Scan(&c.Id, &c.Name, &c.HashPassword)
//			checkErr(err)
//		}
//
//		//stmt, err := db.Prepare("INSERT permissions SET user_id=?, space_id=?, rights=?")
//		checkErr(err)
//
//		//_, err = stmt.Exec(a, utils.HashPassword(b), "read")
//		//checkErr(err)
//		//_, err = stmt.Exec(a, utils.HashPassword(b), "write")
//		//checkErr(err)
//	} else {
//		w.WriteHeader(http.StatusUnauthorized)
//		return
//	}
//}
