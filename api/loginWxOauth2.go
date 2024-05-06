package api

import (
	"encoding/json"
	"yudai/object"
)

const (
	oauth2Pc  = "pc"
	oauth2H5  = "h5"
	oauth2App = "app"
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
	respHandle := func(bytes []byte) (result *loginActionResult, err error) {
		var resp oauthResp
		if err = json.Unmarshal(bytes, &resp); err != nil {
			return
		}

		result = &loginActionResult{
			unionid: resp.UnionId,
			openid:  resp.Openid,
			data:    &resp,
		}
		return
	}
	wxLoginActions[oauth2Pc] = &loginAction{
		url:        oauth2Url,
		respHandle: respHandle,
		configHandle: func() (*object.LoginConfiguration, error) {
			return object.Conf.WxLogin[oauth2Pc], nil
		},
	}
	wxLoginActions[oauth2H5] = &loginAction{
		url:        oauth2Url,
		respHandle: respHandle,
		configHandle: func() (*object.LoginConfiguration, error) {
			return object.Conf.WxLogin[oauth2H5], nil
		},
	}
	wxLoginActions[oauth2App] = &loginAction{
		url:        oauth2Url,
		respHandle: respHandle,
		configHandle: func() (*object.LoginConfiguration, error) {
			return object.Conf.WxLogin[oauth2App], nil
		},
	}
	authActionMap[oauth2Pc] = &authAction{
		action: func(authForm *AuthForm) (user *object.User, err error) {
			return loginWx(oauth2Pc, authForm.Code)
		},
		setting: func() *object.LoginConfiguration {
			return object.Conf.WxLogin[oauth2Pc]
		},
	}
	authActionMap[oauth2H5] = &authAction{
		action: func(authForm *AuthForm) (user *object.User, err error) {
			return loginWx(oauth2H5, authForm.Code)
		},
		setting: func() *object.LoginConfiguration {
			return object.Conf.WxLogin[oauth2H5]
		},
	}
	authActionMap[oauth2App] = &authAction{
		action: func(authForm *AuthForm) (user *object.User, err error) {
			return loginWx(oauth2App, authForm.Code)
		},
		setting: func() *object.LoginConfiguration {
			return object.Conf.WxLogin[oauth2App]
		},
	}
}
