package data

import "example.com/Task-manager-Api/models"

// In-memory task list (seed data)
var TaskList = []models.Tasks{
	{ID: "1", Title: "Finish Go tutorial"},
	{ID: "2", Title: "Write blog post"},
	{ID: "3", Title: "Read Gin documentation"},
}

func FindTaskByID(id string) (*models.Tasks, int) {
	for i, t := range TaskList {
		if t.ID == id {
			return &TaskList[i], i
		}
	}
	return nil, -1
}
