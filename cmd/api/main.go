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
	userRepo := postgres.NewUserRepository(db)
	ratingRepo := postgres.NewRatingRepository(db)

	ratingService := service.NewRatingService(ratingRepo, bookRepo)
	bookService := service.NewBookService(bookRepo, ratingService)
	authorService := service.NewAuthorService(authorRepo, bookRepo)
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)

	bookHandler := v1.NewBookHandler(bookService)
	authorHandler := v1.NewAuthorHandler(authorService)
	authHandler := v1.NewAuthHandler(authService)
	ratingHandler := v1.NewRatingHandler(ratingService)

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

		users := api.Group("/users")
		{
			users.Post("/login", authHandler.Login)
		}
		ratings := api.Group("/ratings")
		{
			ratings.Post("/", ratingHandler.RateBook)
			ratings.Get("/my", ratingHandler.GetMyRatings)
			ratings.Get("/books-with-my-ratings", ratingHandler.GetBooksWithMyRatings)
			ratings.Get("/book/:book_id", ratingHandler.GetMyRating)
			ratings.Delete("/book/:book_id", ratingHandler.RemoveRating)
		}
	}

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	log.Printf("Server starting on port %s", cfg.AppPort)
	log.Fatal(app.Listen(cfg.AppPort))
}
