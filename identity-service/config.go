package main

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	AccessToken   JWTConfig `mapstructure:"access_token"`
	RefreshToken  JWTConfig `mapstructure:"refresh_token"`
	InternalToken JWTConfig `mapstructure:"internal_token"`

	InvalidatedTokenStorage struct {
		Exp time.Duration `mapstructure:"exp"`
	} `mapstructure:"invalidated_token_storage"`

	UserService struct {
		GrpcAddr string `mapstructure:"grpc_addr"`
	} `mapstructure:"userservice"`

	GRPCService struct {
		Port int `mapstructure:"port"`
	} `mapstructure:"grpc_service"`

	Redis struct {
		Addr     string `mapstructure:"addr"`
		Password string `mapstructure:"password"`
		DB       int    `mapstructure:"db"`
	} `mapstructure:"redis"`

	SignInMeta struct {
		Lifetime time.Duration `mapstructure:"lifetime"`
	} `mapstructure:"signin_meta"`

	SignUpMeta struct {
		Lifetime time.Duration `mapstructure:"lifetime"`
	} `mapstructure:"signup_meta"`

	Idempotency struct {
		DataExp time.Duration `mapstructure:"data_exp"`
	} `mapstructure:"idempotency"`

	PhoneCode struct {
		SendFrequency time.Duration `mapstructure:"send_frequency"`
	} `mapstructure:"phone_code"`

	Sms struct {
		Type   string `mapstructure:"type"`
		Email  string `mapstructure:"email"`
		ApiKey string `mapstructure:"api_key"`

		Stub struct {
			Addr string `mapstructure:"addr"`
		} `mapstructure:"stub"`
	} `mapstructure:"sms"`

	Otlp struct {
		GrpcAddr string `mapstructure:"grpc_addr"`
	} `mapstructure:"otlp"`
}

type JWTConfig struct {
	SigningMethod string        `mapstructure:"signing_method"`
	Lifetime      time.Duration `mapstructure:"lifetime"`
	Issuer        string        `mapstructure:"issuer"`
	Audience      []string      `mapstructure:"audience"`
	KeyFilePath   string        `mapstructure:"key_file_path"`
}

func loadConfig(file string) *Config {
	// viper.SetConfigFile("/app/config.yml")

	viper.AutomaticEnv()
	viper.BindEnv("sms.email", "SMSAERO_EMAIL")
	viper.BindEnv("sms.api_key", "SMSAERO_APIKEY")

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
