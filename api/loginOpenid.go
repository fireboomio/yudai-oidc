package api

import (
	"simple-casdoor/object"
)

const openid = "openid"

func init() {
	authActionMap[openid] = &authAction{
		login: func(authForm *AuthForm) (user *object.User, err error) {
			user, err = object.GetUserByWxOpenid(authForm.Code)
			if err != nil {
				return
			}

			if user == nil {
				user = &object.User{
					Name:     "WxUser",
					WxOpenid: authForm.Code,
				}
				if _, err = object.AddUser(user); err != nil {
					return
				}
			}
			return
		},
	}
}
