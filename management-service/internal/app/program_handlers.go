package app

import (
	"cousework/internal"
	"cousework/internal/metrics"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

const loginHeader = "X-User-Login"

// @Summary      Получить программу тренировок
// @Description  Получить программу тренировок
// @Tags         Программы тренировок
// @Param 		 X-User-Login header string true "Логин пользователя"
// @Param        id path int true "id программы тренировок"
// @Success      200 {object} internal.ProgramDb "Программа тренировок"
// @Failure      400 "Ошибка"
// @Router       /programs/{id} [get]
func (app *App) GetProgram(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	statusCode := http.StatusOK

	defer func() {
		metrics.ObserveRequest(time.Since(start), statusCode)
	}()

	idStr := mux.Vars(r)["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.processError(fmt.Errorf("ошибка парсинга id пользователя в GetProgram: %w", err), w, &statusCode)
		return
	}

	login := r.Header.Get(loginHeader)

	program, err := app.db.ProgramGetById(int64(id), login)
	if err != nil {
		app.processError(fmt.Errorf("ошибка получения программы тренировок по userId в GetProgram: %w", err), w, &statusCode)
		return
	}

	marshalledProgram, err := json.Marshal(program)
	if err != nil {
		app.processError(fmt.Errorf("ошибка маршаллинга ответа при получении программы тренировок в GetProgram: %w", err), w, &statusCode)
		return
	}

	if _, err = w.Write(marshalledProgram); err != nil {
		app.processError(fmt.Errorf("В GetUserProgram ошибка записи ответа: %s", err), w, &statusCode)
		return
	}
}

// @Summary      Удалить программу тренировок
// @Description  Удалить программу тренировок
// @Tags         Программы тренировок
// @Param 		 X-User-Login header string true "Логин пользователя"
// @Param        id path int true "id программы тренировок"
// @Success      200 "Успешно удалили"
// @Failure      400 "Ошибка"
// @Router       /programs/{id} [delete]
func (app *App) DeleteProgram(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	statusCode := http.StatusOK

	defer func() {
		metrics.ObserveRequest(time.Since(start), statusCode)
	}()

	idStr := mux.Vars(r)["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.processError(fmt.Errorf("ошибка парсинга id программы тренировок в DeleteProgram: %w", err), w, &statusCode)
		return
	}

	login := r.Header.Get(loginHeader)

	if err := app.db.ProgramDeleteById(int64(id), login); err != nil {
		app.processError(fmt.Errorf("ошибка удаления программы тренировок в DeleteProgram: %w", err), w, &statusCode)
		return
	}

	w.WriteHeader(statusCode)
}

// @Summary      Получить все программы тренировок пользователя
// @Description  Получить все программы тренировок пользователя
// @Tags         Программы тренировок
// @Param 		 X-User-Login header string true "Логин пользователя"
// @Success      200 {object} []internal.ProgramDb "Программы тренировок"
// @Failure      400 "Ошибка"
// @Router       /programs [get]
// GET /programs
func (app *App) GetPrograms(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	statusCode := http.StatusOK

	defer func() {
		metrics.ObserveRequest(time.Since(start), statusCode)
	}()

	login := r.Header.Get(loginHeader)

	programs, err := app.db.ProgramsGetByUserLogin(login)
	if err != nil {
		app.processError(fmt.Errorf("ошибка получения программы тренировок по userId в GetPrograms: %w", err), w, &statusCode)
		return
	}

	marshalledPrograms, err := json.Marshal(programs)
	if err != nil {
		app.processError(fmt.Errorf("ошибка маршаллинга ответа при получении программы тренировок в GetPrograms: %w", err), w, &statusCode)
		return
	}

	w.WriteHeader(statusCode)
	if _, err = w.Write(marshalledPrograms); err != nil {
		app.processError(fmt.Errorf("В GetPrograms ошибка записи ответа: %s", err), w, &statusCode)
		return
	}
}

// @Summary      Создать программу тренировок
// @Description  Создать программу тренировок
// @Tags         Программы тренировок
// @Param        body body internal.CreateProgramReq true "Данные для создания программы"
// @Success      201 "Программа тренировок создана"
// @Failure      400 "Ошибка"
// @Router       /programs [post]
// POST /program
func (app *App) CreateProgram(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	statusCode := http.StatusOK

	defer func() {
		metrics.ObserveRequest(time.Since(start), statusCode)
	}()

	var program internal.CreateProgramReq
	if err := json.NewDecoder(r.Body).Decode(&program); err != nil {
		app.processError(fmt.Errorf("ошибка анмаршаллинга в CreateProgram: %w", err), w, &statusCode)
		return
	}

	if err := app.db.ProgramCreate(internal.ProgramDb{
		Name:      program.Name,
		UserLogin: program.UserLogin,
	}); err != nil {
		app.processError(fmt.Errorf("ошибка создания программы тренировок в CreateProgram: %w", err), w, &statusCode)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (app *App) processError(err error, w http.ResponseWriter, statusCode *int) {
	slog.Error(err.Error())
	w.WriteHeader(http.StatusBadRequest)
	metrics.ErrorMetrics.Inc()
	*statusCode = http.StatusBadRequest
}
