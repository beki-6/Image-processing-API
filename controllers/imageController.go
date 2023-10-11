package controllers

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"storage-api/models"
	"storage-api/services"

	"github.com/gin-gonic/gin"
)

type ImageController struct {
	imageService *services.ImageService
}

func NewImageController(imgService *services.ImageService) *ImageController {
	return &ImageController{
		imageService: imgService,
	}
}

func (ic ImageController) CreateImage(ctx *gin.Context) {
	// Read JSON body
	body, err := io.ReadAll(ctx.Request.Body)
	// Handle parse errors
	if err != nil {
		log.Println("Error while reading JSON data", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	// Unmarshal JSON
	var imageData *models.ImageData
	err = json.Unmarshal(body, &imageData)
	if err != nil {
		log.Println("Error while unmarshaling request body ", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	// Decode the base64-encoded image data
	imageBytes, err := base64.StdEncoding.DecodeString(imageData.Data)
	if err != nil {
		log.Println("Failed to decode image data", err)
		return
	}
	// Call ImageService
	response, responseError := ic.imageService.CreateImage(imageBytes, imageData.Filename)
	// Handle error
	if err != nil {
		ctx.AbortWithStatusJSON(responseError.Status, responseError)
		return
	}
	ctx.JSON(http.StatusOK, response)
}

func (ic ImageController) DeleteImage(ctx *gin.Context) {
	imageID := ctx.Param("imageID")
	responseError := ic.imageService.DeleteImage(imageID)
	if responseError != nil {
		ctx.JSON(responseError.Status, responseError)
	}
	ctx.Status(http.StatusNoContent)
}

func (ic ImageController) GetImageByID(ctx *gin.Context) {
	imageID := ctx.Param("imageID")
	response, responseError := ic.imageService.GetImageByID(imageID)
	if responseError != nil {
		ctx.JSON(responseError.Status, responseError)
	}
	ctx.JSON(http.StatusOK, response)
}
