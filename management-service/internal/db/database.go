package db

import (
	"cousework/internal"
	"database/sql"
	"embed"
	_ "embed"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

var (
	//go:embed scripts/program_add.sql
	addProgramSql string

	//go:embed scripts/training_get_by_id.sql
	getProgramByIdSql string

	//go:embed scripts/programs_get_by_user_login.sql
	getProgramByUserLoginSql string

	//go:embed scripts/program_delete.sql
	deleteProgramSql string

	//go:embed scripts/training_add.sql
	addTrainingSql string

	//go:embed scripts/training_update.sql
	updateTrainingSql string

	//go:embed scripts/training_get_by_id.sql
	getTrainingByIdSql string

	//go:embed scripts/trainings_get_by_user_id.sql
	getTrainingByUserLoginSql string

	//go:embed scripts/training_delete_by_id.sql
	deleteTrainingByIdSql string

	//go:embed scripts/program_training_add.sql
	addProgramTrainingsSql string

	//go:embed scripts/program_training_get.sql
	getProgramTrainingSql string

	//go:embed scripts/program_training_delete.sql
	deleteProgramTrainingSql string
)

type Database struct {
	Db *sqlx.DB
}

func NewDatabase(db *sqlx.DB) *Database {
	return &Database{Db: db}
}

func (r *Database) InitDb() error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("ошибка запуска миграций: %w", err)
	}

	if err := goose.Up(r.Db.DB, "migrations"); err != nil {
		return fmt.Errorf("ошибка запуска миграций: %w", err)
	}

	return nil
}

func (r *Database) ProgramCreate(program internal.ProgramDb) error {
	if _, err := r.Db.Exec(addProgramSql, program.UserLogin, program.Name); err != nil {
		return fmt.Errorf("ошибка добавления программы тренировки в БД: %w", err)
	}

	return nil
}

func (r *Database) ProgramGetById(id int64, login string) (internal.ProgramDb, error) {
	var res internal.ProgramDb
	if err := r.Db.Get(&res, getProgramByIdSql, id, login); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return internal.ProgramDb{}, fmt.Errorf("ошибка получения программы тренировки из БД: %w", err)
	}

	return res, nil
}

func (r *Database) ProgramsGetByUserLogin(login string) ([]internal.ProgramDb, error) {
	var res []internal.ProgramDb
	if err := r.Db.Select(&res, getProgramByUserLoginSql, login); err != nil {
		return nil, fmt.Errorf("ошибка получения программы тренировки из БД: %w", err)
	}

	return res, nil
}

func (r *Database) ProgramDeleteById(id int64, login string) error {
	if _, err := r.Db.Exec(deleteProgramSql, id, login); err != nil {
		return fmt.Errorf("ошибка удаления программы тренировки из БД: %w", err)
	}

	return nil
}

func (r *Database) TrainingCreate(training internal.TrainingDb) error {
	if _, err := r.Db.Exec(addTrainingSql, training.UserLogin, training.Name, training.Exercises); err != nil {
		return fmt.Errorf("ошибка добавления тренировки в БД: %w", err)
	}

	return nil
}

func (r *Database) TrainingGetById(id int64) (internal.TrainingDb, error) {
	var res internal.TrainingDb
	if err := r.Db.Get(&res, getTrainingByIdSql, id); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return internal.TrainingDb{}, fmt.Errorf("ошибка получения тренировки из БД: %w", err)
	}

	return res, nil
}

func (r *Database) TrainingDeleteById(id int64, login string) error {
	if _, err := r.Db.Exec(deleteTrainingByIdSql, id, login); err != nil {
		return fmt.Errorf("ошибка удаления тренировки из БД: %w", err)
	}

	return nil
}

func (r *Database) TrainingsGetByUserLogin(login string) ([]internal.TrainingDb, error) {
	var trainings []internal.TrainingDb

	if err := r.Db.Select(&trainings, getTrainingByUserLoginSql, login); err != nil {
		return nil, fmt.Errorf("ошибка получения тренировок пользователя из БД: %w", err)
	}

	return trainings, nil
}

func (r *Database) GetProgramTrainings(programId int64) ([]internal.ProgramTrainingDb, error) {
	var trainingDays []internal.TrainingDayDb

	var err error
	if err = r.Db.Select(&trainingDays, getProgramTrainingSql, programId); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("ошибка получения тренировок для программы тренировок из БД: %w", err)
	}

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	fmt.Println("trainingDays = ", trainingDays)

	// Формируем из вида который хранился в БД новую форму, которая будет использоваться для вывода.
	res := make([]internal.ProgramTrainingDb, len(trainingDays))
	for i, v := range trainingDays {
		program := internal.ProgramTrainingDb{
			Day: v.Day,
		}

		if v.Id == nil || v.Name == nil || v.Exercises == nil {
			program.Training = nil

		} else {
			program.Training = &internal.TrainingDb{
				Id:        *v.Id,
				Name:      *v.Name,
				Exercises: *v.Exercises,
			}
		}

		res[i] = program
	}

	return res, nil
}

func (r *Database) AddProgramTrainings(programId int64, trainingId *int64, day int64) error {
	if _, err := r.Db.Exec(addProgramTrainingsSql, programId, trainingId, day); err != nil {
		return fmt.Errorf("ошибка добавления связи программа тренировок - тренировка в БД: %w", err)
	}

	return nil
}

func (r *Database) DeleteProgramTrainings(programId int64, trainingId *int64, day int64) error {
	if _, err := r.Db.Exec(deleteProgramTrainingSql, trainingId, programId, day); err != nil {
		return fmt.Errorf("ошибка удаления связи программа тренировок - тренировка в БД: %w", err)
	}

	return nil
}

func (r *Database) UpdateTraining(id int64, training internal.TrainingDb) error {
	if _, err := r.Db.Exec(updateTrainingSql, training.Name, training.Exercises, id, training.UserLogin); err != nil {
		return fmt.Errorf("ошибка добавления связи программа тренировок - тренировка в БД: %w", err)
	}

	return nil
}
