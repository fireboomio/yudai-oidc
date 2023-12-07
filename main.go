package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"log"
	"net/http"
	"simple-casdoor/api"
	"simple-casdoor/object"
	"strings"
)

func main() {
	// 初始化 xorm 映射
	if err := object.Init(); err != nil {
		fmt.Printf("init settings failed, err:%v\n", err)
		return
	}
	object.InitAdapter()

	r := NewRouter()

	r.Logger.Fatal(r.Start(":9825"))
}

func NewRouter() *echo.Echo {
	//logFile, err := os.OpenFile("casdoor.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer logFile.Close()

	e := echo.New()

	// Debug mode
	e.Debug = true

	//e.Logger.SetOutput(logFile)

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

	e.POST("/api/register", api.Register)

	e.POST("/api/isUserExist", api.IsUserExists)

	e.POST("/api/refresh-token", api.RefreshToken)

	e.POST("/api/add-user", api.AddUser)

	e.GET("/.well-known/openid-configuration", api.GetOidcDiscovery)

	e.GET("/.well-known/jwks.json", api.GetJwks)

	e.POST("/api/send-verification-code", api.SendVerificationCode)

	// 以下的接口需要 token 认证
	r := e.Group("/api")
	r.Use(Authorize)

	r.GET("/userinfo", api.GetUserInfo)

	r.POST("/update-provider", api.UpdateProvider)

	r.POST("/update-user", api.UpdateUser)

	return e
}

func Authorize(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		header := c.Request().Header.Get("Authorization")
		if header == "" {
			return c.JSON(http.StatusUnauthorized, "please login first")
		}

		tokenString := formatToken(header)
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

func formatToken(header string) string {
	tokens := strings.Split(header, " ")
	if len(tokens) != 2 {
		return ""
	}
	if tokens[0] != "Bearer" {
		return ""
	}
	return tokens[1]
}
