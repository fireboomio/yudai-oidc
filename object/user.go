package object

import (
	"errors"
	"fmt"
	"regexp"
	"time"
	"yudai/util"

	"github.com/google/uuid"
)

type User struct {
	Id        int       `xorm:"id pk autoincr" json:"id"`
	UserId    string    `xorm:"user_id varchar(36) notnull" json:"userId"`
	Avatar    string    `xorm:"avatar varchar(255)" json:"avatar"`
	CreatedAt time.Time `xorm:"created_at datetime notnull" json:"createdAt"`
	UpdatedAt time.Time `xorm:"updated_at datetime" json:"updatedAt"`

	Name         string `xorm:"name varchar(64) index" json:"name"`
	Password     string `xorm:"password varchar(100)" json:"password"`
	PasswordSalt string `xorm:"password_salt varchar(100)" json:"passwordSalt"`
	PasswordType string `xorm:"password_type varchar(100)" json:"passwordType"`

	Phone       string `xorm:"phone varchar(20) index" json:"phone,omitempty"`
	CountryCode string `xorm:"country_code varchar(6)" json:"countryCode,omitempty"`

	SocialUserId string `xorm:"-" json:"socialUserId,omitempty"`
}

type Userinfo struct {
	UserId       string        `json:"userId"`
	Name         string        `json:"name"`
	Phone        string        `json:"phone,omitempty"`
	Avatar       string        `json:"avatar,omitempty"`
	SocialUserId string        `json:"socialUserId,omitempty"`
	Socials      []*UserSocial `json:"socials,omitempty"`
}

func (u *User) Transform() *Userinfo {
	return &Userinfo{
		UserId:       u.UserId,
		Name:         u.Name,
		Phone:        u.Phone,
		Avatar:       u.Avatar,
		SocialUserId: u.SocialUserId,
	}
}

func AddUser(user *User) (int64, string, error) {
	var err error
	if user.UserId == "" {
		user.UserId = uuid.NewString()
	}
	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now()
	}
	if user.Name == "" {
		user.Name = user.UserId
	}

	fillUserInfo(user)
	affected, err := engine.Insert(user)
	if err != nil {
		return 0, "", err
	}

	return affected, user.UserId, nil
}

func UpdateUser(user *User) (int64, error) {
	if user.UserId == "" {
		return 0, errors.New("用户ID为空")
	}

	if user.UpdatedAt.IsZero() {
		user.UpdatedAt = time.Now()
	}
	fillUserInfo(user)
	return engine.Where("user_id=?", user.UserId).Update(user)
}

func fillUserInfo(user *User) {
	if len(user.Password) > 0 {
		user.PasswordType = "md5"
		user.PasswordSalt = util.RandomString(10)
		user.Password = util.GenMd5(user.PasswordSalt, user.Password)
	}
	if len(user.Phone) > 0 && len(user.CountryCode) == 0 {
		user.CountryCode = "CN"
	}
}

func GetUserByPhone(phone string) (*User, bool, error) {
	if phone == "" {
		return nil, false, errors.New("手机号码为空")
	}

	user := User{Phone: phone}
	existed, err := engine.Get(&user)
	return &user, existed, err
}

func GetUserByUserId(userId string) (*User, bool, error) {
	if userId == "" {
		return nil, false, errors.New("用户ID为空")
	}

	user := User{UserId: userId}
	existed, err := engine.Get(&user)
	return &user, existed, err
}

func GetUserByName(name string) (user *User, existed bool, err error) {
	if name == "" {
		err = errors.New("用户名为空")
		return
	}

	if len(name) > 39 {
		err = errors.New("用户名太长（最多39个字符）")
		return
	}

	exclude, _ := regexp.Compile("^[\u0021-\u007E]+$")
	if !exclude.MatchString(name) {
		err = errors.New("用户名非法字符")
		return
	}

	queryUser := User{Name: name}
	existed, err = engine.Get(&queryUser)
	user = &queryUser
	return
}

func CheckUserPassword(name, password string) (user *User, err error) {
	user, existed, err := GetUserByName(name)
	if err != nil {
		return
	}

	if !existed {
		err = fmt.Errorf("该用户不存在")
		return
	}

	if user.Password == "" {
		err = fmt.Errorf("该用户未设置密码")
		return
	}

	if user.Password != util.GenMd5(user.PasswordSalt, password) {
		err = fmt.Errorf("密码错误")
		return
	}

	return
}
