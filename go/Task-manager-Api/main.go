// main.go
package main

import (
	"example.com/Task-manager-Api/router"
)

func main() {
	r := router.SetupRouter()
	r.Run("localhost:8080")
}
