package api

import (
	"fmt"
	"net/http"
	"time"
	"yudai/object"
	"yudai/util"

	"github.com/labstack/echo/v4"
)

func Register(c echo.Context) (err error) {
	var user object.User

	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, object.Response{
			Msg: err.Error(),
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

	//添加用户
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
