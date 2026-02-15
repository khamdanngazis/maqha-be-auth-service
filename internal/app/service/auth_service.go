// internal/service/auth_service.go

package service

import (
	"context"
	"maqhaa/auth_service/internal/app/entity"
	"maqhaa/auth_service/internal/app/model"
	"maqhaa/auth_service/internal/app/repository"
	"maqhaa/library/helper"
	"maqhaa/library/logging"
	"strings"
	"time"

	"maqhaa/library/middleware"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

// AuthService handles user authentication and authorization.
type AuthService interface {
	Authenticate(ctx context.Context, username, password string) (*entity.User, AppError)
	Authorize(ctx context.Context, token string) (*model.User, AppError)
	AddUser(ctx context.Context, request model.AddUserRequest, token string) AppError
	EditUser(ctx context.Context, request model.EditUserRequest, token string) AppError
	GetAllUser(ctx context.Context, token string) ([]*model.User, AppError)
	DeactivateUser(ctx context.Context, ID uint, token string) AppError
	Logout(ctx context.Context, token string) AppError
	// Add other authentication and authorization methods as needed
}

// authServiceImpl implements the AuthService interface
type authServiceImpl struct {
	userRepository repository.UserRepository
}

// NewAuthService creates a new AuthService instance.
func NewAuthService(userRepository repository.UserRepository) AuthService {
	return &authServiceImpl{
		userRepository: userRepository,
	}
}

// Authenticate performs user authentication based on the provided username and password.
func (a *authServiceImpl) Authenticate(ctx context.Context, username, password string) (*entity.User, AppError) {
	logID, _ := ctx.Value(middleware.RequestIDKey).(string)
	user, err := a.userRepository.GetUserByUsername(ctx, username)
	if err != nil {
		if err.Error() != "record not found" {
			return nil, *NewQueryDBError()
		}
	}

	// Check if the provided password matches the stored password
	if user != nil {
		err := helper.CompareHashAndPassword(user.Password, password)
		if err != nil {
			logging.Log.WithFields(logrus.Fields{"request_id": logID}).Errorf("Error CompareHashAndPassword  %s", err.Error())
			return nil, *NewInvalidPasswordError()
		}
		// Assuming you have a token generation function, replace generateToken with the actual implementation
		token, err := helper.GenerateRandomString(16)
		if err != nil {
			logging.Log.WithFields(logrus.Fields{"request_id": logID}).Errorf("Error GenerateRandomString  %s", err.Error())
			return nil, *NewGeneralSystemError()
		}

		// Update all fields of the user in the database
		user.Token = token
		user.TokenExpired = calculateTokenExpiration()
		user.IsActive = true

		err = a.userRepository.UpdateUser(ctx, user)
		if err != nil {
			return nil, *NewUpdateQueryDBError()
		}

		return user, *NewSuccessError()
	}

	return nil, *NewUserNotFoundError()
}

// Authorize performs user authorization based on the provided token.
func (a *authServiceImpl) Authorize(ctx context.Context, token string) (*model.User, AppError) {
	result, err := a.userRepository.GetUserByToken(ctx, token)
	if err != nil {
		if err.Error() != "record not found" {
			return nil, *NewQueryDBError()
		}
	}

	if result == nil {
		return nil, *NewUserNotFoundError()
	}

	// Check if the user is active or any other authorization logic
	var user *model.User
	if !result.IsActive {
		return user, *NewUserNotActiveError()
	}

	if result.TokenExpired.Before(time.Now()) {
		user = &model.User{
			IsLogin: false,
		}
		return user, *NewSuccessError()
	}

	user = &model.User{
		ID:       result.ID,
		ClientID: result.ClientID,
		Username: result.Username,
		FullName: result.FullName,
		IsLogin:  true,
		IsAdmin:  result.Role == entity.RoleAdminCode,
	}

	return user, *NewSuccessError()

}

func (a *authServiceImpl) AddUser(ctx context.Context, request model.AddUserRequest, token string) AppError {
	logID, _ := ctx.Value(middleware.RequestIDKey).(string)
	user, appError := a.Authorize(ctx, token)
	if appError.Code != SuccessError {
		return appError
	}

	if !user.IsLogin {
		return *NewInvalidTokenError()
	}

	if !user.IsAdmin {
		return *NewUserNotAllowError()
	}

	validate := validator.New()
	if err := validate.Struct(request); err != nil {
		return *NewInvalidRequestError(err.Error())
	}

	// Check for duplicate username
	existingUser, err := a.userRepository.GetUserByUsername(ctx, request.Username)
	if err != nil && err.Error() != "record not found" {
		return *NewQueryDBError()
	}
	if existingUser != nil {
		return *NewDuplicateUserError()
	}

	hashedPassword, err := helper.HashPassword(request.Password)
	if err != nil {
		logging.Log.WithFields(logrus.Fields{"request_id": logID}).Errorf("Error HashPassword  %s", err.Error())
		return *NewGeneralSystemError()
	}
	newUser := &entity.User{
		ClientID: user.ClientID,
		Username: request.Username,
		Password: hashedPassword,
		FullName: request.FullName,
		Role:     request.Role,
	}

	err = a.userRepository.CreateUser(ctx, newUser)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") ||
			strings.Contains(err.Error(), "duplicate key value violates unique constraint") ||
			strings.Contains(err.Error(), "SQLSTATE 23505") {
			return *NewDuplicateUserError()
		}
		return *NewUpdateQueryDBError()
	}

	return *NewSuccessError()
}

