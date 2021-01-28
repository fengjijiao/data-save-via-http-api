package sqliteLib

import (
	"github.com/fengjijiao/data-save-via-http-api/pkg/conf"
	"github.com/fengjijiao/data-save-via-http-api/pkg/logio"
	"github.com/fengjijiao/data-save-via-http-api/pkg/commonio"
	"github.com/fengjijiao/data-save-via-http-api/pkg/coreLib"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"time"
	"path/filepath"
	"go.uber.org/zap"
	"bytes"
)

var sqliteR *SqliteDB

func Init() {
	sqliteR, _ = OpenDataBase(filepath.Join(conf.Config.WorkDir, "dsvha.db"))
}

func InitTable() int64 {
	lockFlag := []byte{0x09, 0x04, 0x06, 0x05}
	lockFilePath := filepath.Join(conf.Config.WorkDir, "db.flag")
	if commonio.IsFileExists(lockFilePath) {
		dat, err := commonio.ReadFile(lockFilePath)
		if err != nil {
			logio.Logger.Error("read db InitTable lock file error", zap.Error(err))
		}
		if bytes.Equal(dat, lockFlag) {
			return 0
		}
	}
	sql := `
		CREATE TABLE "DataSet" (dsid integer PRIMARY KEY, did INTEGER NOT NULL DEFAULT '0', valtype INTEGER NOT NULL DEFAULT '0', value TEXT);
		CREATE TABLE "DataSource" (did integer primary key , createdTimestamp INTEGER, updatedTimestamp INTEGER, valtype INTEGER NOT NULL DEFAULT '0', uid INTEGER NOT NULL DEFAULT '0');
		CREATE TABLE "User" (uid integer primary key, username varchar(20), password text , total INTEGER NOT NULL DEFAULT '0', token TEXT);
		`
	result, err := sqliteR.Exec(sql)
	if err != nil {
		logio.Logger.Error("InitTable Error", zap.Error(err))
		return -1
	}
	id, err := result.LastInsertId()
	if err != nil {
		logio.Logger.Error("InitTable Error", zap.Error(err))
		return -1
	}
	commonio.WriteToFile(lockFilePath, lockFlag)
	return id
}

func AuthUser(username, password string) bool {
	rows, err := sqliteR.PrepareQuery("SELECT COUNT(*) as count FROM User WHERE username = ? AND password = ? LIMIT 1;", username, password)
	if err != nil {
		logio.Logger.Info("AuthUser Fail", zap.Error(err))
		return false
	}
	defer rows.Close()
	count, err := sqliteR.GetCount(rows)
	return count > 0
}

func ExistUser(username string) bool {
	rows, err := sqliteR.PrepareQuery("SELECT COUNT(*) as count FROM User WHERE username = ? LIMIT 1;", username)
	if err != nil {
		logio.Logger.Info("ExistUser Fail", zap.Error(err))
		return false
	}
	defer rows.Close()
	count, err := sqliteR.GetCount(rows)
	return count > 0
}

func AddUser(username, password string) (int64, error) {
	if ExistUser(username) {
		logio.Logger.Info("this username was used.")
		return -1, errors.New("this username was used.")
	}
	result, err := sqliteR.PrepareExec("INSERT INTO User (username, password, token) VALUES (?, ?, ?);", username, password, coreLib.RandString(TokenLength))
	if err != nil {
		logio.Logger.Info("AddUser Fail", zap.Error(err))
		return -1, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		logio.Logger.Info("AddUser Fail", zap.Error(err))
		return -1, err
	}
	return id, nil
}

func DelUser(uid int64) error {
	_, err := sqliteR.PrepareExec("DELETE FROM User WHERE uid = ?;", uid)
	if err != nil {
		logio.Logger.Info("DelUser Fail", zap.Error(err))
		return err
	}
	return nil
}

