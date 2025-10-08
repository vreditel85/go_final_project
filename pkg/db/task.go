package db

import (
	"database/sql"
	"fmt"
	"strconv"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// AddTask добавляет задачу в таблицу scheduler и возвращает ID добавленной записи
func AddTask(task *Task) (string, error) {
	if task == nil {
		return "", fmt.Errorf("task cannot be nil")
	}
	// Проверяем обязательные поля
	if task.Date == "" {
		return "", fmt.Errorf("date is required")
	}
	if task.Title == "" {
		return "", fmt.Errorf("title is required")
	}
	// Проверяем, что база данных инициализирована
	if DB == nil {
		return "", fmt.Errorf("database not initialized")
	}
	// SQL запрос для вставки задачи
	query := `
        INSERT INTO scheduler (date, title, comment, repeat) 
        VALUES (?, ?, ?, ?)
    `

	// Выполняем запрос
	result, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return "", fmt.Errorf("error adding task: %v", err)
	}

	// Получаем ID добавленной записи
	id, err := result.LastInsertId()
	if err != nil {
		return "", fmt.Errorf("error getting ID: %v", err)
	}

	return strconv.FormatInt(id, 10), nil
}

// Tasks возвращает список ближайших задач, отсортированных по дате
func Tasks(limit int) ([]*Task, error) {
	// Проверяем, что база данных инициализирована
	if DB == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	// SQL запрос для получения задач
	query := `
        SELECT id, date, title, comment, repeat 
        FROM scheduler 
        WHERE date >= date('now')
        ORDER BY date ASC, id ASC
    `

	// Если указан лимит, добавляем его в запрос
	if limit > 0 {
		query += " LIMIT ?"
	}

	var rows *sql.Rows
	var err error

	if limit > 0 {
		rows, err = DB.Query(query, limit)
	} else {
		rows, err = DB.Query(query)
	}

	if err != nil {
		return nil, fmt.Errorf("error getting tasks: %v", err)
	}
	defer rows.Close()

	// Считываем результаты
	var tasks []*Task
	for rows.Next() {
		var task Task
		var id int
		err := rows.Scan(&id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, fmt.Errorf("error scanning task: %v", err)
		}
		task.ID = strconv.Itoa(id)
		tasks = append(tasks, &task)
	}

	// Проверяем ошибки итерации
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tasks: %v", err)
	}

	// Если задач нет, возвращаем пустой slice вместо nil
	if tasks == nil {
		tasks = []*Task{}
	}

	return tasks, nil
}

// GetTask возвращает задачу по идентификатору
func GetTask(id string) (*Task, error) {
	// Проверяем, что база данных инициализирована
	if DB == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	// Проверяем, что идентификатор не пустой
	if id == "" {
		return nil, fmt.Errorf("id is required")
	}

	// Конвертируем строковый ID в число
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("invalid task ID")
	}

	var task Task
	var dbID int
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`

	err = DB.QueryRow(query, idInt).Scan(&dbID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task not found")
		}
		return nil, fmt.Errorf("error getting task: %v", err)
	}
	task.ID = strconv.Itoa(dbID)

	return &task, nil
}

// UpdateTask обновляет задачу
func UpdateTask(task *Task) error {
	// Проверяем, что база данных инициализирована
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	if task == nil {
		return fmt.Errorf("task cannot be nil")
	}

	// Проверяем обязательные поля
	if task.ID == "" {
		return fmt.Errorf("id is required")
	}
	if task.Date == "" {
		return fmt.Errorf("date is required")
	}
	if task.Title == "" {
		return fmt.Errorf("title is required")
	}

	// Конвертируем строковый ID в число
	id, err := strconv.Atoi(task.ID)
	if err != nil {
		return fmt.Errorf("invalid task ID")
	}

	// SQL запрос для обновления задачи
	query := `
        UPDATE scheduler 
        SET date = ?, title = ?, comment = ?, repeat = ?
        WHERE id = ?
    `

	// Выполняем запрос
	result, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, id)
	if err != nil {
		return fmt.Errorf("error updating task: %v", err)
	}

	// Проверяем, что запись была обновлена
	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking update: %v", err)
	}
	if count == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}

// UpdateDate обновляет дату задачи
func UpdateDate(id string, date string) error {
	// Проверяем, что база данных инициализирована
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	// Конвертируем строковый ID в число
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("invalid task ID")
	}

	// SQL запрос для обновления даты задачи
	query := `UPDATE scheduler SET date = ? WHERE id = ?`

	// Выполняем запрос
	result, err := DB.Exec(query, date, idInt)
	if err != nil {
		return fmt.Errorf("error updating date: %v", err)
	}

	// Проверяем, что запись была обновлена
	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking update: %v", err)
	}
	if count == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}

// DeleteTask удаляет задачу по идентификатору
func DeleteTask(id string) error {
	// Проверяем, что база данных инициализирована
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	// Проверяем, что идентификатор не пустой
	if id == "" {
		return fmt.Errorf("id is required")
	}

	// Конвертируем строковый ID в число
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("invalid task ID")
	}

	// SQL запрос для удаления задачи
	query := `DELETE FROM scheduler WHERE id = ?`

	// Выполняем запрос
	result, err := DB.Exec(query, idInt)
	if err != nil {
		return fmt.Errorf("error deleting task: %v", err)
	}

	// Проверяем, что запись была удалена
	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking deletion: %v", err)
	}
	if count == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}
