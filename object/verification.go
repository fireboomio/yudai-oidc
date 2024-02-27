package object

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
	"xorm.io/xorm/schemas"
)

const (
	VerificationSuccess = iota
	wrongCodeError
	noRecordError
	timeoutError
)

type VerificationRecord struct {
	Id        int       `xorm:"id pk autoincr" json:"id"`
	CreatedAt time.Time `xorm:"created_at datetime notnull" json:"createdAt"`
	Type      string    `xorm:"type varchar(10)" json:"type,omitempty"`
	UserId    string    `xorm:"user_id varchar(36) notnull" json:"userId"`
	Provider  string    `xorm:"provider varchar(100) notnull" json:"provider"`
	Receiver  string    `xorm:"receiver varchar(100) notnull" json:"receiver"`
	Code      string    `xorm:"code varchar(10) notnull" json:"code"`
	IsUsed    bool      `xorm:"is_used bool" json:"isUsed,omitempty"`
}

type VerifyResult struct {
	Code int
	Msg  string
}

func IsAllowSend(user *User, recordType string) error {
	var record VerificationRecord
	record.Type = recordType
	if user != nil {
		record.UserId = user.UserId
	}
	has, err := engine.Desc("created_at").Get(&record)
	if err != nil {
		return err
	}

	now := time.Now().Unix()
	if has && now-record.CreatedAt.Unix() < 60 {
		return errors.New("you can only send one code in 60s")
	}

	return nil
}

func CheckVerificationCode(dest, code string) *VerifyResult {
	record, err := getVerificationRecord(dest)
	if err != nil {
		panic(err)
	}

	if record == nil {
		return &VerifyResult{noRecordError, "无效的验证码!"}
	}

	var timeout int64 = 10
	if err != nil {
		panic(err)
	}

	now := time.Now().Unix()
	if now-record.CreatedAt.Unix() > timeout*60 {
		return &VerifyResult{timeoutError, fmt.Sprintf("您应该在%d分钟内验证您的验证码!", timeout)}
	}

	if record.Code != code {
		return &VerifyResult{wrongCodeError, "短信验证码错误!"}
	}

	return &VerifyResult{VerificationSuccess, ""}
}

func CheckSignInCode(dest, code string) string {
	result := CheckVerificationCode(dest, code)
	switch result.Code {
	case VerificationSuccess:
		return ""
	case wrongCodeError:
		return fmt.Sprintf("短信验证码错误!")
	default:
		return result.Msg
	}
}

func SendVerificationCodeToPhone(user *User, provider *Provider, remoteAddr string, dest string) error {
	if provider == nil {
		return errors.New("please set a SMS provider first")
	}

	if err := IsAllowSend(user, remoteAddr); err != nil {
		return err
	}

	code := getRandomCode(6)
	if err := SendSms(provider, code, dest); err != nil {
		return err
	}

	if err := AddToVerificationRecord(user, provider, dest, code); err != nil {
		return err
	}

	return nil
}

func AddToVerificationRecord(user *User, provider *Provider, dest, code string) error {
	var record VerificationRecord
	if user != nil {
		record.UserId = user.UserId
	}
	record.CreatedAt = time.Now()

	record.Provider = provider.Name
	record.Receiver = dest
	record.Code = code

	_, err := engine.Insert(record)
	if err != nil {
		return err
	}

	return nil
}

func DisableVerificationCode(dest string) (err error) {
	record, err := getVerificationRecord(dest)
	if record == nil || err != nil {
		return
	}

	record.IsUsed = true
	_, err = engine.ID(schemas.PK{record.Id}).AllCols().Update(record)
	return
}

func getVerificationRecord(dest string) (*VerificationRecord, error) {
	var record VerificationRecord
	record.Receiver = dest
	has, err := engine.Desc("time").Where("is_used = 0").Get(&record)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return &record, nil
}

var stdNums = []byte("0123456789")

func getRandomCode(length int) string {
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, stdNums[r.Intn(len(stdNums))])
	}
	return string(result)
}
