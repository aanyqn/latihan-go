package main

import (
	"context"
	"fmt"
	"latihan-fiber/app/repository"
	"latihan-fiber/config"
	"latihan-fiber/database"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

var metodeBerbody = map[string]bool{
	fiber.MethodPost:  true,
	fiber.MethodPut:   true,
	fiber.MethodPatch: true,
}

func requireJSON(c *fiber.Ctx) error {
	if metodeBerbody[c.Method()] {
		ct := c.Get("Content-Type")
		if !strings.HasPrefix(ct, fiber.MIMEApplicationJSON) {
			return fail(c, fiber.StatusUnsupportedMediaType,
				"Content-Type harus application/json")
		}
	}
	return c.Next()
}

func main() {
	config.LoadEnv()

	pool, err := database.NewPool(context.Background())
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer pool.Close()

	userRepository := repository.NewUserRepository(pool)
	userHandler := NewUserHandler(userRepository)

	studentRepository := repository.NewStudentRepository(pool)
	studentHandler := NewStudentHandler(studentRepository)

	app := fiber.New(fiber.Config{
		AppName: "Praktikum Backend Lanjut - Pertemuan 2",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			status := fiber.StatusInternalServerError
			pesan := "There's an error on the server"
			if e, ok := err.(*fiber.Error); ok {
				status = e.Code
				pesan = e.Message
			}
			return fail(c, status, pesan)
		},
	})
	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${locals:requestid} ${method} ${path} ${status} ${latency}\n",
	}))
	app.Use(cors.New())
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	api := app.Group("/api/v1")

	api.Get("/health", func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.UserContext(), 2*time.Second)
		defer cancel()
		if err := pool.Ping(ctx); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Server berjalan, tapi database terputus",
				"error":   err.Error(),
			})
		}
		return ok(c, "server dan database berjalan normal", fiber.Map{
			"timestamp": time.Now(),
			"database":  "connected",
		})
	})

	// u := api.Group("/users", requireJSON)
	// u.Get("/", listUsers)
	// u.Get("/:id", getUser)
	// u.Post("/", createUser)
	// u.Put("/:id", replaceUser)
	// u.Patch("/:id", patchUser)
	// u.Delete("/:id", deleteUser)

	// s := api.Group("students", requireJSON)
	// s.Get("/", listStudents)
	// s.Get("/:id", getStudent)
	// s.Post("/", createStudent)
	// s.Put("/:id", replaceStudent)
	// s.Patch("/:id", patchStudent)
	// s.Delete("/:id", deleteStudent)
	
	s := api.Group("students", requireJSON)
	s.Get("/", studentHandler.ListStudents)
	s.Get("/:id", studentHandler.GetStudent)
	s.Post("/", studentHandler.CreateStudent)
	s.Put("/:id", studentHandler.ReplaceStudent)
	s.Patch("/:id", studentHandler.PatchStudent)
	s.Delete("/:id", studentHandler.DeleteStudent)

	u := api.Group("users", requireJSON)
	u.Get("/", userHandler.List)
	u.Get("/:id", userHandler.Get)
	u.Post("/", userHandler.Create)
	u.Put("/:id", userHandler.Replace)
	u.Patch("/:id", userHandler.Patch)
	u.Delete("/:id", userHandler.Delete)

	port := config.GetEnv("APP_PORT", "3000")
	log.Fatal(app.Listen(":" + port))

	app.Use(func(c *fiber.Ctx) error {
		return fail(c, fiber.StatusNotFound, "endpoint tidak ditemukan")
	})
	fmt.Println("Server berjalan di http://localhost:3000")
	log.Fatal(app.Listen(":3000"))
}
