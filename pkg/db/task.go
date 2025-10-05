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
	if task == nil {
		return 0, fmt.Errorf("task cannot be nil")
	}
	// Проверяем обязательные поля
	if task.Date == "" {
		return 0, fmt.Errorf("дата не может быть пустой")
	}
	if task.Title == "" {
		return 0, fmt.Errorf("заголовок не может быть пустым")
	}
	// Проверяем, что база данных инициализирована
	if DB == nil {
		return 0, fmt.Errorf("база данных не инициализирована")
	}
	// SQL запрос для вставки задачи
	query := `
        INSERT INTO scheduler (date, title, comment, repeat) 
        VALUES (?, ?, ?, ?)
    `

	// Выполняем запрос
	result, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
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

// Tasks возвращает список ближайших задач, отсортированных по дате
func Tasks(limit int) ([]*Task, error) {
	// Проверяем, что база данных инициализирована
	if DB == nil {
		return nil, fmt.Errorf("база данных не инициализирована")
	}

	// SQL запрос для получения задач
	query := `
        SELECT id, date, title, comment, repeat 
        FROM scheduler 
        WHERE date >= date('now', 'localtime')
        ORDER BY date ASC, id ASC
        LIMIT ?
    `

	// Выполняем запрос
	rows, err := DB.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении задач: %v", err)
	}
	defer rows.Close()

	// Считываем результаты
	var tasks []*Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, fmt.Errorf("ошибка при сканировании задачи: %v", err)
		}
		tasks = append(tasks, &task)
	}

	// Проверяем ошибки итерации
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при итерации по задачам: %v", err)
	}

	// Если задач нет, возвращаем пустой slice вместо nil
	if tasks == nil {
		tasks = []*Task{}
	}

	return tasks, nil
}
