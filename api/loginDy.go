package api

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"yudai/object"

	jsoniter "github.com/json-iterator/go"
)

var dyLoginActions = make(map[string]*loginAction)

func loginDy(actionType, code string) (user *object.User, err error) {
	action, ok := dyLoginActions[actionType]
	if !ok {
		return
	}

	dyConfig := action.configHandle()
	if dyConfig == nil {
		err = fmt.Errorf("no config found for %s", actionType)
		return
	}

	var resp *http.Response
	appid, secret := dyConfig.Appid, dyConfig.Secret
	if action.bodyFormat != "" {
		bodyBytes := []byte(fmt.Sprintf(action.bodyFormat, appid, secret, code))
		resp, err = http.Post(action.url, echo.MIMEApplicationJSON, bytes.NewReader(bodyBytes))
	} else {
		resp, err = http.Get(fmt.Sprintf(action.url, appid, secret, code))
	}
	if err != nil {
		return
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if errmsg := jsoniter.Get(respBody, "err_tips").ToString(); len(errmsg) > 0 && errmsg != "success" {
		err = errors.New(errmsg)
		return
	}

	dyLoginResp, err := action.respHandle([]byte(jsoniter.Get(respBody, "data").ToString()))
	if err != nil {
		return
	}

	userSocial, existed, err := object.GetUserSocialByProviderUserId(dyLoginResp.openid)
	if err != nil {
		return
	}

	if existed {
		if userSocial.UserId != "" {
			user, _, err = object.GetUserByUserId(userSocial.UserId)
		}
	} else {
		userSocial = &object.UserSocial{
			Provider:         "douyin",
			ProviderUserId:   dyLoginResp.openid,
			ProviderUnionid:  dyLoginResp.unionid,
			ProviderPlatform: actionType,
		}
		_, err = object.AddUserUserSocial(userSocial)
	}
	if err != nil {
		return
	}

	if user == nil {
		user = &object.User{UserId: dyLoginResp.openid}
	}
	user.SocialUserId = dyLoginResp.openid
	return
}
