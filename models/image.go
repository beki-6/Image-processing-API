package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Image struct {
	ID       string             `json:"id" bson:"-"`
	ObjectID primitive.ObjectID `json:"-" bson:"_id"`
	Data     []byte             `json:"data"`
	Metadata Metadata           `json:"metadata,omitempty"`
}

type ImageData struct {
	Filename string `json:"filename"`
	Data     string `json:"data"`
}

type Dimension struct {
	Width  int
	Height int
}

type Metadata struct {
	Filename  string    `json:"filename"`
	Format    string    `json:"format"`
	Size      int64     `json:"size"`
	Dimension Dimension `json:"dimension"`
}
