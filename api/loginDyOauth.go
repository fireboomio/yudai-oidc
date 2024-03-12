package api

import (
	"encoding/json"
	"errors"
	"yudai/object"
)

const (
	dyOauthPc   = "dy_pc"
	dyOauthH5   = "dy_h5"
	dyOauthApp  = "dy_app"
	dyOauth2Url = "https://open.douyin.com/oauth/access_token?client_key=%s&client_secret=%s&code=%s&grant_type=authorization_code"
)

type dyOauthResp struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshToken     string `json:"refresh_token"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	Openid           string `json:"open_id"`
	Scope            string `json:"scope"`
	ErrorCode        int    `json:"error_code"`
	Description      string `json:"description"`
}

func init() {
	respHandle := func(bytes []byte) (result *loginActionResult, err error) {
		var resp dyOauthResp
		if err = json.Unmarshal(bytes, &resp); err != nil {
			return
		}

		if resp.ErrorCode != 0 {
			err = errors.New(resp.Description)
			return
		}

		result = &loginActionResult{
			openid: resp.Openid,
			data:   &resp,
		}
		return
	}
	wxLoginActions[dyOauthPc] = &loginAction{
		url:        dyOauth2Url,
		respHandle: respHandle,
		configHandle: func() *object.LoginConfiguration {
			return object.Conf.DyLogin[dyOauthPc]
		},
	}
	wxLoginActions[dyOauthH5] = &loginAction{
		url:        dyOauth2Url,
		respHandle: respHandle,
		configHandle: func() *object.LoginConfiguration {
			return object.Conf.DyLogin[dyOauthH5]
		},
	}
	wxLoginActions[dyOauthApp] = &loginAction{
		url:        dyOauth2Url,
		respHandle: respHandle,
		configHandle: func() *object.LoginConfiguration {
			return object.Conf.DyLogin[dyOauthApp]
		},
	}
	authActionMap[dyOauthPc] = func(authForm *AuthForm) (user *object.User, err error) {
		return loginDy(dyOauthPc, authForm.Code)
	}
	authActionMap[dyOauthH5] = func(authForm *AuthForm) (user *object.User, err error) {
		return loginDy(dyOauthH5, authForm.Code)
	}
	authActionMap[dyOauthApp] = func(authForm *AuthForm) (user *object.User, err error) {
		return loginDy(dyOauthApp, authForm.Code)
	}
}
