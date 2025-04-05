package tests

import (
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/suite"

	db "golang-fiber-postgres-template/db/sqlc"
	. "golang-fiber-postgres-template/setup"
)

type BaseTestSuite struct {
	suite.Suite
	App *fiber.App
}

func (suite *BaseTestSuite) SetupSuite() {
	// Initialize the app
	var dbConn *db.DB
	suite.App, dbConn = SetupApp()

	// Start the app asynchronously
	StartApp(suite.App, dbConn)

	// Wait briefly to ensure the app is running before tests
	time.Sleep(3 * time.Second)
}

func (suite *BaseTestSuite) TearDownSuite() {
	// Shutdown the app after tests
	_ = suite.App.Shutdown()
}

func RunTestSuite(t *testing.T, testSuite suite.TestingSuite) {
	suite.Run(t, testSuite)
}
