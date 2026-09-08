package main

import (
	"context"
	"latihan-fiber/app/handler"
	"latihan-fiber/app/repository"
	"latihan-fiber/app/service/student"
	"latihan-fiber/app/service/user"
	"latihan-fiber/config"
	"latihan-fiber/database"
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

func main() {
	config.LoadEnv()
	logger := config.NewLogger()

	pool, err := database.NewPool(context.Background())
	if err != nil {
		logger.Error("gagal terhubung ke database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer pool.Close()

	userRepository := repository.NewUserRepository(pool)
	userService := user.NewUserService(userRepository)

	studentRepository := repository.NewStudentRepository(pool)
	studentService := student.NewStudentService(studentRepository)
	studentHandler := handler.NewStudentHandler(studentService)


	app := config.NewApp(logger, pool, userService, studentHandler)

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
