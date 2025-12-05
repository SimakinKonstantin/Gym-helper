package internal

import (
	"database/sql/driver"
	"encoding/json"
)

type CreateProgramReq struct {
	UserLogin string `json:"user_login"`
	Name      string `json:"name"`
}

type AddTrainingToProgramReq struct {
	ProgramId  int64  `json:"program_id"`
	TrainingId *int64 `json:"training_id"`
	Day        int64  `json:"day"`
}

type DeleteTrainingToProgramReq struct {
	ProgramId  int64  `json:"program_id"`
	TrainingId *int64 `json:"training_id"`
	Day        int64  `json:"day"`
}

// Программа тренировок.
type ProgramDb struct {
	Id        int64      `db:"id" json:"id"`
	UserLogin string     `db:"user_login" json:"user_login"`
	Name      string     `db:"name" json:"name"`
	Trainings TrainingDb `db:"trainings" json:"trainings"`
}

// Содержит информацию о тренировке и дне этой тренировки.
type TrainingDayDb struct {
	Day       int64          `db:"day"`
	Id        *int64         `db:"id"`
	Name      *string        `db:"name"`
	Exercises *ExercisesJSON `db:"exercises"`
}

type ProgramTrainingDb struct {
	Day      int64       `json:"day"`
	Training *TrainingDb `json:"training"`
}

// Тренировка.
type TrainingDb struct {
	Id        int64         `db:"id" json:"id"`
	UserLogin string        `db:"user_login" json:"user_login"`
	Name      string        `db:"name" json:"name"`
	Exercises ExercisesJSON `db:"exercises" json:"exercises"`
}

type Set struct {
	Weight    float64 `json:"weight"`
	Reps      int64   `json:"reps"`
	CalPerSet int32   `json:"cal_per_set"`
}

type Exercise struct {
	Name string `json:"name"`
	Sets []Set  `json:"sets"`
}

type ExercisesJSON []Exercise

func (e *ExercisesJSON) Scan(src interface{}) error {
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

func (e ExercisesJSON) Value() (driver.Value, error) {
	return json.Marshal(e)
}
