package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"goMail/handler"
	"net/http"
	"github.com/gorilla/mux"
	"goMail/workerpool"
	"go.uber.org/zap"
	"time"
	"goMail/utils"
	"goMail/dbs/mysql"
	"goMail/dbs/tarantool"
)

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	// initial zap logging
	defer utils.Logger.Sync()

	// initial mysql connection
	mysql.MysqlConn = mysql.InitMysql()
	defer mysql.MysqlConn.Close()

	// initial tarantool connection
	tarantool.TarantoolConn = tarantool.InitTarantool()
	defer tarantool.TarantoolConn.Close()

	router := mux.NewRouter()
	handler.InitHandlers(router)

	workerpool.MainPool = workerpool.NewPool(10) //create pool
	utils.Logger.Info("main is running",
		zap.String("url", "localhost:9090"),
		zap.Duration("backoff", time.Second),
	)

	srv := &http.Server{
		Addr:           ":9090",
		Handler:        router,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
	}
	err := srv.ListenAndServe()
	if err != nil {
		utils.Logger.Error("error occured while starting main: ",
			zap.Error(err))
		log.Fatal("ListenAndServe error: ", err)
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
