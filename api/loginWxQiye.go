package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
	"yudai/object"
)

const (
	qyPc         = "qy_pc"
	qyPcTokenUrl = "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s"
	qyPcUrl      = "https://qyapi.weixin.qq.com/cgi-bin/auth/getuserinfo?access_token=%s&code=%s"
)

type (
	qyPcToken struct {
		Errcode     int    `json:"errcode"`
		Errmsg      string `json:"errmsg"`
		AccessToken string `json:"access_token"`
		ExpireIn    int    `json:"expires_in"`

		createdAt time.Time
	}
	QyPcResp struct {
		Userid     string `json:"userid"`
		UserTicket string `json:"user_ticket"`
	}
)

var qcPcTokenInstance = &qyPcToken{}

func init() {
	wxLoginActions[qyPc] = &loginAction{
		url: qyPcUrl,
		configHandle: func() (*object.LoginConfiguration, error) {
			qyPcConf := object.Conf.WxLogin[qyPc]
			if qyPcConf == nil {
				return nil, nil
			}

			now := time.Now()
			if qcPcTokenInstance.createdAt.IsZero() || int(now.Sub(qcPcTokenInstance.createdAt).Seconds()) > qcPcTokenInstance.ExpireIn {
				resp, err := http.Get(fmt.Sprintf(qyPcTokenUrl, qyPcConf.Appid, qyPcConf.Secret))
				if err != nil {
					return nil, err
				}
				defer func() { _ = resp.Body.Close() }()
				respBody, err := io.ReadAll(resp.Body)
				if err != nil {
					return nil, err
				}
				if err = json.Unmarshal(respBody, &qcPcTokenInstance); err != nil {
					return nil, err
				}
				if qcPcTokenInstance.Errcode != 0 {
					return nil, errors.New(qcPcTokenInstance.Errmsg)
				}
				qcPcTokenInstance.createdAt = now
				qyPcConf.AccessToken = qcPcTokenInstance.AccessToken
			}
			return qyPcConf, nil
		},
		respHandle: func(bytes []byte) (result *loginActionResult, err error) {
			var resp QyPcResp
			if err = json.Unmarshal(bytes, &resp); err != nil {
				return
			}

			result = &loginActionResult{
				openid: resp.Userid,
				data:   &resp,
			}
			return
		},
	}
	authActionMap[qyPc] = func(authForm *AuthForm) (user *object.User, err error) {
		return loginWx(qyPc, authForm.Code)
	}
}
