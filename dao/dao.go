package dao

import (
	// "os"
	"database/sql"
	"log"
	"fmt"
)


//DBへの接続を初期化
func NewDBConnection()(*sql.DB,error) {
	// mysqlUser := os.Getenv("MYSQL_USER")
	// mysqlPwd := os.Getenv("MYSQL_PWD")
	// mysqlHost := os.Getenv("MYSQL_HOST")
	// mysqlDatabase := os.Getenv("MYSQL_DATABASE")
	// connStr := fmt.Sprintf("%s:%s@%s/%s", mysqlUser, mysqlPwd, mysqlHost, mysqlDatabase)
	// db, err := sql.Open("mysql", connStr)
	// if err != nil {
	// 	log.Fatalf("fail: sql.Open, %v\n", err)
	// 	return nil,err
	// }
	// if err := db.Ping(); err != nil {
	// 	log.Fatalf("fail: _db.Ping, %v\n", err)
	// 	return nil,err
	// }
	// return db,nil

	////localとつなげるとき
	mysqlUser := "user"
	mysqlUserPwd := "password"
	mysqlDatabase := "mydatabase"
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@(localhost:3306)/%s", mysqlUser, mysqlUserPwd, mysqlDatabase))
	if err != nil {
		log.Fatalf("fail: sql.Open, %v\n", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("fail: _db.Ping, %v\n", err)
	}
	return db,nil

}