package main

import (
	"fmt"
	"log"
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

// @title           Template Go (Fiber) Backend API
// @version         1.0
// @description     This is a template backend REST API written in Go with the Fiber framework using SQLC.

// @contact.name   Petar Vico
// @contact.url    https://google.com
// @contact.email  test@gmail.com

// @host 	  localhost:3500
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter the token with the `Bearer ` prefix, e.g., "Bearer abcde12345"

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	// ==================== Fiber Setup =================
	app := fiber.New()

	allowedOriginRegex := regexp.MustCompile(`^(https?:\/\/)([^/]+\.)?(127\.0\.0\.1|localhost.*)$`)

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
	defer dbConn.Close()
	fmt.Println("Connected to the database successfully!")

	DbConnectionHealtcheck(dbConn)

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

	// ==================== Start Fiber App ====================
	app.Listen(":3500")

	// Close the database connection
	fmt.Println("Database connection closed.")

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
