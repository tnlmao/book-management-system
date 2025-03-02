package handler

import (
	"book-management-system/models"
	"book-management-system/service/driver"
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetBooks godoc
// @Summary Gets all book entries  from redis cache
// @Description Retrieves a list of all books
// @Tags Books
// @Produce json
// @Success 200 {array} models.Book "Books retrieved successfully"
// @Failure 500 {string} string "Internal Server Error"
// @Router /books [get]
func GetBooks(svc driver.BookService) gin.HandlerFunc {

	return func(c *gin.Context) {
		// Extract pagination parameters from query
		limitStr := c.DefaultQuery("limit", "10")  // Default to 10 items per page
		offsetStr := c.DefaultQuery("offset", "0") // Default to start from the first item

		// Convert to integers with validation
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit parameter"})
			return
		}

		offset, err := strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offset parameter"})
			return
		}

		// Call the service layer to get paginated books and total count
		ctx := context.Background()
		books, err := svc.GetBooks(ctx, limit, offset) // Updated service method
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		var resp models.Response
		if books != nil {
			resp = models.Response{
				Code:    200,
				Message: "retrieved",
				Object:  &books,
			}
		}

		c.JSON(http.StatusOK, resp)
	}
}

// GetBook godoc
// @Summary Gets a book by ID
// @Description Retrieve a book's details using its ID
// @Tags Books
// @Produce json
// @Param id path int true "Book ID"
// @Success 200 {object} models.Book "Book retrieved successfully"
// @Failure 404 {string} string "Book not found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /books/{id} [get]
func GetBook(svc driver.BookService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the book ID from the URL parameter
		idStr := c.Param("id")

		// Convert the ID to an integer with validation
		id, err := strconv.Atoi(idStr)
		if err != nil || id < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book ID"})
			return
		}

		// Call the service layer to get the book by ID
		ctx := context.Background()
		book, err := svc.GetBook(ctx, id)
		if err != nil {
			// Handle "not found" error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
				return
			}
			// Handle other errors
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		var resp models.Response
		if book != nil {
			resp = models.Response{
				Code:    200,
				Message: "retrieved",
				Object:  &book,
			}
		}
		// Return the book as JSON
		c.JSON(http.StatusOK, resp)
	}
}

// CreateBook godoc
// @Summary Creates a new book
// @Description Adds a new book to the system
// @Tags Books
// @Accept json
// @Produce json
// @Param book body models.Book true "Book data"
// @Success 201 {boolean} bool "Book created successfully"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /books [post]
func CreateBook(svc driver.BookService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse the request body into a Book struct
		var book models.Book
		if err := c.ShouldBindJSON(&book); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		// Validate the book fields
		if book.Title == "" || book.Author == "" || book.Year == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "title, author, and year are required"})
			return
		}

		// Call the service layer to create the book
		ctx := context.Background()
		createdBook, err := svc.CreateBook(ctx, &book)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		var resp models.Response
		if createdBook {
			resp = models.Response{
				Code:    201,
				Message: "created",
			}
		}
		// Return the created book with a 201 status

		c.JSON(http.StatusCreated, resp)
	}
}

// UpdateBook godoc
// @Summary Update a book's details
// @Description Modify book details based on ID
// @Tags Books
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Param book body models.Book true "Updated book data"
// @Success 200 {boolean} bool "Book updated successfully"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Book not found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /books/{id} [put]
func UpdateBook(svc driver.BookService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the book ID from the URL parameter
		idStr := c.Param("id")

		// Convert the ID to an integer with validation
		id, err := strconv.Atoi(idStr)
		if err != nil || id < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book ID"})
			return
		}

		// Parse the request body into a Book struct
		var book models.Book
		if err := c.ShouldBindJSON(&book); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		// Validate the book fields
		if book.Title == "" || book.Author == "" || book.Year == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "title, author, and year are required"})
			return
		}

		// Call the service layer to update the book
		ctx := context.Background()
		updatedBook, err := svc.UpdateBook(ctx, id, &book)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		var resp models.Response
		if updatedBook {
			resp = models.Response{
				Code:    201,
				Message: "updated",
			}
		}
		// Return the updated book with a 200 status
		c.JSON(http.StatusOK, resp)
	}
}

// DeleteBook godoc
// @Summary Delete a book
// @Description Remove a book using its ID
// @Tags Books
// @Param id path int true "Book ID"
// @Success 204 {string} string "Book deleted successfully"
// @Failure 404 {string} string "Book not found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /books/{id} [delete]
func DeleteBook(svc driver.BookService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the book ID from the URL parameter
		idStr := c.Param("id")

		// Convert the ID to an integer with validation
		id, err := strconv.Atoi(idStr)
		if err != nil || id < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book ID"})
			return
		}

		// Call the service layer to delete the book
		ctx := context.Background()
		err = svc.DeleteBook(ctx, id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return a success message with a 200 status
		c.JSON(http.StatusOK, models.Response{Code: 200, Message: "book deleted successfully"})
	}
}
