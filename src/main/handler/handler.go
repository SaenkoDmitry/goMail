package handler

import (
	"fmt"
	"github.com/gorilla/mux"
)

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		//panic(err)
	}
}



func InitHandlers(router *mux.Router) {

	router.HandleFunc("/users", getAllUsers).Methods("GET") // set router
	router.HandleFunc("/spaces", getAllSpaces).Methods("GET")
	router.HandleFunc("/history", getUserHistory).Methods("GET")
	router.HandleFunc("/space/{name_space}/history", getSpaceHistory).Methods("GET")
	router.HandleFunc("/spaces/{name_space}/permissions", getAllSpacePermissions).Methods("GET")
	router.HandleFunc("/spaces/{name_space}/tuples", getAllTuples).Methods("GET")
	router.HandleFunc("/spaces/{name_space}/tuples/{id_tuple}", getTuple).Methods("GET")

	router.HandleFunc("/users", loginUser).Methods("POST")
	router.HandleFunc("/spaces/{name_space}", addSpace).Methods("POST")
	router.HandleFunc("/spaces/{name_space}/tuples/{id_tuple}", addTuple).Methods("POST")
	router.HandleFunc("/users/{name}/spaces/{name_space}/permissions", addPermission).Methods("POST")

	router.HandleFunc("/spaces/{name_space}/tuples/{id_tuple}", updateTuple).Methods("PUT")
	//router.HandleFunc("/users/{name}/spaces/{name_space}/permissions/{value}", updatePermissions).Methods("PUT")

	//router.HandleFunc("/users", deleteUser).Methods("DELETE")
	//router.HandleFunc("/users/{name}/permissions", deletePermission).Methods("DELETE")
	//router.HandleFunc("/spaces/{name_space}", deleteSpace).Methods("DELETE")
	router.HandleFunc("/spaces/{name_space}/tuples/{id_tuple}", deleteTuple).Methods("DELETE")
}
