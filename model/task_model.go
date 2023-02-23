package model

import "time"

type Task struct {
	Name      string    `bson:"name" json:"name" validate:"required"`
	TotalTime int       `bson:"total_time" json:"total_time,default=10"`
	CreatedAt time.Time `bson:"created_at,omitepty" json:"created_at"`
}
