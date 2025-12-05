package app

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

func (app *App) CreateServer(addr string) *http.Server {
	router := mux.NewRouter()
	router.HandleFunc("/statistics", app.GetResults).Methods(http.MethodGet)
	router.HandleFunc("/statistics/{id}", app.GetResult).Methods(http.MethodGet)
	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(httpSwagger.URL("http://localhost:8083/swagger/doc.json")))
	router.Handle("/metrics", promhttp.Handler())
	return &http.Server{
		Addr:    addr,
		Handler: router,
	}
}
