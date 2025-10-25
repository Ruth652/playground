// main.go
package main

import (
	"example.com/Task-manager-Api/data"
	"example.com/Task-manager-Api/router"
)

func main() {

	data.ConnectDB()

	r := router.SetupRouter()
	r.Run("localhost:8080")
}
