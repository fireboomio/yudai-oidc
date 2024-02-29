package api

import (
	"net/http"
	"yudai/object"

	"github.com/labstack/echo/v4"
)

func GetJwks(c echo.Context) (err error) {
	jwks, err := object.GetJsonWebKeySet()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, jwks)
}

// GetUserInfo ...
//
//	@Title			GetUserInfo
//	@Tag			UserInfo API
//	@Description	Get User Info
//	@Success		200		{object}	object.Userinfo	成功
//	@router			/userinfo [get]
func GetUserInfo(c echo.Context) (err error) {
	userinfo := c.Get("user").(*object.Userinfo)
	userinfo.Socials, _ = object.GetUserSocialsByUserId(userinfo.UserId)
	return c.JSON(http.StatusOK, userinfo)
}

func GetOidcDiscovery(c echo.Context) (err error) {
	host := c.Request().Host
	return c.JSON(http.StatusOK, object.GetOidcDiscovery(host))
}
