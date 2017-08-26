package mysql

import (
	"database/sql"
	"fmt"
	"main/model"
	"main/utils"
)


var mysqlConn sql.DB

func init() {
	mysqlConn = InitMysql()
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		//panic(err)
	}
}

func GetUser(name string) (model.User, bool) {
	rows, err := mysqlConn.Query("SELECT * FROM users where name=?", name)
	checkErr(err)
	var user model.User
	if err == nil {
		b := rows.Next()
		err = rows.Scan(&user.Id, &user.Name, &user.HashPassword)
		checkErr(err)
		return user, b
	}
	return user, false
}

func SelectAllUser() ([]model.User) {
	rows, err := mysqlConn.Query("SELECT * FROM users")
	var a[]model.User
	var temp model.User
	if err == nil {
		for rows.Next() {
			err = rows.Scan(&temp.Id, &temp.Name, &temp.HashPassword)
			a = append(a, temp)
			checkErr(err)
		}
		return a
	}
	return nil
}

func AddUser(name string, hash string) {
	stmt, err := mysqlConn.Prepare("INSERT users SET name=?, hash_password=?")
	checkErr(err)
	_, err = stmt.Exec(name, hash)
	checkErr(err)
}

func DeleteUser(name string) {
	stmt, err := mysqlConn.Prepare("DELETE FROM users WHERE name=?")
	checkErr(err)
	_, err = stmt.Exec(name)
	checkErr(err)
}

func GetSpace(name_space string, user_id uint64) (model.Space, bool) {
	name_space = utils.ConvertToNameInDB(name_space, user_id)
	rows, err := mysqlConn.Query("SELECT * FROM spaces where name=? AND user_id=?", name_space, user_id)
	checkErr(err)
	var space model.Space
	if err == nil {
		b := rows.Next()
		err = rows.Scan(&space.Id, &space.Name, &space.UserId)
		checkErr(err)
		return space, b
	}
	return space, false
}

func GetAllSpaces(user_id uint64) ([]model.Space, bool) {
	rows, err := mysqlConn.Query("SELECT * FROM spaces where user_id=?", user_id)
	checkErr(err)
	var space model.Space
	var spaces []model.Space
	var b bool
	if err == nil {
		for rows.Next() {
			b = true
			err = rows.Scan(&space.Id, &space.Name, &space.UserId)
			space.Name = utils.ConvertToRealName(space.Name,  space.UserId)
			spaces = append(spaces, space)
		}
		return spaces, b
	}
	return spaces, false
}

func AddSpace(name string, user_id uint64) {
	name = utils.ConvertToNameInDB(name, user_id)
	stmt, err := mysqlConn.Prepare("INSERT spaces SET name=?, user_id=?")
	checkErr(err)
	_, err = stmt.Exec(name, user_id)
	checkErr(err)
}

func DeleteSpace(name string, user_id uint64) {
	name = utils.ConvertToNameInDB(name, user_id)
	stmt, err := mysqlConn.Prepare("DELETE FROM spaces WHERE name=?")
	checkErr(err)
	_, err = stmt.Exec(name)
	checkErr(err)
}

func AddPermission(user_id uint64, space_id uint64) {
	stmt, err := mysqlConn.Prepare("INSERT permissions SET user_id=?, space_id=?")
	checkErr(err)
	_, err = stmt.Exec(user_id, space_id)
	checkErr(err)
}

func GetUserPermissions(user_id uint64) (model.Permission, bool) {
	rows, err := mysqlConn.Query("SELECT * FROM permissions where user_id=?", user_id)
	checkErr(err)
	var permissions model.Permission
	if err == nil {
		b := rows.Next()
		err = rows.Scan(&permissions.Id, &permissions.User_id, &permissions.Space_id)
		checkErr(err)
		return permissions, b
	}
	return permissions, false
}

func GetSpacePermissions(space_id uint64) ([]model.Permission, bool) {
	rows, err := mysqlConn.Query("SELECT * FROM permissions where space_id=?", space_id)
	checkErr(err)
	var b bool
	var permissions []model.Permission
	var permission model.Permission
	if err == nil {
		for rows.Next() {
			b = true
			err = rows.Scan(&permission.Id, &permission.User_id, &permission.Space_id)
			checkErr(err)
			permissions = append(permissions, permission)
		}
		return permissions, b
	}
	return permissions, false
}

func CheckPermissionsOnSpace(user_id uint64, space_id uint64) bool {
	rows, err := mysqlConn.Query("SELECT * FROM permissions where user_id=? and space_id=?", user_id, space_id)
	checkErr(err)
	b := rows.Next()
	return b
}

//func DeletePermissionsOnSpace(name_space string, user_id uint64) bool {
//	name_space = utils.ConvertToNameInDB(name_space, user_id)
//	stmt, err := mysqlConn.Prepare("DELETE FROM permissions WHERE name=?")
//	checkErr(err)
//	_, err = stmt.Exec(name_space)
//	checkErr(err)
//}

func GetUserHistory(user_id uint64) ([]model.History, bool) {
	rows, err := mysqlConn.Query("SELECT * FROM history where user_id=?", user_id)
	checkErr(err)
	var hist model.History
	var hists []model.History
	var b bool
	if err == nil {
		for rows.Next() {
			b = true
			err = rows.Scan(&hist.Id, &hist.User_id, &hist.Space_id, &hist.Command, &hist.Result)
			hists = append(hists, hist)
		}
		return hists, b
	}
	return hists, false
}

func GetSpaceHistory(space_id uint64) ([]model.History, bool) {
	rows, err := mysqlConn.Query("SELECT * FROM history where space_id=?", space_id)
	checkErr(err)
	var hist model.History
	var hists []model.History
	var b bool
	if err == nil {
		for rows.Next() {
			b = true
			err = rows.Scan(&hist.Id, &hist.User_id, &hist.Space_id, &hist.Command)
			hists = append(hists, hist)
		}
		return hists, b
	}
	return hists, false
}

func AddHistory(user_id uint64, space_id uint64, command string, result string) {
	stmt, err := mysqlConn.Prepare("INSERT history SET user_id=?, space_id=?, command=?, result=?")
	checkErr(err)
	_, err = stmt.Exec(stmt, user_id, space_id, command, result)
	checkErr(err)
}

func InitMysql() (sql.DB) {
	db, err := sql.Open("mysql", "adminGo:gomail@tcp(localhost:3306)/tarantool_spaces_store?charset=utf8")
	defer db.Close()
	if err != nil {
		fmt.Println(err)
	}
	return *db
}