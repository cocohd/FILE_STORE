package meta

import (
	mydb "FILE_STORE/db"
	"fmt"
)

// 文件元信息结构
type FileMeta struct {
	// 文件唯一标志
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

var fileMetas map[string]FileMeta

// 初始化
func init() {
	fileMetas = make(map[string]FileMeta)
}

// 修改元信息
func UpdateFileMeta(fmeta FileMeta)  {
	fileMetas[fmeta.FileSha1] = fmeta
}


// 更新文件信息到mysql中
func UpdateFileDB(fmeta FileMeta) bool {
	return mydb.OnFileUploadedFinished(fmeta.FileSha1, fmeta.FileName, fmeta.FileSize, fmeta.Location)
}


// 从mysql中查询信息
func GetFileDB(filesha1 string)  (FileMeta, error) {
	file, err := mydb.GetFile(filesha1)
	if err != nil {
		fmt.Println(err.Error())
		return FileMeta{}, err
	}
	fmeta := FileMeta{
		FileSha1: file.FileHash,
		// 数据库里的类型有些不同，需要转换下
		FileName: file.FileName.String,
		FileSize: file.FileSize.Int64,
		Location: file.FileAddr.String,
	}
	return fmeta, nil
}


// 通过sha1值获取元信息
func GetFileMeta(fileSha1 string) FileMeta  {
	return fileMetas[fileSha1]
}


// 删除fileMeta
func RemoveFileMeta(fileSha1 string)  {
	delete(fileMetas, fileSha1)
}