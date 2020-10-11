package sqliteLib

import (
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
		CREATE TABLE DataSourceSet (did integer primary key , createdTimestamp INTEGER, updatedTimestamp INTEGER, valtype INTEGER NOT NULL DEFAULT '0');
		CREATE TABLE "User" (uid integer primary key, username varchar(20), password text , total INTEGER NOT NULL DEFAULT '0');
		`
	result, err := sqliteR.Exec(sql)
	if(err != nil) {
		fmt.Printf("InitTable Error: %s\n", err)
		return -1
	}
	id, err := result.LastInsertId()
	if(err != nil) {
		fmt.Printf("InitTable Error: %s\n", err)
		return -1
	}
	return id
}

func AddUser(username, password string) (int64, error) {
	sql := fmt.Sprintf(`INSERT INTO user (uid,username,password,total) VALUES (
				'',
				'%s',
				'%s',
				''
			); `, username, password)
	result, err := sqliteR.Exec(sql)
	if(err != nil) {
		fmt.Printf("AddUser Error: %s\n", err)
		return -1, err
	}
	id, err := result.LastInsertId()
	if(err != nil) {
		fmt.Printf("AddUser Error: %s\n", err)
		return -1, err
	}
	return id, nil
}

func DelUser(uid int) error {
	sql := fmt.Sprintf(`DELETE FROM User WHERE uid = '%d';`, uid)
	_, err := sqliteR.Exec(sql)
	if(err != nil) {
		fmt.Printf("DelUser Error: %s\n", err)
		return err
	}
	return nil
}

func AddDataSource(valtype int) (int64, error) {
	sql := fmt.Sprintf(`INSERT INTO DataSourceSet (createdTimestamp,updatedTimestamp,valtype) VALUES (
						'%d',
						'%d',
						'%d'
						);`, time.Now().Unix(), time.Now().Unix(), valtype)
	result, err := sqliteR.Exec(sql)
	if(err != nil) {
		fmt.Printf("AddDataSource Error: %s\n", err)
		return 0, err
	}
	id, err := result.LastInsertId()
	if(err != nil) {
		fmt.Printf("AddDataSource Error: %s\n", err)
		return -1, err
	}
	return id, nil
}

func DelDataSource(did int) error {
	sql := fmt.Sprintf(`DELETE FROM DataSourceSet WHERE did='%d';`, did)
	_, err := sqliteR.Exec(sql)
	if(err != nil) {
		fmt.Printf("DelDataSource Error: %s\n", err)
		return err
	}
	return nil
}