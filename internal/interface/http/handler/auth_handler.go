// internal/handler/auth_handler.go

package handler

import (
	"encoding/json"
	"maqhaa/auth_service/internal/app/model"
	"maqhaa/auth_service/internal/app/service"
	"maqhaa/library/logging"
	"net/http"
	"strconv"

	"maqhaa/library/middleware"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// AuthHandler handles HTTP requests related to user authentication and authorization.
type AuthHandler struct {
	authService service.AuthService
}

// NewAuthHandler creates a new AuthHandler instance.
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// LoginHandler handles the HTTP request for user login.
func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var loginRequest model.LoginRequest
	var loginResponse model.LoginResponse
	var appError service.AppError
	logID, _ := r.Context().Value(middleware.RequestIDKey).(string)
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		logging.Log.WithFields(logrus.Fields{"request_id": logID}).Info("Invalid request payload")

		appError = *service.NewInvalidFormatError()
		loginResponse = model.LoginResponse{
			HTTPResponse: *model.NewHTTPResponse(appError.Code, appError.Message, nil),
		}
		sendJSONResponse(w, loginResponse, appError.Code)
		return
	}

	// Perform user authentication
	user, appError := h.authService.Authenticate(r.Context(), loginRequest.Username, loginRequest.Password)
	loginResponse = model.LoginResponse{
		HTTPResponse: *model.NewHTTPResponse(appError.Code, appError.Message, nil),
	}

	if appError.Code != 00 {
		loginResponse.Data = nil
		sendJSONResponse(w, loginResponse, appError.Code)
		return
	}

	loginResponse.Data = &struct {
		Token string `json:"token"`
	}{
		Token: user.Token,
	}

	sendJSONResponse(w, loginResponse, appError.Code)
}

func (h *AuthHandler) AddUserHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var addUserRequest model.AddUserRequest
	var response *model.HTTPResponse
	var appError service.AppError
	logID, _ := r.Context().Value(middleware.RequestIDKey).(string)
	err := json.NewDecoder(r.Body).Decode(&addUserRequest)
	if err != nil {
		logging.Log.WithFields(logrus.Fields{"request_id": logID}).Info("Invalid request payload")

		appError = *service.NewInvalidFormatError()

		response := model.NewHTTPResponse(appError.Code, appError.Message, nil)

		sendJSONResponse(w, response, appError.Code)
		return
	}

	token := r.Header.Get("Token")

	if token == "" {
		appError = *service.NewInvalidTokenError()
		response := model.NewHTTPResponse(appError.Code, appError.Message, nil)
		sendJSONResponse(w, response, appError.Code)
		return
	}

	// Perform user authentication
	appError = h.authService.AddUser(r.Context(), addUserRequest, token)
	response = model.NewHTTPResponse(appError.Code, appError.Message, nil)
	sendJSONResponse(w, response, appError.Code)
}

func (h *AuthHandler) EditUserHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var editUserRequest model.EditUserRequest
	var response *model.HTTPResponse
	var appError service.AppError
	logID, _ := r.Context().Value(middleware.RequestIDKey).(string)
	err := json.NewDecoder(r.Body).Decode(&editUserRequest)
	if err != nil {
		logging.Log.WithFields(logrus.Fields{"request_id": logID}).Info("Invalid request payload")

		appError = *service.NewInvalidFormatError()

		response := model.NewHTTPResponse(appError.Code, appError.Message, nil)

		sendJSONResponse(w, response, appError.Code)
		return
	}

	token := r.Header.Get("Token")

	if token == "" {
		appError = *service.NewInvalidTokenError()
		response := model.NewHTTPResponse(appError.Code, appError.Message, nil)
		sendJSONResponse(w, response, appError.Code)
		return
	}

	// Perform user authentication
	appError = h.authService.EditUser(r.Context(), editUserRequest, token)
	response = model.NewHTTPResponse(appError.Code, appError.Message, nil)
	sendJSONResponse(w, response, appError.Code)
}

func (h *AuthHandler) DeactivateUserHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var response *model.HTTPResponse
	var appError service.AppError
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["userID"])

	if err != nil {
		appError = *service.NewInvalidRequestError("Invalid userID")
		response := model.NewHTTPResponse(appError.Code, appError.Message, nil)
		sendJSONResponse(w, response, appError.Code)
		return
	}

	token := r.Header.Get("Token")

	if token == "" {
		appError = *service.NewInvalidTokenError()
		response := model.NewHTTPResponse(appError.Code, appError.Message, nil)
		sendJSONResponse(w, response, appError.Code)
		return
	}

	// Perform user authentication
	appError = h.authService.DeactivateUser(r.Context(), uint(userID), token)
	response = model.NewHTTPResponse(appError.Code, appError.Message, nil)
	sendJSONResponse(w, response, appError.Code)
}

func (h *AuthHandler) GetAllUserHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var response *model.HTTPResponse
	var appError service.AppError

	token := r.Header.Get("Token")

	if token == "" {
		appError = *service.NewInvalidTokenError()
		response := model.NewHTTPResponse(appError.Code, appError.Message, nil)
		sendJSONResponse(w, response, appError.Code)
		return
	}

	// Perform user authentication
	users, appError := h.authService.GetAllUser(r.Context(), token)

	if appError.Code != 00 {
		response := model.NewHTTPResponse(appError.Code, appError.Message, nil)
		sendJSONResponse(w, response, appError.Code)
		return
	}
	response = model.NewHTTPResponse(appError.Code, appError.Message, users)
	sendJSONResponse(w, response, appError.Code)
}

func (h *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Get the token from the request header
	token := r.Header.Get("Token")

	if token == "" {
		appError := *service.NewInvalidTokenError()
		response := model.NewHTTPResponse(appError.Code, appError.Message, nil)
		sendJSONResponse(w, response, appError.Code)
		return
	}

	// Perform user logout
	appError := h.authService.Logout(r.Context(), token)
	response := model.NewHTTPResponse(appError.Code, appError.Message, nil)
	sendJSONResponse(w, response, appError.Code)
}
