package object

import (
	"fmt"
	"yudai/util"
)

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
