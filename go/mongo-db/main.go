package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Connect to MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	// Check connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("❌ Cannot connect to MongoDB:", err)
	}
	fmt.Println("✅ Connected to MongoDB!")

	// Choose your database and collection
	collection := client.Database("practice_db").Collection("students")

	fmt.Println("Using database:", collection.Database().Name())

	// Insert student
	// student := map[string]interface{}{
	// 	"name":  "John Doe",
	// 	"age":   22,
	// 	"grade": "A",
	// }

	// result, err := collection.InsertOne(ctx, student)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println("Insterted ID:", result.InsertedID)

	// Find all documents
	cursor, err := collection.Find(ctx, map[string]interface{}{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var doc map[string]interface{}
		cursor.Decode(&doc)
		fmt.Println(doc)
	}

	// update a Document
	// 	filter := map[string]interface{}{"name": "John Doe"}
	// 	update := map[string]interface{}{
	// 		"$set": map[string]interface{}{"grade": "A+"},
	// 	}

	// 	result, err := collection.UpdateOne(ctx, filter, update)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	fmt.Println("Modified count:", result.ModifiedCount)
	//

	filter := map[string]interface{}{"name": "John Doe"}
	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Deleted count:", result.DeletedCount)

}
