package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"strings"
)

func NewConfig(p string) *viper.Viper {
	envConf := os.Getenv("APP_CONF")
	if envConf == "" {
		envConf = p
	}
	fmt.Println("load conf file:", envConf)
	return getConfig(envConf)
}

func getConfig(path string) *viper.Viper {
	conf := viper.New()
	conf.SetConfigFile(path)

	// Support environment variable overrides
	conf.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	conf.AutomaticEnv()

	err := conf.ReadInConfig()
	if err != nil {
		panic(err)
	}
	return conf
}
