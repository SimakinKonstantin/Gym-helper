package app

import (
	"coursework_gateway/metrics"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log/slog"
	"net/http"
	"strings"
)

func (app *App) CreateServer(addr string) *http.Server {
	router := mux.NewRouter()
	router.HandleFunc("/sign-in", app.forwardAuth).Methods(http.MethodPost)
	router.HandleFunc("/sign-up", app.forwardAuth).Methods(http.MethodPost)

	router.Handle("/programs", authMiddleware(http.HandlerFunc(app.forwardManage))).Methods(http.MethodGet)
	router.Handle("/programs", authMiddleware(http.HandlerFunc(app.forwardManage))).Methods(http.MethodPost)
	router.Handle("/programs/{id}", authMiddleware(http.HandlerFunc(app.forwardManage))).Methods(http.MethodDelete)

	router.Handle("/trainings/{id}", authMiddleware(http.HandlerFunc(app.forwardManage))).Methods(http.MethodGet)
	router.Handle("/trainings", authMiddleware(http.HandlerFunc(app.forwardManage))).Methods(http.MethodGet)
	router.Handle("/trainings", authMiddleware(http.HandlerFunc(app.forwardManage))).Methods(http.MethodPost)
	router.Handle("/trainings/{id}", authMiddleware(http.HandlerFunc(app.forwardManage))).Methods(http.MethodDelete)
	router.Handle("/trainings/{id}", authMiddleware(http.HandlerFunc(app.forwardManage))).Methods(http.MethodPatch)

	router.Handle("/programs/{id}/trainings", authMiddleware(http.HandlerFunc(app.forwardManage))).Methods(http.MethodGet)
	router.Handle("/programs/trainings", authMiddleware(http.HandlerFunc(app.forwardManage))).Methods(http.MethodPost)
	router.Handle("/programs/trainings", authMiddleware(http.HandlerFunc(app.forwardManage))).Methods(http.MethodDelete)

	router.Handle("/statistics/process-training", authMiddleware(http.HandlerFunc(app.processTraining))).Methods(http.MethodPost)
	router.Handle("/statistics", authMiddleware(http.HandlerFunc(app.forwardStatistics))).Methods(http.MethodGet)
	router.Handle("/statistics/{id}", authMiddleware(http.HandlerFunc(app.forwardStatistics))).Methods(http.MethodGet)

	router.Handle("/metrics", promhttp.Handler())

	return &http.Server{
		Addr:    addr,
		Handler: router,
	}
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		slog.Info("URL: ", r.URL.String())
		token := r.Header.Get("Authorization")
		token = strings.TrimPrefix(token, "Bearer ")

		_, err := checkToken(token)
		if err != nil {
			slog.Error("Ошибка авторизации: ", err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			metrics.ErrorMetrics.Inc()
		}
		next.ServeHTTP(w, r)
	})
}

func checkToken(token string) (string, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	if err != nil {
		return "", fmt.Errorf("ошибка валидации токена: %w", err)
	}
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		return claims["login"].(string), nil
	} else {
		return "", fmt.Errorf("не удалось подтвердить валидность токена")
	}
}
