package db

import (
	"os"

	"github.com/jmoiron/sqlx"

	_ "modernc.org/sqlite"
)

// Schema содержит команды DDL для таблицы scheduler и индекса по date.
const Schema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(256) NOT NULL DEFAULT "",
    comment TEXT NOT NULL DEFAULT "",
    repeat VARCHAR(128) NOT NULL DEFAULT ""
);
CREATE INDEX IF NOT EXISTS scheduler_date ON scheduler(date);
`

// DefaultDbFile — путь по умолчанию к файлу базы данных.
var DefaultDbFile = "scheduler.db"

// DateString — формат представления даты (YYYYMMDD).
var DateString = "20060102"

// DB — глобальный обработчик подключения к БД.
var DB *sqlx.DB

// getDbFile возвращает актуальный путь к файлу БД с учётом переменной окружения.
// Читает переменную окружения TODO_DBFILE каждый раз при вызове
func getDbFile() string {
	if envDbFile := os.Getenv("TODO_DBFILE"); len(envDbFile) > 0 {
		return envDbFile
	}
	return DefaultDbFile
}

// Init открывает БД и при необходимости создаёт таблицу и индекс.
func Init() error {
	dbFile := getDbFile()
	_, err := os.Stat(dbFile)
	install := err != nil

	conn, err := sqlx.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	if install {
		if _, err := conn.Exec(Schema); err != nil {
			_ = conn.Close()
			return err
		}
	}

	DB = conn
	return nil
}

// Close закрывает соединение с БД.
func Close() { _ = DB.Close() }
