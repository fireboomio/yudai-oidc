package api

import (
	"encoding/json"
	"yudai/object"
)

const (
	dyMini             = "dy_mini"
	dyJscodeUrl        = "https://developer.toutiao.com/api/apps/v2/jscode2session"
	dyJscodeBodyFormat = `{"appid": "%s", "secret": "%s", "code": "%s"}`
)

func init() {
	dyLoginActions[dyMini] = &loginAction{
		url:        dyJscodeUrl,
		bodyFormat: dyJscodeBodyFormat,
		configHandle: func() (*object.LoginConfiguration, error) {
			return object.Conf.DyLogin[dyMini], nil
		},
		respHandle: func(bytes []byte) (result *loginActionResult, err error) {
			var resp jscodeResp
			if err = json.Unmarshal(bytes, &resp); err != nil {
				return
			}

			result = &loginActionResult{
				unionid: resp.UnionId,
				openid:  resp.Openid,
				data: map[string]any{
					"session_key": resp.SessionKey,
				},
			}
			return
		},
	}
	authActionMap[dyMini] = func(authForm *AuthForm) (user *object.User, err error) {
		return loginDy(dyMini, authForm.Code)
	}
}
