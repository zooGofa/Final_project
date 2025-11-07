package db

import (
	"fmt"
	"strconv"
	"time"
)

// Task представляет задачу в системе планировщика
type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// AddTask добавляет задачу в таблицу scheduler и возвращает ID созданной записи
func AddTask(task *Task) (int64, error) {
	var id int64

	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, fmt.Errorf("ошибка при добавлении задачи: %w", err)
	}

	id, err = res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("ошибка при получении ID задачи: %w", err)
	}

	return id, nil
}

// Tasks возвращает список ближайших задач, отсортированных по дате
// limit - максимальное количество возвращаемых записей
// search - строка поиска (опционально)
func Tasks(limit int, search string) ([]*Task, error) {
	var query string
	var args []interface{}

	if search == "" {
		// Обычный запрос без поиска
		query = `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC LIMIT ?`
		args = []interface{}{limit}
	} else {
		// Проверяем, является ли search датой в формате 02.01.2006
		if isDateFormat(search) {
			// Поиск по дате
			dateStr, err := convertDateFormat(search)
			if err != nil {
				return nil, fmt.Errorf("некорректный формат даты: %w", err)
			}
			query = `SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? ORDER BY date ASC LIMIT ?`
			args = []interface{}{dateStr, limit}
		} else {
			// Поиск по заголовку и комментарию
			searchPattern := "%" + search + "%"
			query = `SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date ASC LIMIT ?`
			args = []interface{}{searchPattern, searchPattern, limit}
		}
	}

	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении списка задач: %w", err)
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		var task Task
		var id int64

		err := rows.Scan(&id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, fmt.Errorf("ошибка при сканировании задачи: %w", err)
		}

		task.ID = strconv.FormatInt(id, 10)
		tasks = append(tasks, &task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при обработке результатов: %w", err)
	}

	// Если задач нет, возвращаем пустой слайс вместо nil
	if tasks == nil {
		tasks = []*Task{}
	}

	return tasks, nil
}

// isDateFormat проверяет, является ли строка датой в формате 02.01.2006
func isDateFormat(s string) bool {
	_, err := time.Parse("02.01.2006", s)
	return err == nil
}

// convertDateFormat преобразует дату из формата 02.01.2006 в 20060102
func convertDateFormat(dateStr string) (string, error) {
	t, err := time.Parse("02.01.2006", dateStr)
	if err != nil {
		return "", err
	}
	return t.Format("20060102"), nil
}

// GetTask возвращает задачу по указанному ID
func GetTask(id string) (*Task, error) {
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`

	var task Task
	var taskID int64

	err := DB.QueryRow(query, id).Scan(&taskID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return nil, fmt.Errorf("задача не найдена")
	}

	task.ID = strconv.FormatInt(taskID, 10)
	return &task, nil
}

// UpdateTask обновляет существующую задачу
func UpdateTask(task *Task) error {
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`

	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении задачи: %w", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка при проверке количества обновленных записей: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}

	return nil
}

// DeleteTask удаляет задачу по указанному ID
func DeleteTask(id string) error {
	query := `DELETE FROM scheduler WHERE id = ?`

	res, err := DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("ошибка при удалении задачи: %w", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка при проверке количества удаленных записей: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}

	return nil
}

// UpdateDate обновляет только дату задачи
func UpdateDate(next string, id string) error {
	query := `UPDATE scheduler SET date = ? WHERE id = ?`

	res, err := DB.Exec(query, next, id)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении даты задачи: %w", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка при проверке количества обновленных записей: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}

	return nil
}
