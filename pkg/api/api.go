package api

import "net/http"

// Init регистрирует все обработчики HTTP
func Init() {
	http.HandleFunc("/api/nextdate", CORSMiddleware(nextDateHandler))
	http.HandleFunc("/api/task", CORSMiddleware(taskHandler))
	http.HandleFunc("/api/tasks", CORSMiddleware(tasksHandler))
	http.HandleFunc("/api/task/done", CORSMiddleware(taskDoneHandler))
}
