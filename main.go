package main

import (
	"context"
	"latihan-fiber/app/handler"
	"latihan-fiber/app/repository"
	// "latihan-fiber/app/service/achievement"
	"latihan-fiber/app/service/auth"
	"latihan-fiber/app/service/student"
	"latihan-fiber/app/service/user"
	"latihan-fiber/config"
	"latihan-fiber/database"
	"latihan-fiber/helper"
	"latihan-fiber/route"
	"log/slog"
	"os"
	"os/signal"
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
		logger.Error("JWT_SECRET tidak diisi atau terlalu pendek",
			slog.Int("minimal_karakter", minSecretLength))
		os.Exit(1)
	}

	pool, err := database.NewPool(context.Background())
	if err != nil {
		logger.Error("gagal terhubung ke database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer pool.Close()
	jwtManager := helper.NewJWTManager(
		jwtSecret,
		config.GetEnv("JWT_ISSUER", "praktikum-backend"),
		time.Duration(config.GetEnvInt("JWT_ACCESS_TTL_MINUTES", 15))*time.Minute,
	)

	userRepository := repository.NewUserRepository(pool)
	studentRepository := repository.NewStudentRepository(pool)
	studentService := student.NewStudentService(studentRepository)
	tokenRepository := repository.NewTokenRepository(pool)
	userService := user.NewUserService(userRepository)
	studentHandler := handler.NewStudentHandler(studentService)
	authService := auth.NewAuthService(
		userRepository, tokenRepository, jwtManager,
		time.Duration(config.GetEnvInt("JWT_REFRESH_TTL_DAYS", 7))*24*time.Hour,
	)
	app := config.NewApp(logger, route.Dependencies{
		Pool:        pool,
		JWT:         jwtManager,
		UserService: userService,
		AuthService: authService,
		StudentHandler: studentHandler,
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
