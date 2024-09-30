package app

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

type AppConfig struct {
	AppName string
	AppEnv  string
	AppPort string
}

type DBConfig struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
	DBDriver   string
}

func (server *Server) Initialize(appConfig AppConfig, dbConfig DBConfig) {
	fmt.Println("welcome to " + appConfig.AppName)

	server.InitializeDB(dbConfig)
	server.InitializeRoutes()
}

func (server *Server) Run(addr string) {
	fmt.Printf("Listening on http://localhost%s \n", addr)
	log.Fatal((http.ListenAndServe(addr, server.Router)))
}

func (server *Server) InitializeDB(dbConfig DBConfig) {
	var err error

	// Running Database
	if dbConfig.DBDriver == "mysql" {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbConfig.DBUser, dbConfig.DBPassword, dbConfig.DBHost, dbConfig.DBPort, dbConfig.DBName)
		server.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

		fmt.Println("Database MySQL Running")
	} else {
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", dbConfig.DBHost, dbConfig.DBUser, dbConfig.DBPassword, dbConfig.DBName, dbConfig.DBPort)
		server.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

		fmt.Println("Database Postgre Running")
	}

	if err != nil {
		panic("Failed on connecting to the database server")
	}

	for _, model := range RegisterModels() {
		err = server.DB.Debug().AutoMigrate(model.Model)

		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Database migrated Successfully")
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func Run() {
	var server = &Server{}
	var appConfig = AppConfig{}
	var dbConfig = DBConfig{}

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading.env file")
	}

	appConfig.AppName = getEnv("APP_NAME", "Gotoko")
	appConfig.AppEnv = getEnv("APP_ENV", "development")
	appConfig.AppPort = getEnv("APP_PORT", "9000")

	if os.Getenv("DB_DRIVER") == "mysql" {
		dbConfig.DBHost = getEnv("DB_HOST_MYSQL", "localhost")
		dbConfig.DBUser = getEnv("DB_USER_MYSQL", "root")
		dbConfig.DBPassword = getEnv("DB_PASSWORD_MYSQL", "")
		dbConfig.DBName = getEnv("DB_NAME_MYSQL", "gotoko")
		dbConfig.DBPort = getEnv("DB_PORT_MYSQL", "3306")
		dbConfig.DBDriver = getEnv("DB_DRIVER", "mysql")
	} else {
		dbConfig.DBHost = getEnv("DB_HOST_POSTGRES", "localhost")
		dbConfig.DBUser = getEnv("DB_USER_POSTGRES", "root")
		dbConfig.DBPassword = getEnv("DB_PASSWORD_POSTGRES", "")
		dbConfig.DBName = getEnv("DB_NAME_POSTGRES", "gotoko")
		dbConfig.DBPort = getEnv("DB_PORT_POSTGRES", "3306")
		dbConfig.DBDriver = getEnv("DB_DRIVER", "postgres")
	}

	server.Initialize(appConfig, dbConfig)
	server.Run(":" + appConfig.AppPort)
}
