package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/kayprogrammer/socialnet-v6/config"
	"github.com/kayprogrammer/socialnet-v6/database"
	"github.com/kayprogrammer/socialnet-v6/routes"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func CreateSingleTable(db *gorm.DB, model interface{}) {
	db.AutoMigrate(&model)
}

func DropAndCreateSingleTable(db *gorm.DB, model interface{}) {
	db.Migrator().DropTable(&model)
	db.AutoMigrate(&model)
}

func connectToTestDatabase(t *testing.T, dsn string) *gorm.DB {
	t.Logf("Connecting to database....")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatal("Failed to connect to the database! \n", err.Error())
		os.Exit(2)
	}
	t.Logf("Connected to the database successfully")

	// Add UUID extension
	result := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	if result.Error != nil {
		log.Fatal("failed to create extension: " + result.Error.Error())
	}
	return db
}

func SetupTestDatabase(t *testing.T) *gorm.DB {
	cfg := config.GetConfig(true)

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s TimeZone=%s",
		cfg.PostgresServer,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.TestPostgresDB,
		cfg.PostgresPort,
		"UTC",
	)
	return connectToTestDatabase(t, dsn)
}

func CloseTestDatabase(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get database connection: " + err.Error())
	}
	if err := sqlDB.Close(); err != nil {
		log.Fatal("Failed to close database connection: " + err.Error())
	}
}

func Setup(t *testing.T, app *fiber.App) *gorm.DB {
	os.Setenv("ENVIRONMENT", "TESTING")

	// Set up the test database
	db := SetupTestDatabase(t)

	routes.SetupRoutes(app, db)
	t.Logf("Making Database Migrations....")
	database.DropTables(db)
	database.CreateTables(db)
	t.Logf("Database Migrations Made successfully")
	return db
}

func ParseResponseBody(t *testing.T, b io.ReadCloser) interface{} {
	body, _ := io.ReadAll(b)
	// Parse the response body as JSON
	responseBody := make(map[string]interface{})
	err := json.Unmarshal(body, &responseBody)
	if err != nil {
		t.Errorf("error parsing response body as JSON: %s", err)
	}
	return responseBody
}

func ProcessTestBody(t *testing.T, app *fiber.App, url string, method string, body interface{}, access ...string) *http.Response {
	// Marshal the test data to JSON
	requestBytes, err := json.Marshal(body)
	requestBody := bytes.NewReader(requestBytes)
	assert.Nil(t, err)
	req := httptest.NewRequest(method, url, requestBody)
	req.Header.Set("Content-Type", "application/json")
	if access != nil {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access[0]))
	}
	res, err := app.Test(req)
	if err != nil {
		log.Println(err)
	}
	return res
}
