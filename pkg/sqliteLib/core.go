package sqliteLib

import (
	_ "github.com/mattn/go-sqlite3"
	"database/sql"
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
	return sdb.DB.Query("select id, name from foo")
}

func (sdb *SqliteDB) Close() {
	sdb.DB.Close()
}