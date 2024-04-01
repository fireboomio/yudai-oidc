package api

import (
	"fmt"
	"net/http"
	"yudai/object"

	"github.com/labstack/echo/v4"
)

// UpdateSmsProvider
//
//	@Title			UpdateSmsProvider
//	@Tag			Provider API
//	@Description	update provider
//	@Param			clientId		body		string			true	"clientId"
//	@Param			clientSecret	body		string			true	"clientSecret"
//	@Param			signName		body		string			true	"签名"
//	@Param			templateCode	body		string			true	"模板代码"
//	@Success		200				{object}	object.Response	成功
//	@router			/update-provider [post]
func UpdateSmsProvider(c echo.Context) (err error) {
	var provider object.Provider
	if err = c.Bind(&provider); err != nil {
		return c.JSON(http.StatusBadRequest, Response{Msg: err.Error()})
	}

	affected, err := object.UpdateSmsProvider(&provider)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{Msg: err.Error()})
	}

	return c.JSON(http.StatusOK, Response{
		Success: true,
		Code:    http.StatusOK,
		Msg:     fmt.Sprintf("affected:%d ", affected),
	})
}
