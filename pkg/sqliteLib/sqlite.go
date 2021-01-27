package sqliteLib

import (
	"github.com/fengjijiao/data-save-via-http-api/pkg/coreLib"
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
		CREATE TABLE "DataSource" (did integer primary key , createdTimestamp INTEGER, updatedTimestamp INTEGER, valtype INTEGER NOT NULL DEFAULT '0', uid INTEGER NOT NULL DEFAULT '0');
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
	defer rows.Close()
	count, err := sqliteR.GetCount(rows)
	return count > 0
}

func ExistUser(username string) bool {
	rows, err := sqliteR.PrepareQuery("SELECT COUNT(*) as count FROM User WHERE username = ? LIMIT 1;", username)
	if err != nil {
		fmt.Printf("AddUser Error: %s\n", err)
		return false
	}
	defer rows.Close()
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
	Uid int64 `json:"uid"`
	ValueType int `json:"valtype"`
	CreatedTimestamp int64 `json:"createdTimestamp"`
	UpdatedTimestamp int64 `json:"updatedTimestamp"`
}

func UpdateUpdatedTimestamp(did int64) error {
	_, err := sqliteR.PrepareExec("UPDATE DataSource SET updatedTimestamp = ? WHERE did = ?;", time.Now().Unix(), did)
	if err != nil {
		fmt.Printf("UpdateUpdatedTimestamp Error: %s\n", err)
		return err
	}
	return nil
}

func ExistDataSource(did int64) bool {
	rows, err := sqliteR.PrepareQuery("SELECT did FROM DataSource WHERE did = ? LIMIT 1;", did)
	if err != nil {
		fmt.Printf("ExistDataSource Error: %s\n", err)
		return false
	}
	count, err := sqliteR.GetCount(rows)
	if err != nil {
		return false
	}
	return count > 0
}

func GetDataSource(did int64) (*DataSource, error) {
	rows, err := sqliteR.PrepareQuery("SELECT * FROM DataSource WHERE did = ? LIMIT 1;", did)
	if err != nil {
		fmt.Printf("GetDataSource Error: %s\n", err)
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		var datasource DataSource
		err = rows.Scan(&datasource.Did, &datasource.CreatedTimestamp, &datasource.UpdatedTimestamp, &datasource.ValueType, &datasource.Uid)
		if err != nil {
			fmt.Printf("GetDataSource Error: %s\n", err)
			return nil, err
		}
		return &datasource, nil
	}
	fmt.Println("GetDataSource Error: No eligible dataSource fond.")
	return nil, errors.New("No eligible dataSource fond.")
}

func GetDataSourcesViaUid(uid int64) (*[]DataSource, error) {
	rows, err := sqliteR.PrepareQuery("SELECT did, createdTimestamp, updatedTimestamp, valtype FROM DataSource WHERE uid = ?;", uid)
	if err != nil {
		fmt.Printf("GetDataSourcesViaUid Error: %s\n", err)
		return nil, err
	}
	defer rows.Close()
	var result []DataSource
	for rows.Next() {
		var item DataSource
		err = rows.Scan(&item.Did, &item.CreatedTimestamp, &item.UpdatedTimestamp, &item.ValueType)
		if err != nil {
			fmt.Printf("GetDataSourcesViaUid Error: %s\n", err)
			return nil, err
		}
		result = append(result, item)
	}
	return &result, nil
}

