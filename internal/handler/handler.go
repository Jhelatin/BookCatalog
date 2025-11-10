package v1

import (
	"awesomeProject/internal/entity"
	"awesomeProject/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type BookHandler struct {
	bookService *service.BookService
}

func NewBookHandler(bookService *service.BookService) *BookHandler {
	return &BookHandler{bookService: bookService}
}

func (h *BookHandler) CreateBook(c *fiber.Ctx) error {
	var book entity.Book
	if err := c.BodyParser(&book); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	if err := h.bookService.CreateBook(c.Context(), &book); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create book",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(book)
}

func (h *BookHandler) GetBook(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid book ID",
		})
	}

	var userID int
	if c.Locals("userID") != nil {
		userID = c.Locals("userID").(int)
	}

	book, err := h.bookService.GetBook(c.Context(), id, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get book",
		})
	}

	if book == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Book not found",
		})
	}

	return c.JSON(book)
}

func (h *BookHandler) GetAllBooks(c *fiber.Ctx) error {
	books, err := h.bookService.GetAllBooks(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get books",
		})
	}

	return c.JSON(books)
}

func (h *BookHandler) UpdateBook(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid book ID",
		})
	}

	var book entity.Book
	if err := c.BodyParser(&book); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	book.ID = id
	if err := h.bookService.UpdateBook(c.Context(), &book); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update book",
		})
	}

	return c.JSON(book)
}

func (h *BookHandler) DeleteBook(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid book ID",
		})
	}

	if err := h.bookService.DeleteBook(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete book",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
