package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Venue struct {
	ID       primitive.ObjectID `bson:"_id" json:"_id"`
	Location struct {
		Type        string     `bson:"type" json:"type"`
		Coordinates [2]float64 `bson:"coordinates" json:"coordinates"`
	} `bson:"venueLocation" json:"venueLocation"`
	Name       string             `bson:"name" json:"name"`
	Homepage   string             `bson:"homepageURL" json:"homepageURL"`
	OtherLinks []string           `bson:"otherLinks" json:"otherLinks"`
	Verified   bool               `bson:"verified" json:"verified"`
	OwnedBy    primitive.ObjectID `bson:"ownedBy" json:"ownedBy"`
}

func NewVenue() Venue {
	newVenue := Venue{
		ID: primitive.NewObjectID(),
	}
	newVenue.Location.Coordinates = [2]float64{0.0, 0.0}

	return newVenue
}
