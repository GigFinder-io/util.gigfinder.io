package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NotificationSetting struct {
	ID   primitive.ObjectID `bson:"_id" json:"_id"`
	User primitive.ObjectID `bson:"user" json:"user"`
}

func NewNotificationSetting(user primitive.ObjectID) NotificationSetting {
	return NotificationSetting{
		ID:   primitive.NewObjectID(),
		User: user,
	}
}
