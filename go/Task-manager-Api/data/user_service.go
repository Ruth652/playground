package data

import (
	"context"
	"fmt"
	"time"

	"example.com/Task-manager-Api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func UserNameExist(username string) (bool, error) {
	collection := Client.Database("task_manager_db").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()
	var user models.User
	err := collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, fmt.Errorf("DB error: %v", err)
	}
	return true, nil // Return true if user exists
}
func InsertUser(user models.User) (models.User, error) {
	collection := Client.Database("task_manager_db").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if user.Role == "" {
		user.Role = "User"
	}

	_, err := collection.InsertOne(ctx, user)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func UserExists(user models.User) (models.User, error) {
	collection := Client.Database("task_manager_db").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var foundUser models.User
	err := collection.FindOne(ctx, bson.M{"username": user.UserName}).Decode(&foundUser)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.User{}, fmt.Errorf("user not found")
		}
		return models.User{}, fmt.Errorf("database error: %v", err)
	}

	// Check password
	if bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(user.Password)) != nil {
		return models.User{}, fmt.Errorf("invalid password")
	}

	return foundUser, nil
}
