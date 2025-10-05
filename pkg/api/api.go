package api

import (
	"net/http"
)

// Init регистрирует все обработчики HTTP
func Init() {
	http.HandleFunc("/api/nextdate", nextDateHandler)
	http.HandleFunc("/api/task", taskHandler)
	http.HandleFunc("/api/tasks", tasksHandler)
	http.HandleFunc("/api/task/done", taskDoneHandler) // новый обработчик для завершения задач
}
