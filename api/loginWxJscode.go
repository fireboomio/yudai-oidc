package api

import (
	"encoding/json"
	"yudai/object"
)

const (
	mini      = "mini"
	jscodeUrl = "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"
)

type jscodeResp struct {
	SessionKey string `json:"session_key"`
	UnionId    string `json:"unionid"`
	Openid     string `json:"openid"`
}

func init() {
	wxLoginActions[mini] = &loginAction{
		url: jscodeUrl,
		configHandle: func() (*object.LoginConfiguration, error) {
			return object.Conf.WxLogin[mini], nil
		},
		respHandle: func(bytes []byte) (result *loginActionResult, err error) {
			var resp jscodeResp
			if err = json.Unmarshal(bytes, &resp); err != nil {
				return
			}

			result = &loginActionResult{
				unionid: resp.UnionId,
				openid:  resp.Openid,
				data:    &resp,
			}
			return
		},
	}
	authActionMap[mini] = func(authForm *AuthForm) (user *object.User, err error) {
		return loginWx(mini, authForm.Code)
	}
}
