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

	encode_pwd := util.Sha1([]byte(password + password_salt))

	isSuccess := db.UserSignup(username, encode_pwd)
	if isSuccess {
		w.Write([]byte("Create user success"))
	} else {
		w.Write([]byte("Create user failed"))
	}
}
