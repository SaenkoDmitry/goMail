package main

import (
	"fmt"
	//"net/http"
	//"github.com/tarantool/go-tarantool"

	_ "github.com/go-sql-driver/mysql"
	//"main/utils"
	//"github.com/gorilla/mux"
	"log"
	"main/handler"
	"net/http"

	"github.com/gorilla/mux"
	//"main/dbs/mysql"
	//"main/workerpool"
	"main/workerpool"
)

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		//panic(err)
	}
}

func main() {

	router := mux.NewRouter()
	handler.InitHandlers(router)

	workerpool.MainPool = workerpool.NewPool(5) //create pool
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
