package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Service struct {
	ID          primitive.ObjectID `bson:"_id"`
	ServiceName string             `bson:"serviceName"`
	Identifier  string             `bson:"identifier"`
	Version     int                `bson:"version"`
	RunnerArgs  []string           `bson:"runArgs"`
	Exposed     bool               `bson:"exposed"`
}
