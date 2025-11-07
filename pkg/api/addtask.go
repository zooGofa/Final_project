package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"final_project/pkg/db"
)

// addTaskHandler обрабатывает POST-запросы для добавления задач
func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	// Десериализуем JSON
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJson(w, map[string]string{"error": "ошибка десериализации JSON"}, http.StatusBadRequest)
		return
	}

	// Проверяем обязательное поле title
	if task.Title == "" {
		writeJson(w, map[string]string{"error": "Не указан заголовок задачи"}, http.StatusBadRequest)
		return
	}

	// Проверяем и корректируем дату
	if err := checkDate(&task); err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	// Добавляем задачу в базу данных
	id, err := db.AddTask(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, errorStatus(err))
		return
	}

	// Возвращаем ID созданной задачи
	writeJson(w, map[string]string{"id": strconv.FormatInt(id, 10)}, http.StatusCreated)
}
