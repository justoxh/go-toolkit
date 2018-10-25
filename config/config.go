package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// LoadConfig load config
func LoadConfig(file string, output interface{}) error {
	if "" == file {
		return fmt.Errorf("blank file")
	}

	v := viper.New()
	v.SetConfigFile(file) // auto detect the file suffix

	err := v.ReadInConfig()
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	err = v.Unmarshal(&output)
	return err
}
