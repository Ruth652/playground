package router

import (
	"example.com/Task-manager-Api/controllers"
	"example.com/Task-manager-Api/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/tasks", middleware.AuthMiddleware(), controllers.GetAllTasks)
	r.GET("/tasks/:id", middleware.AuthMiddleware(), controllers.GetTaskByID)
	r.POST("/tasks", middleware.AuthMiddleware(), middleware.RoleMiddleWare("admin"), controllers.PostTask)
	r.PUT("/tasks/:id", middleware.AuthMiddleware(), middleware.RoleMiddleWare("admin"), controllers.PutTask)
	r.DELETE("/tasks/:id", middleware.AuthMiddleware(), middleware.RoleMiddleWare("admin"), controllers.DeleteTask)
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	return r
}
