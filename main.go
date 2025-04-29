package main

import (
	_ "golang-fiber-postgres-template/docs"

	. "golang-fiber-postgres-template/setup"
)

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
	// Initialize the app
	app, dbConn := SetupApp()

	// Start the app asynchronously
	StartApp(app, dbConn)
}
