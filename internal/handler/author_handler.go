package v1

import (
	"awesomeProject/internal/entity"
	"awesomeProject/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type AuthorHandler struct {
	authorService *service.AuthorService
}

func NewAuthorHandler(authorService *service.AuthorService) *AuthorHandler {
	return &AuthorHandler{authorService: authorService}
}

func (h *AuthorHandler) CreateAuthor(c *fiber.Ctx) error {
	var author entity.Author
	if err := c.BodyParser(&author); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	if err := h.authorService.CreateAuthor(c.Context(), &author); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create author",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(author)
}

func (h *AuthorHandler) GetAuthor(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid author ID",
		})
	}

	author, err := h.authorService.GetAuthor(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get author",
		})
	}

	if author == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Author not found",
		})
	}

	return c.JSON(author)
}

func (h *AuthorHandler) GetAllAuthors(c *fiber.Ctx) error {
	authors, err := h.authorService.GetAllAuthors(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get authors",
		})
	}

	return c.JSON(authors)
}

func (h *AuthorHandler) UpdateAuthor(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid author ID",
		})
	}

	var author entity.Author
	if err := c.BodyParser(&author); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	author.ID = id
	if err := h.authorService.UpdateAuthor(c.Context(), &author); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update author",
		})
	}

	return c.JSON(author)
}

func (h *AuthorHandler) DeleteAuthor(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid author ID",
		})
	}

	if err := h.authorService.DeleteAuthor(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete author",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *AuthorHandler) GetAuthorWithBooks(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid author ID",
		})
	}

	author, books, err := h.authorService.GetAuthorWithBooks(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get author with books",
		})
	}

	if author == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Author not found",
		})
	}

	return c.JSON(fiber.Map{
		"author": author,
		"books":  books,
	})
}
