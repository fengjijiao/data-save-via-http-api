package sqliteLib

import (
	"../coreLib"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	//"database/sql"
	//"os"
	"fmt"
	"time"
)

var sqliteR *SqliteDB

func Init() {
	sqliteR, _ = OpenDataBase(DBPath)
}

func InitTable() int64 {
	sql := `
		CREATE TABLE "DataSet" (dsid integer PRIMARY KEY, did INTEGER NOT NULL DEFAULT '0', valtype INTEGER NOT NULL DEFAULT '0', value TEXT);
		CREATE TABLE "DataSource" (did integer primary key , createdTimestamp INTEGER, updatedTimestamp INTEGER, valtype INTEGER NOT NULL DEFAULT '0');
		CREATE TABLE "User" (uid integer primary key, username varchar(20), password text , total INTEGER NOT NULL DEFAULT '0', token TEXT);
		`
	result, err := sqliteR.Exec(sql)
	if err != nil {
		fmt.Printf("InitTable Error: %s\n", err)
		return -1
	}
	id, err := result.LastInsertId()
	if err != nil {
		fmt.Printf("InitTable Error: %s\n", err)
		return -1
	}
	return id
}

func AuthUser(username, password string) bool {
	rows, err := sqliteR.PrepareQuery("SELECT COUNT(*) as count FROM User WHERE username = ? AND password = ? LIMIT 1;", username, password)
	if err != nil {
		fmt.Printf("AuthUser Error: %s\n", err)
		return false
	}
	count, err := sqliteR.GetCount(rows)
	return count > 0
}

func ExistUser(username string) bool {
	rows, err := sqliteR.PrepareQuery("SELECT COUNT(*) as count FROM User WHERE username = ? LIMIT 1;", username)
	if err != nil {
		fmt.Printf("AddUser Error: %s\n", err)
		return false
	}
	count, err := sqliteR.GetCount(rows)
	return count > 0
}

func AddUser(username, password string) (int64, error) {
	if ExistUser(username) {
		return -1, errors.New("this username was used.")
	}
	result, err := sqliteR.PrepareExec("INSERT INTO User (username, password, token) VALUES (?, ?, ?);", username, password, coreLib.RandString(TokenLength))
	if err != nil {
		fmt.Printf("AddUser Error: %s\n", err)
		return -1, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		fmt.Printf("AddUser Error: %s\n", err)
		return -1, err
	}
	return id, nil
}

func DelUser(uid int64) error {
	_, err := sqliteR.PrepareExec("DELETE FROM User WHERE uid = ?;", uid)
	if err != nil {
		fmt.Printf("DelUser Error: %s\n", err)
		return err
	}
	return nil
}

func UpdateToken(username, token string) error {
	_, err := sqliteR.PrepareExec("UPDATE User SET token = ? WHERE username = ?;", token, username)
	if err != nil {
		fmt.Printf("UpdateToken Error: %s\n", err)
		return err
	}
	return nil
}

func UpdatePassword(username, password string) error {
	_, err := sqliteR.PrepareExec("UPDATE User SET password = ? WHERE username = ?;", password, username)
	if err != nil {
		fmt.Printf("UpdatePassword Error: %s\n", err)
		return err
	}
	return nil
}

type DataSource struct {
	Did int64 `json:"did"`
	CreatedTimestamp int `json:"createdTimestamp"`
	UpdatedTimestamp int `json:"updatedTimestamp"`
}

func AddDataSource(valtype int) (int64, error) {
	result, err := sqliteR.PrepareExec("INSERT INTO DataSourceSet (createdTimestamp, updatedTimestamp, valtype) VALUES (?, ?, ?);", time.Now().Unix(), time.Now().Unix(), valtype)
	if err != nil {
		fmt.Printf("AddDataSource Error: %s\n", err)
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		fmt.Printf("AddDataSource Error: %s\n", err)
		return -1, err
	}
	return id, nil
}

func DelDataSource(did int) error {
	_, err := sqliteR.PrepareExec("DELETE FROM DataSourceSet WHERE did = ?;", did)
	if err != nil {
		fmt.Printf("DelDataSource Error: %s\n", err)
		return err
	}
	return nil
}

type UserInfo struct {
	Uid int64 `json:"uid"`
	UserName string `json:"username"`
	PassWord string `json:"password"`
	Total int `json:"total"`
	Token string `json:"token"`
}

func GetUserInfo(uid int64) (*UserInfo, error) {
	rows, err := sqliteR.PrepareQuery("SELECT * FROM User WHERE uid = ? LIMIT 1;", uid)
	if err != nil {
		fmt.Printf("GetUserInfo Error: %s\n", err)
		return nil, err
	}
	if rows.Next() {
		var userinfo UserInfo
		err = rows.Scan(&userinfo)
		if err != nil {
			fmt.Printf("GetUserInfo Error: %s\n", err)
			return nil, err
		}
		return &userinfo, nil
	}
	return nil, errors.New("No eligible user fond.")
}

func GetUidViaToken(token string) (int64, error) {
	rows, err := sqliteR.PrepareQuery("SELECT * FROM User WHERE token = ? LIMIT 1;", token)
	if err != nil {
		fmt.Printf("GetUidViaToken Error: %s\n", err)
		return -1, err
	}
	if rows.Next() {
		var userinfo UserInfo
		err = rows.Scan(&userinfo)
		if err != nil {
			fmt.Printf("GetUidViaToken Error: %s\n", err)
			return -1, err
		}
		return userinfo.Uid, nil
	}
	fmt.Println("GetUidViaToken Error: No eligible user fond.")
	return -1, errors.New("No eligible user fond.")
}

type DataSet struct {
	DSId int64 `json:"dsid"`
	Did int64 `json:"did"`
	ValueType int `json:"valtype"`
	Value string `json:"value"`
}

