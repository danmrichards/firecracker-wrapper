package main

import (
	"flag"
	"time"

	"github.com/danmrichards/firecracker-wrapper/internal/api"
	"github.com/danmrichards/firecracker-wrapper/internal/http"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
)

var bind string

const shutdownTimeout = 5 * time.Second

func main() {
	flag.StringVar(&bind, "bind", "0.0.0.0:5000", "the ip:port to bind the API server to")
	flag.Parse()

	logger := logrus.New()

	r := mux.NewRouter()

	srv := http.NewServer(bind, shutdownTimeout, r)

	if err := api.Init(r, logger); err != nil {
		logger.Fatalf("could not setup API handler: %v", err)
	}

	ctx := signals.SetupSignalHandler()

	logger.WithField("bind", bind).Info("starting API server")
	if err := srv.Serve(ctx); err != nil {
		logger.Fatalf("could not start API server: %v", err)
	}
}
