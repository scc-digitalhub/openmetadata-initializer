package main

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func loadConfig() (cfg config, err error) {
	viper.AutomaticEnv()

	pflag.StringP("username", "u", "admin", "admin username used to call openmetadata")
	pflag.StringP("password", "p", "admin", "admin password used to call openmetadata")
	pflag.StringP("hostname", "h", "http://openmetadata", "openmetadata hostname")
	pflag.Int32P("port", "P", 8085, "openmetadata port")
	pflag.Int32P("admin_port", "a", 8086, "openmetadata admin port")
	pflag.BoolP("change_user_password", "c", true, "change username password")
	pflag.BoolP("generate_token", "g", true, "generate token")
	pflag.Int16P("readiness_max_retry", "m", 10, "openmetadata readiness probe max retry")

	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
	err = viper.Unmarshal(&cfg)
	return cfg, err
}
