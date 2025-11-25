package app

import (
	"encoding/json"
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

func (app *App) startTraining(w http.ResponseWriter, r *http.Request) {
	var input ProcessStatsInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		slog.Error("Не удалось получить значение request", err.Error())
		return
	}

	if err := app.Write(input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		slog.Error("Не удалось записать сообщение в kafka", err.Error())
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func forwardReq(w http.ResponseWriter, r *http.Request, newUrl *url.URL) {
	// Обрезаем, чтобы получить независимый от base url.
	//cutFrom := len("gateway/")
	//r.URL.Path = r.URL.Path[cutFrom:]
	proxy := httputil.NewSingleHostReverseProxy(newUrl)
	proxy.ServeHTTP(w, r)
}
