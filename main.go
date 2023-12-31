package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"yudai/api"
	"yudai/object"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func main() {
	// 初始化 xorm 映射
	if err := object.Init(); err != nil {
		fmt.Printf("init settings failed, err:%v\n", err)
		return
	}

	object.InitAdapter()
	r := NewRouter()
	port := 9825
	if object.Conf.SystemConfig != nil && object.Conf.SystemConfig.Port > 0 {
		port = object.Conf.SystemConfig.Port
	}
	r.Logger.Fatal(r.Start(fmt.Sprintf(":%d", port)))
}

func NewRouter() *echo.Echo {
	e := echo.New()
	// Debug mode
	e.Debug = true

	// Middleware
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}, error=${error}\n",
	}))
	e.Use(middleware.Recover())

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	e.POST("/api/login", api.Login)

	e.POST("/api/register", api.AddUser)

	e.POST("/api/add-user", api.AddUser)

	e.POST("/api/isUserExist", api.IsUserExistsByPhone)

	e.POST("/api/refresh-token", api.RefreshToken)

	e.POST("/api/send-verification-code", api.SendVerificationCode)

	e.GET("/.well-known/openid-configuration", api.GetOidcDiscovery)

	e.GET("/.well-known/jwks.json", api.GetJwks)

	// 以下的接口需要 token 认证
	r := e.Group("/api")
	r.Use(Authorize)

	r.GET("/userinfo", api.GetUserInfo)

	r.POST("/update-provider", api.UpdateSmsProvider)

	r.POST("/update-user", api.UpdateUser)

	return e
}

func Authorize(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authorization := c.Request().Header.Get("Authorization")
		if authorization == "" {
			return c.JSON(http.StatusUnauthorized, "token not found")
		}

		tokenString, found := strings.CutPrefix(authorization, "Bearer ")
		if !found {
			return c.JSON(http.StatusUnauthorized, "token format error")
		}

		claims, err := object.ParseToken(tokenString, func() *object.Token {
			return &object.Token{Token: tokenString}
		})
		if err != nil || claims == nil {
			// 验证不通过，不再调用后续的函数处理
			return c.JSON(http.StatusUnauthorized, err.Error())
		}

		c.Set("user", claims.User)
		log.Printf("访问鉴权：{%s}", claims.User.Name)

		return next(c)
	}
}
