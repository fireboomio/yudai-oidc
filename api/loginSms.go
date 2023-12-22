package api

import (
	"fmt"
	"yudai/object"
	"yudai/util"
)

const SMS = "sms"

func init() {
	// 通过手机短信验证码登录
	authActionMap[SMS] = func(authForm *AuthForm) (user *object.User, err error) {
		// 转成E.164格式的电话号码
		dest, ok := util.GetE164Number(authForm.Phone, authForm.CountryCode)
		if !ok {
			return nil, fmt.Errorf("您所在地区:%s的电话号码无效", authForm.CountryCode)
		}
		// 获取最新一条未被使用的验证码进行验证
		checkResult := object.CheckSignInCode(dest, authForm.Code)

		if len(checkResult) != 0 {
			return nil, fmt.Errorf(checkResult)
		}

		// disable the verification code
		go func() {
			err = object.DisableVerificationCode(dest)
		}()

		// get login user
		user, _, err = object.GetUserByPhone(authForm.Phone, false)
		return
	}
}
