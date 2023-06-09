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
	AuthHeader              string `validate:"required"`
	TimeoutSec              int    `validate:"required,number,gte=0"`
	GracefulShutdownTimeout int    `validate:"required,number,gte=0"`
	DisableHTTP             bool   `default:"false"`
}

type AWS struct {
	UseLocal bool   `default:"false"`
	Region   string `validate:"required"`
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
	AWS      AWS
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
			AuthHeader:              getEnv("APP_AUTH_HEADER").String(),
			TimeoutSec:              getEnv("APP_TIMEOUT_SEC").Int(),
			GracefulShutdownTimeout: getEnv("APP_GRACEFUL_SHUTDOWN_TIMEOUT").Int(),
			DisableHTTP:             getEnv("APP_DISABLE_HTTP").Bool(),
		},
		AWS: AWS{
			UseLocal: getEnv("AWS_USE_LOCAL").Bool(),
			Region:   getEnv("AWS_REGION").String(),
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
