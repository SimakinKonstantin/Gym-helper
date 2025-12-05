package app

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
	"statistics-service/metrics"
	"strconv"
	"time"
)

const loginHeader = "X-User-Login"

func (app *App) GetResults(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	statusCode := http.StatusOK

	defer func() {
		metrics.ObserveRequest(time.Since(start), statusCode)
	}()

	login := r.Header.Get(loginHeader)

	results, err := app.db.GetAllResults(login)
	if err != nil {
		statusCode = http.StatusBadRequest
		processError(fmt.Errorf("ошибка получения результатов тренировок в GetResults: %w", err), w, statusCode)
		return
	}

	marshalledResults, err := json.Marshal(results)
	if err != nil {
		statusCode = http.StatusBadRequest
		processError(fmt.Errorf("ошибка маршаллинга результатов тренировок в GetResults: %w", err), w, statusCode)
		return
	}

	if _, err := w.Write(marshalledResults); err != nil {
		statusCode = http.StatusBadRequest
		processError(fmt.Errorf("ошибка записи результатов в GetResults: %w", err), w, statusCode)
		return
	}
}

func (app *App) GetResult(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	statusCode := http.StatusOK

	defer func() {
		metrics.ObserveRequest(time.Since(start), statusCode)
	}()

	login := r.Header.Get(loginHeader)

	idStr := mux.Vars(r)["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		statusCode = http.StatusBadRequest
		processError(fmt.Errorf("ошибка получения id в GetResult: %w", err), w, statusCode)
		return
	}

	result, err := app.db.GetResultById(login, int64(id))
	if err != nil {
		statusCode = http.StatusBadRequest
		processError(fmt.Errorf("ошибка получения результатов тренировок в GetResult: %w", err), w, statusCode)
		return
	}

	marshalledResult, err := json.Marshal(result)
	if err != nil {
		statusCode = http.StatusBadRequest
		processError(fmt.Errorf("ошибка маршаллинга результатов тренировок в GetResult: %w", err), w, statusCode)
		return
	}

	if _, err := w.Write(marshalledResult); err != nil {
		statusCode = http.StatusBadRequest
		processError(fmt.Errorf("ошибка записи результатов в GetResult: %w", err), w, statusCode)
		return
	}
}

func processError(err error, w http.ResponseWriter, statusCode int) {
	slog.Error(err.Error())
	w.WriteHeader(statusCode)
	metrics.ErrorMetrics.Inc()
}
