package repo

import (
	"encoding/json"
	"fmt"

	logger "g.ghn.vn/go-common/zap-logger"
	"github.com/tientp-floware/mgodb-stream/config"
	db "github.com/tientp-floware/mgodb-stream/db/postgres"
	"github.com/tientp-floware/mgodb-stream/lib/request"
	model "github.com/tientp-floware/mgodb-stream/models"
)

var (
	urlDevices   = config.Config.Partner.URLDevices
	urlMetricGPS = config.Config.Partner.URLMetricGPS
	log          = logger.GetLogger("[Device service]")
)

func init() {
	log.Info("Device init service")
}

// Device represents the client
type Device struct {
}

// NewDevice returns a new database instance
func NewDevice() *Device {
	return &Device{}
}

// List devices
func (s *Device) List() []model.IOTDevice {
	return NewCRUD().GetDevices()
}

// ByVehiclePlate get device by vehicles
func (s *Device) ByVehiclePlate(plate string) *model.IOTDeviceLastest {
	result := new(model.IOTDeviceLastest)
	NewCRUD().ByVehiclePlate(plate, result)
	return result.Unmarshal()
}

// CollectByDeviceID collect
func (s *Device) CollectByDeviceID(deviceID int) []model.GPSTRackingJSON {
	btime, etime := db.IntervalTime(-1)
	url := fmt.Sprintf(urlMetricGPS, deviceID, btime, etime)
	makeReq := new(request.HTTPTransport)
	makeReq.GET(url)
	if makeReq.Err != nil {
		log.Info(`[CollectByDeviceID] sent request error:`, makeReq.Err)
	}
	var rawDatas []model.GPSTRackingJSON
	err := json.Unmarshal(makeReq.BodyByte, &rawDatas)
	if err != nil {
		log.Info(`[CollectByDeviceID] error:`, err)
	}
	return rawDatas
}

// CollectByDeviceIDWithTime collect
func (s *Device) CollectByDeviceIDWithTime(deviceID int, t int64) []model.GPSTRackingJSON {
	btime, etime := db.IntervalTime(t)
	url := fmt.Sprintf(urlMetricGPS, deviceID, btime, etime)
	log.Info("Deivice URI:", url)
	makeReq := new(request.HTTPTransport)
	makeReq.GET(url)
	if makeReq.Err != nil {
		log.Info(`[CollectByDeviceID] sent request error:`, makeReq.Err)
	}
	var rawDatas []model.GPSTRackingJSON
	err := json.Unmarshal(makeReq.BodyByte, &rawDatas)
	if err != nil {
		log.Info(`[CollectByDeviceID] error:`, err)
	}
	return rawDatas
}
