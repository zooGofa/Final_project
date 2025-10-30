package api

import (
	"net/http"

	"final_project/pkg/db"
)

// TasksResp представляет ответ API для получения списка задач
type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

// tasksHandler обрабатывает GET-запросы для получения списка задач
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем, что это GET-запрос
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Получаем параметр поиска из URL
	search := r.URL.Query().Get("search")

	// Получаем список задач из базы данных (максимум 50)
	tasks, err := db.Tasks(50, search)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	// Возвращаем список задач в JSON формате
	writeJson(w, TasksResp{
		Tasks: tasks,
	})
}
