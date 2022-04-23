package db

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var Mysql *sql.DB

func InitDB() bool {
	var err error
	Mysql, err = sql.Open("mysql", "poker:fake@tcp(localhost:3306)/lottery")
	if err != nil {
		log.Println("[Error]Failed to connect DB : ", err)
		return false
	}
	Mysql.SetConnMaxLifetime(time.Minute * 3)
	Mysql.SetMaxOpenConns(10)
	Mysql.SetMaxIdleConns(10)
	return true
}

func SelectQueryRows(query string) (*sql.Rows, error) {
	rows, err := Mysql.Query(query)
	if err != nil {
		log.Println("Failed to excute : ", err)
		return nil, err
	}
	return rows, err
}

func SelectQueryRow(query string) *sql.Row {
	row := Mysql.QueryRow(query)
	return row
}

func ExcuteQuery(query string) (int64, error) {
	result, err := Mysql.Exec(query)
	if err != nil {
		log.Println("Failed to excute : ", err)
		return 0, err
	}
	success, err := result.RowsAffected()
	if err != nil {
		log.Println("Failed to excute : ", err)
		return 0, err
	}
	if success == 0 {
		log.Println("Not Excute : 0")
		return 0, err
	}
	return success, nil
}

func CheckPing() error {
	var err error
	if Mysql == nil {
		Mysql, err = sql.Open("mysql", "poker:fake@tcp(localhost:3306)/lottery")
		if err != nil {
			return err
		}
	}
	err = Mysql.Ping()
	if err != nil {
		return err
	}
	return nil
}
