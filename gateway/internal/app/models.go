package app

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type ProcessTrainingInput struct {
	Id           int64       `json:"training_id"`
	UserLogin    string      `json:"user_login"`
	StartTime    time.Time   `json:"start_time"`
	FinishTime   time.Time   `json:"finish_time"`
	ResultValues ResultsJSON `json:"result_values"`
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
