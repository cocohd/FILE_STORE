// 用户控制逻辑
package handler

import (
	dblayer "FILE_STORE/db"
	"FILE_STORE/util"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

//const (
//	// 用于加密的盐值(自定义)
//	saltData = "*#890"
//)

// 注册
func UserSignUpHandler(w http.ResponseWriter, r *http.Request)  {
	//fmt.Printf("get请求")

	if r.Method == http.MethodGet {
		//data, err := ioutil.ReadFile("./static/view/signup.html")
		//if err != nil {
		//	w.WriteHeader(http.StatusInternalServerError)
		//	return
		//}
		//w.Write(data)
		http.Redirect(w, r, "/static/view/signup.html", http.StatusFound)

		return
	}
	//fmt.Printf("post请求")
	r.ParseForm()
	username := r.Form.Get("username")
	passwd := r.Form.Get("password")

	if len(username) < 5 || len(passwd) < 6 {
		w.Write([]byte("Invalid username or passwd"))
		return
	}

	// 密码加密：
	saltData := username
	encPwd := util.Sha1([]byte(passwd + saltData))

	res := dblayer.UserSignUpDb(username, encPwd)
	if res {
		w.Write([]byte("SUCCESS"))
	}else {
		w.Write([]byte("FAILED"))
	}
}


// 登录
func UserSignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data, err := ioutil.ReadFile("./static/view/signin.html")
		if err != nil {
			w.Write([]byte("Internal Error"))
			return
		}
		w.Write(data)
		return
	}

	r.ParseForm()
	userName := r.Form.Get("username")
	passWd := r.Form.Get("password")
	encPwd := util.Sha1([]byte(passWd + userName))
	fmt.Println(encPwd)

	// 校验用户名、密码
	pwdChecked := dblayer.UserSignInDb(userName, encPwd)
	if !pwdChecked {
		w.Write([]byte("PWD FAILED"))
		return
	}

	// 生成访问凭证token
	token := GenToken(userName)
	tokenChecked := dblayer.UserTokenDb(userName, token)
	if !tokenChecked {
		w.Write([]byte("TOKEN FAILED"))
		return
	}

	// 登陆成功后重定向到首页
	// 需要将token等数据传过去
	//w.Write([]byte("http://" + r.Host + "/static/view/home.html"))
	resp := util.RespMsg{
		Code: 0,
		Msg: "OK",
		Data: struct {
			Location string
			Username string
			Token string
		}{
			Location: "http://" + r.Host + "/static/view/home.html",
			Username: userName,
			Token: token,
		},
	}
	fmt.Println("resp:", resp)
	w.Write(resp.JSONBytes())
}


func UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	// 解析请求参数
	r.ParseForm()
	username := r.Form.Get("username")
	token := r.Form.Get("token")
	fmt.Println("userName:" + username)

	// 验证token是否有效
	isValidToken := IsTokenValid(token)
	if !isValidToken {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// 查询用户信息
	user, err := dblayer.GetUserInfoDb(username)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// 组装并响应用户数据
	resp := util.RespMsg{
		Code: 0,
		Msg: "OK",
		Data: user,
	}
	w.Write(resp.JSONBytes())

}


func GenToken(userName string) string {
	// 40位字符， md5() + timestamp[:8]
	ts := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := util.MD5([]byte(userName + ts + "_tokenSalt"))
	return tokenPrefix + ts[:8]
}


func IsTokenValid(token string) bool {
	// 判断token是否过期，根据token后八位，即时间戳

	// tbl_user_token查询username对应的token

	// 判断是否一致

	return true
}