package data

import (
	"context"
	"fmt"
	"log"
	"time"

	"example.com/Task-manager-Api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// // In-memory task list (seed data)
// var TaskList = []models.Tasks{
// 	{ID: "1", Title: "Finish Go tutorial"},
// 	{ID: "2", Title: "Write blog post"},
// 	{ID: "3", Title: "Read Gin documentation"},
// }

var Client *mongo.Client

func ConnectDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal("❌ Cannot connect to MongoDB:", err)
	}

	// Ping to check connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("❌ Cannot ping MongoDB:", err)
	}

	fmt.Println("✅ Connected to MongoDB!")
	Client = client
}

func GetAllTasks() []models.Tasks {
	collection := Client.Database("task_manager_db").Collection("tasks")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	defer cursor.Close(ctx)

	var TaskList []models.Tasks
	for cursor.Next(ctx) {
		var t models.Tasks
		if err := cursor.Decode(&t); err != nil {
			log.Fatal(err)
		}
		TaskList = append(TaskList, t)
	}

	return TaskList
}

func InsertTask(title string) (models.Tasks, error) {
	collection := Client.Database("task_manager_db").Collection("tasks")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create a new task struct
	task := models.Tasks{
		Title: title,
	}

	result, err := collection.InsertOne(ctx, task)
	if err != nil {
		return task, err
	}

	// Set the generated MongoDB ID
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		task.ID = oid.Hex()
	}

	return task, nil
}

func FindTaskByID(id string) (*models.Tasks, error) {
	collection := Client.Database("task_manager_db").Collection("tasks")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, fmt.Errorf("invalid task ID: %v", err)
	}
	var task models.Tasks
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("task not found")
		}
		return nil, err
	}
	return &task, nil
}

func UpdateTask(id string, newTitle string) (models.Tasks, error) {
	collection := Client.Database("task_manager_db").Collection("tasks")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Tasks{}, fmt.Errorf("invalid task ID: %v", err)
	}

	filter := bson.M{"_id": objID}
	update := bson.M{"$set": bson.M{"title": newTitle}}

	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return models.Tasks{}, err
	}

	// Return the updated task
	var updatedTask models.Tasks
	err = collection.FindOne(ctx, filter).Decode(&updatedTask)

	if err != nil {
		return models.Tasks{}, err
	}

	return updatedTask, nil
}

func DeleteTask(id string) error {
	collection := Client.Database("task_manager_db").Collection("tasks")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid task ID: %v", err)
	}

	result, err := collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}

	// To check if the document was actually deleted
	if result.DeletedCount == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}
