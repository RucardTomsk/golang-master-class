package main

import (
	"UGUAPI/api"
	"UGUAPI/auth"
	"UGUAPI/cmd/controllers"
	"UGUAPI/cmd/services"
	"UGUAPI/cmd/storage/dao"
	"UGUAPI/cmd/storage/migration"
	"UGUAPI/server"
	"context"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type AdminCord struct {
	ID       uuid.UUID
	Name     string
	Email    string
	Password string
}

type AuthConfig struct {
	Admin AdminCord
	Salt  string
}

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

type ServerConfig struct {
	Host string
	Port string
}

type JWTManagerConfig struct {
	key        string
	timeToLive time.Duration
}

func main() {

	config := AuthConfig{
		Admin: AdminCord{
			ID:       uuid.MustParse("94a123a8-711d-11ee-a2d8-0242c0a8f005"),
			Name:     "testAdmin",
			Email:    "admin@example.com",
			Password: "admin",
		},
		Salt: "ameiobtiosrn234234sdvwse",
	}

	dbConfig := DBConfig{
		Host:     "92.63.64.241",
		Port:     35446,
		Password: "ugrasu-api",
		User:     "ugrasu-api",
		Name:     "ugrasu-api",
	}

	jwtManagerConfig := JWTManagerConfig{
		key:        "asdfasdfasdf",
		timeToLive: time.Hour,
	}

	serverConfig := ServerConfig{
		Host: "localhost",
		Port: "8080",
	}

	hasher := auth.NewHasher(config.Salt)

	// init connections
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		dbConfig.Host, dbConfig.User, dbConfig.Password, dbConfig.Name, dbConfig.Port)
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		panic(err)
	}

	adminPassword, err := hasher.Hash(config.Admin.Password)
	if err != nil {
		panic(err)
	}

	if err := migration.Migration(
		db,
		config.Admin.ID,
		config.Admin.Name,
		config.Admin.Email,
		adminPassword); err != nil {
		panic(err)
	}

	jwtManager, err := auth.NewJWTManager(jwtManagerConfig.key, jwtManagerConfig.timeToLive)
	if err != nil {
		panic(err)
	}

	userStorage := dao.NewUserStorage(db)
	authService := services.NewAuthService(userStorage, hasher, jwtManager)
	authController := controllers.NewAuthController(authService)

	server := new(server.Server)
	go func() {
		if err := server.Run(serverConfig.Host, serverConfig.Port, api.InitRoutes(authController, jwtManager)); err != nil {
			panic(err)
		}
	}()

	// handle signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	if err := server.Shutdown(context.Background()); err != nil {
		panic(err)
	}
}
