package sqliteLib

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

type SqliteDB struct {
	DBPath string
	DB *sql.DB
}

//创建SqliteDB实体
func OpenDataBase(DBPath string) (*SqliteDB, error) {
	DB, err := sql.Open("sqlite3", DBPath)
	if err != nil {
		return nil, err
	}
	return &SqliteDB{DBPath: DBPath, DB: DB,}, nil
}

func (sdb *SqliteDB) Exec(sql string) (sql.Result, error) {
	result, err := sdb.DB.Exec(sql)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (sdb *SqliteDB) Query(sql string) (*sql.Rows, error) {
	return sdb.DB.Query(sql)
}

func (sdb *SqliteDB) PrepareExec(prepareSql string, args... interface{}) (sql.Result, error) {
	stmt, err := sdb.DB.Prepare(prepareSql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	result, err := stmt.Exec(args...)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (sdb *SqliteDB) PrepareQuery(prepareSql string, args... interface{}) (*sql.Rows, error) {
	stmt, err := sdb.DB.Prepare(prepareSql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (sdb *SqliteDB) GetCount(rows *sql.Rows) (int, error) {
	var count int
	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return -1, err
		}
	}
	return count, nil
}

func (sdb *SqliteDB) Close() {
	sdb.DB.Close()
}

func (sdb *SqliteDB) Remove() {
	os.Remove(sdb.DBPath)
}