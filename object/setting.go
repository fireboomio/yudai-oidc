package object

import (
	"github.com/spf13/viper"
)

type WxLoginConfiguration struct {
	AppID     string
	AppSecret string
}

func Init() (err error) {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("YUDAI")
	viper.SetDefault("port", 9825)
	return
}
