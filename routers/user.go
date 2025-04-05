package routers

import (
	middlewares "golang-fiber-postgres-template/middlewares"
	services "golang-fiber-postgres-template/services"

	"github.com/gofiber/fiber/v2"
)

func UserRouter(api fiber.Router, svc *services.UserService) {
	api.Get("/list", svc.ListUsers)
	api.Post("/create-default", middlewares.AuthMiddleware(), svc.CreateDefaultUsers)
}
