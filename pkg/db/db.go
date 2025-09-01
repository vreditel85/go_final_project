package main

import (
	"database/sql"
	"fmt"
	"log"
	_ "modernc.org/sqlite"
	"os"
)

var db *sql.DB

const schema = `
CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(255) NOT NULL,
    comment TEXT,
    repeat VARCHAR(50)
);
CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler(date);
`

func Init(dbFile string) error {
	// Проверяем существование файла базы данных
	_, err := os.Stat(dbFile)
	install := os.IsNotExist(err)

	// Открываем базу данных
	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("ошибка открытия базы данных: %v", err)
	}

	// Проверяем соединение с базой данных
	if err := db.Ping(); err != nil {
		return fmt.Errorf("ошибка подключения к базе данных: %v", err)
	}

	// Если файл не существовал, создаем таблицу и индекс
	if install {
		if _, err := db.Exec(schema); err != nil {
			return fmt.Errorf("ошибка создания схемы: %v", err)
		}
		fmt.Printf("База данных создана и инициализирована: %s\n", dbFile)
	} else {
		fmt.Printf("База данных подключена: %s\n", dbFile)
	}

	return nil
}
func main() {
	// Инициализация базы данных в начале main()
	if err := Init("scheduler.db"); err != nil {
		log.Fatalf("Ошибка инициализации базы данных: %v", err)
	}

	// Закрытие базы данных при завершении программы
	defer func() {
		if db != nil {
			_ = db.Close()
			fmt.Println("База данных закрыта")
		}
	}()

	fmt.Println("Программа успешно запущена")
}
