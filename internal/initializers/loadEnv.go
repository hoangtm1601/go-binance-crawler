package initializers

import (
	"time"

	"github.com/hoangtm1601/go-binance-crawler/internal/models"
	"github.com/spf13/viper"
)

type Config struct {
	EnableAutoMigrate string `mapstructure:"ENABLE_AUTO_MIGRATE"`
	DBHost            string `mapstructure:"POSTGRES_HOST"`
	DBUserName        string `mapstructure:"POSTGRES_USER"`
	DBUserPassword    string `mapstructure:"POSTGRES_PASSWORD"`
	DBName            string `mapstructure:"POSTGRES_DB"`
	DBPort            string `mapstructure:"POSTGRES_PORT"`
	ServerPort        string `mapstructure:"PORT"`
	DB_TIMEZONE       string `mapstructure:"DB_TIMEZONE"`
	SSLMode           string `mapstructure:"SSL_MODE"`

	GormLogging bool `mapstructure:"GORM_LOGGING"`

	ClientOrigin string `mapstructure:"CLIENT_ORIGIN"`

	AccessTokenPrivateKey  string        `mapstructure:"ACCESS_TOKEN_PRIVATE_KEY"`
	AccessTokenPublicKey   string        `mapstructure:"ACCESS_TOKEN_PUBLIC_KEY"`
	RefreshTokenPrivateKey string        `mapstructure:"REFRESH_TOKEN_PRIVATE_KEY"`
	RefreshTokenPublicKey  string        `mapstructure:"REFRESH_TOKEN_PUBLIC_KEY"`
	AccessTokenExpiresIn   time.Duration `mapstructure:"ACCESS_TOKEN_EXPIRED_IN"`
	RefreshTokenExpiresIn  time.Duration `mapstructure:"REFRESH_TOKEN_EXPIRED_IN"`
	AccessTokenMaxAge      int           `mapstructure:"ACCESS_TOKEN_MAXAGE"`
	RefreshTokenMaxAge     int           `mapstructure:"REFRESH_TOKEN_MAXAGE"`

	RedisHost     string `mapstructure:"REDIS_HOST"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`

	BinanceAPIKey      string `mapstructure:"BINANCE_API_KEY"`
	BinanceAPISecret   string `mapstructure:"BINANCE_API_SECRET"`
	BinanceCandleLimit int    `mapstructure:"BINANCE_CANDLE_LIMIT"`

	CrawlerSymbol                       string                `mapstructure:"CRAWLER_SYMBOL"`
	CrawlerSleepTimeBetweenPastCrawling int                   `mapstructure:"CRAWLER_SLEEP_TIME_BETWEEN_PAST_CRAWLING"`
	CrawlerInitialCrawlingTime          int                   `mapstructure:"CRAWLER_INITIAL_CRAWLING_TIME"`
	CrawlerSleepTimeOnceFailed          int                   `mapstructure:"CRAWLER_SLEEP_TIME_ONCE_FAILED"`
	CrawlerDefaultInterval              models.CandleInterval `mapstructure:"CRAWLER_DEFAULT_INTERVAL"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName("app")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
