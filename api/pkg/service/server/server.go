package server

import (
	"net/http"
	"time"

	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/manufacturer/api/pkg/service"
	"github.com/manufacturer/api/pkg/service/transport"
	"github.com/manufacturer/api/pkg/v1"
)

type Server struct {
	Logger          log.Logger
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
	MetricPrefix    string
}

func makeHandler(svc service.Service) http.Handler {
	var opts = []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(v1.EncodeError),
	}
	router := mux.NewRouter()

	router.Path("/api/v1/manufacturer").Methods(http.MethodPost).Handler(kithttp.NewServer(
		transport.MakePostDoLogicWithManufacturerEndpoint(svc),
		transport.DecodePostDoLogicWithManufacturerRequest,
		transport.EncodePostDoLogicWithManufacturerResponse,
		opts...,
	))

	return router
}

// NewServer creates a new server.
func NewServer(cfg *Server) *http.Server {
	svc := service.NewService()

	svc = service.NewLoggingMiddleware(svc, cfg.Logger)
	svc = service.NewInstrumentingMiddleware(svc, cfg.MetricPrefix+"_api")

	router := http.NewServeMux()
	router.Handle("/metrics", promhttp.Handler())
	router.Handle("/api/v1/", makeHandler(svc))

	return &http.Server{
		Handler:      router,
		Addr:         ":" + cfg.Port,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}
}
