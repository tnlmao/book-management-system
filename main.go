package main

import (
	apihandler "book-management-system/api_handler"
	"book-management-system/app"
	"book-management-system/config"
)

func init() {

}
func main() {
	config.LoadConfig()
	application := app.New()
	apihandler.SetupRoutes(application.Router)

	application.Router.Run(":8080")
}
