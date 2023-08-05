package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Listing struct {
	ID         primitive.ObjectID `bson:"_id" json:"_id"`
	Venue      primitive.ObjectID `bson:"venueID" json:"venueID"`
	Date       time.Time          `bson:"date" json:"date"`
	Lineup     []string           `bson:"lineup" json:"lineup"`
	TicketLink string             `bson:"ticketLink" json:"ticketLink"`
	OtherLinks []string           `bson:"otherLinks" json:"otherLinks"`
}

func NewListing() Listing {
	newListing := Listing{
		ID: primitive.NewObjectID(),
	}

	return newListing
}
