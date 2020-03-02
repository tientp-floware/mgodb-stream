package transport

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/labstack/echo/v4"
	log "github.com/labstack/gommon/log"
	"github.com/tientp-floware/mgodb-stream/db"
	"github.com/tientp-floware/mgodb-stream/db/mgodb"
	model "github.com/tientp-floware/mgodb-stream/models"
	repository "github.com/tientp-floware/mgodb-stream/repositories"
)

var (
	tripstatus  = []string{"fail", "cancel", "finished"}
	flowSetting = mgodb.NewMogoDB(db.Setting)
	//log        = logger.GetLogger("[Device service]")
)

type (
	// FlowareMgoStream background
	FlowareMgoStream struct {
		Worker map[string]DataRow
		repo   *repository.Service
	}
	// DataRow info device
	DataRow struct {
		OID         primitive.ObjectID `bson:"_id"`
		UserSetting string
	}
)

// NewMgoStream create new mgo stream
func NewMgoStream() *FlowareMgoStream {
	return &FlowareMgoStream{repo: repository.New()}
}

// FlowChangeStream stream consumer
// We using to catch event from Trip
func (mgstream *FlowareMgoStream) FlowChangeStream() *FlowareMgoStream {
	log.Info("We run streaming....")
	// Query and can use for scale up
	pipeline := []echo.Map{}
	mgstream.Worker = make(map[string]DataRow)
	// stream func to handle event
	streamer := func(cs *mongo.ChangeStream) {

		for cs.Next(context.Background()) {

			changeDoc := new(model.Mgostream)

			if err := cs.Decode(changeDoc); err != nil {
				log.Info("err:", err)
				break
			}

			fmt.Println("changeDoc:", changeDoc)

			setting := new(model.Setting)
			changeDoc.ToStruct(setting)
			mgstream.createOrupdate(setting)
			fmt.Println("setting:", setting)
		}
	}

	go flowSetting.ChangeStream(pipeline, streamer)

	return mgstream
}

func (mgstream *FlowareMgoStream) createOrupdate(data *model.Setting) {
	data.String()
	err := mgstream.repo.Setting.Crud.SQLCreate(data)
	if err != nil {
		fmt.Println("err:", err)
	}
}
