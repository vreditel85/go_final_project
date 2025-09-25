package main

import (
	"github.com/vreditel85/go_final_project/pkg/api"
	"log"
	"net/http"
)

func main() {
	//Файловый сервер из web
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./web"))))
	api.Init()
	// Запуск сервера
	log.Println("Сервер запущен на http://localhost:7540")
	err := http.ListenAndServe(":7540", nil)
	if err != nil {
		log.Fatal("Ошибка запуска сервера: ", err)
	}
}
