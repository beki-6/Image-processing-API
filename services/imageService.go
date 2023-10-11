package services

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"os"
	"storage-api/models"
	"storage-api/repositories"
)

type ImageService struct {
	imageRepo *repositories.ImageRepo
}

func NewImageService(imgRepo *repositories.ImageRepo) *ImageService {
	return &ImageService{
		imageRepo: imgRepo,
	}
}

func (is ImageService) CreateImage(imageData []byte, filename string) (*models.Image, *models.ResponseError) {
	imageFile, format, size, responseError := processImage(imageData)
	if responseError != nil {
		return nil, responseError
	}
	metaData := models.Metadata{
		Filename: filename,
		Format:   format,
		Size:     size,
		Dimension: models.Dimension{
			Width:  imageFile.Bounds().Dx(),
			Height: imageFile.Bounds().Dy(),
		},
	}
	img := &models.Image{
		Data:     imageData,
		Metadata: metaData,
	}
	return is.imageRepo.CreateImage(img)
}

func (is ImageService) GetImageByID(imageID string) (*models.Image, *models.ResponseError) {
	err := validateImageID(imageID)
	if err != nil {
		return nil, err
	}
	return is.imageRepo.GetImageByID(imageID)
}

func (is ImageService) DeleteImage(imageID string) *models.ResponseError {
	err := validateImageID(imageID)
	if err != nil {
		return err
	}
	return is.imageRepo.DeleteImage(imageID)
}

// Function that takes image data as byte slice and returns decoded image and error, if any
func processImage(imageData []byte) (image.Image, string, int64, *models.ResponseError) {
	// Create temp file
	tempFile, err := os.CreateTemp("", "tempImage.*")
	if err != nil {
		return nil, "", 0, &models.ResponseError{
			Message: "Error while creating temp file",
			Status:  http.StatusInternalServerError,
		}
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()
	// Write image data to temp file
	_, err = tempFile.Write(imageData)
	if err != nil {
		return nil, "", 0, &models.ResponseError{
			Message: "Error while writing image to temp file",
			Status:  http.StatusInternalServerError,
		}
	}
	// open image file for processing
	file, err := os.Open(tempFile.Name())
	if err != nil {
		return nil, "", 0, &models.ResponseError{
			Message: "Error while opening image file",
			Status:  http.StatusInternalServerError,
		}
	}
	// Get the file size
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, "", 0, &models.ResponseError{
			Message: "Error while getting image size",
			Status:  http.StatusInternalServerError,
		}
	}
	fileSize := fileInfo.Size()
	defer file.Close()
	// Decode image
	img, format, err := image.Decode(file)
	if err != nil {
		return nil, "", 0, &models.ResponseError{
			Message: "Error while decoding image",
			Status:  http.StatusInternalServerError,
		}
	}
	return img, format, fileSize, nil
}

func validateImageID(id string) *models.ResponseError {
	if id == "" {
		return &models.ResponseError{
			Message: "Invalid ID",
			Status:  http.StatusBadRequest,
		}
	}
	return nil
}
