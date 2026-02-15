// internal/repository/user_repository.go

package repository

import (
	"context"
	"fmt"
	"maqhaa/auth_service/internal/app/entity"
	"maqhaa/library/logging"

	"maqhaa/library/middleware"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// UserRepository handles database interactions related to users.
type UserRepository interface {
	CreateUser(ctx context.Context, user *entity.User) error
	GetUserByID(ctx context.Context, userID uint) (*entity.User, error)
	GetUserByToken(ctx context.Context, token string) (*entity.User, error)
	GetUserByUsername(ctx context.Context, username string) (*entity.User, error)
	UpdateUser(ctx context.Context, user *entity.User) error
	GetClientByToken(ctx context.Context, token string) (*entity.Client, error)
	GetAllUserByClientID(ctx context.Context, clientID int) ([]*entity.User, error)
	DeactivateUser(ctx context.Context, ID uint) error
	// Add other user-related methods as needed
}

// Implement the interface in the UserRepository struct
type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

// GetAllUserByClientID retrieves all users belonging to a specific client ID from the database.
func (r *userRepository) GetAllUserByClientID(ctx context.Context, clientID int) ([]*entity.User, error) {
	logID, _ := ctx.Value(middleware.RequestIDKey).(string)
	var users []*entity.User
	result := r.db.Where("client_id = ?", clientID).Find(&users)
	if result.Error != nil {
		logging.Log.WithFields(logrus.Fields{"request_id": logID}).Errorf("Error GetAllUserByClientID  %s", result.Error.Error())
		return nil, result.Error
	}
	return users, nil
}

// CreateUser creates a new user in the database.
func (r *userRepository) CreateUser(ctx context.Context, user *entity.User) error {
	logID, _ := ctx.Value(middleware.RequestIDKey).(string)
	result := r.db.Create(user)
	if result.Error != nil {
		logging.Log.WithFields(logrus.Fields{"request_id": logID}).Errorf("Error CreateUser  %s", result.Error.Error())
		return result.Error
	}
	return nil
}

// GetUserByID retrieves a user by ID from the database.
func (r *userRepository) GetUserByID(ctx context.Context, userID uint) (*entity.User, error) {
	logID, _ := ctx.Value(middleware.RequestIDKey).(string)
	var user entity.User
	result := r.db.First(&user, userID)
	if result.Error != nil {
		if result.Error.Error() != "record not found" {
			logging.Log.WithFields(logrus.Fields{"request_id": logID}).Errorf("Error GetUserByID  %s", result.Error.Error())
		}
		return nil, result.Error
	}
	return &user, nil
}

func (r *userRepository) GetUserByToken(ctx context.Context, token string) (*entity.User, error) {
	var user entity.User
	logID, _ := ctx.Value(middleware.RequestIDKey).(string)
	result := r.db.Where("token = ?", token).First(&user)
	if result.Error != nil {
		if result.Error.Error() != "record not found" {
			logging.Log.WithFields(logrus.Fields{"request_id": logID}).Errorf("Error GetUserByToken  %s", result.Error.Error())
		}
		return nil, result.Error
	}
	return &user, nil
}

func (r *userRepository) GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	var user entity.User
	logID, _ := ctx.Value(middleware.RequestIDKey).(string)
	result := r.db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		if result.Error.Error() != "record not found" {
			logging.Log.WithFields(logrus.Fields{"request_id": logID}).Errorf("Error GetUserByUsername  %s", result.Error.Error())
		}
		return nil, result.Error
	}
	return &user, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, user *entity.User) error {
	fmt.Println(user)
	result := r.db.Model(user).Updates(user)
	logID, _ := ctx.Value(middleware.RequestIDKey).(string)
	if result.Error != nil {
		logging.Log.WithFields(logrus.Fields{"request_id": logID}).Errorf("Error UpdateUser  %s", result.Error.Error())
		return result.Error
	}
	return nil
}

func (r *userRepository) DeactivateUser(ctx context.Context, ID uint) error {

	logID, _ := ctx.Value(middleware.RequestIDKey).(string)

	updates := map[string]interface{}{
		"IsActive": false,
	}

	// Perform the update operation
	result := r.db.Model(&entity.User{}).Where("id = ?", ID).Updates(updates)
	if result.Error != nil {
		logging.Log.WithFields(logrus.Fields{"request_id": logID}).Errorf("Error UpdateUser  %s", result.Error.Error())
		return result.Error
	}
	return nil
}

func (r *userRepository) GetClientByToken(ctx context.Context, token string) (*entity.Client, error) {
	var client entity.Client
	logID, _ := ctx.Value(middleware.RequestIDKey).(string)
	result := r.db.Where("token = ?", token).First(&client)
	if result.Error != nil {
		if result.Error.Error() != "record not found" {
			logging.Log.WithFields(logrus.Fields{"request_id": logID}).Errorf("Error GetClientByToken  %s", result.Error.Error())
		}
		return nil, result.Error
	}
	return &client, nil
}
