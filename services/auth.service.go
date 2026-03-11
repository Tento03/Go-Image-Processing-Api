package services

import (
	"errors"
	"go-image-processing-api/models"
	"go-image-processing-api/repositories"
	"go-image-processing-api/utils"
	"time"

	"github.com/google/uuid"
)

var ErrUsernameExist = errors.New("username already exist")
var ErrInvalidCredentials = errors.New("invalid credentials")

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

func Login(username string, password string) (string, string, error) {
	user, err := repositories.FindByUsername(username)
	if err != nil {
		return "", "", err
	}

	if !utils.ComparePassword(user.Password, password) {
		return "", "", ErrInvalidCredentials
	}

	accessToken, err := utils.GenerateAccessToken(user.UserID)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := utils.GenerateRefreshToken(user.UserID)
	if err != nil {
		return "", "", err
	}

	hashRT := utils.HashToken(refreshToken)
	refresh := &models.Refresh{
		UserID:    user.UserID,
		Token:     hashRT,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}

	if err := repositories.SaveRefreshToken(refresh); err != nil {
		return "", "", nil
	}

	return accessToken, refreshToken, nil
}

func Refresh(refreshToken string) (string, string, error) {
	hash := utils.HashToken(refreshToken)
	old, err := repositories.FindValidRefreshToken(hash)
	if err != nil {
		return "", "", nil
	}

	if err := repositories.RevokeToken(old); err != nil {
		return "", "", nil
	}

	newAccessToken, _ := utils.GenerateAccessToken(old.UserID)
	newRefreshToken, _ := utils.GenerateRefreshToken(old.UserID)

	newHashRT := utils.HashToken(newRefreshToken)
	refresh := &models.Refresh{
		UserID:    old.UserID,
		Token:     newHashRT,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}

	if err := repositories.SaveRefreshToken(refresh); err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}

func Logout(refreshToken string) error {
	hashRT := utils.HashToken(refreshToken)
	old, err := repositories.FindValidRefreshToken(hashRT)
	if err != nil {
		return err
	}

	return repositories.RevokeToken(old)
}
