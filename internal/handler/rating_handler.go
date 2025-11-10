package v1

import (
	"awesomeProject/internal/entity"
	"awesomeProject/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type RatingHandler struct {
	ratingService *service.RatingService
}

func NewRatingHandler(ratingService *service.RatingService) *RatingHandler {
	return &RatingHandler{ratingService: ratingService}
}

func (h *RatingHandler) RateBook(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)

	var req entity.RatingRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	if err := h.ratingService.RateBook(c.Context(), userID, &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Rating submitted successfully",
		"book_id": req.BookID,
		"score":   req.Score,
	})
}

func (h *RatingHandler) GetMyRating(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)

	bookID, err := strconv.Atoi(c.Params("book_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid book ID",
		})
	}

	rating, err := h.ratingService.GetUserRating(c.Context(), userID, bookID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get rating",
		})
	}

	if rating == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "No rating found for this book",
		})
	}

	return c.JSON(rating)
}

func (h *RatingHandler) GetMyRatings(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)

	ratings, err := h.ratingService.GetUserRatings(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get ratings",
		})
	}

	return c.JSON(ratings)
}

func (h *RatingHandler) GetBookRatings(c *fiber.Ctx) error {
	bookID, err := strconv.Atoi(c.Params("book_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid book ID",
		})
	}

	ratings, err := h.ratingService.GetBookRatings(c.Context(), bookID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get book ratings",
		})
	}

	return c.JSON(ratings)
}

func (h *RatingHandler) RemoveRating(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)

	bookID, err := strconv.Atoi(c.Params("book_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid book ID",
		})
	}

	if err := h.ratingService.RemoveRating(c.Context(), userID, bookID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to remove rating",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Rating removed successfully",
	})
}

func (h *RatingHandler) GetBooksWithMyRatings(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)

	books, err := h.ratingService.GetBooksWithUserRatings(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get books with ratings",
		})
	}

	return c.JSON(books)
}
