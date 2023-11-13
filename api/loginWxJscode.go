package api

import (
	"encoding/json"
	"simple-casdoor/object"
)

const (
	jscode    = "jscode"
	jscodeUrl = "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"
)

type jscodeResp struct {
	SessionKey string `json:"session_key"`
	UnionId    string `json:"unionid"`
	Openid     string `json:"openid"`
}

func init() {
	loginActions[jscode] = &loginAction{
		url: jscodeUrl,
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
	authActionMap[jscode] = func(authForm *AuthForm) (user *object.User, err error) {
		return loginWx(jscode, authForm.Code)
	}
}
