package main

import (
	"github.com/Ruth652/playground/go/library_management/concurrency"
	"github.com/Ruth652/playground/go/library_management/controllers"
	"github.com/Ruth652/playground/go/library_management/services"
)

func main() {

	lib := services.NewLibrary()

	concurrency.StartReservationWorker(lib)

	controller := controllers.NewController(lib)
	controller.Run()
}
