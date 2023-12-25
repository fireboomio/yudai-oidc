package object

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"time"
	"yudai/util"
)

type User struct {
	Name      string    `xorm:"varchar(100) notnull pk" json:"name"`
	CreatedAt time.Time `xorm:"varchar(100) index" json:"created_at"`

	UserId       string `xorm:"varchar(255)" json:"userId"`
	Avatar       string `xorm:"varchar(255)" json:"avatar"`
	Password     string `xorm:"varchar(100)" json:"password"`
	PasswordSalt string `xorm:"varchar(100)" json:"passwordSalt,omitempty"`
	PasswordType string `xorm:"varchar(100)" json:"passwordType"`
	Phone        string `xorm:"varchar(20) index" json:"phone,omitempty"`
	CountryCode  string `xorm:"varchar(6)" json:"countryCode"`
	WxUnionid    string `xorm:"varchar(100)" json:"WxUnionId,omitempty"`
}

type Userinfo struct {
	Name  string `json:"preferred_username,omitempty"`
	Phone string `json:"phone,omitempty"`
}

func AddUser(user *User) (int64, error) {
	var err error
	if user.UserId == "" {
		user.UserId = uuid.NewString()
	}
	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now()
	}

	affected, err := adapter.Engine.Insert(user)
	if err != nil {
		return 0, err
	}

	return affected, nil
}

func UpdateUser(user *User) (int64, error) {
	var err error
	if user.UserId == "" {
		return 0, errors.New("user not exist")
	}

	affected, err := adapter.Engine.Where("user_id=?", user.UserId).Update(user)
	if err != nil {
		return 0, err
	}

	return affected, nil
}

func GetUserByPhone(phone string, unionidNotEmpty bool) (*User, bool, error) {
	if phone == "" {
		return nil, false, nil
	}

	var unionidExpr string
	if unionidNotEmpty {
		unionidExpr = "<>"
	} else {
		unionidExpr = "="
	}
	user := User{Phone: phone}
	existed, err := adapter.Engine.Where(fmt.Sprintf("wx_unionid %s ''", unionidExpr)).Get(&user)
	if err != nil {
		return nil, false, err
	}

	if existed {
		return &user, true, nil
	} else {
		return nil, false, nil
	}
}

func GetUserByWxUnionid(unionid string) (*User, error) {
	if unionid == "" {
		return nil, nil
	}

	user := User{WxUnionid: unionid}
	existed, err := adapter.Engine.Get(&user)
	if err != nil {
		return nil, err
	}

	if existed {
		return &user, nil
	} else {
		return nil, nil
	}
}

func CheckUserPassword(name, password string) (user *User, err error) {
	// 获取user的密码
	user = &User{}
	existed, err := adapter.Engine.Where("name=?", name).Get(user)
	if err != nil {
		return
	}

	if !existed {
		err = fmt.Errorf("该用户不存在")
		return
	}

	if user.Password != util.GenMd5(user.PasswordSalt, password) {
		err = fmt.Errorf("密码错误")
		return
	}

	return
}