func AddDataSource(uid int64, valtype int) (int64, error) {
	result, err := sqliteR.PrepareExec("INSERT INTO DataSource (uid, createdTimestamp, updatedTimestamp, valtype) VALUES (?, ?, ?, ?);", uid, time.Now().Unix(), time.Now().Unix(), valtype)
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

func DelDataSource(did int64) error {
	_, err := sqliteR.PrepareExec("DELETE FROM DataSource WHERE did = ?;", did)
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
	defer rows.Close()
	if rows.Next() {
		var userinfo UserInfo
		err = rows.Scan(&userinfo.Uid, &userinfo.UserName, &userinfo.PassWord, &userinfo.Total, &userinfo.Token)
		if err != nil {
			fmt.Printf("GetUserInfo Error: %s\n", err)
			return nil, err
		}
		return &userinfo, nil
	}
	return nil, errors.New("No eligible user fond.")
}

func GetUserInfoViaUsername(username string) (*UserInfo, error) {
	rows, err := sqliteR.PrepareQuery("SELECT * FROM User WHERE username = ? LIMIT 1;", username)
	if err != nil {
		fmt.Printf("GetUserInfoViaUsername Error: %s\n", err)
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		var userinfo UserInfo
		err = rows.Scan(&userinfo.Uid, &userinfo.UserName, &userinfo.PassWord, &userinfo.Total, &userinfo.Token)
		if err != nil {
			fmt.Printf("GetUserInfoViaUsername Error: %s\n", err)
			return nil, err
		}
		return &userinfo, nil
	}
	return nil, errors.New("No eligible user fond.")
}

func GetUidViaToken(token string) (int64, error) {
	rows, err := sqliteR.PrepareQuery("SELECT uid FROM User WHERE token = ? LIMIT 1;", token)
	if err != nil {
		fmt.Printf("GetUidViaToken Error: %s\n", err)
		return -1, err
	}
	defer rows.Close()
	if rows.Next() {
		var uid int64
		err = rows.Scan(&uid)
		if err != nil {
			fmt.Printf("GetUidViaToken Error: %s\n", err)
			return -1, err
		}
		return uid, nil
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

func DelDataSetViaDid(did int64) error {
	_, err := sqliteR.PrepareExec("DELETE FROM DataSet WHERE did = ?;", did)
	if err != nil {
		fmt.Printf("DelDataSetByDid Error: %s\n", err)
		return err
	}
	return nil
}

func GetDataSet(dsid int64) (*DataSet, error) {
	rows, err := sqliteR.PrepareQuery("SELECT * FROM DataSet WHERE dsid = ? LIMIT 1;", dsid)
	if err != nil {
		fmt.Printf("GetDataSet Error: %s\n", err)
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		var dataset DataSet
		err = rows.Scan(&dataset.DSId, &dataset.Did, &dataset.ValueType, &dataset.Value)
		if err != nil {
			fmt.Printf("GetDataSet Error: %s\n", err)
			return nil, err
		}
		return &dataset, nil
	}
	fmt.Println("GetDataSet Error: No eligible dataSet fond.")
	return nil, errors.New("No eligible dataSet fond.")
}

func ExistDataSet(dsid int64) bool {
	rows, err := sqliteR.PrepareQuery("SELECT dsid FROM DataSet WHERE dsid = ? LIMIT 1;", dsid)
	if err != nil {
		fmt.Printf("ExistDataSet Error: %s\n", err)
		return false
	}
	count, err := sqliteR.GetCount(rows)
	if err != nil {
		return false
	}
	return count > 0
}

func CheckDataSetPermission(uid, did int64) bool {
	rows, err := sqliteR.PrepareQuery("SELECT COUNT(*) as count FROM DataSource WHERE did = ? AND uid = ? LIMIT 1;", did, uid)
	if err != nil {
		fmt.Printf("CheckDataSetPermission Error: %s\n", err)
		return false
	}
	count, err := sqliteR.GetCount(rows)
	if err != nil {
		fmt.Printf("CheckDataSetPermission Error: %s\n", err)
		return false
	}
	return count > 0
}

func AddDataSet(did int64, valueType int, value string) (int64, error) {
	result, err := sqliteR.PrepareExec("INSERT INTO DataSet (did, valtype, value) VALUES (?, ?, ?);", did, valueType, value)
	if err != nil {
		fmt.Printf("AddDataSet Error: %s\n", err)
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		fmt.Printf("AddDataSet Error: %s\n", err)
		return -1, err
	}
	return id, nil
}

func UpdateDataSetValue(dsid int64, value string) error {
	_, err := sqliteR.PrepareExec("UPDATE DataSet SET value = ? WHERE dsid = ?;", value, dsid)
	if err != nil {
		fmt.Printf("UpdateDataSetValue Error: %s\n", err)
		return err
	}
	return nil
}

func GetDSIdViaDid(did int64) (int64, error) {
	rows, err := sqliteR.PrepareQuery("SELECT dsid FROM DataSet WHERE did = ? LIMIT 1;", did)
	if err != nil {
		fmt.Printf("GetDSIdViaDid Error: %s\n", err)
		return -1, err
	}
	defer rows.Close()
	if rows.Next() {
		var dsid int64
		err = rows.Scan(&dsid)
		if err != nil {
			fmt.Printf("GetDSIdViaDid Error: %s\n", err)
			return -1, err
		}
		return dsid, nil
	}
	fmt.Println("GetDSIdViaDid Error: No eligible dataSet fond.")
	return -1, errors.New("No eligible dataSet fond.")
}
