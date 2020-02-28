package model

import (
	"encoding/json"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func init() {

}

// Base hold general field in model
type Base struct {
	ID        int       `gorm:"primary_key" json:"id" `
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	IsDeleted bool      `json:"is_deleted" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at" gorm:"default:now()"`
	UpdatedAt time.Time `json:"updated_at" gorm:"default:now()"`
}

// InitBase to default value
func (b *Base) InitBase() *Base {
	b.CreatedAt = time.Now()
	b.UpdatedAt = time.Now()
	b.IsDeleted = false
	return b
}

// SetID set id customize
func (b *Base) SetID(id int) *Base {
	b.ID = id
	return b
}

// Mgostream full document stream
type Mgostream struct {
	OPType string `bson:"operationType"`
	Key    struct {
		ID primitive.ObjectID `bson:"_id"`
	} `bson:"documentKey"`
	Detail struct {
		UpdateFiled map[string]interface{} `bson:"updatedFields"`
		RemoveFiled []string               `bson:"removedFields"`
	} `bson:"updateDescription"`
	FullDoc map[string]interface{} `bson:"fullDocument"`
}

// ToStruct convert to model
func (mt *Mgostream) ToStruct(result interface{}) {
	b, err := json.Marshal(mt.FullDoc)
	if err != nil {
		fmt.Println("error:", err)
	}
	err = json.Unmarshal(b, &result)
	if err != nil {
		fmt.Println("error:", err)
	}
}
