package engine

import (
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault("MARLA_ENV", "production")
	if err := loadMarlaConfig(); err != nil {
		panic(err)
	}
}

func loadMarlaConfig() error {
	viper.SetConfigFile("./.env")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}

func LoadConfig(path string) error {
	viper.AddConfigPath(path)
	viper.SetConfigName("marla.config")
	viper.SetConfigType("json")

	if err := viper.MergeInConfig(); err != nil {
		return err
	}
	return nil
}
