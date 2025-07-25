package config

import (
	"sync"

	"github.com/Creative-genius001/Stacklo/services/payment/utils/logger"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
)

type Config struct {
	Env   string `envconfig:"ENV" default:"development"`
	Port  string `envconfig:"PAYMENT_SERVICE_PORT" default:"8080"`
	DBUrl string `envconfig:"DB_URL" required:"true"`
	// PaystackBaseUrl  string `envconfig:"PAYSTACK_BASE_URL" required:"true"`
	JwtKey string `envconfig:"JWT_KEY" required:"true"`
	// PaystackTestKey  string `envconfig:"PAYSTACK_TEST_KEY" required:"true"`
	HMACKey          string `envconfig:"HMAC_KEY" required:"true"`
	BinanceAPIKey    string `envconfig:"BINANCE_API_KEY_TEST" required:"true"`
	BinanceSecretKey string `envconfig:"BINANCE_SECRET_KEY_TEST" required:"true"`
	BinanceBaseUrl   string `envconfig:"BINANCE_BASE_URL" required:"true"`
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
			logger.Logger.Fatal("failed to load environment variables: ", zap.Error(err))
		}
	})
}
