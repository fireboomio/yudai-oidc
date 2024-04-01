package api

import (
	"fmt"
	"net/http"
	"yudai/object"

	"github.com/labstack/echo/v4"
)

type (
	authAction func(authForm *AuthForm) (user *object.User, err error)
	AuthForm   struct {
		object.PlatformConfig
		LoginType string `json:"loginType"`

		Username string `json:"username,omitempty"`
		Password string `json:"password,omitempty"`

		Phone string `json:"phone,omitempty"`
		Code  string `json:"code,omitempty"`
	}

	Response struct {
		Msg     string `json:"msg"`
		Code    int    `json:"code"`
		Success bool   `json:"success"`
		Data    any    `json:"data,omitempty"`
	}
	UserResponse struct {
		Msg     string           `json:"msg"`
		Code    int              `json:"code"`
		Success bool             `json:"success"`
		Data    *object.TokenRes `json:"data,omitempty"`
	}
	refreshInput struct {
		RefreshToken string `json:"refresh_token"`
		object.PlatformConfig
	}
)

var authActionMap map[string]authAction

func init() {
	authActionMap = make(map[string]authAction)
}

// Login ...
//
//	@Title			Login
//	@Tag			Login API
//	@Description	login
//	@Param			username	body		string			true	"用户名"
//	@Param			phone		body		string			true	"号码"
//	@Param			countryCode	body		string			false	"国际区号（默认CN）"
//	@Param			code		body		string			true	"验证码"
//	@Param			password	body		string			true	"密码"
//	@Param			loginType	body		string			true	"登录类型"
//	@Success		200			{object}	UserResponse	成功
//	@router			/login [post]
func Login(c echo.Context) (err error) {
	authForm := new(AuthForm)
	if err = c.Bind(authForm); err != nil {
		return c.JSON(http.StatusBadRequest, Response{Msg: err.Error()})
	}

	action, ok := authActionMap[authForm.LoginType]
	if !ok {
		return c.JSON(http.StatusBadRequest, Response{Msg: fmt.Sprintf("不支持的登录类型：%s", authForm.LoginType)})
	}

	user, err := action(authForm)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{Msg: err.Error()})
	}

	tokenRes, err := object.GenerateToken(user.Transform(), authForm.PlatformConfig)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{Msg: err.Error()})
	}

	return c.JSON(http.StatusOK, &UserResponse{
		Msg:     "Login Success",
		Success: true,
		Data:    tokenRes,
		Code:    http.StatusOK,
	})
}

// RefreshToken
//
//	@Title			RefreshToken
//	@Description	刷新token
//	@Param			refresh-token	body		string			true	"refresh-token"
//	@Success		200				{object}	UserResponse	成功
//	@router			/refresh-token [post]
func RefreshToken(c echo.Context) (err error) {
	var jsonInput refreshInput
	if err = c.Bind(&jsonInput); err != nil {
		return c.JSON(http.StatusBadRequest, Response{Msg: err.Error()})
	}

	if len(jsonInput.RefreshToken) == 0 {
		return c.JSON(http.StatusBadRequest, Response{Msg: "refresh_token为空"})
	}

	claims, err := object.ParseToken(jsonInput.RefreshToken, func() object.Token {
		return object.Token{RefreshToken: jsonInput.RefreshToken}
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{Msg: fmt.Sprintf("token解析错误(%s)", err.Error())})
	}

	tokenRes, err := object.GenerateToken(claims.User, jsonInput.PlatformConfig)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{Msg: fmt.Sprintf("token生成错误(%s)", err.Error())})
	}

	return c.JSON(http.StatusOK, &UserResponse{
		Msg:     "Refresh Success",
		Success: true,
		Data:    tokenRes,
		Code:    http.StatusOK,
	})
}
