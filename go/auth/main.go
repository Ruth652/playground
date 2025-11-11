package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type Login struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

var users = make(map[string]*User)
var jwtSecret = []byte("my_secret_key")

func main() {
	router := gin.Default()
	password := "1234"
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	users["ruthfian74@gmail.com"] = &User{
		ID:       1,
		Email:    "ruthfian74@gmail.com",
		Password: string(hash),
	}

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to the Auth Service",
		})
	})

	router.POST("/register", func(c *gin.Context) {
		var user User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request payload"})
			return
		}
		// Todo

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(500, gin.H{"error": "internal server error"})
			return
		}

		user.Password = string(hashedPassword)
		users[user.Email] = &user
		c.JSON(201, gin.H{"message": "user registered successfully"})

	})
	router.POST("/login", func(c *gin.Context) {
		var loginData Login
		if err := c.ShouldBindJSON(&loginData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}
		user, exists := users[loginData.Email]
		if !exists || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password)) != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or Password"})
			return
		}

		// Generate JWT token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":    user.ID,
			"email": user.Email,
			"exp":   time.Now().Add(time.Hour * 24).Unix(), // expires in 24h
		})

		jwtToken, err := token.SignedString(jwtSecret)

		if err != nil {
			c.JSON(500, gin.H{"error": "Internal server error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": jwtToken})

	})
	router.GET("/secure", AuthMiddleware(), func(c *gin.Context) {
		userID, _ := c.Get("userID")
		userEmail, _ := c.Get("userEmail")

		c.JSON(http.StatusOK, gin.H{
			"message":   "This is a secure route",
			"userID":    userID,
			"userEmail": userEmail,
		})
	})

	router.Run(":8080")
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get the AUthorization header

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// split "Bearer <token>"
		authParts := strings.Split(authHeader, " ")
		if len(authParts) != 2 || strings.ToLower(authParts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header"})
			c.Abort()
			return
		}

		// parse and validate JWT
		token, err := jwt.Parse(authParts[1], func(token *jwt.Token) (interface{}, error) {
			// Ensure signing method is HMAC
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid JWT"})
			c.Abort()
			return
		}

		// set user info to context
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("userID", claims["id"])
			c.Set("userEmail", claims["email"])

		}

		c.Next()
	}
}
