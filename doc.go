// checkDbbackup project doc.go

/*
checkDbbackup document
*/
package main

import (
	"database/sql"
	"fmt"
	//	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func linkSql(dblogin string) {

	fmt.Println("加载数据库信息：" + dblogin)
	//db, _ = sql.Open("mysql", "root:owenshen123@tcp(127.0.0.1:3306)/test?charset=utf8")
	db, _ = sql.Open("mysql", dblogin)
	db.SetMaxOpenConns(2000)
	db.SetMaxIdleConns(1000)
	db.Ping()
}

func ReadData(script string) [][]string { //返回多列数据
	var DBdata [][]string
	rows, _ := db.Query(script)
	columns, _ := rows.Columns()
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	for rows.Next() {
		var tmpArr []string
		_ = rows.Scan(scanArgs...)
		for _, col := range values {
			if col == nil {
				tmpArr = append(tmpArr, "NULL")
			} else {
				tmpArr = append(tmpArr, string(col))
			}
		}
		DBdata = append(DBdata, tmpArr)
	}

	defer rows.Close()
	return DBdata
}
