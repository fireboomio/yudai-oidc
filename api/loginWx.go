package api

import (
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"io"
	"net/http"
	"simple-casdoor/object"
)

type (
	loginAction struct {
		url        string
		respHandle func([]byte) (*loginActionResult, error)
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

	wxConfig := object.Conf.WxLoginConfig
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

	result, err := action.respHandle(respBody)
	if err != nil {
		return
	}

	user, err = object.GetUserByWxUnionid(result.unionid)
	if err != nil {
		return
	}

	if user == nil {
		user = &object.User{
			Name:      "WxUser",
			WxUnionid: result.unionid,
		}
		if _, err = object.AddUser(user); err != nil {
			return
		}
	}

	_, _ = object.AddUserWx(&object.UserWx{Platform: actionType, Unionid: result.unionid, Openid: result.openid})
	user.WxResp = result.data
	return
}
