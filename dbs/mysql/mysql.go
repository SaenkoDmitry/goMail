package mysql

import (
	"database/sql"
	"fmt"
	"goMail/model"
	"goMail/utils"
	"go.uber.org/zap"
)

var MysqlConn *sql.DB

func GetUser(name string) (user model.User, err error) {
	rows, err := MysqlConn.Query("SELECT * FROM users where name=?", name)
	if err != nil {
		utils.Logger.Info("cannot execute query",
			zap.Error(err),
		)
		return
	}
	rows.Next()
	err = rows.Scan(&user.Id, &user.Name, &user.HashPassword)
	if err != nil {
		utils.Logger.Info("rows scan error",
			zap.Error(err),
		)
	}
	return
}

func GetAllUsers() (users []model.User, err error) {
	rows, err := MysqlConn.Query("SELECT * FROM users")
	if err != nil {
		utils.Logger.Info("cannot execute query",
			zap.Error(err),
		)
		return
	}
	var temp model.User
	for rows.Next() {
		err = rows.Scan(&temp.Id, &temp.Name, &temp.HashPassword)
		users = append(users, temp)
		if err != nil {
			utils.Logger.Info("rows scan error",
				zap.Error(err),
			)
			return
		}
	}
	return
}

func AddUser(name string, hash string) (ok bool, err error) {
	stmt, err := MysqlConn.Prepare("INSERT INTO users(name, hash_password) VALUES (?, ?)")
	if err != nil {
		utils.Logger.Info("cannot prepare statement for inserting",
			zap.Error(err),
			zap.String("username", name),
		)
		return
	}
	_, err = stmt.Exec(name, hash)
	if err != nil {
		utils.Logger.Info("cannot insert user",
			zap.Error(err),
			zap.String("username", name),
		)
		return
	}
	ok = true
	return
}

func DeleteUser(name string) (ok bool, err error) {
	stmt, err := MysqlConn.Prepare("DELETE FROM users WHERE name=?")
	if err != nil {
		utils.Logger.Info("cannot prepare statement for deleting",
			zap.Error(err),
			zap.String("username", name),
		)
		return
	}
	_, err = stmt.Exec(name)
	if err != nil {
		utils.Logger.Info("cannot delete user",
			zap.Error(err),
			zap.String("username", name),
		)
		return
	}
	ok = true
	return
}

func GetSpace(name_space string, user_id uint64) (space model.Space, err error) {
	name_space = utils.ConvertToNameInDB(name_space, user_id)
	rows, err := MysqlConn.Query("SELECT * FROM spaces where name=? AND user_id=?", name_space, user_id)
	if err != nil {
		utils.Logger.Info("cannot execute query",
			zap.Error(err),
			zap.String("name_space", name_space),
		)
		return
	}
	rows.Next()
	err = rows.Scan(&space.Id, &space.Name, &space.UserId)
	space.Name = utils.ConvertToRealName(space.Name, space.UserId)
	if err != nil {
		utils.Logger.Info("rows scan error",
			zap.Error(err),
			zap.String("name_space", name_space),
		)
		return
	}
	return
}

func GetAllSpaces(user_id uint64) (spaces []model.Space, err error) {
	rows, err := MysqlConn.Query("SELECT * FROM spaces where user_id=?", user_id)
	if err != nil {
		utils.Logger.Info("cannot execute query",
			zap.Error(err),
		)
		return
	}
	var space model.Space
	for rows.Next() {
		err = rows.Scan(&space.Id, &space.Name, &space.UserId)
		space.Name = utils.ConvertToRealName(space.Name, space.UserId)
		spaces = append(spaces, space)
		if err != nil {
			utils.Logger.Info("rows scan error",
				zap.Error(err),
			)
			return
		}
	}
	return
}

func AddSpace(name_space string, user_id uint64) (ok bool, err error) {
	name_space = utils.ConvertToNameInDB(name_space, user_id)
	stmt, err := MysqlConn.Prepare("INSERT INTO spaces (name, user_id) VALUES (?, ?)")
	if err != nil {
		utils.Logger.Info("cannot prepare statement for inserting",
			zap.Error(err),
			zap.String("name_space", name_space),
		)
		return
	}
	_, err = stmt.Exec(name_space, user_id)
	if err != nil {
		utils.Logger.Info("cannot insert space",
			zap.Error(err),
			zap.String("name_space", name_space),
		)
		return
	}
	ok = true
	return
}

