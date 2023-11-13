package api

import (
	"simple-casdoor/object"
)

const PWD = "password"

func init() {
	// 通过密码登录
	authActionMap[PWD] = &authAction{
		login: func(authForm *AuthForm) (user *object.User, err error) {
			user, err = object.CheckUserPassword(authForm.Phone, authForm.Password)
			if err != nil {
				return nil, err
			}

			return
		},
	}

}