func (a *authServiceImpl) EditUser(ctx context.Context, request model.EditUserRequest, token string) AppError {
	logID, _ := ctx.Value(middleware.RequestIDKey).(string)
	user, appError := a.Authorize(ctx, token)
	if appError.Code != SuccessError {
		return appError
	}

	if !user.IsLogin {
		return *NewInvalidTokenError()
	}

	if !user.IsAdmin {
		return *NewUserNotAllowError()
	}

	validate := validator.New()
	if err := validate.Struct(request); err != nil {
		return *NewInvalidRequestError(err.Error())
	}

	hashedPassword, err := helper.HashPassword(request.Password)
	if err != nil {
		logging.Log.WithFields(logrus.Fields{"request_id": logID}).Errorf("Error HashPassword  %s", err.Error())
		return *NewGeneralSystemError()
	}
	newUser := &entity.User{
		ID:       request.ID,
		ClientID: user.ClientID,
		Username: request.Username,
		Password: hashedPassword,
		FullName: request.FullName,
		Role:     request.Role,
	}

	err = a.userRepository.UpdateUser(ctx, newUser)
	if err != nil {
		return *NewUpdateQueryDBError()
	}

	return *NewSuccessError()
}

func (a *authServiceImpl) GetAllUser(ctx context.Context, token string) ([]*model.User, AppError) {
	user, appError := a.Authorize(ctx, token)
	if appError.Code != SuccessError {
		return nil, appError
	}

	if !user.IsLogin {
		return nil, *NewInvalidTokenError()
	}

	if !user.IsAdmin {
		return nil, *NewUserNotAllowError()
	}

	usersFull, err := a.userRepository.GetAllUserByClientID(ctx, int(user.ClientID))
	if err != nil {
		return nil, *NewQueryDBError()
	}
	var users []*model.User
	for _, u := range usersFull {
		addUser := &model.User{
			ID:       u.ID,
			Username: u.Username,
			FullName: u.FullName,
			Role:     u.Role,
		}
		users = append(users, addUser)
	}

	return users, *NewSuccessError()
}

func (a *authServiceImpl) DeactivateUser(ctx context.Context, ID uint, token string) AppError {
	user, appError := a.Authorize(ctx, token)
	if appError.Code != SuccessError {
		return appError
	}

	if !user.IsLogin {
		return *NewInvalidTokenError()
	}

	if !user.IsAdmin {
		return *NewUserNotAllowError()
	}

	if ID == 0 {
		return *NewInvalidRequestError("Invalid UserID")
	}

	err := a.userRepository.DeactivateUser(ctx, ID)
	if err != nil {
		return *NewUpdateQueryDBError()
	}

	return *NewSuccessError()
}

func (a *authServiceImpl) Logout(ctx context.Context, token string) AppError {
	logID, _ := ctx.Value(middleware.RequestIDKey).(string)
	user, appError := a.Authorize(ctx, token)
	if appError.Code != SuccessError {
		return appError
	}

	token, err := helper.GenerateRandomString(8)
	if err != nil {
		logging.Log.WithFields(logrus.Fields{"request_id": logID}).Errorf("Error GenerateRandomString  %s", err.Error())
		return *NewGeneralSystemError()
	}

	newUser := &entity.User{
		ID:    user.ID,
		Token: token,
	}

	err = a.userRepository.UpdateUser(ctx, newUser)
	if err != nil {
		return *NewUpdateQueryDBError()
	}

	return *NewSuccessError()
}

func calculateTokenExpiration() time.Time {
	// Set token expiration to 15 minutes
	return time.Now().Add(time.Minute * 15)
}
