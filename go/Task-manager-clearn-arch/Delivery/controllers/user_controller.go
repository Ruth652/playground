package controllers

import (
	"net/http"
	"task-manager-clean-arch/domain"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserController struct {
	userUsecase domain.UserUsecase
}

func NewUserController(uu domain.UserUsecase) *UserController {
	return &UserController{
		userUsecase: uu,
	}
}

// POST /register
func (uc *UserController) Register(c *gin.Context) {
	var req domain.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	user := &domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	access, refresh, err := uc.userUsecase.Register(c.Request.Context(), user)
	if err != nil {
		status := http.StatusInternalServerError
		if err == domain.ErrConflict {
			status = http.StatusConflict
		}
		c.JSON(status, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, domain.RegisterResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	})
}

// POST /login
func (uc *UserController) Login(c *gin.Context) {
	var req domain.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	access, refresh, err := uc.userUsecase.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		status := http.StatusUnauthorized
		if err == domain.ErrNotFound {
			status = http.StatusUnauthorized
		}
		c.JSON(status, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.LoginResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	})
}

// POST /promote/:id
func (uc *UserController) Promote(c *gin.Context) {
	targetIDHex := c.Param("id")
	targetID, err := primitive.ObjectIDFromHex(targetIDHex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid user ID"})
		return
	}

	actorIDHex, exists := c.Get("x-user-id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	actorID, err := primitive.ObjectIDFromHex(actorIDHex.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid actor ID"})
		return
	}

	err = uc.userUsecase.Promote(c.Request.Context(), actorID, targetID)
	if err != nil {
		status := http.StatusInternalServerError
		if err == domain.ErrForbidden {
			status = http.StatusForbidden
		} else if err == domain.ErrNotFound {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user promoted to admin"})
}
