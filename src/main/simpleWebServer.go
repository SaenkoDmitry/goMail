package main

import (
	"fmt"
	"net/http"
	"strings"
	"log"
	"unicode/utf8"
	//"github.com/tarantool/go-tarantool"
)

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	fmt.Println(`_____________________________________
	sayHello`)
	r.ParseForm()  // parse arguments, you have to call this by yourself
	fmt.Println(r.Form)  // print form information in server side
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	//fmt.Fprintf(w, "Hello astaxie!") // send data to client side
}

func register(w http.ResponseWriter, r *http.Request) {
	fmt.Println(`_____________________________________
	register`)
	w.Header().Set("Access-Control-Allow-Origin", "*")     //<---------- here!
	//URISegments := strings.Split(r.URL.Path[:utf8.RuneCountInString(r.URL.Path) - 7] + "mainPage", "/")
	//w.Write([]byte(URISegments[1]))
	r.ParseForm()  // parse arguments, you have to call this by yourself
	credentials := r.Form
	var a, _ = credentials["user"]
	fmt.Println("user : ", a)
	var b, _ = credentials["password"]
	fmt.Println("password : ", b)
	var c, _ = credentials["second-password"]
	fmt.Println("second-password : ", c)
	var d, _ = credentials["address"]
	fmt.Println("address : ", d)
	//opts := tarantool.Opts{User: a[0], Pass: b[0]}
	//conn, err := tarantool.Connect("127.0.0.1:3302", opts)
	//if err != nil {
	//	fmt.Println("Connection refused: %s", err.Error())
	//}
	//resp, err := conn.Insert(10, []interface{}{99999, "BB"})
	//if err != nil {
	//	fmt.Println("Error", err)
	//	fmt.Println("Code", resp.Code)
	//}
}

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println(`_____________________________________
	login`)
	w.Header().Set("Access-Control-Allow-Origin", "*")     //<---------- here!
	URISegments := strings.Split(r.URL.Path[:(utf8.RuneCountInString(r.URL.Path) - 7)] + "mainPage", "/")
	w.Write([]byte(URISegments[1]))
	r.ParseForm()
	credentials := r.Form
	var a, _ = credentials["user"]
	fmt.Println("user : ", a)
	var b, _ = credentials["password"]
	fmt.Println("password : ", b)
	//opts := tarantool.Opts{User: a[0], Pass: b[0]}
	//conn, err := tarantool.Connect("127.0.0.1:3302", opts)
	//// conn, err := tarantool.Connect("/path/to/tarantool.socket", opts)
	//if err != nil {
	//	fmt.Println("Connection refused: %s", err.Error())
	//}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", sayhelloName) // set router
	mux.HandleFunc("/sign-in", login) // set router
	mux.HandleFunc("/sign-up", register) // set router
	err := http.ListenAndServe(":9090", mux) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}