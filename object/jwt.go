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
	Platform          string    `xorm:"varchar(36)" json:"platform"`
	UserId            string    `xorm:"varchar(255)" json:"userId"`
	Token             string    `xorm:"varchar(255)" json:"token"`
	ExpireTime        time.Time `xorm:"varchar(100)" json:"expire_time"`
	RefreshToken      string    `xorm:"varchar(255)" json:"refresh_token"`
	RefreshExpireTime time.Time `xorm:"varchar(100)" json:"refresh_expire_time"`
	Banned            bool      `xorm:"bool" json:"banned"`
}

func GenerateToken(user *User, platform PlatformConfig) (res *TokenRes, err error) {
	// Create the Claims
	nowTime := time.Now()
	expireAt := nowTime.Add(24 * 30 * time.Hour)

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
	refreshExpireTime := nowTime.Add(7 * 24 * time.Hour)

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
		Platform:          platform.Platform,
		UserId:            user.UserId,
		Token:             tokenString,
		ExpireTime:        expireAt,
		RefreshToken:      refreshTokenString,
		RefreshExpireTime: refreshExpireTime,
		Banned:            false,
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
			Where("banned=0 and expire_time>?", nowTime.Format(time.DateTime)).
			In("user_id", userIds).
			In("platform", []string{"", platform.Platform}).
			NotIn("token", []string{tokenString}).
			SetExpr("banned", true).
			Update(&Token{}); err != nil {
			return
		}
	}

	return &TokenRes{
		AccessToken:  tokenString,
		RefreshToken: refreshTokenString,
		ExpireIn:     expireAt.Unix(),
	}, err
}

func ParseToken(token string, beanFetch func() *Token) (*Claims, error) {
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

	if claims.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("token expired")
	}

	return claims, validateToken(beanFetch())
}

func validateToken(tokenBean *Token) error {
	exist, err := adapter.Engine.Get(tokenBean)
	if err != nil {
		return err
	}
	if !exist {
		return errors.New("token not exist")
	}
	if tokenBean.Banned {
		return errors.New("token banned")
	}
	return nil
}
