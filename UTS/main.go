package main

import (
	"context"
	"siakad/app/handler"
	"siakad/app/repository"
	"siakad/app/service/auth"
	"siakad/app/service/course"
	"siakad/app/service/enrollment"
	"siakad/app/service/student"

	// "siakad/app/service/user"
	"log/slog"
	"os"
	"os/signal"
	"siakad/config"
	"siakad/database"
	"siakad/helper"
	"siakad/route"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
)

var metodeBerbody = map[string]bool{
	fiber.MethodPost:  true,
	fiber.MethodPut:   true,
	fiber.MethodPatch: true,
}

const minSecretLength = 32

func main() {
	config.LoadEnv()
	logger := config.NewLogger()

	jwtSecret := config.GetEnv("JWT_SECRET", "")
	if len(jwtSecret) < minSecretLength {
		logger.Error("JWT_SECRET not available or too short",
			slog.Int("minimum_character", minSecretLength))
		os.Exit(1)
	}

	pool, err := database.NewPool(context.Background())
	if err != nil {
		logger.Error("failed to connect with database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer pool.Close()
	jwtManager := helper.NewJWTManager(
		jwtSecret,
		config.GetEnv("JWT_ISSUER", "siakad"),
		time.Duration(config.GetEnvInt("JWT_ACCESS_TTL_MINUTES", 15))*time.Minute,
	)

	roleRepository := repository.NewRoleRepository(pool)
	rawPermissions, err := roleRepository.LoadPermissions(context.Background())
	if err != nil {
		logger.Error("Failed to load permissions", slog.String("error", err.Error()))
		os.Exit(1)
	}
	permissions := helper.NewPermissionSet(rawPermissions)
	logger.Info("permission dimuat", slog.Any("roles", permissions.KnownRoles()), slog.Any("all_permissions", permissions.KnownPermissions()))

	userRepository := repository.NewUserRepository(pool)
	studentRepository := repository.NewStudentRepository(pool)
	studentService := student.NewStudentService(studentRepository, userRepository)
	tokenRepository := repository.NewTokenRepository(pool)
	studentHandler := handler.NewStudentHandler(studentService, permissions)
	courseRepository := repository.NewCourseRepository(pool)
	courseService := course.NewCourseService(courseRepository)
	courseHandler := handler.NewCourseHandler(courseService, permissions)
	enrollmentRepository := repository.NewEnrollmentRepository(pool)
	enrollmentService := enrollment.NewEnrollmentService(enrollmentRepository, studentRepository)
	enrollmnetHandler := handler.NewEnrollmentHandler(enrollmentService, permissions)

	authService := auth.NewAuthService(
		userRepository, tokenRepository, jwtManager,
		time.Duration(config.GetEnvInt("JWT_REFRESH_TTL_DAYS", 7))*24*time.Hour, permissions, studentRepository,
	)
	app := config.NewApp(logger, route.Dependencies{
		Pool:           pool,
		Permissions:    permissions,
		JWT:            jwtManager,
		AuthService:    authService,
		StudentHandler: studentHandler,
		CourseHandler: courseHandler,
		EnrollmentHandler: enrollmnetHandler,
	})

	port := config.GetEnv("APP_PORT", "3000")

	go func() {
		if err := app.Listen(":" + port); err != nil {
			logger.Error("Server has stopped", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	logger.Info("Server is Running", slog.String("port", port))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Stop receiving signal, closing")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := app.ShutdownWithContext(ctx); err != nil {
		logger.Error("Failed to stop the server",
			slog.String("error", err.Error()))
	}
	logger.Info("Server has stop flawlessly")
}
