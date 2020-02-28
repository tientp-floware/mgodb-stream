package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	// Setting collection
	Setting struct {
		ID     primitive.ObjectID `bson:"_id"`
		UserID int                `json:"user_id" bson:"user_id"`
		TopNav []TopNav           `json:"top_nav" bson:"top_nav"`
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
