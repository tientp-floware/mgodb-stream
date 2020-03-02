package model

import (
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	// Setting collection
	Setting struct {
		ID           uint               `gorm:"primary_key"`
		MID          primitive.ObjectID `gorm:"-" json:"_id" bson:"_id"`
		UserID       int                `json:"user_id" bson:"user_id"`
		TopNav       TopNav             `json:"top_nav" bson:"top_nav"`
		TopNavString string             `gorm:"column:top_nav;type:varchar(100)"`
		CreatedAt    time.Time
		UpdatedAt    time.Time
		DeletedAt    *time.Time
	}

	// TopNav columns
	TopNav struct {
		UID     string `json:"uid"`
		View    string `json:"view"`
		Class   string `json:"class"`
		Title   string `json:"title"`
		IsView  bool   `json:"isView"`
		Color   string `json:"color"`
		HashKey string `json:"$$hashKey"`
	}
)

// String
func (st *Setting) String() *Setting {
	bolB, _ := json.Marshal(st.TopNav)
	st.TopNavString = (string(bolB))
	return st
}
