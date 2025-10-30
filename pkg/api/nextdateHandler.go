package api

import (
	"fmt"
	"net/http"
	"time"

	"final_project/pkg/nextdate"
)

// Константа для формата даты 20060102
const DateFormat = "20060102"

// NextDateHandler обрабатывает GET-запросы к /api/nextdate
// Принимает параметры:
//   - now: текущая дата в формате 20060102 (опционально, если не указана - используется текущая дата)
//   - date: исходная дата в формате 20060102
//   - repeat: правило повторения
//
// Возвращает следующую дату в формате 20060102 или текст ошибки
func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем, что это GET-запрос
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Получаем параметры из URL
	nowParam := r.FormValue("now")
	dateParam := r.FormValue("date")
	repeatParam := r.FormValue("repeat")

	// Проверяем обязательные параметры
	if dateParam == "" {
		http.Error(w, "Параметр date обязателен", http.StatusBadRequest)
		return
	}

	// Определяем время now
	var now time.Time
	if nowParam == "" {
		// Если параметр now не указан, используем текущую дату
		now = time.Now()
	} else {
		// Парсим переданную дату now
		parsedNow, err := time.Parse(DateFormat, nowParam)
		if err != nil {
			http.Error(w, fmt.Sprintf("Некорректный формат параметра now: %s", nowParam), http.StatusBadRequest)
			return
		}
		now = parsedNow
	}

	// Вызываем функцию NextDate из пакета nextdate
	nextDate, err := nextdate.NextDate(now, dateParam, repeatParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Возвращаем результат в формате 20060102
	// Если nextDate пустой, возвращаем ошибку
	if nextDate == "" {
		http.Error(w, "Некорректное правило повторения", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, nextDate)
}
