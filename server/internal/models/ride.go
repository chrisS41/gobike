package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Ride struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	RouteID   primitive.ObjectID `bson:"route_id,omitempty" json:"route_id"`
	StartTime time.Time          `bson:"start_time" json:"start_time"`
	EndTime   time.Time          `bson:"end_time" json:"end_time"`
	Distance  float64            `bson:"distance" json:"distance"`
	Duration  time.Duration      `bson:"duration" json:"duration"`
	AvgSpeed  float64            `bson:"avg_speed" json:"avg_speed"`
	MaxSpeed  float64            `bson:"max_speed" json:"max_speed"`
	Calories  float64            `bson:"calories" json:"calories"`
	Locations []GeoPoint         `bson:"locations" json:"locations"`
	Weather   WeatherInfo        `bson:"weather" json:"weather"`
}

type WeatherInfo struct {
	Temperature float64 `bson:"temperature" json:"temperature"`
	Humidity    int     `bson:"humidity" json:"humidity"`
	Conditions  string  `bson:"conditions" json:"conditions"`
}
