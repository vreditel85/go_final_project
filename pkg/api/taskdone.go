package api

import (
	"github.com/vreditel85/go_final_project/pkg/db"
	"log"
	"net/http"
	"time"
)

// taskDoneHandler обрабатывает запрос на завершение задачи
func taskDoneHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			log.Printf("panic recovered in taskDoneHandler: %v", rec)
			jsonError(w, "Internal server error", http.StatusInternalServerError)
		}
	}()

	// Проверяем метод запроса
	if r.Method != http.MethodPost {
		jsonError(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

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

	// Если задача не имеет правила повторения - удаляем её
	if task.Repeat == "" {
		err = db.DeleteTask(id)
		if err != nil {
			jsonError(w, "Ошибка при удалении задачи: "+err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		// Если задача периодическая - вычисляем следующую дату
		// Используем текущую дату как now для вычисления следующей даты
		now := time.Now()
		nextDate, err := NextDate(now, task.Date, task.Repeat)

		if err != nil {
			jsonError(w, "Ошибка при вычислении следующей даты: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Обновляем дату задачи
		err = db.UpdateDate(id, nextDate)
		if err != nil {
			jsonError(w, "Ошибка при обновлении даты задачи: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Возвращаем пустой JSON объект
	writeJSON(w, map[string]interface{}{})
}
