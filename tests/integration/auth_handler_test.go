package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"maqhaa/auth_service/internal/app/entity"
	"maqhaa/auth_service/internal/app/model"
	"maqhaa/auth_service/internal/app/service"
	"maqhaa/library/helper"
	"maqhaa/library/logging"

	"maqhaa/library/middleware"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	pb "maqhaa/auth_service/internal/interface/grpc/model"
)

func TestLoginHandler_Positive(t *testing.T) {
	// create mock data
	client := SampleClient()
	db.Create(client)
	userLogin := SampleUser(client.ID)
	hashedPassword, _ := helper.HashPassword(userLogin.Password)
	userLogin.Password = hashedPassword
	db.Create(userLogin)

	tables := []string{"user", "client"}
	defer clearDB(tables)

	user := SampleUser(client.ID)
	// Create a login request
	loginRequest := model.LoginRequest{
		Username: user.Username,
		Password: user.Password,
	}

	// Marshal the login request to JSON
	loginRequestJSON, err := json.Marshal(loginRequest)
	if err != nil {
		t.Fatal(err)
	}

	// Mock HTTP request
	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(loginRequestJSON))
	if err != nil {
		t.Fatal(err)
	}
	requestID := uuid.New().String()
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, requestID)
	req = req.WithContext(ctx)

	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()

	// Call the handler function
	http.HandlerFunc(authHandler.LoginHandler).ServeHTTP(rr, req)
	logging.Log.WithFields(logrus.Fields{
		"RequestID": requestID,
		"Status":    rr.Code,
		"Body":      rr.Body.String(),
	}).Info("Outgoing response")
	// Check the response status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse the response body
	var response model.LoginResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	// Perform assertions based on the expected login response
	assert.Equal(t, service.SuccessMessage, response.Message)
	assert.Equal(t, service.SuccessError, response.Code)
	assert.NotEmpty(t, response.Data.Token)
}

func TestLoginHandler_NotExisUser(t *testing.T) {
	// Clean up the testing environment
	tables := []string{"user", "client"}
	defer clearDB(tables)

	// Create a login request with invalid credentials
	loginRequest := model.LoginRequest{
		Username: "nonExistentUser",
		Password: "invalidPassword",
	}

	// Marshal the login request to JSON
	loginRequestJSON, err := json.Marshal(loginRequest)
	if err != nil {
		t.Fatal(err)
	}

	// Mock HTTP request
	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(loginRequestJSON))
	if err != nil {
		t.Fatal(err)
	}
	requestID := uuid.New().String()
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, requestID)
	req = req.WithContext(ctx)
	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()
	logging.Log.WithFields(logrus.Fields{
		"RequestID": requestID,
		"Status":    rr.Code,
		"Body":      rr.Body.String(),
	}).Info("Outgoing response")
	// Call the handler function
	http.HandlerFunc(authHandler.LoginHandler).ServeHTTP(rr, req)

	// Check the response status code
	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	// Parse the response body
	var response model.HTTPResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	// Perform assertions based on the expected error response
	assert.Equal(t, service.InvalidUsernameMessage, response.Message)
	assert.Equal(t, service.InvalidUsername, response.Code)
	assert.Nil(t, response.Data) // Data should be nil for failed login attempts
}

