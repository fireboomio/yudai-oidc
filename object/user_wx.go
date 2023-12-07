package object

import (
	"time"
)

type UserWx struct {
	CreatedAt time.Time `xorm:"varchar(100) index" json:"createdAt"`
	Unionid   string    `xorm:"varchar(255) not null" json:"unionid"`
	Openid    string    `xorm:"varchar(100) not null" json:"openid"`
	Platform  string    `xorm:"varchar(100) not null" json:"platform"`
}

func AddUserWx(userWx *UserWx) (int64, error) {
	_, existed, err := GetUserWxByOpenid(userWx.Openid, userWx.Platform)
	if err != nil || existed {
		return 0, err
	}

	userWx.CreatedAt = time.Now()
	affected, err := adapter.Engine.Insert(userWx)
	if err != nil {
		return 0, err
	}

	return affected, nil
}

func GetUserWxByOpenid(openid, platform string) (*UserWx, bool, error) {
	if len(openid) == 0 || len(platform) == 0 {
		return nil, false, nil
	}

	userWx := UserWx{Openid: openid, Platform: platform}
	existed, err := adapter.Engine.Get(&userWx)
	if err != nil {
		return nil, false, err
	}

	if existed {
		return &userWx, true, nil
	} else {
		return nil, false, nil
	}
}

func GetUserWxsByUnionid(unionid string) (data []*UserWx, err error) {
	if len(unionid) == 0 {
		return
	}

	err = adapter.Engine.Where("unionid=?", unionid).Find(&data)
	return
}
