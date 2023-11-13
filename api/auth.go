package api

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"simple-casdoor/object"
)

var (
	authActionMap map[string]*authAction
)

type (
	authAction struct {
		login func(authForm *AuthForm) (user *object.User, err error)
	}
	AuthForm struct {
		LoginType       string `json:"loginType"`
		Username        string `json:"username,omitempty"`
		Password        string `json:"password,omitempty"`
		Phone           string `json:"phone,omitempty"`
		Code            string `json:"code,omitempty"`
		CountryCode     string `json:"countryCode,omitempty"`
		MiniProGramCode string `json:"miniProGramCode,omitempty"`
		AppId           string `json:"appId,omitempty"`
		Secret          string `json:"secret,omitempty"`
	}
	UserResponse struct {
		Success bool             `json:"success"`
		Msg     string           `json:"msg"`
		Data    *object.TokenRes `json:"data,omitempty"`
		Code    int              `json:"code"`
	}
)

func init() {
	authActionMap = make(map[string]*authAction, 0)
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
	if err := c.Bind(authForm); err != nil {
		return c.JSON(http.StatusBadRequest, object.Response{
			Msg:  err.Error(),
			Code: http.StatusBadRequest,
		})
	}

	if authForm.CountryCode == "" {
		authForm.CountryCode = "CN"
	}

	var user *object.User

	user, err = authActionMap[authForm.LoginType].login(authForm)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, object.Response{
			Msg:  err.Error(),
			Code: http.StatusInternalServerError,
		})
	}

	resp := handleLoggedIn(user)
	return c.JSON(http.StatusOK, resp)
}

func handleLoggedIn(user *object.User) (resp *UserResponse) {
	tokenRes, err := object.GenerateToken(user)

	if err != nil {
		resp = &UserResponse{
			Msg:     fmt.Sprintf("generate token failed: %v", err),
			Success: false,
			Data:    &object.TokenRes{},
			Code:    http.StatusInternalServerError,
		}
	}

	resp = &UserResponse{
		Msg:     "generate token success",
		Success: true,
		Data:    tokenRes,
		Code:    http.StatusOK,
	}
	return resp
}

// RefreshToken
//
//	@Title			RefreshToken
//	@Description	刷新token
//	@Param			refresh-token	body		string			true	"refresh-token"
//	@Success		200				{object}	UserResponse	成功
//	@router			/refresh-token [post]
func RefreshToken(c echo.Context) (err error) {
	var jsonInput map[string]string

	if err := c.Bind(&jsonInput); err != nil {
		return c.JSON(http.StatusBadRequest, object.Response{
			Msg: err.Error(),
		})
	}

	refreshToken, exist := jsonInput["refresh_token"]
	if !exist {
		return c.JSON(400, object.Response{
			Msg: "refresh-token不存在",
		})
	}

	claims, err := object.ParseToken(refreshToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, object.Response{
			Msg: "token解析错误，请重新登录",
		})
	}

	resp := handleLoggedIn(claims.User)
	return c.JSON(http.StatusOK, resp)
}
