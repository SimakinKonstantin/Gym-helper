package app

import (
	"fmt"
	"log/slog"
	"statistics-service/db"
)

func (app *App) ProcessStats(input db.ProcessStatsInput) error {
	setsCount := getSetsCount(input.ResultValues)

	slog.Info(fmt.Sprintf("Обработка соощения в ProcessStats: %+v", input))

	id, err := app.saveStats(input)

	if err != nil {
		return fmt.Errorf("ошибка обработки статистики в ProcessStats: %w", err)
	}

	var totalCal int32

	goodSet := 0
	badSet := 0

	// Подсчет кол-ва
	for i := 0; i < len(input.ResultValues); i++ { // Проходим по всем упражнениям
		for j := 0; j < len(input.ResultValues[i].Sets); j++ { // Проходим по всем подходам
			setsCount++
			setCal := calculateCal(input.ResultValues[i].Sets[j].CalPerSet, int32(input.ResultValues[i].Sets[j].RealReps))
			totalCal += setCal
			if input.ResultValues[i].Sets[j].RealReps >= input.ResultValues[i].Sets[j].OriginalReps {
				goodSet++
			} else {
				badSet++
			}
		}
	}

	slog.Info("cal = ", totalCal)

	if err = app.updateStats(id, input.UserLogin, totalCal, goodSet, badSet); err != nil {
		return fmt.Errorf("ошибка сохранения результатов тренировки в ProcessStats: %v", err)
	}

	return nil
	// Считать кол-во сженных ккал x (добавлять в БД)
	// Считать кол-во хорошо выполненных упражнений. Дальше расчитываем по процентам и пишем в комментарий.
	// На фронте отображаем все подходы и план, а на нижней планке - время тренировки, общее кол-во ккал, комментарий
	//
	// На фронте нажимаем "завершить тренироку" - на gateway дергаем "start processing" (кладем в БД), затем на фронт кидаем что-то вроде accepted.
	// Далее либо на фронте не отрисовываем элементы, которые еще не обработаны, либо дергаем short polling (достаем только дата тренировки).
	// После того, как обработали через kafka, изменяем статус.
}

// Возвращаем id сохранененной статистики.
func (app *App) saveStats(stats db.ProcessStatsInput) (int64, error) {
	id, err := app.db.SaveBeforeProcessing(stats)
	if err != nil {
		return -1, fmt.Errorf("ошибка сохранения статистики перед обработкой: %w", err)
	}

	return id, nil
}

func (app *App) updateStats(id int64, userLogin string, cal int32, goodSet, badSet int) error {
	err := app.db.SaveKcal(userLogin, id, cal/1000)
	if err != nil {
		return fmt.Errorf("ошибка сохранения Kcal в БД: %w", err)
	}

	comment := generateComment(goodSet, badSet)
	if err = app.db.SaveComment(userLogin, id, comment); err != nil {
		return fmt.Errorf("ошибка сохранения комментария тренировки: %w", err)
	}

	if err = app.db.SetProcessedStatus(userLogin, id, db.DoneStatus); err != nil {
		return fmt.Errorf("ошибка установки статуса статистики: %w", err)
	}

	return nil
}

func calculateCal(calPerSet int32, reps int32) int32 {
	return calPerSet * reps
}

func getSetsCount(training db.ResultsJSON) int {
	setsCount := 0

	for i := 0; i < len(training); i++ { // Проходим по всем упражнениям
		for j := 0; j < len(training[i].Sets); j++ {
			setsCount++
		}
	}

	return setsCount
}

func generateComment(goodSet, badSet int) string {
	totalSets := goodSet + badSet
	if totalSets == 0 {
		return "Нет данных для анализа"
	}

	badPercentage := float64(badSet) / float64(totalSets) * 100

	// Вариативные комментарии в зависимости от процента неудачных подходов
	if badPercentage > 10 {
		comments := []string{
			"Слишком много неудачных подходов. Рекомендуется уменьшить рабочий вес.",
			"Слишком много неудачных подходов. Попробуйте снизить нагрузку и уделить больше внимания технике.",
			"Некоторые подходы выполнены неудачно. Рассмотрите возможность корректировки веса и техники выполнения.",
			"Тренировка получилась тяжелой? Слишком много неудачных попыток. Попробуйте немного снизить вес.",
		}
		return comments[badSet%len(comments)]
	} else if badPercentage > 5 {
		comments := []string{
			"Умеренное количество неудачных подходов (>5%). Рекомендуется обратить внимание на технику выполнения.",
			"Некоторые подходы не получились. Уделите внимание технике выполнения упражнений.",
			"Небольшой процент неудачных подходов. Попробуйте немного скорректировать технику.",
			"Тренировка прошла неплохо, но есть несколько неудачных подходов. Обратите внимание на выполнение.",
			"Почти все подходы выполнены, но несколько оказались сложными. Повторите технику выполнения.",
		}
		return comments[badSet%len(comments)]
	} else {
		comments := []string{
			"Отличный результат! Продолжайте в том же духе!",
			"Потрясающе! Все подходы выполнены на отлично. Так держать!",
			"Превосходный результат! Техника на высоте, нагрузка подобрана идеально.",
			"Замечательная тренировка! Молодец!",
			"Отличная работа! Все получилось - продолжайте в том же духе!",
			"Все подходы выполнены безупречно. Вы на правильном пути!",
		}
		return comments[goodSet%len(comments)]
	}
}
