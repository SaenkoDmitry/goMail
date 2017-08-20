package mysql

import (
	"database/sql"
	"fmt"
	"main/model"
)

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		//panic(err)
	}
}

func GetUser(name string) (model.User, bool) {
	db := GetDb()
	rows, err := db.Query("SELECT * FROM users where name=?", name)
	checkErr(err)
	b := rows.Next()
	var user model.User
	err = rows.Scan(&user.Id, &user.Name, &user.HashPassword)
	checkErr(err)
	return user, b
}

func AddUser(name string, hash string) {
	db := GetDb()
	stmt, err := db.Prepare("INSERT users SET name=?, hash_password=?")
	checkErr(err)
	_, err = stmt.Exec(name, hash)
	checkErr(err)
}

func DeleteUser(name string) {
	db := GetDb()
	stmt, err := db.Prepare("DELETE FROM users WHERE name=?")
	checkErr(err)
	_, err = stmt.Exec(name)
	checkErr(err)
}

func GetSpace(name string, user_id uint64) (model.Space, bool) {
	db := GetDb()
	rows, err := db.Query("SELECT * FROM spaces where name=? AND user_id=?", name, user_id)
	checkErr(err)
	b := rows.Next()
	var space model.Space
	err = rows.Scan(&space.Id, &space.Name, &space.UserId)
	checkErr(err)
	return space, b
}

func GetAllSpaces(user_id uint64) ([]model.Space, bool) {
	db := GetDb()
	rows, err := db.Query("SELECT * FROM spaces where user_id=?", user_id)
	checkErr(err)
	var space model.Space
	var spaces []model.Space
	var b bool
	for rows.Next() {
		b = true
		err = rows.Scan(&space.Id, &space.Name, &space.UserId)
		spaces = append(spaces, space)
	}
	return spaces, b
}

func AddSpace(name string, user_id uint64) {
	db := GetDb()
	stmt, err := db.Prepare("INSERT spaces SET name=?, user_id=?")
	checkErr(err)
	_, err = stmt.Exec(name, user_id)
	checkErr(err)
}

func DeleteSpace(name string) {
	db := GetDb()
	stmt, err := db.Prepare("DELETE FROM spaces WHERE name=?")
	checkErr(err)
	_, err = stmt.Exec(name)
	checkErr(err)
}

func GetUserPermissions(user_id uint64) (model.Permission, bool) {
	db := GetDb()
	rows, err := db.Query("SELECT * FROM permissions where user_id=?", user_id)
	checkErr(err)
	b := rows.Next()
	var permissions model.Permission
	err = rows.Scan(&permissions.Id, &permissions.User_id, &permissions.Space_id)
	checkErr(err)
	return permissions, b
}

func GetSpacePermissions(space_id uint64) ([]model.Permission, bool) {
	db := GetDb()
	rows, err := db.Query("SELECT * FROM permissions where space_id=?", space_id)
	checkErr(err)
	var b bool
	var permissions []model.Permission
	var permission model.Permission
	for rows.Next() {
		b = true
		err = rows.Scan(&permission.Id, &permission.User_id, &permission.Space_id)
		checkErr(err)
		permissions = append(permissions, permission)
	}
	return permissions, b
}

func CheckPermissionsOnSpace(user_id uint64, space_id uint64) bool {
	db := GetDb()
	rows, err := db.Query("SELECT * FROM permissions where user_id=? and space_id=?", user_id, space_id)
	checkErr(err)
	b := rows.Next()
	return b
}

func GetUserHistory(user_id uint64) ([]model.History, bool) {
	db := GetDb()
	rows, err := db.Query("SELECT * FROM history where user_id=?", user_id)
	checkErr(err)
	var hist model.History
	var hists []model.History
	var b bool
	for rows.Next() {
		b = true
		err = rows.Scan(&hist.Id, &hist.User_id, &hist.Space_id, &hist.Command)
		hists = append(hists, hist)
	}
	return hists, b
}

func GetSpaceHistory(space_id uint64) ([]model.History, bool) {
	db := GetDb()
	rows, err := db.Query("SELECT * FROM history where space_id=?", space_id)
	checkErr(err)
	var hist model.History
	var hists []model.History
	var b bool
	for rows.Next() {
		b = true
		err = rows.Scan(&hist.Id, &hist.User_id, &hist.Space_id, &hist.Command)
		hists = append(hists, hist)
	}
	return hists, b
}

func AddHistory(user_id uint64, space_id uint64, command string, result string) {
	db := GetDb()
	stmt, err := db.Prepare("INSERT spaces SET user_id=?, space_id=?, command=?, result=?")
	checkErr(err)
	_, err = stmt.Exec(stmt, user_id, space_id, command, result)
	checkErr(err)
}

func GetDb() (sql.DB) {
	db, err := sql.Open("mysql", "adminGo:gomail@tcp(localhost:3306)/tarantool_spaces_store?charset=utf8")
	defer db.Close()
	if err != nil {
		fmt.Println(err)
	}
	return *db
}