package main

import (
	"log"

	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/kayprogrammer/socialnet-v6/config"
	"github.com/kayprogrammer/socialnet-v6/database"
	"github.com/kayprogrammer/socialnet-v6/initials"
	"github.com/kayprogrammer/socialnet-v6/routes"

	_ "github.com/kayprogrammer/socialnet-v6/docs"
)

// @title SOCIALNET API
// @version 6.0
// @description.markdown api
// @Accept json
// @Produce json
// @BasePath  /api/v6
// @Security BearerAuth
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type 'Bearer jwt_string' to correctly set the API Key
func main() {
	cfg := config.GetConfig()
	db := database.ConnectDb(cfg)
	sqlDb, _ := db.DB()
	initials.CreateInitialData(cfg, db)

	app := fiber.New()

	// CORS config
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORSAllowedOrigins,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, Access-Control-Allow-Origin, Content-Disposition",
		AllowCredentials: true,
		AllowMethods:     "GET, POST, PUT, PATCH, DELETE, OPTIONS",
	}))

	// Swagger Config
	swaggerCfg := swagger.Config{
		FilePath: "./docs/swagger.json",
		Path:     "/",
		Title:    "SOCIALNET API Documentation",
		CacheAge: 1,
	}

	app.Use(swagger.New(swaggerCfg))

	// Register Routes & Sockets
	app.Use("/ws", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	routes.SetupRoutes(app, db)
	defer sqlDb.Close()
	log.Fatal(app.Listen(":8000"))
}
