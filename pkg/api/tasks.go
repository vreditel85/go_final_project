package api

import (
	"encoding/json"
	"github.com/vreditel85/go_final_project/pkg/db"
	"net/http"
)

// TasksResp структура для ответа с задачами
type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

// tasksHandler обрабатывает запрос на получение списка задач
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод запроса
	if r.Method != http.MethodGet {
		jsonError(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Получаем задачи из базы данных
	tasks, err := db.Tasks(50) // ограничиваем 50 задачами
	if err != nil {
		jsonError(w, "Ошибка при получении задач: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем задачи в формате JSON
	writeJSON(w, TasksResp{
		Tasks: tasks,
	})
}

// jsonError возвращает ошибку в формате JSON
func jsonError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	response := map[string]string{"error": message}
	json.NewEncoder(w).Encode(response)
}
