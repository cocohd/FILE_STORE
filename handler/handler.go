package handler

import (
	dblayer "FILE_STORE/db"
	"FILE_STORE/meta"
	"FILE_STORE/util"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)


// 文件上传
func UploadHandler(w http.ResponseWriter, r *http.Request)  {
	if r.Method == "GET" {
		// 返回上传html页面
		data, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			io.WriteString(w, "internal server error html!")
			return
		}
		// 成功将data传回去
		io.WriteString(w, string(data))
	} else if r.Method == "POST" {
		// 接收文件流、存储到本地目录
		// 表单form，返回三参数file，header, err
		file, header, err := r.FormFile("file")
		if err != nil {
			fmt.Printf("Failed to get data, err: %s\n", err.Error())
			return
		}
		defer file.Close()


		// 获取元信息
		fileMeta := meta.FileMeta{
			FileName: header.Filename,
			Location: "./tmp/" + header.Filename,
			UploadAt: time.Now().Format("2006-01-02 15:04:05"),
		}

		// 创建
		//newFile, err := os.Create("./tmp/" + header.Filename)
		newFile, err := os.Create("./tmp/" + header.Filename)
		if err != nil {
			fmt.Printf("Failed to create file, err: %s\n", err.Error())
			return
		}
		defer  newFile.Close()

		// 将内存中的文件拷贝到新的文件的buff区
		// 返回文件的字节信息，err
		fileMeta.FileSize, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Printf("Failed to save data into file, err: %s\n", err.Error())
			return
		}

		// 生成fileSha1，先定位到文件头部，在传入生成哈希

		newFile.Seek(0, 0)
		fileMeta.FileSha1 = util.FileSha1(newFile)
		meta.UpdateFileMeta(fileMeta)

		//存储到tbl_user_file表中
		r.ParseForm()
		username := r.Form.Get("username")
		fmt.Println(username)
		suc := dblayer.OnFileUploadFinishedDb(username, fileMeta.FileSha1, fileMeta.FileName, fileMeta.FileSize)
		if suc {
			http.Redirect(w, r, "/static/view/home.html", http.StatusFound)
		}else {
			w.Write([]byte("Upload Failed"))
		}

		// 更行到mysql中
		meta.UpdateFileDB(fileMeta)

		//重定向到上传成功页
		http.Redirect(w, r,"/file/upload/suc", http.StatusFound)
		
	}
}


// 上传成功页面
func UploadSucHandler(w http.ResponseWriter, t *http.Request)  {
	io.WriteString(w, "Upload finish")
}


// 通过fileSha1查看先骨干meta
func GetFileMetaHandler(w http.ResponseWriter, r *http.Request)  {
	r.ParseForm()
	// 取出请求中的参数(url中的)
	//filehash := r.Form["filehash"][0]
	//filehash := r.Form.Get("filehash")

	//fMeta := meta.GetFileMeta(filehash)
	// 从数据库中查询
	//fMeta, err := meta.GetFileDB((filehash))

	userName := r.Form.Get("username")
	// 从tbl_user_file表查询数据
	fMeta, err := dblayer.QueryUserFileDb(userName, 5)
	if err != nil {
		fmt.Printf(err.Error())
	}

	//将其转换为JSON
	data, err := json.Marshal(fMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}


// 下载文件
func DownloadFileHandler(w http.ResponseWriter, r *http.Request)  {
	r.ParseForm()
	fileSha1 := r.Form.Get("filehash")
	fm := meta.GetFileMeta(fileSha1)

	f, err := os.Open(fm.Location)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	// 客户端读取出来即为下载
	// 小文件可以直接用readAll读出来，否则用buff
	data, err := ioutil.ReadAll(f)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// app的话，到这已经ok了；浏览器需加上返回头信息
	w.Header().Set("Content-Type", "application/octect-stream")
	w.Header().Set("Content-Descrption", "attachment;filename=\"" + fm.FileName + "\"")
	w.Write(data)

}


// update操作：重命名
func UpdateFileHandler(w http.ResponseWriter, r *http.Request)  {
	r.ParseForm()
	fileSha1 := r.Form.Get("fileSha1")
	newFileName := r.Form.Get("newFileName")

	fileMeta := meta.GetFileMeta(fileSha1)
	fileMeta.FileName = newFileName
	meta.UpdateFileMeta(fileMeta)

	data, err := json.Marshal(fileMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func DeleteFileHandler(w http.ResponseWriter, r *http.Request)  {
	r.ParseForm()
	fileSha1 := r.Form.Get("fileSha1")
	meta.RemoveFileMeta(fileSha1)
	fileMeta := meta.GetFileMeta(fileSha1)
	os.Remove(fileMeta.Location)

	w.WriteHeader(http.StatusOK)
}

func TryFastUploadHandler(w http.ResponseWriter, r *http.Request)  {
	// 1. 解析请求参数
	r.ParseForm()
	userName := r.Form.Get("username")
	fileHash := r.Form.Get("filehash")
	fileName := r.Form.Get("fileName")
	// 有err返回值
	fileSize, err:= strconv.Atoi(r.Form.Get("fileSize"))
	// 2. 从文件表中查询相同hash的文件记录
	fileMeta, err := meta.GetFileDB(fileHash)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
	// 3. 查不到记录则返回秒传失败
	//fileMeta == nil 这样不行，定义的结构体这样判断不行
	if fileMeta == (meta.FileMeta{}) {
		resp := util.RespMsg{
			Code: 1,
			Msg: "秒传失败， 请访问普通接口",
		}
		w.Write(resp.JSONBytes())
		return
	}
	// 4. 上传过则将文件信息写入用户文件表，返回成功
	// 要将string转换为int64的,虚线转换为int型，通过strconv.Atoi()
	suc := dblayer.OnFileUploadFinishedDb(userName, fileHash, fileName, int64(fileSize))
	if suc {
		resp := util.RespMsg{
			Code: 0,
			Msg: "秒传成功",
		}
		w.Write(resp.JSONBytes())
		return
	}
}










