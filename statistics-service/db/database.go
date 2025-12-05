package db

import (
	_ "embed"
	"fmt"

	"github.com/jmoiron/sqlx"
)

var (
	//go:embed scripts/get_all_results.sql
	getAllResults string

	//go:embed scripts/get_result_by_id.sql
	getResultById string

	//go:embed scripts/set_stats.sql
	setStats string

	//go:embed scripts/set_kcal.sql
	setKcal string

	//go:embed scripts/set_comment.sql
	setComment string

	//go:embed scripts/get_status.sql
	getStatus string

	//go:embed scripts/set_status.sql
	setStatus string
)

type Database struct {
	Db *sqlx.DB
}

func NewDatabase(db *sqlx.DB) *Database {
	return &Database{Db: db}
}

func (d *Database) GetAllResults(login string) ([]TrainingResultDb, error) {
	var results []TrainingResultDb
	if err := d.Db.Select(&results, getAllResults, login); err != nil {
		return nil, fmt.Errorf("ошибка получения в GetAllResults: %w", err)
	}

	return results, nil
}

func (d *Database) GetResultById(login string, id int64) (TrainingResultDb, error) {
	var result TrainingResultDb
	if err := d.Db.Get(&result, getResultById, login, id); err != nil {
		return TrainingResultDb{}, fmt.Errorf("ошибка получения в GetResultById: %w", err)
	}

	return result, nil
}

func (d *Database) SaveBeforeProcessing(stats ProcessStatsInput) (int64, error) {
	var status StatusDb
	err := d.Db.Get(&status, getStatus, ProcessingStatus)
	if err != nil {
		return -1, fmt.Errorf("ошибка получения статуса в SaveBeforeProcessing: %w", err)
	}

	resultValues, err := stats.ResultValues.Value()
	if err != nil {
		return -1, fmt.Errorf("ошибка сериализации result_values в SaveBeforeProcessing: %w", err)
	}

	var id int64
	err = d.Db.Get(&id, setStats, stats.UserLogin, stats.TrainingId, stats.StartTime, stats.FinishTime, resultValues, status.Id)
	if err != nil {
		return -1, fmt.Errorf("ошибка обновления в SaveBeforeProcessing: %w", err)
	}

	return id, nil
}

func (d *Database) SaveKcal(login string, id int64, kcal int32) error {
	if _, err := d.Db.Exec(setKcal, kcal, id, login); err != nil {
		return fmt.Errorf("ошибка обновления в SaveKcal: %w", err)
	}

	return nil
}

func (d *Database) SaveComment(login string, id int64, comment string) error {
	if _, err := d.Db.Exec(setComment, comment, id, login); err != nil {
		return fmt.Errorf("ошибка обновления в SaveComment: %w", err)
	}

	return nil
}

func (d *Database) SetProcessedStatus(login string, id int64, statusName string) error {
	var status StatusDb
	err := d.Db.Get(&status, getStatus, statusName)
	if err != nil {
		return fmt.Errorf("ошибка получения статуса в SetProcessedStatus: %w", err)
	}

	_, err = d.Db.Exec(setStatus, status.Id, id, login)
	if err != nil {
		return fmt.Errorf("ошибка установки статуса в SetProcessedStatus: %w", err)
	}

	return nil
}
