package api

import (
	"encoding/json"
	"fmt"
	"github.com/vreditel85/go_final_project/pkg/db"
	"log"
	"net/http"
	"time"
)

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		jsonError(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

// addTaskHandler обрабатывает запрос на добавление задачи
func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			log.Printf("panic recovered in addTaskHandler: %v", rec)
			jsonError(w, "Internal server error", http.StatusInternalServerError)
		}
	}()

	// Проверяем метод запроса
	if r.Method != http.MethodPost {
		jsonError(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Декодируем JSON из тела запроса
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		jsonError(w, "Неверный формат JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Проверяем, что поле Title не пустое
	if task.Title == "" {
		jsonError(w, "Поле 'title' не может быть пустым", http.StatusBadRequest)
		return
	}

	// Проверяем и корректируем дату
	if err := checkDate(&task); err != nil {
		jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Добавляем задачу в базу данных
	id, err := db.AddTask(&task)
	if err != nil {
		jsonError(w, "Ошибка при добавлении задачи: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем идентификатор добавленной задачи в правильном формате
	response := map[string]interface{}{"id": id}
	writeJSON(w, response)
}

// getTaskHandler обрабатывает запрос на получение задачи по ID
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем параметр id из query string
	id := r.URL.Query().Get("id")
	if id == "" {
		jsonError(w, "Не указан идентификатор", http.StatusBadRequest)
		return
	}

	// Получаем задачу из базы данных
	task, err := db.GetTask(id)
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Возвращаем задачу в формате JSON
	writeJSON(w, task)
}

// updateTaskHandler обрабатывает запрос на обновление задачи
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			log.Printf("panic recovered in updateTaskHandler: %v", rec)
			jsonError(w, "Internal server error", http.StatusInternalServerError)
		}
	}()

	// Декодируем JSON из тела запроса
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		jsonError(w, "Неверный формат JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Проверяем, что поле ID не пустое
	if task.ID == 0 {
		jsonError(w, "Поле 'id' не может быть пустым", http.StatusBadRequest)
		return
	}

	// Проверяем, что поле Title не пустое
	if task.Title == "" {
		jsonError(w, "Поле 'title' не может быть пустым", http.StatusBadRequest)
		return
	}

	// Проверяем и корректируем дату
	if err := checkDate(&task); err != nil {
		jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Обновляем задачу в базе данных
	err := db.UpdateTask(&task)
	if err != nil {
		jsonError(w, "Ошибка при обновлении задачи: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем пустой JSON объект
	writeJSON(w, map[string]interface{}{})
}

// deleteTaskHandler обрабатывает запрос на удаление задачи
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			log.Printf("panic recovered in deleteTaskHandler: %v", rec)
			jsonError(w, "Internal server error", http.StatusInternalServerError)
		}
	}()

	// Получаем параметр id из query string
	id := r.URL.Query().Get("id")
	if id == "" {
		jsonError(w, "Не указан идентификатор", http.StatusBadRequest)
		return
	}

	// Удаляем задачу из базы данных
	err := db.DeleteTask(id)
	if err != nil {
		jsonError(w, "Ошибка при удалении задачи: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем пустой JSON объект
	writeJSON(w, map[string]interface{}{})
}

// checkDate проверяет и корректирует дату задачи
func checkDate(task *db.Task) error {
	if task == nil {
		return fmt.Errorf("task is nil")
	}
	now := time.Now()

	// Если дата пустая, устанавливаем текущую дату
	if task.Date == "" {
		task.Date = now.Format("20060102")
	}

	// Проверяем корректность формата даты
	t, err := time.Parse("20060102", task.Date)
	if err != nil {
		return fmt.Errorf("некорректный формат даты: %s", task.Date)
	}

	// Если определено правило повторения
	if task.Repeat != "" {
		// Проверяем корректность правила и получаем следующую дату
		next, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return fmt.Errorf("некорректное правило повторения: %v", err)
		}
		// Обновляем дату на следующую
		task.Date = next
	} else {
		// Если правила повторения нет, проверяем что дата не в прошлом
		if !afterNow(t, now) {
			return fmt.Errorf("дата не может быть в прошлом")
		}
	}

	return nil
}

// afterNow проверяет, что дата t после now (включая сегодняшний день)
func afterNow(t time.Time, now time.Time) bool {
	// Приводим к началу дня для корректного сравнения дат без времени
	t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	return !t.Before(now)
}

// writeJSON записывает данные в ответ в формате JSON
func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err := json.NewEncoder(w).Encode(data); err != nil {
		jsonError(w, "Ошибка при кодировании JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// jsonError возвращает ошибку в формате JSON
func jsonError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	response := map[string]string{"error": message}
	json.NewEncoder(w).Encode(response)
}
