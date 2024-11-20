package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id" json:"_id"`
	Email        string             `bson:"email" json:"email"`
	Password     string             `bson:"password" json:"password"`
	Name         string             `bson:"name" json:"name"`
	Phone        string             `bson:"phone" json:"phone"`
	ProfileImage string             `bson:"profile_image" json:"profile_image"`
	Friends      []string           `bson:"friends" json:"friends"`
	Subscription *Subscription      `bson:"subscription" json:"subscription"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
	LastLoginAt  time.Time          `bson:"last_login_at" json:"last_login_at"`
	Status       string             `bson:"status" json:"status"`
	Role         string             `bson:"role" json:"role"`
}

type Subscription struct {
	PlanID    string    `bson:"plan_id" json:"plan_id"`
	Name      string    `bson:"name" json:"name"`
	Price     float64   `bson:"price" json:"price"`
	StartDate time.Time `bson:"start_date" json:"start_date"`
	EndDate   time.Time `bson:"end_date" json:"end_date"`
	Features  []string  `bson:"features" json:"features"`
	IsActive  bool      `bson:"is_active" json:"is_active"`
}
