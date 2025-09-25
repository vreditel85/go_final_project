package db

import (
	"fmt"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// AddTask добавляет задачу в таблицу scheduler и возвращает ID добавленной записи
func AddTask(task *Task) (int64, error) {
	// Проверяем обязательные поля
	if task.Date == "" {
		return 0, fmt.Errorf("дата не может быть пустой")
	}
	if task.Title == "" {
		return 0, fmt.Errorf("заголовок не может быть пустым")
	}

	// SQL запрос для вставки задачи
	query := `
        INSERT INTO scheduler (date, title, comment, repeat) 
        VALUES (?, ?, ?, ?)
    `

	// Выполняем запрос
	result, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, fmt.Errorf("ошибка при добавлении задачи: %v", err)
	}

	// Получаем ID добавленной записи
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("ошибка при получении ID: %v", err)
	}

	return id, nil
}
