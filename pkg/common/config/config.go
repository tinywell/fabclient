package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

const (
	// ENVPREFIX config eng prefix
	ENVPREFIX = "CLIENT"
)

var (
	myViper *viper.Viper
)

// InitConfig config initialize
func InitConfig(configfile string) {
	if len(configfile) == 0 {
		panic("config file is none")
	}

	myViper = viper.New()
	myViper.SetEnvPrefix(ENVPREFIX)
	myViper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	myViper.SetEnvKeyReplacer(replacer)
	myViper.SetConfigFile(configfile)
	err := myViper.ReadInConfig()
	if err == nil {
		fmt.Printf("Using config file: %s\n", myViper.ConfigFileUsed())
	} else {
		fmt.Println("Read config file error:", err)
	}
}

func GetSDKConfig() string {
	return myViper.GetString("client.sdkconfig")
}
