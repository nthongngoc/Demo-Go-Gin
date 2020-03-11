package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/khoa5773/go-server/docs"
	"github.com/khoa5773/go-server/src/configs"
	"github.com/khoa5773/go-server/src/domains/auth"
	"github.com/khoa5773/go-server/src/domains/documents"
	"github.com/khoa5773/go-server/src/domains/projects"
	"github.com/khoa5773/go-server/src/domains/repositories"
	"github.com/khoa5773/go-server/src/domains/users"
	"github.com/khoa5773/go-server/src/shared"
)

// @title Swagger Example API server
// @version 2.0
// @description This is a sample server.
// @termsOfService http://swagger.io/terms/

func init() {
	gin.SetMode(os.Getenv("GIN_MODE"))
}

func main() {
	app := gin.Default()

	app.Use(cors.Default())
	app.Use(shared.ErrorHandler)

	binding.Validator = &shared.DefaultValidator{}

	swaggerURL := ginSwagger.URL(fmt.Sprintf("http://%s:%d/swagger/doc.json", configs.ConfigsService.Host, configs.ConfigsService.Port))
	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, swaggerURL))

	auth.ApplyRoutes(app)
	users.ApplyRoutes(app)
	projects.ApplyRoutes(app)
	repositories.ApplyRoutes(app)
	documents.ApplyRoutes(app)

	err := app.Run(fmt.Sprintf(":%d", configs.ConfigsService.Port))
	if err != nil {
		log.Fatalln(err)
	}
}
