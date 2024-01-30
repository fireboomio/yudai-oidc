package object

import (
	"yudai/object"

	"github.com/spf13/viper"
)

type WxLoginConfiguration struct {
	AppID     string
	AppSecret string
}

func Init() (err error) {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("YUDAI")
	viper.SetDefault(object.ENV_PORT, 9825)
	return
}
