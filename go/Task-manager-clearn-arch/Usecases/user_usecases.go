package usecases

import (
	"context"
	"errors"

	"task-manager-clean-arch/domain"
	"task-manager-clean-arch/infrastructure"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type userUsecase struct {
	userRepo    domain.UserRepository
	passwordSvc infrastructure.PasswordService
	jwtSvc      infrastructure.JWTService
}

func NewUserUsecase(userRepo domain.UserRepository, ps infrastructure.PasswordService, js infrastructure.JWTService) domain.UserUsecase {
	return &userUsecase{
		userRepo:    userRepo,
		passwordSvc: ps,
		jwtSvc:      js,
	}
}

func (u *userUsecase) Register(ctx context.Context, user *domain.User) (string, string, error) {
	// Check if user exists
	_, err := u.userRepo.GetByEmail(ctx, user.Email)

	// User exists -> conflict
	if err == nil {
		return "", "", domain.ErrConflict
	}

	// If error is not "not found"
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return "", "", err
	}

	// Hash password
	hashed, err := u.passwordSvc.Hash(user.Password)
	if err != nil {
		return "", "", err
	}

	user.Password = hashed

	// Assign role if the first user
	count, err := u.userRepo.Count(ctx)
	if err != nil {
		return "", "", err
	}

	if count == 0 {
		user.Role = domain.RoleAdmin
	} else {
		user.Role = domain.RoleUser
	}

	// Save user
	err = u.userRepo.Create(ctx, user)
	if err != nil {
		return "", "", err
	}

	// Generate tokens
	access, refresh, err := u.jwtSvc.GenerateTokens(user)
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

func (u *userUsecase) Login(ctx context.Context, email, password string) (string, string, error) {
	user, err := u.userRepo.GetByEmail(ctx, email)

	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return "", "", domain.ErrUnauthorized
		}
		return "", "", err
	}

	if !u.passwordSvc.Compare(password, user.Password) {
		return "", "", domain.ErrUnauthorized
	}

	access, refresh, err := u.jwtSvc.GenerateTokens(&user)
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

func (u *userUsecase) Promote(ctx context.Context, actorID, targetID primitive.ObjectID) error {
	actor, err := u.userRepo.GetByID(ctx, actorID)
	if err != nil {
		return err
	}

	if actor.Role != domain.RoleAdmin {
		return domain.ErrForbidden
	}

	return u.userRepo.PromoteToAdmin(ctx, targetID)
}
