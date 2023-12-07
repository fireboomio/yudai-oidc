package object

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type Claims struct {
	User      *User  `json:"username"`
	TokenType string `json:"tokenType,omitempty"`
	jwt.RegisteredClaims
}

type TokenRes struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpireIn     int64  `json:"expireIn"`
}

type PlatformConfig struct {
	Platform  string `json:"platform"`
	Exclusive bool   `json:"exclusive"`
}

type Token struct {
	Platform   string    `xorm:"varchar(36)" json:"platform"`
	UserId     string    `xorm:"varchar(255)" json:"userId"`
	Token      string    `xorm:"varchar(255)" json:"token"`
	ExpireTime time.Time `xorm:"varchar(100)" json:"expire_time"`
	FlushTime  time.Time `xorm:"varchar(100)" json:"flush_time"`
	Banned     bool      `xorm:"bool" json:"banned"`
}

func GenerateToken(user *User, platform PlatformConfig) (res *TokenRes, err error) {
	// Create the Claims
	nowTime := time.Now()
	expireAt := nowTime.Add(120 * time.Hour)

	claims := Claims{
		User:      user,
		TokenType: "access-token",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.UserId,
			NotBefore: jwt.NewNumericDate(nowTime),
			IssuedAt:  jwt.NewNumericDate(nowTime),
			ExpiresAt: jwt.NewNumericDate(expireAt),
			Issuer:    "fireboom",
		},
	}

	var token *jwt.Token
	var refreshToken *jwt.Token
	refreshExpireTime := nowTime.Add(24 * time.Hour)

	token = jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	claims.TokenType = "refresh-token"
	claims.ExpiresAt = jwt.NewNumericDate(refreshExpireTime)
	refreshToken = jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// RSA private key
	// cert通常代表着公私钥对中的私钥，用于对JWT进行签名，验证Token时使用公钥进行解密和验证
	key, err := jwt.ParseRSAPrivateKeyFromPEM(cert.PrivateKey)
	if err != nil {
		return
	}

	token.Header["kid"] = "fireboom"
	tokenString, err := token.SignedString(key)
	if err != nil {
		return
	}

	refreshTokenString, err := refreshToken.SignedString(key)

	at := &Token{
		UserId:     user.UserId,
		Token:      tokenString,
		ExpireTime: expireAt,
		FlushTime:  refreshExpireTime,
		Banned:     false,
		Platform:   platform.Platform,
	}

	adminToken := &Token{Token: tokenString}
	exist, err := adapter.Engine.Get(adminToken)
	if err != nil {
		return
	}

	if !exist {
		if _, err = adapter.Engine.Insert(at); err != nil {
			return
		}
	}

	if platform.Exclusive {
		var samePhoneUsers []*User
		_ = adapter.Engine.Where("phone=?", user.Phone).Find(&samePhoneUsers)
		var userIds []string
		for _, v := range samePhoneUsers {
			userIds = append(userIds, v.UserId)
		}
		if _, err = adapter.Engine.
			Where("platform=? and expire_time>?", platform.Platform, nowTime.Format(time.DateTime)).
			In("user_id", userIds).
			NotIn("token", []string{tokenString}).
			Update(&Token{ExpireTime: nowTime}); err != nil {
			return
		}
	}

	return &TokenRes{
		AccessToken:  tokenString,
		RefreshToken: refreshTokenString,
		ExpireIn:     expireAt.Unix(),
	}, err
}

func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		certificate, err := jwt.ParseRSAPublicKeyFromPEM(cert.Certificate)

		if err != nil {
			return nil, err
		}

		return certificate, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := tokenClaims.Claims.(*Claims)
	if !ok {
		return nil, errors.New("expected point of Claims, but not found")
	}

	adminToken := &Token{Token: token}
	exist, err := adapter.Engine.Get(adminToken)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errors.New("token not exist")
	}
	if adminToken.Banned {
		return nil, errors.New("token banned")
	}
	if adminToken.ExpireTime.Before(time.Now()) {
		return nil, errors.New("token expired")
	}
	return claims, nil
}
