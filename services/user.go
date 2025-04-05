package services

import (
	"context"
	db "golang-fiber-postgres-template/db/sqlc"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type UserService struct {
	db *db.DB
}

func NewUserService(db *db.DB) *UserService {
	return &UserService{db: db}
}

// ListUsers lists all existing users
//
//	@Summary      List users
//	@Description  get users
//	@Tags         users
//	@Accept       json
//	@Produce      json
//	@Success      200  {array} db.User
//	@Failure      500
//	@Router       /user/list [get]
//	@Security     BearerAuth
func (p *UserService) ListUsers(c *fiber.Ctx) error {
	users, err := p.db.Queries.ListUsers(context.Background())
	if err != nil {
		return c.Status(500).SendString("Failed to list users")
	}
	return c.Status(http.StatusOK).JSON(users)
}

// CreateDefaultUsers creates default users
//
//	@Summary      Create default users
//	@Description  create default users
//	@Tags         users
//	@Accept       json
//	@Produce      json
//	@Success      201
//	@Failure      500
//	@Router       /user/create-default [post]
//	@Security     BearerAuth
func (p *UserService) CreateDefaultUsers(c *fiber.Ctx) error {
	err := p.db.Queries.CreateDefaultUsers(context.Background())
	if err != nil {
		return c.Status(500).SendString("Failed to create default users")
	}
	return c.Status(http.StatusCreated).SendString("Default users created successfully")
}
