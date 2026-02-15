// auth_handler_test.go

package handler_test

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"testing"
	"time"

	"maqhaa/auth_service/internal/app/entity"
	"maqhaa/auth_service/internal/app/repository"
	"maqhaa/auth_service/internal/app/service"
	"maqhaa/auth_service/internal/config"
	"maqhaa/auth_service/internal/database"
	"maqhaa/auth_service/internal/interface/http/handler"
	"maqhaa/library/helper"
	"maqhaa/library/logging"
	"maqhaa/library/middleware"

	gRPCHandler "maqhaa/auth_service/internal/interface/grpc/handler"
	pb "maqhaa/auth_service/internal/interface/grpc/model"

	"google.golang.org/grpc"
	"gorm.io/gorm"
)

var db *gorm.DB
var authHandler *handler.AuthHandler
var userHandlerGrpc *gRPCHandler.UserHandler

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	tearDown()
	os.Exit(code)
}

func setup() {
	// Load testing configuration
	cfg, err := config.LoadConfig("../../cmd/config/config-test.yaml")
	if err != nil {
		panic(err)
	}

	// Set up a connection to the testing database
	db, err = database.NewDB(&cfg.Database)
	if err != nil {
		panic(err)
	}

	// Apply database migrations for tests
	_ = db.Migrator().DropTable(&entity.User{}, &entity.Client{})
	if err := db.AutoMigrate(&entity.Client{}, &entity.User{}); err != nil {
		panic(err)
	}

	logFolder := flag.String("log.file", "../../logs", "Logging file")

	flag.Parse()

	// set logging file
	//logging.OutputScreen = true
	//logging.Filename = *logFile
	logging.InitLogger()
	currentDate := time.Now().Format("2006-01-02")
	logFilePath := fmt.Sprintf("%s/app_test_%s.log", *logFolder, currentDate)

	// Create the log file if it doesn't exist
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Error creating log file:", err)
	}

	// Set the logrus output to the log file
	logging.Log.SetOutput(logFile)

	// Apply database migrations (if any)
	// You can use db.AutoMigrate(&YourModel{}) to automatically apply migrations

	userRepository := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepository)
	authHandler = handler.NewAuthHandler(authService)
	userHandlerGrpc = gRPCHandler.NewUserGRPCHandler(authService)

	go func() {
		// Create a gRPC server
		grpcServer := grpc.NewServer(grpc.UnaryInterceptor(middleware.LoggingInterceptor))

		// Register your gRPC service implementation
		pb.RegisterUserServer(grpcServer, userHandlerGrpc)

		// Start gRPC server on a specific port
		listen, err := net.Listen("tcp", "localhost:50051")
		if err != nil {
			fmt.Println("Error starting gRPC server: ", err)
			//log.Fatal("Error starting gRPC server: %s", err.Error())
		}

		defer listen.Close()

		if err := grpcServer.Serve(listen); err != nil {
			fmt.Println("Error serving gRPC server ", err)
			//log.Fatal("Error serving gRPC server: %s", err.Error())
		}
	}()
}

func clearDB(tables []string) {
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	for _, v := range tables {
		sqlDB.Exec("delete from " + v)
	}

}

func tearDown() {
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	defer sqlDB.Close()

}

func SampleUser(clientID uint) *entity.User {
	token, _ := helper.GenerateRandomString(16)
	return &entity.User{
		ID:           0,
		ClientID:     clientID,
		Username:     "sample",
		Password:     "rahasia",
		FullName:     "Sample User",
		Token:        token,
		Role:         1,
		TokenExpired: time.Now().Add(time.Hour), // Set a token expiration time
		IsActive:     true,
		CreatedAt:    time.Now(),
	}
}

func SampleUserCS(clientID uint, username string) *entity.User {
	token, _ := helper.GenerateRandomString(16)
	return &entity.User{
		ID:           0,
		ClientID:     clientID,
		Username:     username,
		Password:     "rahasia",
		FullName:     "Sample User",
		Token:        token,
		Role:         2,
		TokenExpired: time.Now().Add(time.Hour), // Set a token expiration time
		IsActive:     true,
		CreatedAt:    time.Now(),
	}
}

func SampleClient() *entity.Client {
	return &entity.Client{
		ID:          0,
		CompanyName: "Example Coffee",
		Email:       "info@examplecoffee.com",
		PhoneNumber: "+1234567890",
		Address:     "123 Main St, Cityville",
		OwnerName:   "John Doe",
		IsActive:    true,
		Token:       "JYA60sj03G6ii0LR3BfF", // Add a sample token if needed
		CreatedAt:   time.Now(),
	}
}
