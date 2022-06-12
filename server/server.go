package server

import (
	"log"
	"piepay/config"
	"piepay/routes"
	"piepay/services/es"
	"piepay/utils"
	"time"

	"github.com/getsentry/sentry-go"
)

func Init() {

	config := config.Get()

	err := sentry.Init(sentry.ClientOptions{
		Dsn:              config.SentryDSN,
		Debug:            true,
		Environment:      config.AppEnv,
		TracesSampleRate: float64(config.SentrySamplingRate),
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	defer sentry.Flush(2 * time.Second)
	es.Init()
	go utils.Uploader()
	r := routes.NewRouter()
	r.Run(":" + "4000")
}