func UpdateToken(username, token string) error {
	_, err := sqliteR.PrepareExec("UPDATE User SET token = ? WHERE username = ?;", token, username)
	if err != nil {
		logio.Logger.Info("UpdateToken Fail", zap.Error(err))
		return err
	}
	return nil
}

func UpdatePassword(username, password string) error {
	_, err := sqliteR.PrepareExec("UPDATE User SET password = ? WHERE username = ?;", password, username)
	if err != nil {
		logio.Logger.Info("UpdatePassword Fail", zap.Error(err))
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
		logio.Logger.Info("UpdateUpdatedTimestamp Fail", zap.Error(err))
		return err
	}
	return nil
}

func ExistDataSource(did int64) bool {
	rows, err := sqliteR.PrepareQuery("SELECT did FROM DataSource WHERE did = ? LIMIT 1;", did)
	if err != nil {
		logio.Logger.Info("ExistDataSource Fail", zap.Error(err))
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
		logio.Logger.Info("GetDataSource Fail", zap.Error(err))
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		var datasource DataSource
		err = rows.Scan(&datasource.Did, &datasource.CreatedTimestamp, &datasource.UpdatedTimestamp, &datasource.ValueType, &datasource.Uid)
		if err != nil {
			logio.Logger.Info("GetDataSource Fail", zap.Error(err))
			return nil, err
		}
		return &datasource, nil
	}
	logio.Logger.Info("GetDataSource FailNo eligible dataSource fond.")
	return nil, errors.New("No eligible dataSource fond.")
}

func GetDataSourcesViaUid(uid int64) (*[]DataSource, error) {
	rows, err := sqliteR.PrepareQuery("SELECT did, createdTimestamp, updatedTimestamp, valtype FROM DataSource WHERE uid = ?;", uid)
	if err != nil {
		logio.Logger.Info("GetDataSourcesViaUid Fail", zap.Error(err))
		return nil, err
	}
	defer rows.Close()
	var result []DataSource
	for rows.Next() {
		var item DataSource
		err = rows.Scan(&item.Did, &item.CreatedTimestamp, &item.UpdatedTimestamp, &item.ValueType)
		if err != nil {
			logio.Logger.Info("GetDataSourcesViaUid Fail", zap.Error(err))
			return nil, err
		}
		result = append(result, item)
	}
	return &result, nil
}

func AddDataSource(uid int64, valtype int) (int64, error) {
	result, err := sqliteR.PrepareExec("INSERT INTO DataSource (uid, createdTimestamp, updatedTimestamp, valtype) VALUES (?, ?, ?, ?);", uid, time.Now().Unix(), time.Now().Unix(), valtype)
	if err != nil {
		logio.Logger.Info("AddDataSource Fail", zap.Error(err))
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		logio.Logger.Info("AddDataSource Fail", zap.Error(err))
		return -1, err
	}
	return id, nil
}

func DelDataSource(did int64) error {
	_, err := sqliteR.PrepareExec("DELETE FROM DataSource WHERE did = ?;", did)
	if err != nil {
		logio.Logger.Info("DelDataSource Fail", zap.Error(err))
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
		logio.Logger.Info("GetUserInfo Fail", zap.Error(err))
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		var userinfo UserInfo
		err = rows.Scan(&userinfo.Uid, &userinfo.UserName, &userinfo.PassWord, &userinfo.Total, &userinfo.Token)
		if err != nil {
			logio.Logger.Info("GetUserInfo Fail", zap.Error(err))
			return nil, err
		}
		return &userinfo, nil
	}
	return nil, errors.New("No eligible user fond.")
}

func GetUserInfoViaUsername(username string) (*UserInfo, error) {
	rows, err := sqliteR.PrepareQuery("SELECT * FROM User WHERE username = ? LIMIT 1;", username)
	if err != nil {
		logio.Logger.Info("GetUserInfoViaUsername Fail", zap.Error(err))
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		var userinfo UserInfo
		err = rows.Scan(&userinfo.Uid, &userinfo.UserName, &userinfo.PassWord, &userinfo.Total, &userinfo.Token)
		if err != nil {
			logio.Logger.Info("GetUserInfoViaUsername Fail", zap.Error(err))
			return nil, err
		}
		return &userinfo, nil
	}
	return nil, errors.New("No eligible user fond.")
}

