package main

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Jwt JWTConfig `mapstructure:"jwt"`

	Redis struct {
		Addr     string `mapstructure:"addr"`
		Password string `mapstructure:"password"`
		DB       int    `mapstructure:"db"`
	} `mapstructure:"redis"`

	Idempotency struct {
		DataExp time.Duration `mapstructure:"data_exp"`
	} `mapstructure:"idempotency"`

	Upload struct {
		FileSizeLimit int64 `mapstructure:"file_size_limit"`
	} `mapstructure:"upload"`

	MultipartUpload struct {
		MinFileSize int64 `mapstructure:"min_file_size"`
		MaxPartSize int64 `mapstructure:"max_part_size"`
	} `mapstructure:"multipart_upload"`

	S3 struct {
		Bucket    string `mapstructure:"bucket"`
		UrlPrefix string `mapstructure:"url_prefix"`
	} `mapstructure:"s3"`

	AWS struct {
		AccessKeyId     string `mapstructure:"access_key_id"`
		SecretAccessKey string `mapstructure:"secret_access_key"`
		Region          string `mapstructure:"region"`
		EndpointUrl     string `mapstructure:"endpoint_url"`
	} `mapstructure:"aws"`

	DB struct {
		DSN string `mapstructure:"dsn"`
	} `mapstructure:"db"`
}

type JWTConfig struct {
	SigningMethod string        `mapstructure:"signing_method"`
	Lifetime      time.Duration `mapstructure:"lifetime"`
	Issuer        string        `mapstructure:"issuer"`
	Audience      []string      `mapstructure:"audience"`
	KeyFilePath   string        `mapstructure:"key_file_path"`
}

func loadConfig(file string) *Config {
	viper.AutomaticEnv()

	viper.BindEnv("aws.access_key_id", "FILE_STORAGE_AWS_ACCESS_KEY_ID")
	viper.BindEnv("aws.secret_access_key", "FILE_STORAGE_AWS_SECRET_ACCESS_KEY")
	viper.BindEnv("aws.region", "FILE_STORAGE_AWS_REGION")
	viper.BindEnv("aws.endpoint_url", "FILE_STORAGE_AWS_ENDPOINT_URL")
	viper.BindEnv("db.dsn", "FILE_STORAGE_DB_DSN")

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
