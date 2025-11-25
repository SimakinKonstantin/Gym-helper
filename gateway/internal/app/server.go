package app

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
)

func (app *App) CreateServer(addr string) *http.Server {
	router := mux.NewRouter()
	router.HandleFunc("/sign-in", app.forwardAuth).Methods(http.MethodPost)
	router.HandleFunc("/sign-up", app.forwardAuth).Methods(http.MethodPost)

	router.Handle("/program", authMiddleware(http.HandlerFunc(app.forwardManage))).Methods(http.MethodGet)
	router.Handle("/programs", authMiddleware(http.HandlerFunc(app.forwardManage))).Methods(http.MethodGet)
	router.Handle("/program", authMiddleware(http.HandlerFunc(app.forwardManage))).Methods(http.MethodPost)
	router.Handle("/program", authMiddleware(http.HandlerFunc(app.forwardManage))).Methods(http.MethodDelete)

	router.Handle("/training", authMiddleware(http.HandlerFunc(app.forwardManage))).Methods(http.MethodGet)
	router.Handle("/trainings", authMiddleware(http.HandlerFunc(app.forwardManage))).Methods(http.MethodGet)
	router.Handle("/training", authMiddleware(http.HandlerFunc(app.forwardManage))).Methods(http.MethodPost)
	router.Handle("/training", authMiddleware(http.HandlerFunc(app.forwardManage))).Methods(http.MethodDelete)
	router.Handle("/training", authMiddleware(http.HandlerFunc(app.forwardManage))).Methods(http.MethodPatch)

	router.Handle("/program/trainings", authMiddleware(http.HandlerFunc(app.forwardManage))).Methods(http.MethodGet)
	router.Handle("/program/trainings", authMiddleware(http.HandlerFunc(app.forwardManage))).Methods(http.MethodPost)
	router.Handle("/program/trainings", authMiddleware(http.HandlerFunc(app.forwardManage))).Methods(http.MethodDelete)

	router.Handle("/program/start-training", authMiddleware(http.HandlerFunc(app.startTraining))).Methods(http.MethodPost)

	return &http.Server{
		Addr:    addr,
		Handler: router,
	}
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		_, err := checkToken(token)
		if err != nil {
			slog.Error("Ошибка авторизации: ", err.Error())
			w.WriteHeader(http.StatusUnauthorized)
		}
		next.ServeHTTP(w, r)
	})
}

func checkToken(token string) (string, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	if err != nil {
		return "", fmt.Errorf("ошибка валидации токена: %v", err)
	}
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		return claims["login"].(string), nil
	} else {
		return "", fmt.Errorf("не удалось подтвердить валидность токена")
	}
}
