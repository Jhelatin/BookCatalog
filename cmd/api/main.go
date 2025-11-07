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
	authorRepo := postgres.NewAuthorRepository(db)

	bookService := service.NewBookService(bookRepo)
	authorService := service.NewAuthorService(authorRepo, bookRepo)

	bookHandler := v1.NewBookHandler(bookService)
	authorHandler := v1.NewAuthorHandler(authorService)

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

		authors := api.Group("/authors")
		{
			authors.Post("/", authorHandler.CreateAuthor)
			authors.Get("/", authorHandler.GetAllAuthors)
			authors.Get("/:id", authorHandler.GetAuthor)
			authors.Get("/:id/books", authorHandler.GetAuthorWithBooks) // GET author with all books
			authors.Put("/:id", authorHandler.UpdateAuthor)
			authors.Delete("/:id", authorHandler.DeleteAuthor)
		}
	}

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	log.Printf("Server starting on port %s", cfg.AppPort)
	log.Fatal(app.Listen(cfg.AppPort))
}
