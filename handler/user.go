package handler

import (
	"net-disk-server/db"
	"net-disk-server/util"
	"net/http"
)

const password_salt = "**$$110"

// 用户注册
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")

	if len(username) < 3 || len(password) < 5 {
		w.Write([]byte("Invalid parameters"))
		return
	}

	encodePwd := util.Sha1([]byte(password + password_salt))

	isSuccess := db.UserSignup(username, encodePwd)
	if isSuccess {
		w.Write([]byte("Create user success"))
	} else {
		w.Write([]byte("Create user failed"))
	}
}

// 登录接口
func SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	// 校验账号密码
	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	encodePwd := util.Sha1([]byte(password + password_salt))

	checkdPassword := db.UserSignIn(username, encodePwd)

	if !checkdPassword {
		w.Write([]byte("signin failed"))
		return
	}

	// 生成 token
	token := db.GenerateToken(username)
	updateRes := db.UpdateToken(username, token)
	if !updateRes {
		w.Write([]byte("Signin failed"))
		return
	}

	// 登录成功
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: struct {
			Username string
			Token    string
		}{
			Username: username,
			Token:    token,
		},
	}

	w.Write(resp.JSONBytes())
}

// 查询用户信息
func UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	// 解析参数
	r.ParseForm()
	username := r.Form.Get("username")

	// 使用拦截器校验 token

	// 查询用户信息
	user, err := db.GetUserInfo(username)

	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	// 响应
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: user,
	}

	w.Write(resp.JSONBytes())
}

// 判断 token 有效性
func IsTokenValid(token string) bool {
	// 校验长度
	if len(token) != 40 {
		return false
	}
	// 判断 token 是否过期

	// 从数据库 token_table 查询是否有用户的 token 信息

	// 直接对比两个 token 是否一致

	return true
}
