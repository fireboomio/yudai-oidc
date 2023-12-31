package api

import (
	"yudai/object"
)

const PWD = "password"

func init() {
	// 通过密码登录
	authActionMap[PWD] = func(authForm *AuthForm) (*object.User, error) {
		return object.CheckUserPassword(authForm.Username, authForm.Password)
	}
}
