package api

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"golang.org/x/exp/slices"
	"net/http"
	"yudai/object"
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
	if err = c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, Response{Msg: err.Error()})
	}

	_, existed, err := object.GetUserByName(user.Name)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{Msg: err.Error()})
	}
	if existed {
		return c.JSON(http.StatusBadRequest, Response{Msg: fmt.Sprintf("用户名已存在: %s", user.Name)})
	}

	if len(user.Phone) > 0 {
		_, existed, err = object.GetUserByPhone(user.Phone)
		if err != nil {
			return c.JSON(http.StatusBadRequest, Response{Msg: err.Error()})
		}
		if existed {
			return c.JSON(http.StatusBadRequest, Response{Msg: fmt.Sprintf("手机号已存在: %s", user.Phone)})
		}
	}

	affected, userId, err := object.AddUser(&user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{Msg: err.Error()})
	}

	return c.JSON(http.StatusOK, Response{
		Success: true,
		Code:    http.StatusOK,
		Msg:     fmt.Sprintf("affected:%d ", affected),
		Data:    userId,
	})
}

type updateUserInput struct {
	object.User
	object.PlatformConfig
	Code string `json:"code"`
}

func UpdateUser(c echo.Context) (err error) {
	var updatedUser updateUserInput
	if err = c.Bind(&updatedUser); err != nil {
		return c.JSON(http.StatusBadRequest, Response{Msg: err.Error()})
	}

	var prepareUpdateUser *object.Userinfo
	if len(updatedUser.UserId) > 0 {
		existedUser, existed, _ := object.GetUserByUserId(updatedUser.UserId)
		if !existed {
			return c.JSON(http.StatusBadRequest, Response{Msg: "用户不存在"})
		}
		prepareUpdateUser = existedUser.Transform()
	} else {
		prepareUpdateUser = c.Get("user").(*object.Userinfo)
		updatedUser.UserId = prepareUpdateUser.UserId
	}
	if len(updatedUser.Name) > 0 {
		if updatedUser.Name == prepareUpdateUser.Name {
			updatedUser.Name = ""
		} else {
			var repeated bool
			_, repeated, err = object.GetUserByName(updatedUser.Name)
			if err != nil {
				return c.JSON(http.StatusBadRequest, Response{Msg: err.Error()})
			}
			if repeated {
				return c.JSON(http.StatusBadRequest, Response{Msg: "用户名已存在，请更换用户名！"})
			}
		}
	}

	var changeUserToken *object.TokenRes
	var createRequired bool
	if len(updatedUser.Phone) > 0 {
		if updatedUser.Phone == prepareUpdateUser.Phone {
			updatedUser.Phone = ""
		} else {
			var (
				existedPhoneUser *object.User
				existed          bool
			)
			existedPhoneUser, existed, err = object.GetUserByPhone(updatedUser.Phone)
			if err != nil {
				return c.JSON(http.StatusBadRequest, Response{Msg: err.Error()})
			}

			isSocialUser := len(prepareUpdateUser.SocialUserId) > 0
			if existed && !isSocialUser {
				return c.JSON(http.StatusBadRequest, Response{Msg: "手机号码已被使用，请更换手机号码！"})
			}

			if isSocialUser {
				if createRequired = !existed; existed {
					socialUser, socialExisted, _ := object.GetUserSocialByProviderUserId(prepareUpdateUser.SocialUserId)
					if !socialExisted {
						return c.JSON(http.StatusBadRequest, Response{Msg: "社交用户未找到"})
					}
					// 如果手机号用户存在，则检查对应的social用户是否存在（providerUserId不同但provider和platform相同）
					var socials []*object.UserSocial
					if socials, err = object.GetUserSocialsByUserId(existedPhoneUser.UserId); err != nil {
						return c.JSON(http.StatusBadRequest, Response{Msg: err.Error()})
					}
					if len(socials) > 0 && slices.ContainsFunc(socials, func(social *object.UserSocial) bool {
						return social.ProviderUserId != socialUser.ProviderUserId && social.Provider+social.ProviderPlatform == socialUser.Provider+social.ProviderPlatform
					}) {
						return c.JSON(http.StatusBadRequest, Response{Msg: "手机号码已被使用，请更换手机号码！"})
					}
					if _, err = object.UpdateUserSocial(existedPhoneUser.UserId, socialUser.ProviderUserId); err != nil {
						return c.JSON(http.StatusBadRequest, Response{Msg: err.Error()})
					}
					changeUserToken, _ = object.GenerateToken(existedPhoneUser.Transform(), updatedUser.PlatformConfig)
				}
			}
			if updatedUser.Code != "" {
				if err = checkAndDisableCode(updatedUser.Phone, updatedUser.Code, updatedUser.CountryCode); err != nil {
					return c.JSON(http.StatusBadRequest, Response{Msg: err.Error()})
				}
			}
		}
	}

	var affected int64
	if createRequired {
		affected, _, err = object.AddUser(&updatedUser.User)
	} else {
		affected, err = object.UpdateUser(&updatedUser.User)
	}
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{Msg: err.Error()})
	}

	return c.JSON(http.StatusOK, UserResponse{
		Code: http.StatusOK,
		Msg:  fmt.Sprintf("affected:%d ", affected),
		Data: changeUserToken,
	})

}

func IsUserExistsByPhone(c echo.Context) (err error) {
	var user object.User
	if err = c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, Response{Msg: err.Error()})
	}

	_, exist, err := object.GetUserByPhone(user.Phone)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{Msg: err.Error()})
	}

	if !exist {
		return c.JSON(http.StatusBadRequest, Response{Msg: "用户不存在"})
	}

	return c.JSON(http.StatusOK, Response{
		Success: true,
		Code:    http.StatusOK,
		Data:    user.UserId,
	})
}
