package api

import (
	"encoding/json"
	"net/http"
	"strings"
)

// writeJson сериализует данные в JSON и отправляет HTTP-ответ с заданным статусом
func writeJson(w http.ResponseWriter, data any, status int) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

// errorStatus сопоставляет ошибку с HTTP статусом
// если в тексте ошибки есть "не найдена" — возвращаем 404, иначе 500
func errorStatus(err error) int {
	if err == nil {
		return http.StatusOK
	}
	if strings.Contains(err.Error(), "не найдена") {
		return http.StatusNotFound
	}
	return http.StatusInternalServerError
}
