//go:build wireinject

package app

import (
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/wire"
	"github.com/haandol/devops-academy-eda-demo/pkg/adapter/primary/consumer"
	"github.com/haandol/devops-academy-eda-demo/pkg/adapter/primary/router"
	"github.com/haandol/devops-academy-eda-demo/pkg/adapter/secondary/producer"
	"github.com/haandol/devops-academy-eda-demo/pkg/adapter/secondary/repository"
	"github.com/haandol/devops-academy-eda-demo/pkg/adapter/secondary/rest"
	"github.com/haandol/devops-academy-eda-demo/pkg/config"
	cloudconnector "github.com/haandol/devops-academy-eda-demo/pkg/connector/cloud"
	kafkaproducer "github.com/haandol/devops-academy-eda-demo/pkg/connector/producer"
	"github.com/haandol/devops-academy-eda-demo/pkg/port"
	"github.com/haandol/devops-academy-eda-demo/pkg/port/primaryport/routerport"
	"github.com/haandol/devops-academy-eda-demo/pkg/port/secondaryport/producerport"
	"github.com/haandol/devops-academy-eda-demo/pkg/port/secondaryport/repositoryport"
	"github.com/haandol/devops-academy-eda-demo/pkg/port/secondaryport/restport"
	"github.com/haandol/devops-academy-eda-demo/pkg/service"
)

// Common
func provideDbClient(cfg *config.Config) *dynamodb.Client {
	awsCfg, err := cloudconnector.GetAWSConfig(&cfg.Database)
	if err != nil {
		panic(err)
	}

	return dynamodb.NewFromConfig(awsCfg.Cfg)
}

func provideProducer(cfg *config.Config) *kafkaproducer.KafkaProducer {
	kafkaProducer, err := kafkaproducer.Connect(&cfg.Kafka)
	if err != nil {
		panic(err)
	}
	return kafkaProducer
}

// TripApp
func provideTripProducer(
	kafkaProducer *kafkaproducer.KafkaProducer,
) *producer.TripProducer {
	return producer.NewTripProducer(kafkaProducer)
}

func provideTripConsumer(
	cfg *config.Config,
	tripService *service.TripService,
) *consumer.TripConsumer {
	kafkaConsumer := consumer.NewKafkaConsumer(&cfg.Kafka, "trip", "trip-service")
	return consumer.NewTripConsumer(kafkaConsumer, tripService)
}

var provideTripRouters = wire.NewSet(
	router.NewGinRouter,
	wire.Bind(new(http.Handler), new(*router.GinRouter)),
	router.NewServerForce,
	wire.Bind(new(routerport.RouterGroup), new(*router.GinRouter)),
	router.NewTripRouter,
)

var provideRouters = wire.NewSet(
	router.NewGinRouter,
	wire.Bind(new(http.Handler), new(*router.GinRouter)),
	router.NewServer,
	wire.Bind(new(routerport.RouterGroup), new(*router.GinRouter)),
)

func InitTripApp(cfg *config.Config) port.App {
	wire.Build(
		provideDbClient,
		provideProducer,
		provideTripRouters,
		provideTripConsumer,
		service.NewTripService,
		provideTripProducer,
		wire.Bind(new(producerport.TripProducer), new(*producer.TripProducer)),
		repository.NewTripRepository,
		wire.Bind(new(repositoryport.TripRepository), new(*repository.TripRepository)),
		rest.NewTripRestAdapter,
		wire.Bind(new(restport.TripRestAdapter), new(*rest.TripRestAdapter)),
		NewTripApp,
		wire.Bind(new(port.App), new(*TripApp)),
	)
	return nil
}

// CarApp
func provideCarProducer(
	kafkaProducer *kafkaproducer.KafkaProducer,
) *producer.CarProducer {
	return producer.NewCarProducer(kafkaProducer)
}

func provideCarConsumer(
	cfg *config.Config,
	carService *service.CarService,
) *consumer.CarConsumer {
	kafkaConsumer := consumer.NewKafkaConsumer(&cfg.Kafka, "car", "car-service")
	return consumer.NewCarConsumer(kafkaConsumer, carService)
}

func InitCarApp(cfg *config.Config) port.App {
	wire.Build(
		provideDbClient,
		provideProducer,
		provideRouters,
		provideCarConsumer,
		service.NewCarService,
		provideCarProducer,
		wire.Bind(new(producerport.CarProducer), new(*producer.CarProducer)),
		repository.NewCarRepository,
		wire.Bind(new(repositoryport.CarRepository), new(*repository.CarRepository)),
		NewCarApp,
		wire.Bind(new(port.App), new(*CarApp)),
	)
	return nil
}

// HotelApp
func provideHotelProducer(
	kafkaProducer *kafkaproducer.KafkaProducer,
) *producer.HotelProducer {
	return producer.NewHotelProducer(kafkaProducer)
}

func provideHotelConsumer(
	cfg *config.Config,
	hotelService *service.HotelService,
) *consumer.HotelConsumer {
	kafkaConsumer := consumer.NewKafkaConsumer(&cfg.Kafka, "hotel", "hotel-service")
	return consumer.NewHotelConsumer(kafkaConsumer, hotelService)
}

func InitHotelApp(cfg *config.Config) port.App {
	wire.Build(
		provideDbClient,
		provideProducer,
		provideRouters,
		router.NewHotelRouter,
		provideHotelConsumer,
		service.NewHotelService,
		provideHotelProducer,
		wire.Bind(new(producerport.HotelProducer), new(*producer.HotelProducer)),
		repository.NewHotelRepository,
		wire.Bind(new(repositoryport.HotelRepository), new(*repository.HotelRepository)),
		NewHotelApp,
		wire.Bind(new(port.App), new(*HotelApp)),
	)
	return nil
}

// FlightApp
func provideFlightProducer(
	kafkaProducer *kafkaproducer.KafkaProducer,
) *producer.FlightProducer {
	return producer.NewFlightProducer(kafkaProducer)
}

func provideFlightConsumer(
	cfg *config.Config,
	flightService *service.FlightService,
) *consumer.FlightConsumer {
	kafkaConsumer := consumer.NewKafkaConsumer(&cfg.Kafka, "flight", "flight-service")
	return consumer.NewFlightConsumer(kafkaConsumer, flightService)
}

func InitFlightApp(cfg *config.Config) port.App {
	wire.Build(
		provideDbClient,
		provideProducer,
		provideRouters,
		provideFlightConsumer,
		service.NewFlightService,
		provideFlightProducer,
		wire.Bind(new(producerport.FlightProducer), new(*producer.FlightProducer)),
		repository.NewFlightRepository,
		wire.Bind(new(repositoryport.FlightRepository), new(*repository.FlightRepository)),
		NewFlightApp,
		wire.Bind(new(port.App), new(*FlightApp)),
	)
	return nil
}
