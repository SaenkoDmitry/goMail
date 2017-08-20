package main

import (
	"fmt"
	//"net/http"
	//"github.com/tarantool/go-tarantool"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	//"main/utils"
	//"github.com/gorilla/mux"
	"github.com/gorilla/mux"
	"main/handler"
	"net/http"
	"log"
	//"main/dbs/mysql"
)

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		//panic(err)
	}
}

func selectFromUsers(db *sql.DB, dbname string) {
	rows, err := db.Query("SELECT * FROM users")
	checkErr(err)

	for rows.Next() {
		var username string
		var token string
		err = rows.Scan(&username, &token)
		checkErr(err)
		fmt.Println("user: " + username + "; token: " + token)
	}
}

func insertTo(db *sql.DB, dbname string) {
	stmt, err := db.Prepare("INSERT users SET username=?, token=?")
	checkErr(err)

	res, err := stmt.Exec("user5", "12345")
	checkErr(err)

	id, err := res.LastInsertId()
	checkErr(err)

	fmt.Println(id)
}

func main() {

	router := mux.NewRouter()
	handler.InitHandlers(router)

	err := http.ListenAndServe(":9090", router) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	//token, _ := utils.CreateToken("user", "12345")
	//fmt.Println(token)
	//res := utils.ParseToken(token)
	//fmt.Println(res)

	//db, err := sql.Open("mysql", "adminGo:gomail@tcp(localhost:3306)/tarantooldbs?charset=utf8")
	//defer db.Close()
	//if err != nil {
	//	//handling errors------------------------------------------------------!!!
	//}
	//selectFromUsers(db, "users")
	////insertTo(db, "users")
	//selectFromUsers(db, "users")
}