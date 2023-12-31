package api

import (
	"fmt"
	"net/http"
	"yudai/object"
	"yudai/util"

	"github.com/labstack/echo/v4"
)

type VerificationForm struct {
	Dest        string `form:"dest"`
	CountryCode string `form:"countryCode"`
	Provider    string `form:"provider"`
}

// SendVerificationCode ...
//
//	@Title			SendVerificationCode
//	@Tag			Verification API
//	@Description	发送验证码
//	@Param			dest		body		string			true	"发送手机号"
//	@Param			countryCode	body		string			false	"国际区号（默认CN）"
//	@Success		200			{object}	object.Response	成功
//	@router			/send-verification-code [post]
func SendVerificationCode(c echo.Context) (err error) {
	var vForm VerificationForm
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
		user = &object.User{
			Name:        vForm.Dest,
			Phone:       vForm.Dest,
			CountryCode: vForm.CountryCode,
			Password:    vForm.Dest,
		}
		defer func() {
			if c.Response().Status == http.StatusOK {
				_, _ = object.AddUser(user)
			}
		}()
	}

	smsProvider := vForm.Provider
	if len(smsProvider) == 0 {
		smsProvider = "fireboom/provider_sms"
	}
	// 获取短信提供商
	provider, existed, _ := object.GetSmsProvider(smsProvider)
	if !existed {
		return c.JSON(http.StatusBadRequest, Response{
			Msg: fmt.Sprintf("短信提供商:%s不存在", smsProvider),
		})
	}

	phone, ok := util.GetE164Number(vForm.Dest, vForm.CountryCode)
	if !ok {
		return c.JSON(http.StatusBadRequest, Response{
			Msg: fmt.Sprintf("您所在地区:%s的电话号码无效", vForm.CountryCode),
		})
	}

	remoteAddr := util.GetIPFromRequest(c.Request())
	if err = object.SendVerificationCodeToPhone(user, provider, remoteAddr, phone); err != nil {
		return c.JSON(http.StatusBadRequest, Response{Msg: err.Error()})
	}

	return c.JSON(http.StatusOK, Response{Msg: "ok", Code: http.StatusOK})
}
