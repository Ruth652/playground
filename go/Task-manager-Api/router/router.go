package router

import (
	"example.com/Task-manager-Api/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/tasks", controllers.GetTasks)
	r.GET("/tasks/:id", controllers.GetTaskByID)
	r.POST("/tasks", controllers.PostTask)
	r.PUT("/tasks/:id", controllers.PutTask)
	r.DELETE("/tasks/:id", controllers.DeleteTask)

	return r
}
