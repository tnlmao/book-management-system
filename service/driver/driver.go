package driver

import (
	"book-management-system/models"
	"context"
)

type BookService interface {
	GetBooks(ctx context.Context, limit, offset int) ([]models.Book, error)
	GetBook(ctx context.Context, id int) (*models.Book, error)
	CreateBook(ctx context.Context, book *models.Book) (bool, error)
	UpdateBook(ctx context.Context, id int, book *models.Book) (bool, error)
	DeleteBook(ctx context.Context, id int) error
}
