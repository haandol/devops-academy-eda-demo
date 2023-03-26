package config

import (
	"log"

	"github.com/haandol/devops-academy-eda-demo/pkg/util"
	"github.com/joho/godotenv"
)

type App struct {
	Name                    string `validate:"required"`
	Stage                   string `validate:"required"`
	Port                    int    `validate:"required,number"`
	RPS                     int    `validate:"required,number"`
	TimeoutSec              int    `validate:"required,number,gte=0"`
	GracefulShutdownTimeout int    `validate:"required,number,gte=0"`
	DisableHTTP             bool   `default:"false"`
}

type Kafka struct {
	Seeds            []string `validate:"required"`
	MessageExpirySec int      `validate:"required,number"`
	BatchSize        int      `validate:"required,number"`
}

type Database struct {
	TableName string
}

type Rest struct {
	HotelHost string
}

type Config struct {
	App      App
	Kafka    Kafka
	Database Database
	Rest     Rest
}

// Load config.Config from environment variables for each stage
func Load() Config {
	stage := getEnv("APP_STAGE").String()
	log.Printf("Loading %s config\n", stage)

	if err := godotenv.Load(); err != nil {
		log.Panic("Error loading .env file")
	}

	cfg := Config{
		App: App{
			Name:                    getEnv("APP_NAME").String(),
			Stage:                   getEnv("APP_STAGE").String(),
			Port:                    getEnv("APP_PORT").Int(),
			RPS:                     getEnv("APP_RPS").Int(),
			TimeoutSec:              getEnv("APP_TIMEOUT_SEC").Int(),
			GracefulShutdownTimeout: getEnv("APP_GRACEFUL_SHUTDOWN_TIMEOUT").Int(),
			DisableHTTP:             getEnv("APP_DISABLE_HTTP").Bool(),
		},
		Kafka: Kafka{
			Seeds:            getEnv("KAFKA_SEEDS").Split(","),
			MessageExpirySec: getEnv("KAFKA_MESSAGE_EXPIRY_SEC").Int(),
			BatchSize:        getEnv("KAFKA_BATCH_SIZE").Int(),
		},
		Database: Database{
			TableName: getEnv("DB_TABLE_NAME").String(),
		},
		Rest: Rest{
			HotelHost: getEnv("REST_HOTEL_HOST").String(),
		},
	}

	if err := util.ValidateStruct(cfg); err != nil {
		log.Panicf("Error validating config: %s", err)
	}

	return cfg
}
