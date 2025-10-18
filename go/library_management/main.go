package main

import (
	"fmt"

	"github.com/Ruth652/playground/go/library_management/controllers"
	"github.com/Ruth652/playground/go/library_management/services"
)

func main() {
	lib := services.NewLibrary()

	lib.AddBook(services.BookFromArgs(1, "The Go Programming Language", "Alan Donovan"))
	lib.AddBook(services.BookFromArgs(2, "Introduction to Algorithms", "CLRS"))
	lib.AddMember(services.MemberFromArgs(1, "Alice"))
	lib.AddMember(services.MemberFromArgs(2, "Bob"))

	ctrl := controllers.NewController(lib)
	fmt.Println("Welcome to the Library Management System")
	ctrl.Run()
}
