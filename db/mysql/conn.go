package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	//第三方包匿名方式导入上面sql的path里，这样课通过sql直接调用
	_"github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	// 连接
	db, _ = sql.Open("mysql", "root:756979099@tcp(127.0.0.1:3306)/go_fileStore?charset=utf8")
	// 最大连接数
	db.SetMaxOpenConns(1000)
	err := db.Ping()
	if err != nil {
		fmt.Println("Failed to connect to mysql, err:" + err.Error())
		os.Exit(1)
	}
}

func DBConn() *sql.DB {
	return db
}


// 返回数据的处理
func ParseRows(rows *sql.Rows) []map[string]interface{} {
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for j := range values {
		scanArgs[j] = &values[j]
	}
	//fmt.Println(scanArgs)
	record := make(map[string]interface{})
	records := make([]map[string]interface{}, 0)
	for rows.Next() {
		//将行数据保存到record字典
		err := rows.Scan(scanArgs...)

		checkErr(err)

		for i, col := range values {
			//fmt.Println(col)

			if col != nil {
				record[columns[i]] = col
			}
		}
		records = append(records, record)
	}
	//fmt.Println(records)

	return records
}


func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}