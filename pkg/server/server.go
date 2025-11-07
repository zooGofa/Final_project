package server

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"final_project/pkg/api"
)

// Константы для настройки порта сервера
const (
	defaultPort = 7540        // порт по умолчанию, используется когда TODO_PORT не задан или невалиден
	envPortKey  = "TODO_PORT" // имя переменной окружения для переопределения порта
)

// resolvePort определяет порт для запуска сервера
// Сначала проверяет переменную окружения TODO_PORT
// Если переменная задана и содержит валидный номер порта (1-65535), использует её
// Иначе возвращает порт по умолчанию (7540)
func resolvePort() int {
	// Получаем значение переменной окружения TODO_PORT
	if v, ok := os.LookupEnv(envPortKey); ok {
		// Пытаемся преобразовать строку в число и проверить диапазон портов
		if p, err := strconv.Atoi(v); err == nil && p > 0 && p < 65536 {
			return p
		}
	}
	// Возвращаем порт по умолчанию, если переменная не задана или невалидна
	return defaultPort
}

// addr формирует строку адреса для сервера в формате ":порт"
func addr() string { return ":" + strconv.Itoa(resolvePort()) }

// New создает и настраивает HTTP сервер для обслуживания статических файлов
// Принимает путь к директории с веб-файлами (webDir)
// Возвращает настроенный http.Server
func New(webDir string) *http.Server {
	// Инициализируем API обработчики
	api.Init()

	// Создаем файловый сервер для обслуживания статических файлов из webDir
	fileServer := http.FileServer(http.Dir(webDir))

	// Регистрируем обработчик для всех остальных запросов к корню "/"
	// Все запросы будут направляться к файловому серверу
	http.Handle("/", fileServer)

	// Возвращаем настроенный сервер с адресом и обработчиком логирования
	return &http.Server{Addr: addr(), Handler: logRequests(http.DefaultServeMux)}
}

// Run запускает HTTP сервер для обслуживания файлов из webDir
// Выводит сообщение о запуске и начинает прослушивание порта
// Возвращает ошибку, если сервер не может быть запущен
func Run(webDir string) error {
	// Создаем новый сервер
	s := New(webDir)

	// Выводим сообщение о том, на каком адресе запущен сервер
	log.Printf("listening on http://localhost%s", s.Addr)

	// Запускаем сервер и начинаем прослушивание входящих соединений
	return s.ListenAndServe()
}

// logRequests создает middleware для логирования всех HTTP запросов
// Обертывает следующий обработчик и выводит в лог метод и путь каждого запроса
func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Логируем метод HTTP (GET, POST, etc.) и путь запроса
		log.Printf("%s %s", r.Method, r.URL.Path)

		// Передаем управление следующему обработчику в цепочке
		next.ServeHTTP(w, r)
	})
}
