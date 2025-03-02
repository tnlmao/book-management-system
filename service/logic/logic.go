package logic

import (
	"book-management-system/config/database"
	"book-management-system/config/kafka"
	"book-management-system/config/redis"
	"book-management-system/models"
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

// BookService handles book operations
type BookService struct {
	producer *kafka.KafkaProducer
}

// NewBookService initializes BookService with Kafka producer
func NewBookService() *BookService {
	broker := "localhost:9092" // Kafka broker
	producer := kafka.NewKafkaProducer(broker)
	return &BookService{producer: producer}
}

func (b *BookService) CreateBook(ctx context.Context, book *models.Book) (bool, error) {
	//Create book entry in Database
	result := database.DB.WithContext(ctx).Create(book)
	if result.Error != nil {
		return false, result.Error
	}
	//Marshal to store entry in Redis Cache
	bookJSON, err := json.Marshal(book)
	if err == nil {
		redis.Client.Set(ctx, fmt.Sprintf("books:%d", book.ID), bookJSON, 10*time.Minute)
	} else {
		fmt.Println("Failed to update book in Redis:", err)
		return false, err
	}
	//Publish created event
	b.producer.PublishBookEvent("create", book)
	fmt.Println("New book added to Redis cache")
	return true, nil
}

func (b *BookService) DeleteBook(ctx context.Context, id int) error {
	redisKey := fmt.Sprintf("books:%d", id)

	// Delete from Database
	result := database.DB.WithContext(ctx).Delete(&models.Book{}, id)
	if result.Error != nil {
		return result.Error
	}

	// Delete from Redis
	err := redis.Client.Del(ctx, redisKey).Err()
	if err != nil {
		fmt.Println("Failed to delete book from Redis:", err)
		return err
	}
	//Publish Deleted  event
	b.producer.PublishBookEvent("delete", id)
	fmt.Println("Book deleted from both DB and Redis")
	return nil
}

func (b *BookService) GetBook(ctx context.Context, id int) (*models.Book, error) {
	redisKey := fmt.Sprintf("books:%d", id)

	// Check Redis cache first
	val, err := redis.Client.Get(ctx, redisKey).Result()
	if err == nil {
		// Book found in Redis, return it
		var book models.Book
		if jsonErr := json.Unmarshal([]byte(val), &book); jsonErr == nil {
			fmt.Println("Cache hit: Returning book from Redis")
			return &book, nil
		}
		fmt.Println("Error unmarshalling book from Redis:", err)
	}

	// If not found in Redis, return an error (No DB lookup)
	return nil, fmt.Errorf("book with ID %d not found in cache", id)
}

func (b *BookService) GetBooks(ctx context.Context, limit int, offset int) ([]models.Book, error) {
	var books []models.Book

	// Get all keys matching "books:*" from Redis.
	bookKeys, err := redis.Client.Keys(ctx, "books:*").Result()
	if err != nil || len(bookKeys) == 0 {
		// Cache miss: fetch paginated data from DB.
		fmt.Println("Cache miss: Fetching books from DB...")
		if err := database.DB.Limit(limit).Offset(offset).Find(&books).Error; err != nil {
			return nil, err
		}
		// Cache each book individually.
		for _, book := range books {
			bookJSON, _ := json.Marshal(book)
			redis.Client.Set(ctx, fmt.Sprintf("books:%d", book.ID), bookJSON, 10*time.Minute)
		}
		return books, nil
	}

	// Define a helper type to sort keys by their numeric ID.
	type bookKey struct {
		id  int
		key string
	}
	var bk []bookKey

	// Parse the keys and extract numeric IDs.
	for _, key := range bookKeys {
		parts := strings.Split(key, ":")
		if len(parts) != 2 {
			continue
		}
		id, err := strconv.Atoi(parts[1])
		if err != nil {
			continue
		}
		bk = append(bk, bookKey{id: id, key: key})
	}

	// Sort keys numerically.
	sort.Slice(bk, func(i, j int) bool {
		return bk[i].id < bk[j].id
	})

	// Apply pagination: compute start and end indices.
	start := offset
	end := offset + limit
	if start >= len(bk) {
		return []models.Book{}, nil // No results for this page.
	}
	if end > len(bk) {
		end = len(bk)
	}

	// Retrieve books for the specified range.
	for _, k := range bk[start:end] {
		bookJSON, err := redis.Client.Get(ctx, k.key).Result()
		if err == nil {
			var book models.Book
			if jsonErr := json.Unmarshal([]byte(bookJSON), &book); jsonErr == nil {
				books = append(books, book)
			}
		}
	}
	fmt.Println("Serving paginated results from Redis cache")
	return books, nil
}
func (b *BookService) UpdateBook(ctx context.Context, id int, updatedBook *models.Book) (bool, error) {
	redisKey := fmt.Sprintf("books:%d", id)

	// Find the existing book
	var existingBook models.Book
	if err := database.DB.WithContext(ctx).First(&existingBook, id).Error; err != nil {
		return false, err // Return if the book is not found
	}

	// Update fields
	existingBook.Title = updatedBook.Title
	existingBook.Author = updatedBook.Author
	existingBook.Year = updatedBook.Year

	// Save the updated book in DB
	if err := database.DB.WithContext(ctx).Save(&existingBook).Error; err != nil {
		return false, err
	}

	// Update Redis Cache
	bookJSON, err := json.Marshal(existingBook)
	if err == nil {
		redis.Client.Set(ctx, redisKey, bookJSON, 0) // Update Redis
	} else {
		fmt.Println("Failed to update book in Redis:", err)
		return false, err
	}
	//Publish updated event
	b.producer.PublishBookEvent("update", id)
	fmt.Println("Book updated in both DB and Redis")
	return true, nil
}
