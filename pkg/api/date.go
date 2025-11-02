package api

import (
	"fmt"
	"time"

	"final_project/pkg/db"
	"final_project/pkg/nextdate"
)

// checkDate проверяет и корректирует дату задачи
// Если дата пустая - устанавливает текущую дату
// Если дата в прошлом и есть правило повторения - вычисляет следующую дату
// Если дата в прошлом и нет правила повторения - устанавливает текущую дату
func checkDate(task *db.Task) error {
	now := time.Now()

	// Если дата не указана, используем текущую дату
	if task.Date == "" {
		task.Date = now.Format(DateFormat)
		return nil
	}

	// Проверяем корректность формата даты
	t, err := time.Parse(DateFormat, task.Date)
	if err != nil {
		return fmt.Errorf("дата представлена в формате, отличном от 20060102")
	}

	// Если дата в прошлом
	if afterNow(now, t) {
		if task.Repeat == "" {
			// Если правила повторения нет, используем сегодняшнюю дату
			task.Date = now.Format(DateFormat)
		} else {
			// Если есть правило повторения, вычисляем следующую дату
			next, err := nextdate.NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return fmt.Errorf("правило повторения указано в неправильном формате: %w", err)
			}
			task.Date = next
		}
	}

	return nil
}

// afterNow проверяет, что первая дата больше второй (игнорируя время)
func afterNow(date, now time.Time) bool {
	// Нормализуем даты, убирая время
	dateOnly := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	nowOnly := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	return dateOnly.After(nowOnly)
}
