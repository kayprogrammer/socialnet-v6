package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	ProjectName               string `mapstructure:"PROJECT_NAME"`
	Debug                     bool   `mapstructure:"DEBUG"`
	EmailOtpExpireSeconds     int64  `mapstructure:"EMAIL_OTP_EXPIRE_SECONDS"`
	AccessTokenExpireMinutes  int    `mapstructure:"ACCESS_TOKEN_EXPIRE_MINUTES"`
	RefreshTokenExpireMinutes int    `mapstructure:"REFRESH_TOKEN_EXPIRE_MINUTES"`
	Port                      string `mapstructure:"PORT"`
	SecretKey                 string `mapstructure:"SECRET_KEY"`
	FirstSuperuserEmail       string `mapstructure:"FIRST_SUPERUSER_EMAIL"`
	FirstSuperUserPassword    string `mapstructure:"FIRST_SUPERUSER_PASSWORD"`
	FirstClientEmail          string `mapstructure:"FIRST_CLIENT_EMAIL"`
	FirstClientPassword       string `mapstructure:"FIRST_CLIENT_PASSWORD"`
	PostgresUser              string `mapstructure:"POSTGRES_USER"`
	PostgresPassword          string `mapstructure:"POSTGRES_PASSWORD"`
	PostgresServer            string `mapstructure:"POSTGRES_SERVER"`
	PostgresPort              string `mapstructure:"POSTGRES_PORT"`
	PostgresDB                string `mapstructure:"POSTGRES_DB"`
	TestPostgresDB            string `mapstructure:"TEST_POSTGRES_DB"`
	MailSenderEmail           string `mapstructure:"MAIL_SENDER_EMAIL"`
	MailSenderPassword        string `mapstructure:"MAIL_SENDER_PASSWORD"`
	MailSenderHost            string `mapstructure:"MAIL_SENDER_HOST"`
	MailSenderPort            int    `mapstructure:"MAIL_SENDER_PORT"`
	CORSAllowedOrigins        string `mapstructure:"CORS_ALLOWED_ORIGINS"`
	CloudinaryCloudName       string `mapstructure:"CLOUDINARY_CLOUD_NAME"`
	CloudinaryAPIKey          string `mapstructure:"CLOUDINARY_API_KEY"`
	CloudinaryAPISecret       string `mapstructure:"CLOUDINARY_API_SECRET"`
	SocketSecret		      string `mapstructure:"SOCKET_SECRET"`
}

func GetConfig() (config Config) {
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	var err error
	if err = viper.ReadInConfig(); err != nil {
		panic(err)
	}
	viper.Unmarshal(&config)
	return
}
