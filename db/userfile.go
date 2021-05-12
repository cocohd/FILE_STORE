package db

import (
	"FILE_STORE/db/mysql"
	"fmt"
	"time"
)

type UserFile struct{
	UserName string
	FileHash string
	FileName string
	FileSize int64
	UploadedAt string
	LastUpdated string
}


// 插入tbl_user_file表
func OnFileUploadFinishedDb(userName string, filesha1 string, fileName string, fileSize int64,) bool {
	stmt, err := mysql.DBConn().Prepare("insert ignore into tbl_user_file (`user_name`, " +
		"`file_sha1`, `file_name`, `file_size`, `upload_at`) values (?,?,?,?,?)")
	if err != nil {
		return  false
	}
	defer stmt.Close()

	_, err = stmt.Exec(userName, filesha1, fileName, fileSize, time.Now())
	if err != nil {
		fmt.Println("Error insert into tbl_user_file: ", err)
		return false
	}
	return true
}


// 返回tbl_user_file表中的数据
func QueryUserFileDb(userName string, limit int) ([]UserFile, error) {
	stmt, err := mysql.DBConn().Prepare("select file_sha1,file_name,file_size,upload_at,last_update from" +
		"tbl_user_file where user_name=? limit ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(userName, limit)
	if err != nil {
		return nil, err
	}

	var userFiles []UserFile
	for rows.Next() {
		uFile := UserFile{}
		rows.Scan(&uFile.FileHash, &uFile.FileName, &uFile.UploadedAt, &uFile.LastUpdated)
		userFiles = append(userFiles, uFile)
	}

	return userFiles, err
}