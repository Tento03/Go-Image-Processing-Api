package services

import (
	"errors"
	"go-image-processing-api/models"
	"go-image-processing-api/repositories"
	"go-image-processing-api/utils"

	"github.com/google/uuid"
)

var ErrUsernameExist = errors.New("username already exist")

func Register(username string, password string) (*models.Auth, error) {
	exist, err := repositories.IsUsernameExists(username)
	if err != nil {
		return nil, err
	}

	if exist {
		return nil, ErrUsernameExist
	}

	hashed, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &models.Auth{
		UserID:   uuid.NewString(),
		Username: username,
		Password: string(hashed),
	}

	if err := repositories.CreateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}
