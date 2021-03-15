package config

import (
	"fmt"
	"github.com/spf13/viper"
)

func LoadConfig() {
	viper.SetConfigFile("./config/rchain-dev.toml")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("viper read file error ", err)
		return
	}
	// fmt.Fprintln(os.Stderr, "using config file:", viper.ConfigFileUsed())
}
