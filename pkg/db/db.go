package db

import (
	"database/sql"
	"fmt"
	"log"
	_ "modernc.org/sqlite"
	"os"
)

var DB *sql.DB

const schema string = `
CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(255) NOT NULL,
    comment TEXT,
    repeat VARCHAR(50)
);
`

func Init(dbFile string) error {
	// Проверяем существование файла базы данных
	_, err := os.Stat(dbFile)
	install := os.IsNotExist(err)
	if install {
		fmt.Printf("Файл %s не существует, создаем новую базу данных\n", dbFile)
	}

	// Открываем (или создаем) базу данных
	DB, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("ошибка открытия базы данных: %v", err)
	}

	// Проверяем соединение с базой данных
	if err := DB.Ping(); err != nil {
		return fmt.Errorf("ошибка подключения к базе данных: %v", err)
	}

	// Если файл не существовал, создаем таблицу и индекс
	if install {
		if _, err := DB.Exec(schema); err != nil {
			return fmt.Errorf("ошибка создания схемы: %v", err)
		}
		fmt.Printf("Таблица scheduler создана в базе данных: %s\n", dbFile)
	} else {
		fmt.Printf("База данных подключена: %s\n", dbFile)
	}

	return nil
}

func main() {
	// Инициализация базы данных
	dbFile := "scheduler.db"
	if err := Init(dbFile); err != nil {
		log.Fatalf("Ошибка инициализации базы данных: %v", err)
	}

	// Закрытие базы данных при завершении программы
	defer func() {
		if DB != nil {
			err := DB.Close()
			if err != nil {
				return
			}
			fmt.Println("База данных закрыта")
		}
	}()

	fmt.Println("Программа успешно запущена")
}
