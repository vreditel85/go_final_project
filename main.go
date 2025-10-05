package main

import (
	"log"
	"net/http"

	"github.com/vreditel85/go_final_project/pkg/api"
	"github.com/vreditel85/go_final_project/pkg/db"
)

func main() {
	// Инициализация базы данных
	dbFile := "scheduler.db"
	if err := db.Init(dbFile); err != nil {
		log.Fatalf("Ошибка инициализации базы данных: %v", err)
	}
	//defer db.Close()

	// Файловый сервер для статических файлов из папки web
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./web"))))

	// Инициализация API маршрутов
	api.Init()

	// Запуск сервера
	log.Println("Сервер запущен на http://localhost:7540")
	err := http.ListenAndServe(":7540", nil)
	if err != nil {
		log.Fatal("Ошибка запуска сервера: ", err)
	}
}
