package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type handler struct {
	logger *logrus.Logger
}

// Init initialises a new API handler using the given router.
func Init(r *mux.Router, logger *logrus.Logger) error {
	sr := r.PathPrefix("/v1").Subrouter()

	h := handler{logger: logger}

	sr.Path("/start").Methods(http.MethodPost).HandlerFunc(h.startProcessHandler)

	return nil
}
