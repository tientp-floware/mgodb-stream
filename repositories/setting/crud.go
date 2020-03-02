package repo

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/tientp-floware/mgodb-stream/config"
	"github.com/tientp-floware/mgodb-stream/db"
	"github.com/tientp-floware/mgodb-stream/db/mgodb"
	"github.com/tientp-floware/mgodb-stream/db/mysql"
	model "github.com/tientp-floware/mgodb-stream/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	collSetting = mgodb.NewSwitchMogoDB(config.Config.Mongodb.Database, db.Setting)
	dbmysql     = mysql.GetDB()
)

type (
	// SettingCRUD store data
	SettingCRUD struct {
		Data *model.Setting
	}
)

// NewCRUD collection
func NewCRUD() *SettingCRUD {

	setting := &model.Setting{}
	// After db connection is created.
	dbmysql.AutoMigrate(setting)

	return &SettingCRUD{
		Data: setting,
	}
}

// Find trip
func (tr *SettingCRUD) Find(code string) error {
	return nil
}

// Create trip
func (tr *SettingCRUD) Create() error {
	tr.Data.MID = primitive.NewObjectID()
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

// SQLCreate to mysql
func (tr *SettingCRUD) SQLCreate(dataSetting *model.Setting) error {

	err := dbmysql.Where(model.Setting{UserID: dataSetting.UserID}).Assign(dataSetting).FirstOrCreate(dataSetting).Error
	if err != nil {
		log.Infof("[Setting] CreateOrUpdate error is: %s", err)
	} else {
		log.Infof("[Setting] CreateOrUpdate done ID is: %d", tr.Data.UserID)
	}
	return err
}

// SQLUpdate to mysql
func (tr *SettingCRUD) SQLUpdate() error {
	return dbmysql.Model(&model.Setting{}).Where("user_id = ?", tr.Data.UserID).Updates(tr.Data).Error
}

// SQLDelete to mysql
func (tr *SettingCRUD) SQLDelete() error {
	return nil
}
