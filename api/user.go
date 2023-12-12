package api

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"regexp"
	"simple-casdoor/object"
	"simple-casdoor/util"
	"time"
)

// AddUser ...
//
//	@Title			AddUser
//	@Tag			User API
//	@Description	add user
//	@Param			name			body		string			true	"名称"
//	@Param			password		body		string			true	"密码"
//	@Param			passwordType	body		string			false	"密码类型"
//	@Param			phone			body		string			true	"电话号码"
//	@Param			countryCode		body		string			false	"国际区号（默认CN）"
//	@Success		200				{object}	object.Response	成功
//	@router			/add-user [post]
func AddUser(c echo.Context) (err error) {
	var user object.User

	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, object.Response{
			Msg: err.Error(),
		})
	}

	if user.Phone == "" {
		return c.JSON(http.StatusBadRequest, object.Response{
			Msg: "手机号不能为空!",
		})
	}
	if user.CountryCode == "" {
		user.CountryCode = "CN"
	}

	if user.PasswordType == "" {
		user.PasswordType = "md5"
	}
	user.PasswordSalt = util.RandomString(12)

	user.Password = util.GenMd5(user.PasswordSalt, user.Password)
	msg := checkUsername(user.Name)
	if msg != "" {
		return c.JSON(http.StatusBadRequest, object.Response{
			Msg: msg,
		})
	}

	user.CreatedAt = time.Now()
	affected, err := object.AddUser(&user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, object.Response{
			Msg: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, object.Response{
		Code: http.StatusOK,
		Msg:  fmt.Sprintf("affected:%d ", affected),
	})
}

func UpdateUser(c echo.Context) (err error) {
	var user object.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, object.Response{
			Msg: err.Error(),
		})
	}
	if user.CountryCode == "" {
		user.CountryCode = "CN"
	}

	if len(user.Password) > 0 {
		user.PasswordType = "md5"
		user.PasswordSalt = util.RandomString(12)
		user.Password = util.GenMd5(user.PasswordSalt, user.Password)
	}

	if user.Name != "" {
		if msg := checkUsername(user.Name); len(msg) > 0 {
			return c.JSON(http.StatusBadRequest, object.Response{
				Msg: msg,
			})
		}
	}

	if len(user.Phone) > 0 {
		if _, existed, _ := object.GetUserByPhone(user.Phone, true); existed {
			return c.JSON(http.StatusBadRequest, object.Response{
				Msg: "手机号码已被使用，请更换手机号码！",
			})
		}
	}

	if smsCode := c.QueryParam("smsCode"); len(smsCode) > 0 {
		checkPhone, ok := util.GetE164Number(user.Phone, user.CountryCode)
		if !ok {
			return c.JSON(http.StatusBadRequest, object.Response{
				Msg: fmt.Sprintf("verification:Phone %s is invalid in your region %s", user.Phone, user.CountryCode),
			})
		}

		checkResult := object.CheckSignInCode(checkPhone, smsCode)
		if len(checkResult) > 0 {
			return c.JSON(http.StatusBadRequest, object.Response{Msg: checkResult})
		}

		_ = object.DisableVerificationCode(checkPhone)
	}

	affected, err := object.UpdateUser(&user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, object.Response{
			Msg: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, object.Response{
		Code: http.StatusOK,
		Msg:  fmt.Sprintf("affected:%d ", affected),
	})

}

func checkUsername(username string) string {
	if username == "" {
		return "检查:空用户名."
	} else if len(username) > 39 {
		return "检查:用户名太长（最多39个字符）."
	}

	exclude, _ := regexp.Compile("^[\u0021-\u007E]+$")
	if !exclude.MatchString(username) {
		return ""
	}
	return ""
}

func IsUserExists(c echo.Context) (err error) {
	var user object.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, UserResponse{
			Msg:     "",
			Success: false,
		})
	}
	_, exist, err := object.GetUserByPhone(user.Phone, false)

	if err != nil {
		return err
	}

	if !exist {
		return c.JSON(http.StatusBadRequest, UserResponse{
			Msg:     "用户不存在",
			Success: exist,
		})
	}
	return c.JSON(http.StatusOK, UserResponse{
		Msg:     "success",
		Success: exist,
		Code:    http.StatusOK,
	})
}
