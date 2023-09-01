package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/lmzuccarelli/golang-redis-publisher/pkg/connectors"
	"github.com/lmzuccarelli/golang-redis-publisher/pkg/handlers"
	"github.com/lmzuccarelli/golang-redis-publisher/pkg/validator"
	"github.com/microlib/simple"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	CONTENTTYPE     string = "Content-Type"
	APPLICATIONJSON string = "application/json"
)

var (
	logger       *simple.Logger
	httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "redis_publisher_http_duration_seconds",
		Help: "Duration of HTTP requests.",
	}, []string{"path"})
)

// prometheusMiddleware implements mux.MiddlewareFunc.
func prometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(CONTENTTYPE, APPLICATIONJSON)
		// use this for cors
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept-Language")
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()
		timer := prometheus.NewTimer(httpDuration.WithLabelValues(path))
		next.ServeHTTP(w, r)
		timer.ObserveDuration()
	})
}

func startHttpServer(con connectors.Clients) *http.Server {
	srv := &http.Server{Addr: ":" + os.Getenv("SERVER_PORT")}
	logger.Info("Starting server on port " + srv.Addr)

	r := mux.NewRouter()
	r.Use(prometheusMiddleware)
	r.Path("/metrics").Handler(promhttp.Handler())

	r.HandleFunc("/api/v1/publish", func(w http.ResponseWriter, req *http.Request) {
		handlers.SendPayloadHandler(w, req, con)
	}).Methods("POST", "OPTIONS")

	r.HandleFunc("/api/v1/isalive", handlers.IsAlive).Methods("GET")

	http.Handle("/", r)

	if err := srv.ListenAndServe(); err != nil {
		con.Error("Httpserver: ListenAndServe() error: " + err.Error())
	}

	return srv
}

func main() {

	if os.Getenv("LOG_LEVEL") == "" {
		logger = &simple.Logger{Level: "info"}
	} else {
		logger = &simple.Logger{Level: os.Getenv("LOG_LEVEL")}
	}

	err := validator.ValidateEnvars(logger)
	if err != nil {
		os.Exit(-1)
	}

	conn := connectors.NewClientConnections(logger)
	startHttpServer(conn)
}
