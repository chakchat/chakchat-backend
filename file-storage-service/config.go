package main

import (
	"fmt"
	"log"
	"os"
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

	GRPCService struct {
		Port int `mapstructure:"port"`
	} `mapstructure:"grpc_service"`

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
	viper.AutomaticEnv()

	viper.MustBindEnv("s3.bucket", "FILE_STORAGE_S3_BUCKET")
	viper.MustBindEnv("s3.url_prefix", "FILE_STORAGE_S3_URL_PREFIX")

	viper.MustBindEnv("aws.access_key_id", "FILE_STORAGE_AWS_ACCESS_KEY_ID")
	viper.MustBindEnv("aws.secret_access_key", "FILE_STORAGE_AWS_SECRET_ACCESS_KEY")
	viper.MustBindEnv("aws.region", "FILE_STORAGE_AWS_REGION")
	viper.MustBindEnv("aws.endpoint_url", "FILE_STORAGE_AWS_ENDPOINT_URL")
	viper.MustBindEnv("db.dsn", "FILE_STORAGE_DB_DSN")

	for k, v := range os.Environ() {
		fmt.Println(k, v)
	}

	viper.SetConfigFile(file)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("viper reading config failed: %s", err)
	}

	conf := new(Config)
	if err := viper.UnmarshalExact(&conf); err != nil {
		log.Fatalf("viper config unmarshalling failed: %s", err)
	}

	fmt.Printf("%#v\n", conf)

	return conf
}
