package tarantool

import (
	"fmt"
	"github.com/tarantool/go-tarantool"
	"main/utils"
	"strconv"
)

var tarantoolConn *tarantool.Connection

func init() {
	tarantoolConn = InitTarantool()
}


func InitTarantool() *tarantool.Connection {
	a := utils.GetTarantool()
	server := a.Host + ":" + strconv.Itoa(a.Port)
	opts := tarantool.Opts{
		Timeout:       a.Timeout,
		Reconnect:     a.Reconnect,
		MaxReconnects: a.Maxreconnects,
		User:          a.Username,
		Pass:          a.Password,
	}
	conn, err := tarantool.Connect(server, opts)
	if err != nil {
		fmt.Println("Connection refused: %s", err.Error())
	}
	return conn
}

func CreateSpace(name_space string, user_id uint64) ([]interface{}, bool) {
	name_space = utils.ConvertToNameInDB(name_space, user_id)
	fmt.Println("name : " + name_space)
	fmt.Println("string : " + "box.schema.space.create('" + name_space + "')\n" +
		"box.space." + name_space + ":create_index('primary', {type = 'hash', parts = {1, 'NUM'}})\n")
	resp, err := tarantoolConn.Eval("box.schema.space.create('" + name_space + "')\n" +
		"box.space." + name_space + ":create_index('primary', {type = 'hash', parts = {1, 'NUM'}})\n",[] interface{}{})
	if err != nil {
		fmt.Println(err)
		return nil, false
	}
	return resp.Data, true
}

//func SelectAllTuples(name_space string, user_id uint64) ([]interface{}, bool) {
//	name_space = utils.ConvertToNameInDB(name_space, user_id)
//	resp, err := tarantoolConn.Eval("box.space." + name_space + ":select{25}",[] interface{}{})
//	if err != nil {
//		fmt.Println(err)
//		return nil, false
//	}
//	return resp.Data, true
//}

//func SelectSpace(name_space string, user_id uint64) ([]interface{}, bool)  {
//	name_space = utils.ConvertToNameInDB(name_space, user_id)
//	resp, err := tarantoolConn.Select(name_space, "primary", 0, 1, tarantool.IterAll, []interface{}{})
//	if err != nil {
//		fmt.Println(err.Error())
//		return nil, false
//	}
//	return resp.Data, true
//}

func DeleteSpace(name_space string, user_id uint64) ([]interface{}, bool)  {
	name_space = utils.ConvertToNameInDB(name_space, user_id)
	resp, err := tarantoolConn.Delete(name_space, "primary", []interface{}{})
	if err != nil {
		fmt.Println(err.Error())
		return nil, false
	}
	return resp.Data, true
}

func DeleteTuple(tuple_id uint64, name_space string, user_id uint64) ([]interface{}, bool) {
	name_space = utils.ConvertToNameInDB(name_space, user_id)
	resp, err := InitTarantool().Delete(name_space, "primary", []interface{}{uint(tuple_id)})
	if err != nil {
		fmt.Println(err.Error())
		return nil, false
	}
	return resp.Data, true
}

func SelectTuple(tarantoolConn *tarantool.Connection, tuple_id uint64, name_space string, user_id uint64) ([]interface{}, bool) {
	name_space = utils.ConvertToNameInDB(name_space, user_id)
	resp, err := tarantoolConn.Select(name_space, "primary", 0, 1, tarantool.IterEq, []interface{}{uint(tuple_id)})
	if err != nil {
		fmt.Println(err.Error())
		return nil, false
	}
	return resp.Data, true
}

func InsertTuple(tarantoolConn *tarantool.Connection, tuple_id uint64, name_space string, user_id uint64, data interface{}) ([]interface{}, bool) {
	name_space = utils.ConvertToNameInDB(name_space, user_id)
	resp, err := tarantoolConn.Insert(name_space, []interface{}{tuple_id, data})
	if err != nil {
		fmt.Println(err.Error())
		return nil, false
	}
	return resp.Data, true
}

func UpdateTuple(tuple_id uint64, name_space string, user_id uint64, data interface{}) ([]interface{}, bool) {
	name_space = utils.ConvertToNameInDB(name_space, user_id)
	resp, err := tarantoolConn.Replace(name_space, []interface{}{tuple_id, data})
	if err != nil {
		fmt.Println(err.Error())
		return nil, false
	}
	return resp.Data, true
}

func SelectAllTuples(name_space string, user_id uint64) ([]interface{}, bool) {
	name_space = utils.ConvertToNameInDB(name_space, user_id)
	resp, err := tarantoolConn.Select(name_space, "primary", 0, 10, tarantool.IterAll, []interface{}{})
	if err != nil {
		fmt.Println(err.Error())
		return nil, false
	}
	return resp.Data, true
}