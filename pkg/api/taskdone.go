package api

import (
	"net/http"
	"time"

	"final_project/pkg/db"
	"final_project/pkg/nextdate"
)

// taskDoneHandler обрабатывает POST-запросы для отметки задачи как выполненной
func taskDoneHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем, что это POST-запрос
	if r.Method != http.MethodPost {
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

	// Если правило повторения отсутствует, удаляем задачу
	if task.Repeat == "" {
		err = db.DeleteTask(id)
		if err != nil {
			writeJson(w, map[string]string{"error": err.Error()})
			return
		}
	} else {
		// Если задача периодическая, вычисляем следующую дату
		now := time.Now()
		nextDate, err := nextdate.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeJson(w, map[string]string{"error": err.Error()})
			return
		}

		// Обновляем дату задачи
		err = db.UpdateDate(nextDate, id)
		if err != nil {
			writeJson(w, map[string]string{"error": err.Error()})
			return
		}
	}

	// Возвращаем пустой JSON при успешном выполнении
	writeJson(w, map[string]interface{}{})
}
