package db

import (
	"database/sql"
	"fmt"
	"gin-todolist/library"
	"sync"

	"github.com/didi/gendry/builder"
	"github.com/didi/gendry/scanner"
	_ "github.com/go-sql-driver/mysql"
)

type dbConfig struct {
	IP       string
	Port     int
	Username string
	Password string
	DBName   string
}

type Mysql struct {
	DB *sql.DB
}

var mysqlDB *Mysql
var dbMu sync.Mutex

// GetDB: 获取数据库, 单例模式
func GetMysqlDB(configName string) *Mysql {
	if mysqlDB == nil {
		dbMu.Lock()
		defer dbMu.Unlock()
		if mysqlDB == nil {
			var dbconf dbConfig

			if _, err := library.GetConf(configName, &dbconf); err != nil {
				library.CheckErr(err)
			}

			db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", dbconf.Username, dbconf.Password, dbconf.IP, dbconf.Port, dbconf.DBName))
			library.CheckErr(err)
			mysqlDB = &Mysql{DB: db}
		}
	}
	return mysqlDB
	// var dbconf dbConfig

	// if _, err := library.GetConf(configName, &dbconf); err != nil {
	// 	library.CheckErr(err)
	// }

	// db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", dbconf.Username, dbconf.Password, dbconf.IP, dbconf.Port, dbconf.DBName))
	// library.CheckErr(err)
	// return &Mysql{
	// 	DB: db,
	// }
}

// Select: 获取数据
func (mysql *Mysql) Select(tableName string, conds Where, fields []string, target interface{}) error {
	query, args, err := builder.BuildSelect(tableName, conds, fields)
	if err != nil {
		return err
	}

	return mysql.query(target, query, args, nil)
}

// SelectTx: 事务中的查询
func (mysql *Mysql) SelectTx(tx *sql.Tx, tableName string, conds Where, fields []string, target interface{}) error {
	query, args, err := builder.BuildSelect(tableName, conds, fields)
	if err != nil {
		return err
	}

	return mysql.query(target, query, args, tx)
}

// NamedQuery: 复杂查询
func (mysql *Mysql) NamedQuery(namedQuery string, namedConds Where, target interface{}) error {
	query, args, err := builder.NamedQuery(namedQuery, namedConds)
	if err != nil {
		return err
	}
	return mysql.query(target, query, args, nil)
}

func (mysql *Mysql) query(target interface{}, query string, args []interface{}, tx *sql.Tx) error {
	var rows *sql.Rows
	var err error
	// 根据是否传入tx决定是否使用事务
	if tx == nil {
		rows, err = mysql.DB.Query(query, args...)
	} else {
		rows, err = tx.Query(query, args...)
	}
	if err != nil {
		return err
	}

	if rows != nil {
		defer rows.Close()
	}

	err = scanner.Scan(rows, target)
	return err
}

func (mysql *Mysql) Update(tableName string, conds Where, data RowData) (int64, error) {
	query, args, err := builder.BuildUpdate(tableName, conds, data)
	if err != nil {
		return 0, err
	}

	result, err := mysql.DB.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (mysql *Mysql) UpdateTx(tx *sql.Tx, tableName string, conds Where, data RowData) (int64, error) {
	query, args, err := builder.BuildUpdate(tableName, conds, data)
	if err != nil {
		return 0, err
	}

	result, err := tx.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (mysql *Mysql) Insert(tableName string, data RowDataArr) (int64, error) {
	query, args, err := builder.BuildInsert(tableName, data)
	if err != nil {
		return -1, err
	}
	result, err := mysql.DB.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (mysql *Mysql) InsertTx(tx *sql.Tx, tableName string, data RowDataArr) (int64, error) {
	query, args, err := builder.BuildInsert(tableName, data)
	if err != nil {
		return -1, err
	}
	result, err := tx.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (mysql *Mysql) Delete(tableName string, conds Where) error {
	query, args, err := builder.BuildDelete(tableName, conds)
	if err != nil {
		return err
	}
	_, err = mysql.DB.Exec(query, args...)
	return err
}

func (mysql *Mysql) DeleteTx(tx *sql.Tx, tableName string, conds Where) error {
	query, args, err := builder.BuildDelete(tableName, conds)
	if err != nil {
		return err
	}
	_, err = tx.Exec(query, args...)
	return err
}

func (mysql *Mysql) BeginTx() (*sql.Tx, error) {
	return mysql.DB.Begin()
}

func (mysql *Mysql) Commit(tx *sql.Tx) error {
	return tx.Commit()
}

func (mysql *Mysql) Rollback(tx *sql.Tx) error {
	return tx.Rollback()
}
