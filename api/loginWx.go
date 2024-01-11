package api

import (
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"io"
	"net/http"
	"yudai/object"
)

type (
	loginAction struct {
		url          string
		configHandle func(*object.WxLoginConfig) *object.WxLoginDetail
		respHandle   func([]byte) (*loginActionResult, error)
	}
	loginActionResult struct {
		openid  string
		unionid string
		data    map[string]any
	}
)

var loginActions = make(map[string]*loginAction)

func loginWx(actionType, code string) (user *object.User, err error) {
	action, ok := loginActions[actionType]
	if !ok {
		return
	}

	wxConfig := action.configHandle(object.Conf.WxLoginConfig)
	appid, secret := wxConfig.Appid, wxConfig.Secret
	resp, err := http.Get(fmt.Sprintf(action.url, appid, secret, code))
	if err != nil {
		return
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if errmsg := jsoniter.Get(respBody, "errmsg").ToString(); len(errmsg) > 0 {
		err = errors.New(errmsg)
		return
	}

	wxLoginResp, err := action.respHandle(respBody)
	if err != nil {
		return
	}

	userSocial, existed, err := object.GetUserSocialByProviderUserId(wxLoginResp.openid)
	if err != nil {
		return
	}

	if existed {
		if userSocial.UserId != "" {
			user, _, err = object.GetUserByUserId(userSocial.UserId)
		}
	} else {
		userSocial = &object.UserSocial{
			Provider:         "weixin",
			ProviderUserId:   wxLoginResp.openid,
			ProviderUnionid:  wxLoginResp.unionid,
			ProviderPlatform: actionType,
		}
		_, err = object.AddUserUserSocial(userSocial)
	}
	if err != nil {
		return
	}

	if user == nil {
		user = &object.User{UserId: wxLoginResp.openid}
	}
	user.SocialUserId = wxLoginResp.openid
	return
}
