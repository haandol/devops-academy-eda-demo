package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"sync"

	"golang.org/x/sync/errgroup"

	"github.com/haandol/devops-academy-eda-demo/pkg/adapter/primary/consumer"
	"github.com/haandol/devops-academy-eda-demo/pkg/port/primaryport/consumerport"
	"github.com/haandol/devops-academy-eda-demo/pkg/util"
)

type CarApp struct {
	server   *http.Server
	consumer consumerport.Consumer
}

func NewCarApp(
	server *http.Server,
	carConsumer *consumer.CarConsumer,
) *CarApp {
	return &CarApp{
		server:   server,
		consumer: carConsumer,
	}
}

func (a *CarApp) Init() {
	logger := util.GetLogger().With(
		"module", "CarApp",
		"func", "Init",
	)
	logger.Info("Initializing App...")

	a.consumer.Init()

	logger.Info("App Initialized")
}

func (a *CarApp) Start(ctx context.Context) error {
	logger := util.GetLogger().With(
		"module", "CarApp",
		"func", "Start",
	)
	logger.Info("Starting App...")

	g := new(errgroup.Group)
	if a.server != nil {
		g.Go(func() error {
			logger.Infow("Started and serving HTTP", "addr", a.server.Addr, "pid", os.Getpid())
			if err := a.server.ListenAndServe(); err != nil {
				if errors.Is(err, http.ErrServerClosed) {
					logger.Info("server closed.")
					return err
				} else {
					logger.Errorw("ListenAndServe fail", "error", err)
					return err
				}
			}
			return nil
		})
	}
	g.Go(func() error {
		return a.consumer.Consume(ctx)
	})

	logger.Info("App Started")

	return g.Wait()
}

func (a *CarApp) Cleanup(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	logger := util.GetLogger().With(
		"module", "CarApp",
		"func", "Cleanup",
	)
	logger.Info("Cleaning App...")

	if a.server != nil {
		logger.Info("Shutting down server...")
		if err := a.server.Shutdown(ctx); err != nil {
			logger.Error("Error on server shutdown:", err)
		}
		logger.Info("Server shutdown.")
	}

	if err := a.consumer.Close(ctx); err != nil {
		logger.Errorw("failed to close consumer", "err", err)
	} else {
		logger.Info("Consumer closed.")
	}

	logger.Info("App Cleaned Up")
}
