package app

import (
	"cousework_auth/internal"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

// @Summary      Зарегистрироваться
// @Description  Создает нового пользователя
// @Tags         Аутентификация
// @Param        body body internal.SignUpReq true "Данные для регистрации"
// @Success      200 {string} string "Успешно зарегистрировался"
// @Failure      400 {string} string "Ошибки в передаваемых параметрах"
// @Failure      409 {string} string "Пользователь с таким логином уже существует"
// @Router       /sign-up [post]
func (app *App) SignUp(w http.ResponseWriter, r *http.Request) {
	var input internal.SignUpReq
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		processError(fmt.Errorf("ошибка анмаршаллинга запроса в SignUp: %w", err), http.StatusBadRequest, w)
		return
	}

	user, err := app.db.GetUser(input.Login)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		processError(fmt.Errorf("ошибка получения пользователя из БД в SignUp: %w", err), http.StatusBadRequest, w)
		return
	}

	// Проверка, что в БД нашли пользователя с таким логином.
	if user.Login != "" {
		processError(fmt.Errorf("попытка зарегестрироваться с уже занятым логином: %w", err), http.StatusConflict, w)
		return
	}

	// Кладем в БД логин, хеш
	hash, err := hashPassword(input.Password)
	if err != nil {
		processError(fmt.Errorf("ошибка формирования хеша в SignUp: %w", err), http.StatusBadRequest, w)
		return
	}

	if err := app.db.AddUser(input.Login, hash); err != nil {
		processError(fmt.Errorf("ошибка сохранения пользователя в БД в SignUp: %w", err), http.StatusBadRequest, w)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary      Войти
// @Description  Войти в существующий аккаунт
// @Tags         Аутентификация
// @Param        body body internal.SignInReq true "Данные для входа"
// @Success      200 {object} internal.SignInResp "JWT-токен"
// @Failure      400 {string} string "Ошибки в передаваемых параметрах"
// @Failure      401 {string} string "Неверный логин или пароль"
// @Router       /sign-in [post]
func (app *App) SignIn(w http.ResponseWriter, r *http.Request) {
	var input internal.SignInReq
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		processError(fmt.Errorf("ошибка анмаршаллинга запроса в SignIn: %w", err), http.StatusBadRequest, w)
		return
	}

	user, err := app.db.GetUser(input.Login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			processError(fmt.Errorf("такого пользователя не существует в SignIn: %w", err), http.StatusUnauthorized, w)
			return
		}

		processError(fmt.Errorf("ошибка получения пользователя из БД в SignIn: %w", err), http.StatusBadRequest, w)
		return
	}

	if !checkPasswordHash(input.Password, user.Hash) {
		processError(fmt.Errorf("неверный логин или пароль в SignIn: %w", err), http.StatusUnauthorized, w)
		return
	}

	token, err := buildToken(input.Login)
	if err != nil {
		processError(fmt.Errorf("ошибка получения токена в SignIn: %w", err), http.StatusBadRequest, w)
		return
	}

	resp := internal.SignInResp{
		Token: token,
	}
	marshalledResp, err := json.Marshal(resp)
	if err != nil {
		processError(fmt.Errorf("ошибка маршаллинга ответа в SignIn: %w", err), http.StatusBadRequest, w)
		return
	}

	if _, err = w.Write(marshalledResp); err != nil {
		processError(fmt.Errorf("ошибка записи ответа в SignIn: %w", err), http.StatusBadRequest, w)
		return
	}
}

func buildToken(login string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"login": login,
		"exp":   time.Now().Add(time.Hour * 336).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return "", fmt.Errorf("ошибка формирования токена: %w", err)
	}

	return tokenString, nil
}

func processError(err error, statusCode int, w http.ResponseWriter) {
	slog.Error(err.Error())
	w.WriteHeader(statusCode)
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
