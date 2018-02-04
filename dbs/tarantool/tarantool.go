package tarantool

import (
	"github.com/tarantool/go-tarantool"
	"goMail/utils"
	"strconv"
	"go.uber.org/zap"
)

var TarantoolConn *tarantool.Connection

func init() {
	TarantoolConn = InitTarantool()
}


func InitTarantool() (db *tarantool.Connection) {
	a := utils.GetTarantool()
	server := a.Host + ":" + strconv.Itoa(a.Port)
	opts := tarantool.Opts{
		Timeout:       a.Timeout,
		Reconnect:     a.Reconnect,
		MaxReconnects: a.Maxreconnects,
		User:          a.Username,
		Pass:          a.Password,
	}
	db, err := tarantool.Connect(server, opts)
	if err != nil {
		utils.Logger.Info("cannot open tarantool connection",
			zap.Error(err),
		)
		return
	}
	return
}

func CreateSpace(name_space string, user_id uint64) (resp []interface{}, err error) {
	name_space = utils.ConvertToNameInDB(name_space, user_id)
	query1 := "box.schema.space.create('" + name_space + "')\n"
	query2 := "box.space." + name_space + ":create_index('primary', {type = 'hash', parts = {1, 'NUM'}})"

	res, err := TarantoolConn.Eval(query1, []interface{}{})
	if err != nil {
		utils.Logger.Info("cannot eval query",
			zap.String("query", query1),
			zap.Error(err),
		)
		return
	}

	res, err = TarantoolConn.Eval(query2, []interface{}{})
	if err != nil {
		utils.Logger.Info("cannot eval query",
			zap.String("query", query2),
			zap.Error(err),
		)
		return
	}
	resp = res.Data
	return
}

func DeleteTuple(tuple_id uint64, name_space string, user_id uint64) (resp []interface{}, err error) {
	name_space = utils.ConvertToNameInDB(name_space, user_id)
	res, err := TarantoolConn.Delete(name_space, "primary", []interface{}{uint(tuple_id)})
	if err != nil {
		utils.Logger.Info("cannot delete tuple",
			zap.Uint64("tuple_id", tuple_id),
			zap.Error(err),
		)
		return
	}
	resp = res.Data
	return
}

func DeleteSpace(name_space string, user_id uint64) (resp []interface{}, err error) {
	name_space = utils.ConvertToNameInDB(name_space, user_id)
	res, err := TarantoolConn.Eval("box.schema.space." + name_space + ":drop()", []interface{}{})
	if err != nil {
		utils.Logger.Info("cannot delete space",
			zap.String("name_space", name_space),
			zap.Error(err),
		)
		return
	}
	resp = res.Data
	return
}

func SelectTuple(tuple_id uint64, name_space string, user_id uint64) (resp []interface{}, err error) {
	name_space = utils.ConvertToNameInDB(name_space, user_id)
	res, err := TarantoolConn.Select(name_space, "primary", 0, 1, tarantool.IterEq, []interface{}{uint(tuple_id)})
	if err != nil {
		utils.Logger.Info("cannot select tuple",
			zap.Uint64("tuple_id", tuple_id),
			zap.Error(err),
		)
		return
	}
	resp = res.Data
	return
}

func InsertTuple(tuple_id uint64, name_space string, user_id uint64, data interface{}) (resp []interface{}, err error) {
	name_space = utils.ConvertToNameInDB(name_space, user_id)
	res, err := TarantoolConn.Insert(name_space, []interface{}{tuple_id, data})
	if err != nil {
		utils.Logger.Info("cannot insert tuple",
			zap.Uint64("tuple_id", tuple_id),
			zap.Error(err),
		)
		return
	}
	resp = res.Data
	return
}

func UpdateTuple(tuple_id uint64, name_space string, user_id uint64, data interface{}) (resp []interface{}, ok bool, err error) {
	name_space = utils.ConvertToNameInDB(name_space, user_id)
	res, err := TarantoolConn.Replace(name_space, []interface{}{tuple_id, data})
	if err != nil {
		utils.Logger.Info("cannot update tuple",
			zap.Uint64("tuple_id", tuple_id),
			zap.Error(err),
		)
		return
	}
	resp = res.Data
	ok = true
	return
}

func SelectAllTuples(name_space string, user_id uint64) (resp []interface{}, err error) {
	name_space = utils.ConvertToNameInDB(name_space, user_id)
	res, err := TarantoolConn.Select(name_space, "primary", 0, 10, tarantool.IterAll, []interface{}{})
	if err != nil {
		utils.Logger.Info("cannot update tuple",
			zap.String("name_space", name_space),
			zap.Error(err),
		)
		return
	}
	resp = res.Data
	return
}