package repositories

import (
	"context"
	"net/http"
	"storage-api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ImageRepo interface {
	CreateImage(models.Image) (models.Image, models.ResponseError)
	GetImageByID(string) (models.Image, models.ResponseError)
	DeleteImage(string) models.ResponseError
}

type MongoDBUserRepository struct {
	client *mongo.Client
}

func NewMongoDBUserRepository(client *mongo.Client) *MongoDBUserRepository {
	return &MongoDBUserRepository{
		client: client,
	}
}

func (ir MongoDBUserRepository) CreateImage(image models.Image) (*models.Image, *models.ResponseError) {
	collection := ir.client.Database("lakipay_db").Collection("images")
	result, err := collection.InsertOne(context.TODO(), image)
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	id := result.InsertedID.(primitive.ObjectID)
	return &models.Image{
		ID:       id.Hex(),
		Data:     image.Data,
		Metadata: image.Metadata,
	}, nil
}

func (ir MongoDBUserRepository) GetImageByID(imageID string) (*models.Image, *models.ResponseError) {
	objectId, err := primitive.ObjectIDFromHex(imageID)
	if err != nil {
		return nil, &models.ResponseError{
			Message: "Invalid Image ID",
			Status:  http.StatusBadRequest,
		}
	}
	collection := ir.client.Database("lakipay_db").Collection("images")
	filter := bson.D{{Key: "_id", Value: objectId}}
	var image *models.Image
	err = collection.FindOne(context.TODO(), filter).Decode(&image)
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	return image, nil
}

func (ir MongoDBUserRepository) DeleteImage(imageID string) *models.ResponseError {
	objectId, err := primitive.ObjectIDFromHex(imageID)
	if err != nil {
		return &models.ResponseError{
			Message: "Invalid Image ID",
			Status:  http.StatusBadRequest,
		}
	}
	collection := ir.client.Database("lakipay_db").Collection("images")
	filter := bson.D{{Key: "_id", Value: objectId}}
	result, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	if result.DeletedCount == 0 {
		return &models.ResponseError{
			Message: "Image not found",
			Status:  http.StatusInternalServerError,
		}
	}
	return nil
}
