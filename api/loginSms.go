package api

import (
	"fmt"
	"yudai/object"
	"yudai/util"
)

const SMS = "sms"

func init() {
	// 通过手机短信验证码登录
	authActionMap[SMS] = &authAction{
		action: func(authForm *AuthForm) (user *object.User, err error) {
			user, existed, err := object.GetUserByPhone(authForm.Phone)
			if err != nil {
				return
			}

			if !existed {
				err = fmt.Errorf("手机号不存在: %s", authForm.Phone)
				return
			}

			err = checkAndDisableCode(authForm.Phone, authForm.Code, user.CountryCode)
			return
		},
	}
}

func checkAndDisableCode(phone, code, countryCode string) (err error) {
	if len(phone) == 0 {
		return fmt.Errorf("手机号码不能为空")
	}

	if len(countryCode) == 0 {
		countryCode = "CN"
	}

	// 转成E.164格式的电话号码
	dest, ok := util.GetE164Number(phone, countryCode)
	if !ok {
		err = fmt.Errorf("您所在地区:%s的电话号码无效", countryCode)
		return
	}

	// 获取最新一条未被使用的验证码进行验证
	if checkResult := object.CheckSignInCode(dest, code); len(checkResult) != 0 {
		err = fmt.Errorf(checkResult)
		return
	}

	// disable the verification code
	go func() { _ = object.DisableVerificationCode(dest) }()
	return
}
