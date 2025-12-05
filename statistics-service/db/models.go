package db

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

const (
	ProcessingStatus = "processing"
	DoneStatus       = "done"
)

type ProcessStatsInput struct {
	TrainingId   int64       `db:"training_id" json:"training_id"`
	UserLogin    string      `db:"user_login" json:"user_login"`
	StartTime    time.Time   `db:"start_time" json:"start_time"`
	FinishTime   time.Time   `db:"finish_time" json:"finish_time"`
	ResultValues ResultsJSON `db:"result_values" json:"result_values"`
}

type StatusDb struct {
	Id   int64  `db:"id"`
	Name string `db:"name"`
}

type TrainingResultDb struct {
	Id           int64       `db:"id" json:"id"`
	TrainingId   int64       `db:"training_id" json:"training_id"`
	UserLogin    string      `db:"user_name" json:"user_name"`
	StartTime    time.Time   `db:"start_time" json:"start_time"`
	FinishTime   time.Time   `db:"finish_time" json:"finish_time"`
	ResultValues ResultsJSON `db:"result_values" json:"result_values"`
	Status       string      `db:"status" json:"status"`
	Comment      *string     `db:"comment" json:"comment"`
	Kcal         *string     `db:"kcal" json:"kcal"`
}

type ResultsJSON []Exercise

type Exercise struct {
	Name string `json:"name"`
	Sets []Set  `json:"sets"`
}

type Set struct {
	Weight       float64 `json:"weight"`
	OriginalReps int64   `json:"original_reps"`
	RealReps     int64   `json:"real_reps"`
	CalPerSet    int32   `json:"cal_per_set"`
}

func (e *ResultsJSON) Scan(src interface{}) error {
	var data []byte
	switch v := src.(type) {
	case string:
		data = []byte(v)
	case []byte:
		data = v
	default:
		return nil
	}
	return json.Unmarshal(data, e)
}

func (e *ResultsJSON) Value() (driver.Value, error) {
	return json.Marshal(e)
}
