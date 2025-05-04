package repository

import (
	"authentication-service/pkg/models"
	"context"
	"errors"
	"gorm.io/gorm"
)

// UserRepository - структура для работы с пользователями
type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {

	if err := r.db.WithContext(ctx).Create(&user).Error; err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil // Пользователь не найден
	}
	return &user, err
}

func (r *UserRepository) SetUserVerified(ctx context.Context, userID int64) error {
	return r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", userID).Update("is_verified", true).Error
}

func (r *UserRepository) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}
