package controllers

import (
	"net/http"

	"example.com/Task-manager-Api/data"
	"example.com/Task-manager-Api/models"
	"github.com/gin-gonic/gin"
)

// Get all tasks
func GetTasks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, data.TaskList)
}

// Get task by ID
func GetTaskByID(c *gin.Context) {
	id := c.Param("id")
	task, _ := data.FindTaskByID(id)
	if task != nil {
		c.IndentedJSON(http.StatusOK, task)
		return
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "task not found"})
}

// Create a new task
func PostTask(c *gin.Context) {
	var newTask models.Tasks
	if err := c.BindJSON(&newTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if newTask.ID == "" || newTask.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id and title are required"})
		return
	}
	data.TaskList = append(data.TaskList, newTask)
	c.IndentedJSON(http.StatusCreated, newTask)
}

// Update task
func PutTask(c *gin.Context) {
	id := c.Param("id")
	var editTask struct {
		Title string `json:"title"`
	}
	if err := c.BindJSON(&editTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if editTask.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title is required"})
		return
	}
	task, index := data.FindTaskByID(id)
	if task != nil {
		data.TaskList[index].Title = editTask.Title
		c.IndentedJSON(http.StatusOK, data.TaskList[index])
		return
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "task not found"})
}

// Delete task
func DeleteTask(c *gin.Context) {
	id := c.Param("id")
	task, index := data.FindTaskByID(id)
	if task != nil {
		data.TaskList = append(data.TaskList[:index], data.TaskList[index+1:]...)
		c.IndentedJSON(http.StatusOK, gin.H{"message": "task deleted"})
		return
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "task not found"})
}
