package worker

import (
	"context"
	"log"

	"github.com/jailtonjunior94/order/pkg/bundle"

	"github.com/robfig/cron/v3"
)

type worker struct {
}

func NewWorkers() *worker {
	return &worker{}
}

func (w *worker) Run() {
	ctx := context.Background()
	ioc := bundle.NewContainer(ctx)

	/* Observability */
	tracerProvider := ioc.Observability.TracerProvider()
	defer func() {
		if err := tracerProvider.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	meterProvider := ioc.Observability.MeterProvider()
	defer func() {
		if err := meterProvider.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	/* Close DBConnection */
	defer func() {
		if err := ioc.DB.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	c := cron.New()

	_, err := c.AddFunc(ioc.Config.WorkerConfig.CronExpression, func() {
		// Do something
		log.Println("Running...")
	})
	if err != nil {
		log.Fatal(err)
	}

	c.Run()
}
