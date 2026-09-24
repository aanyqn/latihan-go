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
	Pool           *pgxpool.Pool
	Permissions    *helper.PermissionSet
	JWT            *helper.JWTManager
	UserService    *user.UserService
	AuthService    *auth.AuthService
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

	perms := deps.Permissions
	users.Get("/", middleware.RequirePermission(perms, "user:list"), deps.UserService.List)
	users.Post("/", middleware.RequirePermission(perms, "user:update:any"), deps.UserService.Create)
	users.Delete("/:id", middleware.RequirePermission(perms, "user:delete"), deps.UserService.Delete)
	users.Patch("/:id/role", middleware.RequirePermission(perms, "role:assign"), deps.UserService.AssignRole)
	users.Get("/:id", deps.UserService.Get)
	users.Put("/:id", deps.UserService.Replace)
	users.Patch("/:id", deps.UserService.Patch)

	student := api.Group("/students", middleware.RequireJSON, middleware.RequireAuth(deps.JWT))
	student.Get("/", middleware.RequirePermission(perms, "student:list"), deps.StudentHandler.List)
	student.Post("/", middleware.RequirePermission(perms, "student:create"), deps.StudentHandler.Create)
	student.Delete("/:id", middleware.RequirePermission(perms, "student:delete"), deps.StudentHandler.Delete)
	student.Get("/:id", deps.StudentHandler.Get)
	student.Put("/:id", deps.StudentHandler.Replace)
	student.Patch("/:id", deps.StudentHandler.Patch)

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
