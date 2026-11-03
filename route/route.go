package route

import (
	"context"
	"latihan-fiber/app/handler"

	// "latihan-fiber/app/service/achievement"
	"latihan-fiber/app/service/auth"
	"latihan-fiber/app/service/course"
	"latihan-fiber/app/service/portfolio"
	"latihan-fiber/app/service/review"
	"latihan-fiber/app/service/user"
	"latihan-fiber/helper"
	"latihan-fiber/middleware"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Dependencies struct {
	Pool             *pgxpool.Pool
	Permissions      *helper.PermissionSet
	JWT              *helper.JWTManager
	UserService      *user.UserService
	AuthService      *auth.AuthService
	StudentHandler   *handler.StudentHandler
	CourseService    *course.CourseService
	Mongo            *mongo.Client
	ReviewService    *review.ReviewService
	PortfolioService *portfolio.PortfolioService
}

func Register(app *fiber.App, deps Dependencies) {
	api := app.Group("/api/v1")

	api.Get("/health", healthCheck(deps.Pool, deps.Mongo))

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
	student.Get("/:id/portfolio", deps.PortfolioService.Get)
	student.Put("/:id/portfolio", deps.PortfolioService.Replace)
	student.Post("/:id/portfolio/projects", deps.PortfolioService.AddProject)
	student.Delete("/:id/portfolio/projects/:projectId", deps.PortfolioService.RemoveProject)

	// achievement := api.Group("/achievements", middleware.RequireJSON)
	// achievement.Get("/:id", achievementService.Get)

	courses := api.Group("/courses",
		middleware.RequireJSON,
		middleware.RequireAuth(deps.JWT))
	courses.Get("/", deps.CourseService.List)
	courses.Get("/mine", deps.CourseService.Mine)
	courses.Get("/:id", deps.CourseService.Get)
	courses.Post("/:id/enroll", deps.CourseService.Enroll)
	courses.Delete("/:id/enroll", deps.CourseService.Drop)
	courses.Get("/:id/reviews", deps.ReviewService.ListByCourse)
	courses.Get("/:id/reviews/summary", deps.ReviewService.Summary)
	courses.Post("/:id/reviews", deps.ReviewService.Create)

	enrollments := api.Group("/enrollments",
		middleware.RequireJSON,
		middleware.RequireAuth(deps.JWT))
	enrollments.Post("/swap", deps.CourseService.Swap)

	reviews := api.Group("/reviews",
		middleware.RequireJSON,
		middleware.RequireAuth(deps.JWT))
	reviews.Get("/:id", deps.ReviewService.Get)
	reviews.Patch("/:id", deps.ReviewService.Update)
	reviews.Delete("/:id", deps.ReviewService.Delete)
	reviews.Post("/:id/replies", deps.ReviewService.Reply)
	reviews.Post("/:id/helpful", deps.ReviewService.MarkHelpful)

	// --- portofolio mahasiswa (MongoDB) ---
	portfolios := api.Group("/portfolios",
		middleware.RequireJSON,
		middleware.RequireAuth(deps.JWT))
	portfolios.Get("/stats/technologies", deps.PortfolioService.TopTechnologies)
}

func healthCheck(pool *pgxpool.Pool, mongoClient *mongo.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.UserContext(), 2*time.Second)
		defer cancel()
		if err := pool.Ping(ctx); err != nil {
			return helper.ServiceUnavailable("Failed to reach PostgreSQL")
		}
		if err := mongoClient.Ping(ctx, nil); err != nil {
			return helper.ServiceUnavailable("Failed to reach MongoDB")
		}
		return helper.Success(c, fiber.StatusOK, "server, PostgreSQL, dan MongoDB running", nil)
	}
}
