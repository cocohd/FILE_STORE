package db

import (
	mydb "FILE_STORE/db/mysql"
	"fmt"
)

// 注册用户信息加入用户表
func UserSignUpDb(username string, passwd string) bool {
	stmt, err := mydb.DBConn().Prepare("insert ignore into tbl_user (`user_name`, `user_pwd`) values (?, ?)")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer stmt.Close()

	ret, err := stmt.Exec(username, passwd)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	if rowAffect, err := ret.RowsAffected(); nil == err && rowAffect > 0 {
		return true
	}
	return false
}


// 从用户表查询用户及密码是否正确
func UserSignInDb(userName string, encPwd string) bool  {
	stmt, err := mydb.DBConn().Prepare("select * from tbl_user where user_name=? limit 1")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	defer stmt.Close()

	rows, err := stmt.Query(userName)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}else if rows == nil {
		fmt.Println("username not found:" + userName)
		return false
	}

	pRows := mydb.ParseRows(rows)
	// 此处.([]byte)不懂
	if len(pRows) > 0 && string(pRows[0]["user_pwd"].([]byte)) == encPwd {
		return true
	}
	return false
}


// 存储token
func UserTokenDb(userName string, userToken string) bool {
	// replace指遇到相同的直接替换，token可重复替换
	stmt, err := mydb.DBConn().Prepare("replace into tbl_user_token (`user_name`, `user_token`) values (?, ?)")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	defer stmt.Close()

	// 此处不许处理返回值，仅需保证token更新即可
	_, err = stmt.Exec(userName, userToken)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}


type User struct {
	Username string
	Email string
	Phone string
	SignupAt string
	LastActiveAt string
	status int
}


func GetUserInfoDb(userName string) (User, error) {
	user := User{}

	stmt, err := mydb.DBConn().Prepare("select user_name, signup_at from tbl_user where user_name=? limit 1")
	if err != nil {
		fmt.Println(err.Error())
		return user, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(userName).Scan(&user.Username, &user.SignupAt)
	if err != nil {
		fmt.Println(err.Error())
		return user, err
	}
	return user, nil
}