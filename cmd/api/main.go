package main

import (
	"awesomeProject/internal/config"
	v1 "awesomeProject/internal/handler"
	postgres "awesomeProject/internal/repository"
	"awesomeProject/internal/service"
	"awesomeProject/pkg/database"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	cfg := config.Load()

	db, err := database.NewPostgresConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	bookRepo := postgres.NewBookRepository(db)

	bookService := service.NewBookService(bookRepo)

	bookHandler := v1.NewBookHandler(bookService)

	app := fiber.New()

	app.Use(logger.New())

	api := app.Group("/api/v1")
	{
		books := api.Group("/books")
		{
			books.Post("/", bookHandler.CreateBook)
			books.Get("/", bookHandler.GetAllBooks)
			books.Get("/:id", bookHandler.GetBook)
			books.Put("/:id", bookHandler.UpdateBook)
			books.Delete("/:id", bookHandler.DeleteBook)
		}
	}

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	log.Printf("Server starting on port %s", cfg.AppPort)
	log.Fatal(app.Listen(cfg.AppPort))
}
