package object

import (
	"fmt"
	"simple-casdoor/util"
)

func CheckUserPassword(phone string, password string) (user *User, err error) {
	// 获取user的密码
	user = &User{}
	existed, err := adapter.Engine.Where("phone=?", phone).Get(user)
	if err != nil {
		return nil, err
	}

	if !existed {
		return nil, fmt.Errorf("该用户不存在")
	}

	password = util.GenMd5(user.PasswordSalt, password)

	if user.Password == password {
		return user, nil
	}

	return nil, fmt.Errorf("密码错误")
}
