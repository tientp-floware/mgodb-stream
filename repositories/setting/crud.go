package repo

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/tientp-floware/mgodb-stream/config"
	"github.com/tientp-floware/mgodb-stream/db"
	"github.com/tientp-floware/mgodb-stream/db/mgodb"
	model "github.com/tientp-floware/mgodb-stream/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	collSetting = mgodb.NewSwitchMogoDB(config.Config.Mongodb.Database, db.Setting)
)

type (
	// SettingCRUD store data
	SettingCRUD struct {
		Data *model.Setting
	}
)

// NewCRUD collection
func NewCRUD() *SettingCRUD {
	return &SettingCRUD{
		Data: &model.Setting{},
	}
}

// Find trip
func (tr *SettingCRUD) Find(code string) error {
	return nil
}

// Create trip
func (tr *SettingCRUD) Create() error {
	tr.Data.ID = primitive.NewObjectID()
	_, err := collSetting.InsertOne(tr.Data)
	if err != nil {
		log.Error("Can not insert:", err)
		return err
	}
	return nil
}

// Update  trip
func (tr *SettingCRUD) Update(id string) error {
	matched, err := collSetting.UpdateOne(echo.Map{"user_id": id}, echo.Map{"$set": echo.Map{
		"updated_at": time.Now(),
	}})
	if err != nil {
		log.Error("Can not updated:", err)
		return err
	}
	log.Info("[Tracking] MatchedCount:", matched.MatchedCount, " ModifiedCount:", matched.ModifiedCount)
	return err
}
