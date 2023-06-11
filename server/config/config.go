package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Jwt struct {
	SecretKey   string `mapstructure:"secretKey" yaml:"secretKey"`
	ExpiresTime int64  `mapstructure:"expiresTime" yaml:"expiresTime"`
	Issuer      string `mapstructure:"issuer" yaml:"issuer"`
}

type Snowflake struct {
	MachineId int64 `mapstructure:"machineId" yaml:"machineId"`
}

type Config struct {
	Jwt       *Jwt       `mapstructure:"jwt" yaml:"jwt"`
	Snowflake *Snowflake `mapstructure:"snowflake" yaml:"snowflake"`
}

var GlobalConfig = &Config{}

func InitConfig() {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./server/")
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	if err := v.Unmarshal(&GlobalConfig); err != nil {
		panic(fmt.Errorf("unmarshal config failed, err:%v", err))
	}
}
