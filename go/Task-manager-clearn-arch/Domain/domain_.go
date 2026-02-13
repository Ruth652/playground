package domain

import (
	"context"
	"errors"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Entities
type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name     string             `bson:"name" json:"name"`
	Email    string             `bson:"email" json:"email"`
	Password string             `bson:"password,omitempty" json:"-"`
	Role     Role               `bson:"role" json:"role"`
}

type Task struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title  string             `bson:"title" json:"title"`
	UserID primitive.ObjectID `bson:"userID" json:"userId"`
}

// Domain-level errors - to map to HTTP codes in controllers
var (
	ErrNotFound       = errors.New("not found")
	ErrConflict       = errors.New("conflict")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrForbidden      = errors.New("forbidden")
	ErrInvalidRequest = errors.New("invalid request")
)

type UserRepository interface {
	Create(ctx context.Context, u *User) error
	GetByEmail(ctx context.Context, email string) (User, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (User, error)
	Count(ctx context.Context) (int64, error) // useful to decide first admin
	PromoteToAdmin(ctx context.Context, id primitive.ObjectID) error
}

type TaskRepository interface {
	Create(ctx context.Context, t *Task) error
	FetchByUser(ctx context.Context, userID primitive.ObjectID) ([]Task, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (Task, error)
	Update(ctx context.Context, t *Task) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}

// Usecase interfaces (optional; useful for wiring & testing)
type UserUsecase interface {
	Register(ctx context.Context, u *User) (accessToken string, refreshToken string, err error)
	Login(ctx context.Context, email, password string) (accessToken string, refreshToken string, err error)
	Promote(ctx context.Context, actorID primitive.ObjectID, targetID primitive.ObjectID) error
}

type TaskUsecase interface {
	CreateTask(ctx context.Context, t *Task, actorID primitive.ObjectID) error
	FetchTasks(ctx context.Context, actorID primitive.ObjectID) ([]Task, error)
}

type JwtCustomClaims struct {
	ID   string `json:"id"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

type JwtCustomRefreshClaims struct {
	ID string `json:"id"`
	jwt.RegisteredClaims
}
