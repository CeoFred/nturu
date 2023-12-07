package main

import (
	"github.com/nturu/microservice-template/constants"
	"github.com/nturu/microservice-template/database"
	"github.com/nturu/microservice-template/internal/handlers"
	"github.com/nturu/microservice-template/internal/otp"
	"github.com/nturu/microservice-template/internal/routes"

	"context"
	"flag"
	"log"
	"os"

	apitoolkit "github.com/apitoolkit/apitoolkit-go"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	_ "github.com/nturu/microservice-template/docs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	_ "golang.org/x/text"
)

var (
	prod = flag.Bool("prod", false, "Enable prefork in Production")
)

// @title goFiber App
// @version 1.0
// @description Swagger API documentation for goFiber API
// @termsOfService http://swagger.io/terms/
// @contact.name Your Name
// @contact.email fiber@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:3009
// @BasePath /api/v1
func main() {

	logger, err := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(zapcore.InfoLevel),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig:    zap.NewProductionEncoderConfig(),
	}.Build()
	if err != nil {
		panic(err)
	}
	defer func() {
		err := logger.Sync()
		if err != nil {
			// Handle the error appropriately, such as logging or returning an error
			log.Println("Failed to sync logger:", err)
		}
	}()

	constant := constants.New()
	_ = otp.NewOTPManager()

	// Parse command-line flags
	flag.Parse()

	// Create fiber app
	app := fiber.New(fiber.Config{
		Prefork: *prod, // go run app.go -prod
	})

	ctx := context.Background()

	// Initialize the client using your apitoolkit.io generated apikey
	apitoolkitClient, err := apitoolkit.NewClient(ctx, apitoolkit.Config{APIKey: constant.APIToolkitKey})
	if err != nil {
		// Handle the error
		panic(err)
	}

	app.Get("/swagger/*", swagger.HandlerDefault) // default

	app.Get("/swagger/*", swagger.New(swagger.Config{ // custom
		URL:         "http://example.com/doc.json",
		DeepLinking: false,
		// Expand ("list") or Collapse ("none") tag groups by default
		DocExpansion: "none",
		// Prefill OAuth ClientId on Authorize popup
		OAuth: &swagger.OAuthConfig{
			AppName:  "OAuth Provider",
			ClientId: "21bb4edc-05a7-4afc-86f1-2e151e4ba6e2",
		},
		// Ability to change OAuth2 redirect uri location
		OAuth2RedirectUrl: "http://localhost:8080/swagger/oauth2-redirect.html",
	}))
	app.Static("/", "./static/public")

	// Middleware
	app.Use(recover.New())
	// Set the Zap logger as the Fiber logger
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("logger", logger)
		return c.Next()
	})

	app.Use(apitoolkitClient.FiberMiddleware)

	app.Use(func(c *fiber.Ctx) error {

		logger := c.Locals("logger").(*zap.Logger)

		// Log request details
		logger.Info("Request received",
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.String("ip", c.IP()),
			zap.Any("headers", c.Request()),
		)

		// Proceed to the next middleware or route handler
		err := c.Next()

		// Log response details
		logger.Info("Response sent",
			zap.Int("status", c.Response().StatusCode()),
			zap.Any("headers", c.Response()),
		)

		return err
	})
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, Authentication",
	}))

	dbConfig := database.Config{
		Host:     constant.DbHost,
		Port:     constant.DbPort,
		Password: constant.DbPassword,
		User:     constant.DbUser,
		DBName:   constant.DbName,
	}

	database.Connect(&dbConfig)

	database.RunManualMigration(database.DB)

	// Bind routes
	routes.Routes(app, database.DB)

	// Handle not founds
	app.Use(handlers.NotFound)

	port := os.Getenv("PORT")
	if port == "" {
		port = constant.Port
	}

	// Listen on port set in .env
	log.Fatal(app.Listen(":" + port))
}
