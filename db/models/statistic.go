package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Statistic struct {
	ID                 primitive.ObjectID  `bson:"_id" json:"_id"`
	Date               time.Time           `bson:"date" json:"date"`
	MessagesSent       int                 `bson:"messagesSent" json:"messagesSent"`
	Logins             int                 `bson:"logins" json:"logins"`
	Signups            int                 `bson:"signups" json:"signups"`
	Searches           int                 `bson:"searches" json:"searches"`
	RootViews          int                 `bson:"rootViews" json:"rootViews"`
	Referrers          []Referrer          `bson:"referrers" json:"referrers"`
	Reports            int                 `bson:"reports" json:"reports"`
	ServerErrors       int                 `bson:"serverErrors" json:"serverErrors"`
	UserErrors         int                 `bson:"userErrors" json:"userErrors"`
	AlertsSent         int                 `bson:"alertsSent" json:"alertsSent"`
	UnauthorizedErrors int                 `bson:"unauthorizedErrors" json:"unauthorizedErrors"`
	MatchCount         MatchCount          `bson:"matchCount" json:"matchCount"`
	Timing             map[string]TimeData `bson:"endpointTiming" json:"endpointTiming"`
	Errors             []ErrorInfo         `bson:"errors" json:"errors"`
}

type MatchCount struct {
	Average int `bson:"avg" json:"avg"`
	Total   int `bson:"total" json:"total"`
}

type Referrer struct {
	URL   string `bson:"url" json:"url"`
	Count int    `bson:"count" json:"count"`
}

type ErrorInfo struct {
	Message   string    `bson:"message" json:"message"`
	Origin    string    `bson:"origin" json:"origin"`
	Timestamp time.Time `bson:"timestamp" json:"timestamp"`
}

type TimeData struct {
	Min   float64 `bson:"min" json:"min"`
	Avg   float64 `bson:"avg" json:"avg"`
	Max   float64 `bson:"max" json:"max"`
	Count int     `bson:"count" json:"count"`
}

func NewStatistic() Statistic {
	date := time.Now()
	strDate := date.Format("02-01-2006")
	date, _ = time.Parse("02-01-2006", strDate)

	newStat := Statistic{
		ID:        primitive.NewObjectID(),
		Date:      date,
		Timing:    map[string]TimeData{},
		Referrers: []Referrer{},
		Errors:    []ErrorInfo{},
	}

	return newStat
}
