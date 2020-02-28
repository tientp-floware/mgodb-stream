package repo

import (
	"sync"

	db "github.com/tientp-floware/mgodb-stream/db/postgres"
	"github.com/tientp-floware/mgodb-stream/lib/util"
	model "github.com/tientp-floware/mgodb-stream/models"
)

// DeviceCRUD represents the client
type DeviceCRUD struct {
}

// NewCRUD returns a new Fuel database instance
func NewCRUD() *DeviceCRUD {
	return &DeviceCRUD{}
}

// CreateOrUpdate creates or update new one
func (u *DeviceCRUD) CreateOrUpdate(vdv model.IOTDeviceJSON, wg *sync.WaitGroup) {
	err := db.GetDB().Where(model.IOTDevice{IDv: vdv.IDv}).Assign(vdv).FirstOrCreate(&vdv).Error
	if err != nil {
		log.Infof("[Device] CreateOrUpdate error is: %s", err)
	} else {
		log.Infof("[Device] CreateOrUpdate done ID is: %d", vdv.IDv)
	}
	wg.Done()
}

// BulkCreate bulk insert
/* func (u *Device) BulkCreate(devices []interface{}) error {
	return gormbulk.BulkInsert(db.GetDB(), devices, 3000)
} */

// BulkUpdate update
func (u *DeviceCRUD) BulkUpdate(devices []model.IOTDeviceJSON) error {
	return db.GetDB().Table("device").Updates(&devices).Error
}

// GetOneByDeviceID get device by device_id
func (u *DeviceCRUD) GetOneByDeviceID(id int) *model.IOTDevice {
	dv := new(model.IOTDevice)
	err := db.GetDB().Where(model.IOTDevice{IDv: id}).First(&dv).Error
	if err != nil {
		log.Infof("[Device] GetOneByDeviceID error is: %s - id: %d", err, id)
	}
	return dv
}

// GetDevices list devices
func (u *DeviceCRUD) GetDevices() []model.IOTDevice {
	var devices []model.IOTDevice
	err := db.GetDB().Select(`device_id`).Where("is_active = ? AND is_deleted = ?", true, false).Find(&devices).Error
	if err != nil {
		log.Infof("[Device] GetDevices error is: %s ", err)
	}
	return devices
}

// ByVehiclePlate get device by vehicles
func (u *DeviceCRUD) ByVehiclePlate(plate string, result *model.IOTDeviceLastest) *DeviceCRUD {
	plainPlate := util.RemoveSpecialChar(plate)
	db.GetDB().Raw("SELECT * from gps_tracking_lastest_func(?)", plainPlate).Scan(&result)
	return u
}
