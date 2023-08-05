package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	SearchTypeForm    = "Form"
	SearchTypeJoin    = "Join"
	SearchTypeEither  = "Either"
	SearchTypeRecruit = "Recruit"
)

type User struct {
	ID             primitive.ObjectID `bson:"_id" json:"_id"`
	SearchLocation struct {
		Type        string     `bson:"type" json:"type"`
		Coordinates [2]float64 `bson:"coordinates" json:"coordinates"`
	} `bson:"searchLocation" json:"searchLocation"`
	SearchType      string   `bson:"searchtype" json:"searchType"`
	Genres          []string `bson:"genres" json:"genres"`
	Instruments     []string `bson:"instruments" json:"instruments"`
	SearchRadius    float32  `bson:"searchRadius" json:"searchRadius"`
	Description     string   `bson:"description" json:"description"`
	Active          bool     `bson:"active" json:"active"`
	Admin           bool     `bson:"admin" json:"admin"`
	EmailConfirmed  bool     `bson:"emailConfirmed" json:"emailConfirmed"`
	Email           string   `bson:"email" json:"email"`
	DisplayName     string   `bson:"displayName" json:"displayName"`
	PasswordHash    string   `bson:"passwordHash" json:"-"`
	ConfirmString   string   `bson:"confirmString" json:"-"`
	PassResetString string   `bson:"passResetString" json:"-"`
	PassRest        struct {
		Token     string    `bson:"token" json:"-"`
		Timestamp time.Time `bson:"timestamp" json:"-"`
	} `bson:"passReset"  json:"-"`
	Timestamps struct {
		LastLogin   time.Time `bson:"last_login" json:"last_login"`
		SignupAt    time.Time `bson:"signup_at" json:"signup_at"`
		LastUpdated time.Time `bson:"last_updated" json:"-"`
	} `bson:"timestamps" json:"timestamps"`
	AudioURL string `bson:"audioURL" json:"audioURL"`
}

func NewUser(email string, name string, passHash string, confirmString string) User {
	newUser := User{
		ID:            primitive.NewObjectID(),
		Email:         email,
		DisplayName:   name,
		ConfirmString: confirmString,
		PasswordHash:  passHash,
		Genres:        []string{},
		Instruments:   []string{},
		Active:        true,
		Admin:         false,
		SearchType:    SearchTypeForm,
	}
	newUser.Timestamps.LastLogin = time.Now()
	newUser.Timestamps.SignupAt = time.Now()
	newUser.SearchLocation.Type = "Point"
	newUser.SearchLocation.Coordinates = [2]float64{0.0, 0.0}

	return newUser
}
