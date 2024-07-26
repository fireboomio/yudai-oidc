package api

import (
	"fmt"
	"net/http"
	"yudai/object"
	"yudai/util"

	"github.com/labstack/echo/v4"
)

type SmsMessageForm struct {
	Dest        string            `form:"dest"`
	CountryCode string            `form:"countryCode"`
	Provider    string            `form:"provider"`
	Params      map[string]string `form:"params"`
}

// SendSmsMessage ...
//
//	@Title			SendSmsMessage
//	@Tag			SendSmsMessage API
//	@Description	发送普通短信
//	@Param			dest		body		string			true	"发送手机号"
//	@Param			countryCode	body		string			false	"国际区号（默认CN）"
//	@Param			params		body		string			false	"模板参数"
//	@Success		200			{object}	object.Response	成功
//	@router			/send-verification-code [post]
func SendSmsMessage(c echo.Context) (err error) {
	var vForm SmsMessageForm
	if err = c.Bind(&vForm); err != nil {
		return c.JSON(http.StatusBadRequest, Response{Msg: err.Error()})
	}
	if len(vForm.Dest) == 0 {
		return c.JSON(http.StatusBadRequest, Response{Msg: "用户手机号未提供"})
	}

	// 通过号码获取用户
	user, existed, err := object.GetUserByPhone(vForm.Dest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{Msg: err.Error()})
	}

	if !existed {
		return c.JSON(http.StatusBadRequest, Response{
			Msg: fmt.Sprintf("手机号%s不存在", vForm.Dest),
		})
	}

	// 获取短信提供商
	provider, existed, _ := object.GetSmsProvider(vForm.Provider)
	if !existed {
		return c.JSON(http.StatusBadRequest, Response{
			Msg: fmt.Sprintf("短信提供商:%s不存在", vForm.Provider),
		})
	}

	phone, ok := util.GetE164Number(vForm.Dest, vForm.CountryCode)
	if !ok {
		return c.JSON(http.StatusBadRequest, Response{
			Msg: fmt.Sprintf("您所在地区:%s的电话号码无效", vForm.CountryCode),
		})
	}

	remoteAddr := util.GetIPFromRequest(c.Request())
	if err = object.IsAllowSend(user, remoteAddr); err != nil {
		return c.JSON(http.StatusBadRequest, Response{Msg: err.Error()})
	}

	if err = object.SendSms(provider, vForm.Params, phone); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, Response{Success: true, Msg: "ok", Code: http.StatusOK})
}
