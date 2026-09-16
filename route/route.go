package route

import (
	"context"
	"latihan-fiber/app/handler"
	// "latihan-fiber/app/service/achievement"
	"latihan-fiber/app/service/auth"
	"latihan-fiber/app/service/user"
	"latihan-fiber/helper"
	"latihan-fiber/middleware"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Dependencies struct {
	Pool        *pgxpool.Pool
	JWT         *helper.JWTManager
	UserService *user.UserService
	AuthService *auth.AuthService
	StudentHandler *handler.StudentHandler
}

func Register(app *fiber.App, deps Dependencies) {
	api := app.Group("/api/v1")

	api.Get("/health", healthCheck(deps.Pool))

	auth := api.Group("/auth", middleware.RequireJSON)
	auth.Post("/register", deps.AuthService.Register)
	auth.Post("/login", middleware.LoginRateLimiter(), deps.AuthService.Login)
	auth.Post("/refresh", deps.AuthService.Refresh)
	auth.Post("/logout", deps.AuthService.Logout)
	auth.Get("/me", middleware.RequireAuth(deps.JWT), deps.AuthService.Me)

	users := api.Group("/users",
		middleware.RequireJSON, middleware.RequireAuth(deps.JWT))
	users.Get("/", deps.UserService.List)
	users.Get("/:id", deps.UserService.Get)
	users.Post("/", deps.UserService.Create)
	users.Put("/:id", deps.UserService.Replace)
	users.Patch("/:id", deps.UserService.Patch)
	users.Delete("/:id", deps.UserService.Delete)

	student := api.Group("/students", middleware.RequireJSON)
	student.Get("/", deps.StudentHandler.List)
	student.Get("/:id", deps.StudentHandler.Get)
	student.Post("/", deps.StudentHandler.Create)
	student.Put("/:id", deps.StudentHandler.Replace)
	student.Patch("/:id", deps.StudentHandler.Patch)
	student.Delete("/:id", deps.StudentHandler.Delete)

	// achievement := api.Group("/achievements", middleware.RequireJSON)
	// achievement.Get("/:id", achievementService.Get)
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
