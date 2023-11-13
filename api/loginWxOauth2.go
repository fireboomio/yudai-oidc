package api

import (
	"encoding/json"
	"simple-casdoor/object"
)

const (
	oauth2Pc  = "oauth2_pc"
	oauth2H5  = "oauth2_h5"
	oauth2App = "oauth2_app"
	oauth2Url = "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code"
)

type oauthResp struct {
	AccessToken    string `json:"access_token"`
	ExpiresIn      int    `json:"expires_in"`
	RefreshToken   string `json:"refresh_token"`
	Openid         string `json:"openid"`
	Scope          string `json:"scope"`
	IsSnapshotUser int    `json:"is_snapshotuser"`
	UnionId        string `json:"unionid"`
}

func init() {
	oauth2Action := &loginAction{
		url: oauth2Url,
		respHandle: func(bytes []byte) (result *loginActionResult, err error) {
			var resp oauthResp
			if err = json.Unmarshal(bytes, &resp); err != nil {
				return
			}

			result = &loginActionResult{
				unionid: resp.UnionId,
				openid:  resp.Openid,
				data: map[string]any{
					"access_token":  resp.AccessToken,
					"refresh_token": resp.RefreshToken,
					"expires_in":    resp.ExpiresIn,
				},
			}
			return
		},
	}
	loginActions[oauth2Pc] = oauth2Action
	loginActions[oauth2H5] = oauth2Action
	loginActions[oauth2App] = oauth2Action
	authActionMap[oauth2Pc] = func(authForm *AuthForm) (user *object.User, err error) {
		return loginWx(oauth2Pc, authForm.Code)
	}
	authActionMap[oauth2H5] = func(authForm *AuthForm) (user *object.User, err error) {
		return loginWx(oauth2H5, authForm.Code)
	}
	authActionMap[oauth2App] = func(authForm *AuthForm) (user *object.User, err error) {
		return loginWx(oauth2App, authForm.Code)
	}
}
