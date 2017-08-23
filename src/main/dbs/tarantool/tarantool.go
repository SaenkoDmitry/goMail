package tarantool

import (
	"fmt"
	"github.com/tarantool/go-tarantool"
	"time"
)

var spaceNo uint32 = uint32(999)
//var indexNo uint32 = uint32(0)


var tarantoolConn *tarantool.Connection

func init() {
	tarantoolConn = InitTarantool()
}


func InitTarantool() *tarantool.Connection {

	server := "127.0.0.1:3302"
	opts := tarantool.Opts{
		Timeout:       500 * time.Millisecond,
		Reconnect:     1 * time.Second,
		MaxReconnects: 3,
		User:          "test",
		Pass:          "12345",
	}
	conn, err := tarantool.Connect(server, opts)
	if err != nil {
		fmt.Println("Connection refused: %s", err.Error())
	}
	return conn
}

func convertToNameInTarantool(name string, user_id uint64) string {
	return name + fmt.Sprintf("%v", user_id)
}

func DeleteTuple(tuple_id uint32, name_space string, user_id uint64) {
	name_spaceT := convertToNameInTarantool(name_space, user_id)
	name_spaceT = "examples" // temporary -------------------------------------------
	resp, err := InitTarantool().Delete(name_spaceT, "primary", []interface{}{uint(tuple_id)})
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(resp.Data)
}

func SelectTuple(tuple_id uint32, name_space string, user_id uint64) {
	name_spaceT := convertToNameInTarantool(name_space, user_id)
	name_spaceT = "examples" // temporary -------------------------------------------

	resp, err := tarantoolConn.Select(name_spaceT, "primary", 0, 1, tarantool.IterEq, []interface{}{uint(tuple_id)})
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(resp.Data)
}

func InsertTuple(tuple_id uint32, name_space string, user_id uint64, data interface{}) {
	name_spaceT := convertToNameInTarantool(name_space, user_id)
	name_spaceT = "examples" // temporary -------------------------------------------

	_, err := tarantoolConn.Insert(name_spaceT, []interface{}{tuple_id, data})
	if err != nil {
		fmt.Println(err.Error())
	}
}
