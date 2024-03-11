package object

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

type (
	Configuration struct {
		System  SystemConfig      `mapstructure:"system"`
		WxLogin DatasourceConfigs `mapstructure:"wxlogin"`
		DyLogin DatasourceConfigs `mapstructure:"dylogin"`

		Mysql         *DatasourceConfig `mapstructure:"mysql"`
		Postgres      *DatasourceConfig `mapstructure:"postgres"`
		DbTablePrefix string            `mapstructure:"db_table_prefix"`

		DbDriver  string `mapstructure:"-"`
		DbConnStr string `mapstructure:"-"`
	}
	SystemConfig struct {
		Port int `mapstructure:"port"`
	}
	DatasourceConfigs map[string]*LoginConfiguration
	DatasourceConfig  struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Dbname   string `mapstructure:"dbname"`
	}
	LoginConfiguration struct {
		Appid  string `mapstructure:"appid"`
		Secret string `mapstructure:"secret"`
	}
)

func (c *Configuration) fromEnv() (err error) {
	if v := viper.GetInt("system_port"); v != 0 {
		c.System.Port = v
	}
	if c.System.Port == 0 {
		c.System.Port = 9825
	}
	if v := viper.GetString("db_table_prefix"); v != "" {
		c.DbTablePrefix = v
	}
	if Conf.DbTablePrefix != "" {
		// remove last _ before append _
		Conf.DbTablePrefix = strings.TrimSuffix(Conf.DbTablePrefix, "_") + "_"
	}
	if v := viper.GetString("db_url"); v != "" {
		spits := strings.Split(v, "://")
		switch c.DbDriver = spits[0]; c.DbDriver {
		case "postgres":
			c.DbConnStr = v
		case "mysql":
			c.DbConnStr = spits[1]
		default:
			err = errors.New("YUDAI_DB_URL environment variable require postgres or mysql url")
			return
		}
	}
	if c.DbDriver == "" && c.Mysql != nil {
		c.DbDriver = "mysql"
		c.DbConnStr = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8",
			c.Mysql.User, c.Mysql.Password, c.Mysql.Host, c.Mysql.Port, c.Mysql.Dbname)
	}
	if c.DbDriver == "" && c.Postgres != nil {
		c.DbDriver = "postgres"
		c.DbConnStr = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
			c.Postgres.User, c.Postgres.Password, c.Postgres.Host, c.Postgres.Port, c.Postgres.Dbname)
	}
	if c.DbDriver == "" {
		err = errors.New("please supply db config with env or config.yaml")
		return
	}

	loginFlags := []string{"mini", "pc", "h5", "app"}
	if c.WxLogin == nil {
		c.WxLogin = make(map[string]*LoginConfiguration)
	}
	if c.DyLogin == nil {
		c.DyLogin = make(map[string]*LoginConfiguration)
	}
	c.WxLogin.fromEnv("wx", false, loginFlags...)
	c.DyLogin.fromEnv("dy", true, loginFlags...)
	return
}

func (c DatasourceConfigs) fromEnv(prefix string, prependPrefix bool, flag ...string) {
	for _, item := range flag {
		conf, ok := makeLoginConfiguration(prefix, item)
		if !ok {
			continue
		}

		if prependPrefix {
			item = fmt.Sprintf("%s_%s", prefix, item)
		}
		c[item] = conf
	}
}

func makeLoginConfiguration(prefix, flag string) (conf *LoginConfiguration, ok bool) {
	appid := viper.GetString(fmt.Sprintf("%s_%s_appid", prefix, flag))
	secret := viper.GetString(fmt.Sprintf("%s_%s_secret", prefix, flag))
	if appid != "" && secret != "" {
		conf, ok = &LoginConfiguration{
			Appid:  appid,
			Secret: secret,
		}, true
	}
	return
}

var Conf *Configuration

func Init() (err error) {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("YUDAI")
	viper.SetConfigName("config") // 指定配置文件名（不带后缀）
	viper.AddConfigPath("./conf") // 指定查找配置文件的路径（这里使用相对路径）
	// 读取配置信息
	if err = viper.ReadInConfig(); err != nil {
		fmt.Printf("viper.ReadInConfig failed, err: %v\n", err)
		return
	}

	// 把读取到的配置信息反序列化到 Conf 变量中
	if err = viper.Unmarshal(&Conf); err != nil {
		fmt.Printf("viper.Unmarshal failed, err:%v\n", err)
		return
	}

	err = Conf.fromEnv()
	return
}
