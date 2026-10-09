package route

import (
	"context"
	"siakad/app/handler"
	"siakad/app/service/auth"
	"siakad/helper"
	"siakad/middleware"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Dependencies struct {
	Pool              *pgxpool.Pool
	Permissions       *helper.PermissionSet
	JWT               *helper.JWTManager
	AuthService       *auth.AuthService
	StudentHandler    *handler.StudentHandler
	CourseHandler     *handler.CourseHandler
	EnrollmentHandler *handler.EnrollmentHandler
}

func Register(app *fiber.App, deps Dependencies) {
	api := app.Group("/api/v1")

	api.Get("/health", healthCheck(deps.Pool))

	auth := api.Group("/auth", middleware.RequireJSON)
	auth.Post("/login", middleware.LoginRateLimiter(), deps.AuthService.Login)
	auth.Post("/refresh", deps.AuthService.Refresh)
	auth.Post("/logout", deps.AuthService.Logout)
	auth.Get("/me", middleware.RequireAuth(deps.JWT), deps.AuthService.Me)

	perms := deps.Permissions

	student := api.Group("/students", middleware.RequireJSON, middleware.RequireAuth(deps.JWT))
	student.Get("/", middleware.RequirePermission(perms, "student:list"), deps.StudentHandler.List)
	student.Post("/", middleware.RequirePermission(perms, "student:create"), deps.StudentHandler.Create)
	student.Delete("/:id", middleware.RequirePermission(perms, "student:delete:any"), deps.StudentHandler.Delete)
	student.Get("/:id", deps.StudentHandler.Get)
	student.Put("/:id", deps.StudentHandler.Replace)

	course := api.Group("/courses", middleware.RequireJSON, middleware.RequireAuth(deps.JWT))
	course.Get("/", middleware.RequirePermission(perms, "course:list"), deps.CourseHandler.List)

	enrollment := api.Group("/enrollments", middleware.RequireJSON, middleware.RequireAuth(deps.JWT))
	enrollment.Post("/", middleware.RequirePermission(perms, "enrollment:create"), deps.EnrollmentHandler.Create)
	enrollment.Delete("/:id", middleware.RequirePermission(perms, "enrollment:delete:own"), deps.EnrollmentHandler.Delete)
}

func healthCheck(pool *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.UserContext(), 2*time.Second)
		defer cancel()
		if err := pool.Ping(ctx); err != nil {
			return helper.Fail(c, fiber.StatusServiceUnavailable,
				"Can't connext database")
		}
		return helper.Success(c, fiber.StatusOK, "server and database is running", nil)
	}
}
