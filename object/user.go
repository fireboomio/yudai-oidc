package object

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"regexp"
	"time"
	"yudai/util"
)

type User struct {
	UserId    string    `xorm:"varchar(255)" json:"userId"`
	Avatar    string    `xorm:"varchar(255)" json:"avatar"`
	CreatedAt time.Time `xorm:"varchar(100)" json:"created_at"`
	UpdatedAt time.Time `xorm:"varchar(100)" json:"updated_at"`

	Name         string `xorm:"varchar(100) index" json:"name"`
	Password     string `xorm:"varchar(100)" json:"password"`
	PasswordSalt string `xorm:"varchar(100)" json:"passwordSalt"`
	PasswordType string `xorm:"varchar(100)" json:"passwordType"`

	Phone       string `xorm:"varchar(20) index" json:"phone,omitempty"`
	CountryCode string `xorm:"varchar(6)" json:"countryCode,omitempty"`
}

func AddUser(user *User) (int64, error) {
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
	affected, err := adapter.Engine.Insert(user)
	if err != nil {
		return 0, err
	}

	return affected, nil
}

func UpdateUser(user *User) (int64, error) {
	if user.UserId == "" {
		return 0, errors.New("用户ID为空")
	}

	if user.UpdatedAt.IsZero() {
		user.UpdatedAt = time.Now()
	}
	fillUserInfo(user)
	return adapter.Engine.Where("user_id=?", user.UserId).Update(user)
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

func GetUserByPhone(phone string) (user *User, existed bool, err error) {
	if phone == "" {
		err = errors.New("手机号码为空")
		return
	}

	user = &User{Phone: phone}
	existed, err = adapter.Engine.Get(&user)
	return
}

func GetUserByUserId(userId string) (user *User, existed bool, err error) {
	if userId == "" {
		err = errors.New("用户ID为空")
		return
	}

	user = &User{}
	existed, err = adapter.Engine.Where("user_id=?", user).Get(user)
	return
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

	user = &User{}
	existed, err = adapter.Engine.Where("name=?", name).Get(user)
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
