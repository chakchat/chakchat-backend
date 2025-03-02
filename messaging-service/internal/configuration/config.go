package configuration

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	JWT struct {
		SigningMethod string   `mapstructure:"signing_method"`
		Issuer        string   `mapstructure:"issuer"`
		Audience      []string `mapstructure:"audience"`
		KeyFilePath   string   `mapstructure:"key_file_path"`
	} `mapstructure:"jwt"`
}

func LoadConfig(file string) *Config {
	viper.SetConfigFile(file)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("viper reading config failed: %s", err)
	}

	conf := new(Config)
	if err := viper.UnmarshalExact(&conf); err != nil {
		log.Fatalf("viper config unmarshalling failed: %s", err)
	}

	return conf
}
