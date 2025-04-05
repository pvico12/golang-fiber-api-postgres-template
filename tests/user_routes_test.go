package tests

import (
	"encoding/json"
	db "golang-fiber-postgres-template/db/sqlc"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserRoutesTestSuite(t *testing.T) {
	RunTestSuite(t, new(UserRoutesTestSuite))
}

type UserRoutesTestSuite struct {
	BaseTestSuite
}

func (suite *UserRoutesTestSuite) TestExecutionInOrder() {
	suite.T().Run("TestCreateDefaultUsers_noAuth", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/user/create-default", nil)
		resp, err := suite.App.Test(req)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), http.StatusUnauthorized, resp.StatusCode)
	})

	suite.T().Run("TestCreateDefaultUsers_withAuth", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/user/create-default", nil)
		req.Header.Set("Authorization", "Bearer test")
		resp, err := suite.App.Test(req)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)
	})

	suite.T().Run("TestListUsers", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/user/list", nil)
		resp, err := suite.App.Test(req)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

		// Check that the user list length is 0
		body, _ := io.ReadAll(resp.Body)
		var users []db.User
		json.Unmarshal(body, &users)
		assert.Greater(suite.T(), len(users), 0)
	})
}
