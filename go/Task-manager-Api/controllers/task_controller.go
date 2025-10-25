package controllers

import (
	"net/http"

	"example.com/Task-manager-Api/data"
	"github.com/gin-gonic/gin"
)

// Get all tasks
func GetTasks(c *gin.Context) {
	tasks := data.GetAllTasks()
	c.IndentedJSON(http.StatusOK, tasks)
}

// Get task by ID
func GetTaskByID(c *gin.Context) {
	id := c.Param("id")
	task, err := data.FindTaskByID(id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "task not found"})
		return
	}
	c.IndentedJSON(http.StatusOK, task)
}

// Create a new task
func PostTask(c *gin.Context) {
	var input struct {
		Title string `json:"title"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if input.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
		return
	}
	createdTask, err := data.InsertTask(input.Title)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusCreated, createdTask)
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
	task, err := data.UpdateTask(id, editTask.Title)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, task)
}

// Delete task
func DeleteTask(c *gin.Context) {
	id := c.Param("id")
	err := data.DeleteTask(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "task deleted successfully"})
}
