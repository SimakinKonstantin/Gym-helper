package app

import (
	"github.com/gorilla/mux"
	//httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

func (app *App) CreateServer(addr string) *http.Server {
	router := mux.NewRouter()

	router.HandleFunc("/programs/{id}", app.GetProgram).Methods(http.MethodGet)
	router.HandleFunc("/programs", app.GetPrograms).Methods(http.MethodGet)
	router.HandleFunc("/programs", app.CreateProgram).Methods(http.MethodPost)
	router.HandleFunc("/programs/{id}", app.DeleteProgram).Methods(http.MethodDelete)

	router.HandleFunc("/trainings/{id}", app.GetTraining).Methods(http.MethodGet)
	router.HandleFunc("/trainings", app.GetTrainings).Methods(http.MethodGet)
	router.HandleFunc("/trainings", app.CreateTraining).Methods(http.MethodPost)
	router.HandleFunc("/trainings/{id}", app.DeleteTraining).Methods(http.MethodDelete)
	router.HandleFunc("/trainings/{id}", app.UpdateTraining).Methods(http.MethodPatch)

	router.HandleFunc("/programs/{id}/trainings", app.GetProgramTrainings).Methods(http.MethodGet)
	router.HandleFunc("/programs/trainings", app.AddTrainingToProgram).Methods(http.MethodPost)
	router.HandleFunc("/programs/trainings", app.DeleteTrainingToProgram).Methods(http.MethodDelete)

	//router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(httpSwagger.URL("http://localhost:8082/swagger/doc.json")))

	return &http.Server{
		Addr:    addr,
		Handler: router,
	}
}
