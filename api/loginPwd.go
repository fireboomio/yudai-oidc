package api

import (
	"yudai/object"
)

const PWD = "password"

func init() {
	// 通过密码登录
	authActionMap[PWD] = func(authForm *AuthForm) (user *object.User, err error) {
		user, err = object.CheckUserPassword(authForm.Username, authForm.Password)
		if err != nil {
			return nil, err
		}

		return
	}

}
