package route

import (
	"context"
	"latihan-fiber/app/handler"
	"latihan-fiber/app/service/user"
	"latihan-fiber/helper"
	"latihan-fiber/middleware"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Register(app *fiber.App, pool *pgxpool.Pool, userService *user.UserService, studentHandler *handler.StudentHandler) {
	api := app.Group("/api/v1")

	api.Get("/health", healthCheck(pool))

	users := api.Group("/users", middleware.RequireJSON)
	users.Get("/", userService.List)
	users.Get("/:id", userService.Get)
	users.Post("/", userService.Create)
	users.Put("/:id", userService.Replace)
	users.Patch("/:id", userService.Patch)
	users.Delete("/:id", userService.Delete)

	student := api.Group("/students", middleware.RequireJSON)
	student.Get("/", studentHandler.List)
	student.Get("/:id", studentHandler.Get)
	student.Post("/", studentHandler.Create)
	student.Put("/:id", studentHandler.Replace)
	student.Patch("/:id", studentHandler.Patch)
	student.Delete("/:id", studentHandler.Delete)
}

func healthCheck(pool *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.UserContext(), 2*time.Second)
		defer cancel()
		if err := pool.Ping(ctx); err != nil {
			return helper.Fail(c, fiber.StatusServiceUnavailable,
				"database tidak dapat dihubungi")
		}
		return helper.Success(c, fiber.StatusOK, "server dan database berjalan", nil)
	}
}