func TestAddUserHandler_Positive(t *testing.T) {
	// create mock data
	//tables := []string{"user", "client"}
	//defer clearDB(tables)

	client := SampleClient()
	db.Create(client)
	userLogin := SampleUser(client.ID)
	hashedPassword, _ := helper.HashPassword(userLogin.Password)
	userLogin.Password = hashedPassword
	db.Create(userLogin)

	// Create a login request
	addUserRequest := model.AddUserRequest{
		Username: "New User",
		Password: "Password",
		FullName: "New User",
		Role:     2,
	}

	// Marshal the login request to JSON
	loginRequestJSON, err := json.Marshal(addUserRequest)
	if err != nil {
		t.Fatal(err)
	}

	// Mock HTTP request
	req, err := http.NewRequest("POST", "/user", bytes.NewBuffer(loginRequestJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Token", userLogin.Token)
	requestID := uuid.New().String()
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, requestID)
	req = req.WithContext(ctx)

	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()

	// Call the handler function
	http.HandlerFunc(authHandler.AddUserHandler).ServeHTTP(rr, req)
	logging.Log.WithFields(logrus.Fields{
		"RequestID": requestID,
		"Status":    rr.Code,
		"Body":      rr.Body.String(),
	}).Info("Outgoing response")
	// Check the response status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse the response body
	var response model.LoginResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	// Perform assertions based on the expected login response
	assert.Equal(t, service.SuccessMessage, response.Message)
	assert.Equal(t, service.SuccessError, response.Code)
}

func TestAddUserHandler_UserNotAllowed(t *testing.T) {
	// create mock data
	tables := []string{"user", "client"}
	defer clearDB(tables)

	client := SampleClient()
	db.Create(client)
	userLogin := SampleUser(client.ID)
	hashedPassword, _ := helper.HashPassword(userLogin.Password)
	userLogin.Password = hashedPassword
	userLogin.Role = 2
	db.Create(userLogin)

	// Create a login request
	addUserRequest := model.AddUserRequest{
		Username: "New User",
		Password: "Password",
		FullName: "New User",
		Role:     2,
	}

	// Marshal the login request to JSON
	loginRequestJSON, err := json.Marshal(addUserRequest)
	if err != nil {
		t.Fatal(err)
	}

	// Mock HTTP request
	req, err := http.NewRequest("POST", "/user", bytes.NewBuffer(loginRequestJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Token", userLogin.Token)
	requestID := uuid.New().String()
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, requestID)
	req = req.WithContext(ctx)

	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()

	// Call the handler function
	http.HandlerFunc(authHandler.AddUserHandler).ServeHTTP(rr, req)
	logging.Log.WithFields(logrus.Fields{
		"RequestID": requestID,
		"Status":    rr.Code,
		"Body":      rr.Body.String(),
	}).Info("Outgoing response")
	// Check the response status code
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// Parse the response body
	var response model.LoginResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	// Perform assertions based on the expected login response
	assert.Equal(t, service.UserNotAllowMessage, response.Message)
	assert.Equal(t, service.UserNotAllowError, response.Code)
}

func TestAddUserHandler_UserNotActive(t *testing.T) {
	// create mock data
	tables := []string{"user", "client"}
	defer clearDB(tables)

	client := SampleClient()
	db.Create(client)
	userLogin := SampleUser(client.ID)
	hashedPassword, _ := helper.HashPassword(userLogin.Password)
	userLogin.Password = hashedPassword
	userLogin.IsActive = false
	db.Create(userLogin)

	// Create a login request
	addUserRequest := model.AddUserRequest{
		Username: "New User",
		Password: "Password",
		FullName: "New User",
		Role:     2,
	}

	// Marshal the login request to JSON
	loginRequestJSON, err := json.Marshal(addUserRequest)
	if err != nil {
		t.Fatal(err)
	}

	// Mock HTTP request
	req, err := http.NewRequest("POST", "/user", bytes.NewBuffer(loginRequestJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Token", userLogin.Token)
	requestID := uuid.New().String()
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, requestID)
	req = req.WithContext(ctx)

	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()

	// Call the handler function
	http.HandlerFunc(authHandler.AddUserHandler).ServeHTTP(rr, req)
	logging.Log.WithFields(logrus.Fields{
		"RequestID": requestID,
		"Status":    rr.Code,
		"Body":      rr.Body.String(),
	}).Info("Outgoing response")
	// Check the response status code
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// Parse the response body
	var response model.LoginResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	// Perform assertions based on the expected login response
	assert.Equal(t, service.UserNotActiveMessage, response.Message)
	assert.Equal(t, service.UserNotActiveError, response.Code)
}

func TestAddUserHandler_DuplicateUser(t *testing.T) {
	// create mock data
	tables := []string{"user", "client"}
	defer clearDB(tables)

	client := SampleClient()
	db.Create(client)
	userLogin := SampleUser(client.ID)
	hashedPassword, _ := helper.HashPassword(userLogin.Password)
	userLogin.Password = hashedPassword
	db.Create(userLogin)

	// Create a login request
	addUserRequest := model.AddUserRequest{
		Username: userLogin.Username,
		Password: "Password",
		FullName: "New User",
		Role:     2,
	}

	// Marshal the login request to JSON
	loginRequestJSON, err := json.Marshal(addUserRequest)
	if err != nil {
		t.Fatal(err)
	}

	// Mock HTTP request
	req, err := http.NewRequest("POST", "/user", bytes.NewBuffer(loginRequestJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Token", userLogin.Token)
	requestID := uuid.New().String()
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, requestID)
	req = req.WithContext(ctx)

	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()

	// Call the handler function
	http.HandlerFunc(authHandler.AddUserHandler).ServeHTTP(rr, req)
	logging.Log.WithFields(logrus.Fields{
		"RequestID": requestID,
		"Status":    rr.Code,
		"Body":      rr.Body.String(),
	}).Info("Outgoing response")
	// Check the response status code
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// Parse the response body
	var response model.LoginResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	// Perform assertions based on the expected login response
	assert.Equal(t, service.DuplicateUserMessage, response.Message)
	assert.Equal(t, service.DuplicateUserError, response.Code)
}

func TestAddUserHandler_UserNotFound(t *testing.T) {
	// create mock data
	tables := []string{"user", "client"}
	defer clearDB(tables)

	client := SampleClient()
	db.Create(client)
	userLogin := SampleUser(client.ID)
	hashedPassword, _ := helper.HashPassword(userLogin.Password)
	userLogin.Password = hashedPassword
	db.Create(userLogin)

	// Create a login request
	addUserRequest := model.AddUserRequest{
		Username: "userLogin.Username",
		Password: "Password",
		FullName: "THIS STRING IS 500 CHARACTERS xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
		Role:     2,
	}

	// Marshal the login request to JSON
	loginRequestJSON, err := json.Marshal(addUserRequest)
	if err != nil {
		t.Fatal(err)
	}

	// Mock HTTP request
	req, err := http.NewRequest("POST", "/user", bytes.NewBuffer(loginRequestJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Token", "invalid")
	requestID := uuid.New().String()
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, requestID)
	req = req.WithContext(ctx)

	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()

	// Call the handler function
	http.HandlerFunc(authHandler.AddUserHandler).ServeHTTP(rr, req)
	logging.Log.WithFields(logrus.Fields{
		"RequestID": requestID,
		"Status":    rr.Code,
		"Body":      rr.Body.String(),
	}).Info("Outgoing response")
	// Check the response status code
	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	// Parse the response body
	var response model.LoginResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	// Perform assertions based on the expected login response
	assert.Equal(t, service.InvalidUsernameMessage, response.Message)
	assert.Equal(t, service.InvalidUsername, response.Code)
}

func TestAddUserHandler_InvalidToken(t *testing.T) {
	// create mock data
	tables := []string{"user", "client"}
	defer clearDB(tables)

	client := SampleClient()
	db.Create(client)
	userLogin := SampleUser(client.ID)
	hashedPassword, _ := helper.HashPassword(userLogin.Password)
	userLogin.Password = hashedPassword
	db.Create(userLogin)

	// Create a login request
	addUserRequest := model.AddUserRequest{
		Username: "userLogin.Username",
		Password: "Password",
		FullName: "THIS STRING IS 500 CHARACTERS xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
		Role:     2,
	}

	// Marshal the login request to JSON
	loginRequestJSON, err := json.Marshal(addUserRequest)
	if err != nil {
		t.Fatal(err)
	}

	// Mock HTTP request
	req, err := http.NewRequest("POST", "/user", bytes.NewBuffer(loginRequestJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Token", "")
	requestID := uuid.New().String()
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, requestID)
	req = req.WithContext(ctx)

	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()

	// Call the handler function
	http.HandlerFunc(authHandler.AddUserHandler).ServeHTTP(rr, req)
	logging.Log.WithFields(logrus.Fields{
		"RequestID": requestID,
		"Status":    rr.Code,
		"Body":      rr.Body.String(),
	}).Info("Outgoing response")
	// Check the response status code
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// Parse the response body
	var response model.LoginResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	// Perform assertions based on the expected login response
	assert.Equal(t, service.InvalidTokendMessage, response.Message)
	assert.Equal(t, service.InvalidToken, response.Code)
}

func TestEditUserHandler_Positive(t *testing.T) {
	// create mock data
	tables := []string{"user", "client"}
	defer clearDB(tables)

	client := SampleClient()
	db.Create(client)
	userLogin := SampleUser(client.ID)
	hashedPassword, _ := helper.HashPassword(userLogin.Password)
	userLogin.Password = hashedPassword
	db.Create(userLogin)

	user := SampleUser(client.ID)
	hashedPassword, _ = helper.HashPassword(user.Password)
	user.Password = hashedPassword
	db.Create(user)

	user.FullName = "New Edit User"
	user.Username = "New.Edit.User"
	user.Password = "rahasiaBanget"
	user.Role = 2

	// Create a login request
	UserRequest := model.AddUserRequest{
		FullName: user.FullName,
		Username: user.Username,
		Password: user.Password,
		Role:     user.Role,
	}

	EditUserRequest := model.EditUserRequest{
		ID:             userLogin.ID,
		AddUserRequest: UserRequest,
	}

	// Marshal the login request to JSON
	requestJSON, err := json.Marshal(EditUserRequest)
	if err != nil {
		t.Fatal(err)
	}
	requestID := uuid.New().String()
	logging.Log.WithFields(logrus.Fields{
		"RequestID": requestID,
		"Body":      string(requestJSON),
	}).Info("Incoming request")
	// Mock HTTP request
	req, err := http.NewRequest("PUT", "/user", bytes.NewBuffer(requestJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Token", userLogin.Token)
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, requestID)
	req = req.WithContext(ctx)

	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()

	// Call the handler function
	http.HandlerFunc(authHandler.EditUserHandler).ServeHTTP(rr, req)
	logging.Log.WithFields(logrus.Fields{
		"RequestID": requestID,
		"Status":    rr.Code,
		"Body":      rr.Body.String(),
	}).Info("Outgoing response")
	// Check the response status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse the response body
	var response model.LoginResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	// Perform assertions based on the expected login response
	assert.Equal(t, service.SuccessMessage, response.Message)
	assert.Equal(t, service.SuccessError, response.Code)
	var userNew entity.User
	result := db.First(&userNew, user.ID).Error
	if result != nil {
		t.Fatal(result)
	}

	assert.Equal(t, user.FullName, userNew.FullName)
	assert.Equal(t, user.Username, userNew.Username)
	errPass := helper.CompareHashAndPassword(userNew.Password, user.Password)
	assert.Nil(t, errPass)
	assert.Equal(t, user.Role, userNew.Role)
}

func TestEditUserHandler_InvalidToken(t *testing.T) {
	// create mock data
	tables := []string{"user", "client"}
	defer clearDB(tables)

	client := SampleClient()
	db.Create(client)
	userLogin := SampleUser(client.ID)
	hashedPassword, _ := helper.HashPassword(userLogin.Password)
	userLogin.Password = hashedPassword
	db.Create(userLogin)

	user := SampleUser(client.ID)
	hashedPassword, _ = helper.HashPassword(user.Password)
	user.Password = hashedPassword
	db.Create(user)

	user.FullName = "New Edit User"
	user.Username = "New.Edit.User"
	user.Password = "rahasiaBanget"
	user.Role = 2

	// Create a login request
	UserRequest := model.AddUserRequest{
		FullName: user.FullName,
		Username: user.Username,
		Password: user.Password,
		Role:     user.Role,
	}

	EditUserRequest := model.EditUserRequest{
		ID:             userLogin.ID,
		AddUserRequest: UserRequest,
	}

	// Marshal the login request to JSON
	requestJSON, err := json.Marshal(EditUserRequest)
	if err != nil {
		t.Fatal(err)
	}
	requestID := uuid.New().String()
	logging.Log.WithFields(logrus.Fields{
		"RequestID": requestID,
		"Body":      string(requestJSON),
	}).Info("Incoming request")
	// Mock HTTP request
	req, err := http.NewRequest("PUT", "/user", bytes.NewBuffer(requestJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Token", "")
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, requestID)
	req = req.WithContext(ctx)

	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()

	// Call the handler function
	http.HandlerFunc(authHandler.EditUserHandler).ServeHTTP(rr, req)
	logging.Log.WithFields(logrus.Fields{
		"RequestID": requestID,
		"Status":    rr.Code,
		"Body":      rr.Body.String(),
	}).Info("Outgoing response")
	// Check the response status code
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// Parse the response body
	var response model.LoginResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	// Perform assertions based on the expected login response
	assert.Equal(t, service.InvalidTokendMessage, response.Message)
	assert.Equal(t, service.InvalidToken, response.Code)

}

func TestEditUserHandler_InvalidRequest(t *testing.T) {
	// create mock data
	tables := []string{"user", "client"}
	defer clearDB(tables)

	client := SampleClient()
	db.Create(client)
	userLogin := SampleUser(client.ID)
	hashedPassword, _ := helper.HashPassword(userLogin.Password)
	userLogin.Password = hashedPassword
	db.Create(userLogin)

	user := SampleUser(client.ID)
	hashedPassword, _ = helper.HashPassword(user.Password)
	user.Password = hashedPassword
	db.Create(user)

	user.FullName = "New Edit User"
	user.Username = "New.Edit.User"
	user.Password = "rahasiaBanget"
	user.Role = 2

	// Create a login request
	UserRequest := model.AddUserRequest{
		FullName: "",
		Username: user.Username,
		Password: user.Password,
		Role:     user.Role,
	}

	EditUserRequest := model.EditUserRequest{
		ID:             userLogin.ID,
		AddUserRequest: UserRequest,
	}

	// Marshal the login request to JSON
	requestJSON, err := json.Marshal(EditUserRequest)
	if err != nil {
		t.Fatal(err)
	}
	requestID := uuid.New().String()
	logging.Log.WithFields(logrus.Fields{
		"RequestID": requestID,
		"Body":      string(requestJSON),
	}).Info("Incoming request")
	// Mock HTTP request
	req, err := http.NewRequest("PUT", "/user", bytes.NewBuffer(requestJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Token", userLogin.Token)
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, requestID)
	req = req.WithContext(ctx)

	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()

	// Call the handler function
	http.HandlerFunc(authHandler.EditUserHandler).ServeHTTP(rr, req)
	logging.Log.WithFields(logrus.Fields{
		"RequestID": requestID,
		"Status":    rr.Code,
		"Body":      rr.Body.String(),
	}).Info("Outgoing response")
	// Check the response status code
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// Parse the response body
	var response model.LoginResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	// Perform assertions based on the expected login response
	assert.Equal(t, service.InvalidRequestError, response.Code)

}

func TestDeactivateUserHandler_Positive(t *testing.T) {
	// create mock data
	tables := []string{"user", "client"}
	defer clearDB(tables)

	client := SampleClient()
	db.Create(client)

	userLogin := SampleUser(client.ID)
	hashedPassword, _ := helper.HashPassword(userLogin.Password)
	userLogin.Password = hashedPassword
	db.Create(userLogin)

	user := SampleUserCS(client.ID, "jovan.yoooo")
	hashedPassword, _ = helper.HashPassword(user.Password)
	user.Password = hashedPassword
	db.Create(user)

	fmt.Println(user.ID)

	router := mux.NewRouter()
	router.HandleFunc("/user/{userID}", authHandler.DeactivateUserHandler).Methods("DELETE")

	// Mock HTTP request
	req, err := http.NewRequest("DELETE", "/user/"+strconv.Itoa(int(user.ID)), nil)
	if err != nil {
		t.Fatal(err)
	}
	requestID := uuid.New().String()
	req.Header.Set("Token", userLogin.Token)
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, requestID)
	req = req.WithContext(ctx)

	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()

	// Call the handler function
	router.ServeHTTP(rr, req)
	logging.Log.WithFields(logrus.Fields{
		"RequestID": requestID,
		"Status":    rr.Code,
		"Body":      rr.Body.String(),
	}).Info("Outgoing response")
	// Check the response status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse the response body
	var response model.HTTPResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	// Perform assertions based on the expected login response
	assert.Equal(t, service.SuccessMessage, response.Message)
	assert.Equal(t, service.SuccessError, response.Code)
	var userNew entity.User
	result := db.First(&userNew, user.ID).Error
	if result != nil {
		t.Fatal(result)
	}

	assert.Equal(t, userNew.IsActive, false)
}

func TestLogoutHandler_Positive(t *testing.T) {
	// create mock data
	tables := []string{"user", "client"}
	defer clearDB(tables)

	client := SampleClient()
	db.Create(client)

	userLogin := SampleUser(client.ID)
	hashedPassword, _ := helper.HashPassword(userLogin.Password)
	userLogin.Password = hashedPassword
	db.Create(userLogin)

	router := mux.NewRouter()
	router.HandleFunc("/logout", authHandler.LogoutHandler).Methods("DELETE")

	// Mock HTTP request
	req, err := http.NewRequest("DELETE", "/logout", nil)
	if err != nil {
		t.Fatal(err)
	}
	requestID := uuid.New().String()
	req.Header.Set("Token", userLogin.Token)
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, requestID)
	req = req.WithContext(ctx)

	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()

	// Call the handler function
	router.ServeHTTP(rr, req)
	logging.Log.WithFields(logrus.Fields{
		"RequestID": requestID,
		"Status":    rr.Code,
		"Body":      rr.Body.String(),
	}).Info("Outgoing response")
	// Check the response status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse the response body
	var response model.HTTPResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	// Perform assertions based on the expected login response
	assert.Equal(t, service.SuccessMessage, response.Message)
	assert.Equal(t, service.SuccessError, response.Code)
	var userNew entity.User
	result := db.First(&userNew, userLogin.ID).Error
	if result != nil {
		t.Fatal(result)
	}

	assert.NotEqual(t, userLogin.Token, userNew.Token)
}

func TestGetAllUserHandler_Positive(t *testing.T) {
	// create mock data
	tables := []string{"user", "client"}
	defer clearDB(tables)

	client := SampleClient()
	db.Create(client)

	userLogin := SampleUser(client.ID)
	hashedPassword, _ := helper.HashPassword(userLogin.Password)
	userLogin.Password = hashedPassword
	db.Create(userLogin)

	user := SampleUserCS(client.ID, "user1")
	db.Create(user)
	user = SampleUserCS(client.ID, "user12")
	db.Create(user)
	user = SampleUserCS(client.ID, "user3")
	db.Create(user)
	user = SampleUserCS(client.ID, "user4")
	db.Create(user)
	user = SampleUserCS(client.ID, "user5")
	db.Create(user)

	router := mux.NewRouter()
	router.HandleFunc("/user", authHandler.GetAllUserHandler).Methods("GET")

	// Mock HTTP request
	req, err := http.NewRequest("GET", "/user", nil)
	if err != nil {
		t.Fatal(err)
	}
	requestID := uuid.New().String()
	req.Header.Set("Token", userLogin.Token)
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, requestID)
	req = req.WithContext(ctx)

	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()

	// Call the handler function
	router.ServeHTTP(rr, req)
	logging.Log.WithFields(logrus.Fields{
		"RequestID": requestID,
		"Status":    rr.Code,
		"Body":      rr.Body.String(),
	}).Info("Outgoing response")
	// Check the response status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse the response body
	var response model.GetUserResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	// Perform assertions based on the expected login response
	assert.Equal(t, service.SuccessMessage, response.Message)
	assert.Equal(t, service.SuccessError, response.Code)

	assert.Equal(t, len(*response.Data), 6)
}

func TestGetUserGRPCHandler_Positive(t *testing.T) {
	// create mock data
	tables := []string{"user", "client"}
	defer clearDB(tables)

	client := SampleClient()
	db.Create(client)

	userLogin := SampleUser(client.ID)
	hashedPassword, _ := helper.HashPassword(userLogin.Password)
	userLogin.Password = hashedPassword
	db.Create(userLogin)

	// Set up a gRPC connection to the server
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Error creating gRPC client connection: %v", err)
	}
	defer conn.Close()

	clientServer := pb.NewUserClient(conn)

	// Prepare a request
	reqRpc := &pb.GetUserRequest{
		Token: userLogin.Token, // Replace with a valid product ID for your test data
	}

	// Call the gRPC method
	resp, err := clientServer.GetUser(context.Background(), reqRpc)
	if err != nil {
		t.Fatalf("Error calling GetUser gRPC method: %v", err)
	}

	// Assertions
	assert.NotNil(t, resp)
	assert.Equal(t, int32(service.SuccessError), resp.Code) // Assuming 0 is the success code
	assert.Equal(t, service.SuccessMessage, resp.Message)
	// Add more assertions based on your response structure

	// Cleanup if needed
	// ...

}
