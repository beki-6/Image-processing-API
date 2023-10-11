package server

import (
	"database/sql"
	"log"
	"storage-api/controllers"
	"storage-api/repositories"
	"storage-api/services"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type HttpServer struct {
	config          *viper.Viper
	router          *gin.Engine
	imageController *controllers.ImageController
}

func InitHttpServer(config *viper.Viper, dbHandler *sql.DB) HttpServer {
	imageRepo := repositories.NewImageRepo(dbHandler)
	imageService := services.NewImageService(imageRepo)
	imageController := controllers.NewImageController(imageService)

	router := gin.Default()
	router.POST("/image", imageController.CreateImage)
	router.DELETE("/image/:imageID", imageController.DeleteImage)
	router.GET("/image/:imageID", imageController.GetImageByID)

	return HttpServer{
		config:          config,
		router:          router,
		imageController: imageController,
	}
}

func (hs HttpServer) Start() {
	err := hs.router.Run(hs.config.GetString("http.address"))
	if err != nil {
		log.Fatalf("Error while starting HTTP server: %v", err)
	}
}
