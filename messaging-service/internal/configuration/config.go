package configuration

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	JWT struct {
		SigningMethod string   `mapstructure:"signing_method"`
		Issuer        string   `mapstructure:"issuer"`
		Audience      []string `mapstructure:"audience"`
		KeyFilePath   string   `mapstructure:"key_file_path"`
	} `mapstructure:"jwt"`

	DB struct {
		ConnString string `mapstructure:"conn_string"`
	} `mapstructure:"db"`

	Redis struct {
		Addr     string `mapstructure:"addr"`
		Password string `mapstructure:"password"`
		DB       int    `mapstructure:"db"`
	} `mapstructure:"redis"`

	Otlp struct {
		GrpcAddr string `mapstructure:"grpc_addr"`
	} `mapstructure:"otlp"`
}

func LoadConfig(file string) (*Config, error) {
	viper.AutomaticEnv()

	viper.MustBindEnv("db.conn_string", "DB_CONN_STRING")

	viper.SetConfigFile(file)

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("viper reading config failed: %s", err)
	}

	conf := new(Config)
	if err := viper.UnmarshalExact(&conf); err != nil {
		return nil, fmt.Errorf("viper config unmarshalling failed: %s", err)
	}

	return conf, nil
}
