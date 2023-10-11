package server

import (
	"context"
	"log"
	"storage-api/controllers"
	"storage-api/repositories"
	"storage-api/services"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type HttpServer struct {
	config          *viper.Viper
	router          *gin.Engine
	imageController *controllers.ImageController
}

func InitDatabase(config *viper.Viper) *mongo.Client {
	// Set MongoDB connection options
	clientOptions := options.Client().ApplyURI(config.GetString("database.connection_string"))

	// Create a new MongoDB client
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		// Handle connection error
		log.Fatalln(err)
	}
	return client
}

func InitHttpServer(config *viper.Viper, client *mongo.Client) HttpServer {
	imageRepo := repositories.NewMongoDBUserRepository(client)
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
