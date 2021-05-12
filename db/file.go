package db

import (
	mydb"FILE_STORE/db/mysql"
	"database/sql"
	"fmt"
)


// 存储文件到数据库中
func OnFileUploadedFinished(filehash string, filename string, filesize int64, fileaddr string ) bool {
	stmt, err := mydb.DBConn().Prepare("insert ignore into tbl_file (`file_sha1`, `file_name`, `file_size`, `file_addr`, `status`) values (?,?,?,?,1)")
	if err != nil {
		fmt.Println("Failed to prepare statement,err:", err.Error())
		return false
	}
	defer stmt.Close()
	ret, err := stmt.Exec(filehash, filename, filesize, fileaddr)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	// 判断是否重复插入
	// 是否返回一条新的表记录
	if rf, err := ret.RowsAffected(); nil==err {
		if rf <= 0 {
			fmt.Println("File with hash:%s has been uploaded before!", filehash)
		}
		return true
	}
	return false
}


type FileTable struct {
	FileHash string
	FileName sql.NullString
	FileSize sql.NullInt64
	FileAddr sql.NullString
}
func GetFile(filehash string) (*FileTable, error) {
	stmt, err := mydb.DBConn().Prepare("select file_sha1, file_name, file_size, file_addr from tbl_file " +
		"where file_sha1=? and status=1 limit 1")
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	defer stmt.Close()

	fileTable := FileTable{}
	err = stmt.QueryRow(filehash).Scan(&fileTable.FileHash, &fileTable.FileName, &fileTable.FileSize, &fileTable.FileAddr)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return &fileTable, nil
}