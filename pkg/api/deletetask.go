package api

import (
	"net/http"

	"final_project/pkg/db"
)

// deleteTaskHandler обрабатывает DELETE-запросы для удаления задачи
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем, что это DELETE-запрос
	if r.Method != http.MethodDelete {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Получаем параметр id из URL
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	// Удаляем задачу из базы данных
	err := db.DeleteTask(id)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	// Возвращаем пустой JSON при успешном удалении
	writeJson(w, map[string]interface{}{})
}
