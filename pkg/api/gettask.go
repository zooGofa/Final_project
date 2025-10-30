package api

import (
	"net/http"

	"final_project/pkg/db"
)

// getTaskHandler обрабатывает GET-запросы для получения задачи по ID
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем, что это GET-запрос
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Получаем параметр id из URL
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	// Получаем задачу из базы данных
	task, err := db.GetTask(id)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	// Возвращаем задачу в JSON формате
	writeJson(w, task)
}
