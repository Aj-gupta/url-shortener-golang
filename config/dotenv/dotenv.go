package dotenv

import (
	"log"

	"github.com/spf13/viper"
)

type config struct {
	GoEnv           string `mapstructure:"GO_ENV"`
	GoPORT          string `mapstructure:"GO_PORT"`
	PgHost          string `mapstructure:"PG_HOST"`
	PgDB            string `mapstructure:"PG_DB"`
	PgPort          string `mapstructure:"PG_PORT"`
	PgPassword      string `mapstructure:"PG_PASS"`
	PgUser          string `mapstructure:"PG_USER"`
	ApiKey          string `mapstructure:"URL_SHORTNER_API_KEY"`
	ShortURLBaseURI string `mapstructure:"SHORT_URL_BASE_URI"`
	RequestTimeout  int64  `mapstructure:"REQUEST_TIMEOUT"`
}

var Global config

func LoadConfig(path string) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatalf("⛔ %s file not found.", path)
		} else {
			log.Fatalf("⛔ Loading Configuration failed: %s", err.Error())
		}
	}
	if err := viper.Unmarshal(&Global); err != nil {
		log.Fatalf("⛔ unmarshal configuration failed: %s", err.Error())
	}
}
