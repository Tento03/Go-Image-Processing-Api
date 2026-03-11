package repositories

import (
	"go-image-processing-api/config"
	"go-image-processing-api/models"
)

func CreateUser(user *models.Auth) error {
	return config.DB.Create(user).Error
}

func IsUsernameExists(username string) (bool, error) {
	var count int64
	err := config.DB.Model(&models.Auth{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}

func FindByUsername(username string) (*models.Auth, error) {
	var user models.Auth
	err := config.DB.Model(&models.Auth{}).Where("username = ?", username).First(&user).Error
	return &user, err
}
