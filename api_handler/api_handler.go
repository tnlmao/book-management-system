package apihandler

import (
	"book-management-system/api_handler/handler"
	"book-management-system/middleware"
	"book-management-system/service/logic"

	_ "book-management-system/docs"

	"github.com/gin-gonic/gin"
	files "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(router *gin.Engine) {
	svc := logic.NewBookService()
	router.Use(middleware.RedisMiddleware)
	router.Use(middleware.DbMiddleware)
	v1 := router.Group("/api/v1")
	{
		books := v1.Group("/books")
		{
			books.GET("", handler.GetBooks(svc))
			books.GET("/:id", handler.GetBook(svc))
			books.POST("", handler.CreateBook(svc))
			books.PUT("/:id", handler.UpdateBook(svc))
			books.DELETE("/:id", handler.DeleteBook(svc))
		}
	}
	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(files.Handler))

}
