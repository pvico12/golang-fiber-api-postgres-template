package setup

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"

	db "golang-fiber-postgres-template/db/sqlc"
	"golang-fiber-postgres-template/docs"
	_ "golang-fiber-postgres-template/docs"
	routers "golang-fiber-postgres-template/routers"
	services "golang-fiber-postgres-template/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	_ "github.com/lib/pq"
)

type BackendApp struct {
	UserService *services.UserService
}

func SetupApp() (*fiber.App, *db.DB) {
	// ==================== Fiber Setup =================
	app := fiber.New()

	allowedOriginRegex := regexp.MustCompile(`^(https?:\/\/)([^/]+\.)?(localhost.*)$`)

	app.Use(cors.New(cors.Config{
		AllowHeaders: "Origin,Content-Type,Accept,Content-Length,Accept-Language,Authorization,Accept-Encoding,Connection,Access-Control-Allow-Origin",
		AllowOriginsFunc: func(origin string) bool {
			// If there's no origin or if the origin matches our regex, allow it
			return origin == "" || allowedOriginRegex.MatchString(origin)
		},
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))

	fmt.Println("Launching server!")

	// ==================== Database Setup =================
	dbConn, err := SetupDatabaseConnection()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	fmt.Println("Connected to the database successfully!")

	DbConnectionHealtcheck(dbConn)

	// cleanup database if in testing mode
	testing_mode := os.Getenv("TESTING_MODE")
	if testing_mode == "true" {
		fmt.Println("Running in testing mode, re-initializing database.")
		ReInitDatabaseForTesting(dbConn)
	} else {
		fmt.Println("Running in production mode, skipping database re-initialization.")
	}

	// ==================== Backend Setup =================
	backendApplication := BackendApp{
		UserService: services.NewUserService(dbConn),
	}

	// healthcheck endpoint
	app.Get("/healthcheck", healthcheck)

	// ==================== Routes Setup =================
	app.Route("/user", func(router fiber.Router) {
		routers.UserRouter(router, backendApplication.UserService)
	})

	// ====================== Swagger Setup =================
	docs.SwaggerInfo.Host = "localhost:3500"
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Remove the blocking app.Listen call
	fmt.Println("App setup complete.")
	return app, dbConn
}

func StartApp(app *fiber.App, dbConn *db.DB) {
	// Start the app asynchronously
	fmt.Println("Starting server on port 3500...")
	testingMode := os.Getenv("TESTING_MODE")

	if testingMode == "true" {
		// Run asynchronously for testing
		go func() {
			defer dbConn.Close()
			if err := app.Listen(":3500"); err != nil {
				log.Fatalf("Failed to start the server: %v", err)
			}
		}()
	} else {
		// Run indefinitely for production
		defer dbConn.Close()
		if err := app.Listen(":3500"); err != nil {
			log.Fatalf("Failed to start the server: %v", err)
		}
	}
}

func SwaggerHandler(c *fiber.Ctx) error {
	return swagger.HandlerDefault(c)
}

func SetupDatabaseConnection() (*db.DB, error) {
	dbConn, err := db.NewDB()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %w", err)
	}
	return dbConn, nil
}

func DbConnectionHealtcheck(db *db.DB) {
	err := db.Conn.Ping()
	if err != nil {
		log.Fatalf("Failed to ping the database: %v", err)
	}
	fmt.Println("Pinged database successfully!")
}

func healthcheck(c *fiber.Ctx) error {
	return c.SendString("OK")
}

func ReInitDatabaseForTesting(dbConn *db.DB) {
	// Create a context
	ctx := context.Background()

	// Retrieve all table names from the database
	tableNames, err := dbConn.Queries.GetAllTableNames(ctx)
	if err != nil {
		log.Fatalf("Failed to get table names: %v", err)
	}
	log.Printf("Retrieved table names after re-initialization: %v", fmt.Sprintf("%q", tableNames))

	// Drop all tables
	for _, tableName := range tableNames {
		if _, err := dbConn.Conn.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", tableName)); err != nil {
			log.Fatalf("Failed to drop table %s: %v", tableName, err)
		}
	}
	log.Println("All tables dropped successfully!")

	// Create all tables
	schemaFile, err := os.ReadFile("../db/schema.sql")
	if err != nil {
		log.Fatalf("Failed to read schema.sql file: %v", err)
	}
	// execute the file contents
	_, err = dbConn.Conn.Exec(string(schemaFile))
	if err != nil {
		log.Fatalf("Failed to execute schema.sql: %v", err)
	}
	log.Println("Database re-initialized successfully!")

	// get all table names again
	tableNames, err = dbConn.Queries.GetAllTableNames(ctx)
	if err != nil {
		log.Fatalf("Failed to get table names: %v", err)
	}
	log.Printf("Retrieved table names after re-initialization: %v", fmt.Sprintf("%q", tableNames))

}