func GetUidViaToken(token string) (int64, error) {
	rows, err := sqliteR.PrepareQuery("SELECT uid FROM User WHERE token = ? LIMIT 1;", token)
	if err != nil {
		logio.Logger.Info("GetUidViaToken Fail", zap.Error(err))
		return -1, err
	}
	defer rows.Close()
	if rows.Next() {
		var uid int64
		err = rows.Scan(&uid)
		if err != nil {
			logio.Logger.Info("GetUidViaToken Fail", zap.Error(err))
			return -1, err
		}
		return uid, nil
	}
	logio.Logger.Info("GetUidViaToken FailNo eligible user fond.")
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
		logio.Logger.Info("DelDataSetByDid Fail", zap.Error(err))
		return err
	}
	return nil
}

func GetDataSet(dsid int64) (*DataSet, error) {
	rows, err := sqliteR.PrepareQuery("SELECT * FROM DataSet WHERE dsid = ? LIMIT 1;", dsid)
	if err != nil {
		logio.Logger.Info("GetDataSet Fail", zap.Error(err))
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		var dataset DataSet
		err = rows.Scan(&dataset.DSId, &dataset.Did, &dataset.ValueType, &dataset.Value)
		if err != nil {
			logio.Logger.Info("GetDataSet Fail", zap.Error(err))
			return nil, err
		}
		return &dataset, nil
	}
	logio.Logger.Info("GetDataSet FailNo eligible dataSet fond.")
	return nil, errors.New("No eligible dataSet fond.")
}

func ExistDataSet(dsid int64) bool {
	rows, err := sqliteR.PrepareQuery("SELECT dsid FROM DataSet WHERE dsid = ? LIMIT 1;", dsid)
	if err != nil {
		logio.Logger.Info("ExistDataSet Fail", zap.Error(err))
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
		logio.Logger.Info("CheckDataSetPermission Fail", zap.Error(err))
		return false
	}
	count, err := sqliteR.GetCount(rows)
	if err != nil {
		logio.Logger.Info("CheckDataSetPermission Fail", zap.Error(err))
		return false
	}
	return count > 0
}

func AddDataSet(did int64, valueType int, value string) (int64, error) {
	result, err := sqliteR.PrepareExec("INSERT INTO DataSet (did, valtype, value) VALUES (?, ?, ?);", did, valueType, value)
	if err != nil {
		logio.Logger.Info("AddDataSet Fail", zap.Error(err))
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		logio.Logger.Info("AddDataSet Fail", zap.Error(err))
		return -1, err
	}
	return id, nil
}

func UpdateDataSetValue(dsid int64, value string) error {
	_, err := sqliteR.PrepareExec("UPDATE DataSet SET value = ? WHERE dsid = ?;", value, dsid)
	if err != nil {
		logio.Logger.Info("UpdateDataSetValue Fail", zap.Error(err))
		return err
	}
	return nil
}

func GetDSIdViaDid(did int64) (int64, error) {
	rows, err := sqliteR.PrepareQuery("SELECT dsid FROM DataSet WHERE did = ? LIMIT 1;", did)
	if err != nil {
		logio.Logger.Info("GetDSIdViaDid Fail", zap.Error(err))
		return -1, err
	}
	defer rows.Close()
	if rows.Next() {
		var dsid int64
		err = rows.Scan(&dsid)
		if err != nil {
			logio.Logger.Info("GetDSIdViaDid Fail", zap.Error(err))
			return -1, err
		}
		return dsid, nil
	}
	logio.Logger.Info("GetDSIdViaDid FailNo eligible dataSet fond.")
	return -1, errors.New("No eligible dataSet fond.")
}
