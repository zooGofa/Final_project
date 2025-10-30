package api

import (
	"encoding/json"
	"net/http"
)

// writeJson сериализует данные в JSON и отправляет HTTP-ответ
func writeJson(w http.ResponseWriter, data any) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	return json.NewEncoder(w).Encode(data)
}
