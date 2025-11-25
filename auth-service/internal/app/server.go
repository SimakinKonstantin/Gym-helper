package app

import (
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

func (app *App) CreateServer(addr string) *http.Server {
	router := mux.NewRouter()

	router.HandleFunc("/sign-up", app.SignUp).Methods(http.MethodPost)
	router.HandleFunc("/sign-in", app.SignIn).Methods(http.MethodPost)
	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(httpSwagger.URL("http://localhost:8081/swagger/doc.json")))
	return &http.Server{
		Addr:    addr,
		Handler: router,
	}
}
