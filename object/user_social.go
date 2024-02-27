package object

import (
	"errors"
	"time"
)

type UserSocial struct {
	Id               int       `xorm:"id pk autoincr" json:"id"`
	UserId           string    `xorm:"user_id varchar(36)" json:"userId,omitempty"`
	Provider         string    `xorm:"provider varchar(64) not null" json:"provider"`
	ProviderUserId   string    `xorm:"provider_user_id varchar(64) not null" json:"providerUserId"`
	ProviderPlatform string    `xorm:"provider_platform varchar(64)" json:"providerPlatform,omitempty"`
	ProviderUnionid  string    `xorm:"provider_unionid varchar(64)" json:"providerUnionid,omitempty"`
	CreatedAt        time.Time `xorm:"created_at datetime index" json:"createdAt"`
}

func AddUserUserSocial(social *UserSocial) (int64, error) {
	_, existed, err := GetUserSocialByProviderUserId(social.ProviderUserId)
	if err != nil || existed {
		return 0, err
	}

	social.CreatedAt = time.Now()
	return engine.Insert(social)
}

func UpdateUserSocial(userId, providerUserId string) (int64, error) {
	if providerUserId == "" {
		return 0, errors.New("Social用户ID为空")
	}

	return engine.Where("provider_user_id=?", providerUserId).Update(&UserSocial{UserId: userId})
}

func GetUserSocialByProviderUserId(providerUserId string) (*UserSocial, bool, error) {
	if len(providerUserId) == 0 {
		return nil, false, errors.New("providerUserId is empty")
	}

	userSocial := UserSocial{ProviderUserId: providerUserId}
	existed, err := engine.Get(&userSocial)
	return &userSocial, existed, err
}

func GetUserSocialsByUserId(userId string, provider ...string) (data []*UserSocial, err error) {
	if len(userId) == 0 {
		err = errors.New("userId is empty")
		return
	}

	session := engine.Where("user_id=?", userId)
	if len(provider) > 0 {
		session.In("provider", provider)
	}
	err = session.Find(&data)
	return
}
