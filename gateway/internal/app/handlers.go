package app

import (
	"coursework_gateway/metrics"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func (app *App) forwardAuth(w http.ResponseWriter, r *http.Request) {
	newURL, _ := url.Parse(os.Getenv("AUTH_SERVICE_URL"))
	forwardReq(w, r, newURL)
}

func (app *App) forwardManage(w http.ResponseWriter, r *http.Request) {
	newUrl, _ := url.Parse(os.Getenv("MANAGEMENT_SERVICE_URL"))
	forwardReq(w, r, newUrl)
}

func (app *App) forwardStatistics(w http.ResponseWriter, r *http.Request) {
	newUrl, _ := url.Parse(os.Getenv("STATISTICS_SERVICE_URL"))
	forwardReq(w, r, newUrl)
}

func (app *App) processTraining(w http.ResponseWriter, r *http.Request) {
	var input ProcessTrainingInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		slog.Error("Не удалось получить значение request", err.Error())
		metrics.ErrorMetrics.Inc()
		return
	}

	slog.Info(fmt.Sprintf("Отправили сообщение в Kafka: %+v", input))

	if err := app.Write(input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		slog.Error("Не удалось записать сообщение в kafka", err.Error())
		metrics.ErrorMetrics.Inc()
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func forwardReq(w http.ResponseWriter, r *http.Request, newUrl *url.URL) {
	proxy := httputil.NewSingleHostReverseProxy(newUrl)
	proxy.ServeHTTP(w, r)
}
