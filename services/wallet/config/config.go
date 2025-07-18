package config

import (
	"sync"

	"github.com/Creative-genius001/go-logger"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Env             string `envconfig:"ENV" default:"development"`
	Port            string `envconfig:"PORT" default:"8080"`
	DBUrl           string `envconfig:"DB_URL" required:"true"`
	PaystackBaseUrl string `envconfig:"PAYSTACK_BASE_URL"`
	PaystackTestKey string `envconfig:"PAYSTACK_TEST_KEY"`
}

var (
	Cfg  *Config
	once sync.Once
)

func Init() {
	once.Do(func() {
		_ = godotenv.Load() // Loads .env file

		Cfg = &Config{}
		if err := envconfig.Process("", Cfg); err != nil {
			logger.Fatal("failed to load environment variables: ", err)
		}
	})
}
