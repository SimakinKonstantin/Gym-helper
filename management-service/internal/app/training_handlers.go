package app

import (
	"cousework/internal"
	"cousework/internal/metrics"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

// @Summary      Получить тренировку по id
// @Description  Получить тренировку по id
// @Tags         Тренировки
// @Param        id path int true "id тренировки"
// @Param 		 login header string true "Логин пользователя"
// @Success      200 {object} internal.TrainingDb "Тренировка"
// @Failure      400 "Ошибка"
// @Router       /trainings/{id} [get]
func (app *App) GetTraining(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	statusCode := http.StatusOK

	defer func() {
		metrics.ObserveRequest(time.Since(start), statusCode)
	}()

	idStr := mux.Vars(r)["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.processError(fmt.Errorf("ошибка парсинга id программы тренировок в GetTraining: %w", err), w, &statusCode)
		return
	}

	trainings, err := app.db.TrainingGetById(int64(id))
	if err != nil {
		app.processError(fmt.Errorf("ошибка получения тренировок в GetTraining: %w", err), w, &statusCode)
		return
	}

	marshalledTrainings, err := json.Marshal(trainings)
	if err != nil {
		app.processError(fmt.Errorf("ошибка получения тренировок в GetTraining: %w", err), w, &statusCode)
		return
	}

	w.WriteHeader(http.StatusOK)

	if _, err = w.Write(marshalledTrainings); err != nil {
		app.processError(fmt.Errorf("ошибка записи ответа в GetProgramTrainings: %w", err), w, &statusCode)
		return
	}
}

// @Summary      Получить все тренировки пользователя
// @Description  Получить все тренировки пользователя
// @Tags         Тренировки
// @Param        login path string true "Логин пользователя"
// @Success      200 {object} []internal.TrainingDb "Тренировки пользователя"
// @Failure      400 "Ошибка"
// @Router       /trainings [get]
func (app *App) GetTrainings(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	statusCode := http.StatusOK

	defer func() {
		metrics.ObserveRequest(time.Since(start), statusCode)
	}()

	login := r.Header.Get(loginHeader)

	fmt.Println("LOGIN = ", login)

	trainings, err := app.db.TrainingsGetByUserLogin(login)
	if err != nil {
		app.processError(fmt.Errorf("ошибка получения тренировок в GetTrainings: %w", err), w, &statusCode)
		return
	}

	marshalledTrainings, err := json.Marshal(trainings)
	if err != nil {
		app.processError(fmt.Errorf("ошибка получения тренировок в GetTrainings: %w", err), w, &statusCode)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(marshalledTrainings); err != nil {
		app.processError(fmt.Errorf("ошибка записи ответа в GetTrainings: %w", err), w, &statusCode)
		return
	}
}

// @Summary      Создать тренировку
// @Description  Создать тренировку
// @Tags         Тренировки
// @Param        body body internal.TrainingDb true "Данные для создания тренировки"
// @Success      200 "Тренировка создана"
// @Failure      400 "Ошибка"
// @Router       /trainings [post]
func (app *App) CreateTraining(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	statusCode := http.StatusOK

	defer func() {
		metrics.ObserveRequest(time.Since(start), statusCode)
	}()

	var training internal.TrainingDb
	if err := json.NewDecoder(r.Body).Decode(&training); err != nil {
		app.processError(fmt.Errorf("ошибка анмаршаллинга запроса на создание тренировки в CreateTraining: %w", err), w, &statusCode)
		return
	}

	if err := app.db.TrainingCreate(training); err != nil {
		app.processError(fmt.Errorf("ошибка сохранения новой тренировки в БД в CreateTraining: %w", err), w, &statusCode)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary      Удалить тренировку
// @Description  Удалить тренировку
// @Tags         Тренировки
// @Param        id path int true "id тренировки"
// @Param 		 X-User-Login header string true "Логин пользователя"
// @Success      200 "Тренировка удалена"
// @Failure      400 "Ошибка"
// @Router       /trainings/{id} [delete]
func (app *App) DeleteTraining(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	statusCode := http.StatusOK

	defer func() {
		metrics.ObserveRequest(time.Since(start), statusCode)
	}()

	idStr := mux.Vars(r)["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.processError(fmt.Errorf("ошибка парсинга id в DeleteTraining: %w", err), w, &statusCode)
		return
	}

	login := r.Header.Get(loginHeader)

	if err = app.db.TrainingDeleteById(int64(id), login); err != nil {
		app.processError(fmt.Errorf("ошибка удаления тренировки в DeleteTraining: %w", err), w, &statusCode)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary      Обновить тренировку
// @Description  Обновить тренировку
// @Tags         Тренировки
// @Param        id path int true "id тренировки"
// @Param        body body internal.TrainingDb true "Данные для обновления тренировки"
// @Success      200 "Тренировка обновлена"
// @Failure      400 "Ошибка"
// @Router       /trainings/{id} [patch]
func (app *App) UpdateTraining(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	statusCode := http.StatusOK

	defer func() {
		metrics.ObserveRequest(time.Since(start), statusCode)
	}()

	var training internal.TrainingDb
	if err := json.NewDecoder(r.Body).Decode(&training); err != nil {
		app.processError(fmt.Errorf("ошибка анмаршаллинга запроса на создание тренировки в UpdateTraining: %w", err), w, &statusCode)
		return
	}

	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.processError(fmt.Errorf("ошибка парсинга id в UpdateTraining: %w", err), w, &statusCode)
		return
	}

	if err = app.db.UpdateTraining(int64(id), training); err != nil {
		app.processError(fmt.Errorf("ошибка обновления тренировки в БД: %w", err), w, &statusCode)
		return
	}
}
