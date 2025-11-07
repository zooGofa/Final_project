package api

import (
	"encoding/json"
	"net/http"

	"final_project/pkg/db"
)

// updateTaskHandler обрабатывает PUT-запросы для обновления задачи
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем, что это PUT-запрос
	if r.Method != http.MethodPut {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

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

	// Проверяем обязательное поле id
	if task.ID == "" {
		writeJson(w, map[string]string{"error": "Не указан идентификатор задачи"}, http.StatusBadRequest)
		return
	}

	// Проверяем и корректируем дату
	if err := checkDate(&task); err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	// Обновляем задачу в базе данных
	if err := db.UpdateTask(&task); err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, errorStatus(err))
		return
	}

	// Возвращаем пустой JSON при успешном обновлении
	writeJson(w, map[string]interface{}{}, http.StatusOK)
}