func DeleteSpace(name_space string, user_id uint64) (ok bool, err error) {
	name_space = utils.ConvertToNameInDB(name_space, user_id)
	stmt, err := MysqlConn.Prepare("DELETE FROM spaces WHERE name=?")
	if err != nil {
		utils.Logger.Info("cannot prepare statement for deleting",
			zap.Error(err),
			zap.String("name_space", name_space),
		)
		return
	}
	_, err = stmt.Exec(name_space)
	if err != nil {
		utils.Logger.Info("cannot delete space",
			zap.Error(err),
			zap.String("name_space", name_space),
		)
		return
	}
	ok = true
	return
}

func AddPermission(user_id uint64, space_id uint64) (ok bool, err error) {
	stmt, err := MysqlConn.Prepare("INSERT INTO permissions (user_id, space_id) VALUES (?, ?)")
	if err != nil {
		utils.Logger.Info("cannot prepare statement for inserting",
			zap.Error(err),
		)
		return
	}
	_, err = stmt.Exec(user_id, space_id)
	if err != nil {
		utils.Logger.Info("cannot add permission",
			zap.Error(err),
		)
		return
	}
	ok = true
	return
}

func GetSpacePermissions(space_id uint64) (permissions []model.Permission, err error) {
	rows, err := MysqlConn.Query("SELECT * FROM permissions where space_id=?", space_id)
	if err != nil {
		utils.Logger.Info("cannot execute query",
			zap.Error(err),
		)
		return
	}
	var permission model.Permission
	for rows.Next() {
		err = rows.Scan(&permission.Id, &permission.User_id, &permission.Space_id)
		if err != nil {
			utils.Logger.Info("rows scan error",
				zap.Error(err),
			)
			return
		}
		permissions = append(permissions, permission)
	}
	return
}

func CheckPermissionsOnSpace(user_id uint64, space_id uint64) (exists bool, err error) {
	rows, err := MysqlConn.Query("SELECT * FROM permissions where user_id=? and space_id=?", user_id, space_id)
	if err != nil {
		utils.Logger.Info("cannot execute query",
			zap.Error(err),
		)
		return
	}
	exists = rows.Next()
	return
}

func GetUserHistory(user_id uint64) (history []model.History, err error) {
	rows, err := MysqlConn.Query("SELECT * FROM history where user_id=?", user_id)
	if err != nil {
		utils.Logger.Info("cannot execute query",
			zap.Error(err),
		)
		return
	}
	var hist model.History
	for rows.Next() {
		err = rows.Scan(&hist.Id, &hist.User_id, &hist.Space_id, &hist.Command, &hist.Result)
		if err != nil {
			utils.Logger.Info("rows scan error",
				zap.Error(err),
			)
			return
		}
		history = append(history, hist)
	}
	return
}

func GetSpaceHistory(space_id uint64) (history []model.History, err error) {
	rows, err := MysqlConn.Query("SELECT * FROM history where space_id=?", space_id)
	if err != nil {
		utils.Logger.Info("cannot execute query",
			zap.Error(err),
		)
		return
	}
	var hist model.History
	for rows.Next() {
		err = rows.Scan(&hist.Id, &hist.User_id, &hist.Space_id, &hist.Command)
		if err != nil {
			utils.Logger.Info("rows scan error",
				zap.Error(err),
			)
			return
		}
		history = append(history, hist)
	}
	return
}

func AddHistory(user_id uint64, space_id uint64, command string, result string) (ok bool, err error) {
	stmt, err := MysqlConn.Prepare("INSERT INTO history (user_id, space_id, command, result) VALUES (?, ?, ?, ?)")
	if err != nil {
		utils.Logger.Info("cannot prepare statement for inserting",
			zap.Error(err),
		)
		return
	}
	_, err = stmt.Exec(user_id, space_id, command, result)
	if err != nil {
		utils.Logger.Info("cannot insert history",
			zap.Error(err),
		)
		return
	}
	ok = true
	return
}

func InitMysql() (db *sql.DB) {
	a := utils.GetMysql()
	conn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s?%s", a.Username, a.Password, a.Protocol, a.Host, a.Port, a.Database, a.Encoding)
	db, err := sql.Open("mysql", conn)
	if err != nil {
		utils.Logger.Info("cannot open mysql connection",
			zap.Error(err),
		)
		return
	}
	return db
}
