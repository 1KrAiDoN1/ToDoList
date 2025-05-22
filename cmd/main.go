package main

import (
	todo "helloapp"
	"helloapp/pkg/handler"
	"helloapp/pkg/repository"
	"helloapp/pkg/service"
	"log"
)

func main() {

	repos := repository.NewRepository()
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	var srv todo.Server
	err := srv.Run("8080", handlers.InitRoutes())
	if err != nil {
		log.Fatal("Ошибка при запуске сервера: ", err.Error())
	}
}
