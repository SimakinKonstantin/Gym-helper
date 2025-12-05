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

// @Summary      Получить тренировки конкретной программы
// @Description  Получить тренировки конкретной программы
// @Tags         Программы тренировок
// @Param        id path int true "id программы тренировок"
// @Param 		 X-User-Login header string true "Логин пользователя"
// @Success      200 {object} []internal.TrainingDb "Тренировки программы"
// @Failure      400 "Ошибка"
// @Router       /programs/{id}/trainings [get]
// GET /program/trainings (получить тренировки конкретной программы)
func (app *App) GetProgramTrainings(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	statusCode := http.StatusOK

	defer func() {
		metrics.ObserveRequest(time.Since(start), statusCode)
	}()

	idStr := mux.Vars(r)["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.processError(fmt.Errorf("ошибка парсинга id программы тренировок в GetProgramTrainings: %w", err), w, &statusCode)
		return
	}

	trainings, err := app.db.GetProgramTrainings(int64(id))
	if err != nil {
		app.processError(fmt.Errorf("ошибка получения тренировок в GetProgramTrainings: %w", err), w, &statusCode)
		return
	}

	marshalledTrainings, err := json.Marshal(trainings)
	if err != nil {
		app.processError(fmt.Errorf("ошибка получения тренировок в GetProgramTrainings: %w", err), w, &statusCode)
		return
	}

	w.WriteHeader(http.StatusOK)

	if _, err = w.Write(marshalledTrainings); err != nil {
		app.processError(fmt.Errorf("ошибка записи ответа в GetProgramTrainings: %w", err), w, &statusCode)
		return
	}
}

// @Summary      Добавить тренировку в программу тренировок
// @Description  Добавить тренировку в программу тренировок
// @Tags         Программы тренировок
// @Param        body body internal.AddTrainingToProgramReq true "Данные для добавления тренировки в программу"
// @Success      200 "Тренировка добавлена в программу"
// @Failure      400 "Ошибка"
// @Router       /programs/trainings [post]
// POST /program/training (добавить тренировку в программу тренировок).
func (app *App) AddTrainingToProgram(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	statusCode := http.StatusOK

	defer func() {
		metrics.ObserveRequest(time.Since(start), statusCode)
	}()

	var req internal.AddTrainingToProgramReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		app.processError(fmt.Errorf("ошибка анмаршаллинга запроса в AddTrainingToProgram: %w", err), w, &statusCode)
		return
	}

	if err := app.db.AddProgramTrainings(req.ProgramId, req.TrainingId, req.Day); err != nil {
		app.processError(fmt.Errorf("ошибка добавления связи программа-программа тренировок в AddTrainingToProgram: %w", err), w, &statusCode)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary      Удалить тренировку из программы тренировок
// @Description  Удалить тренировку из программы тренировок
// @Tags         Программы тренировок
// @Param        body body internal.DeleteTrainingToProgramReq true "Данные для удаления тренировки из программы"
// @Success      200 "Тренировка удалена из программы"
// @Failure      400 "Ошибка"
// @Router       /programs/trainings [delete]
func (app *App) DeleteTrainingToProgram(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	statusCode := http.StatusOK

	defer func() {
		metrics.ObserveRequest(time.Since(start), statusCode)
	}()

	var req internal.DeleteTrainingToProgramReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		app.processError(fmt.Errorf("ошибка анмаршаллинга запроса в DeleteTraininToProgram: %w", err), w, &statusCode)
		return
	}

	if err := app.db.DeleteProgramTrainings(req.ProgramId, req.TrainingId, req.Day); err != nil {
		app.processError(fmt.Errorf("ошибка удаления связи программа - тренировка в DeleteTrainingToProgram: %w", err), w, &statusCode)
		return
	}

	w.WriteHeader(http.StatusOK)
}
