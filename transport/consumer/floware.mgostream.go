package transport

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/tientp-floware/mgodb-stream/config"
	"github.com/tientp-floware/mgodb-stream/db"
	"github.com/tientp-floware/mgodb-stream/db/mgodb"
	model "github.com/tientp-floware/mgodb-stream/models"
	repository "github.com/tientp-floware/mgodb-stream/repositories"
	logger "go.uber.org/zap"
)

var (
	tripstatus = []string{"fail", "cancel", "finished"}
	dbtrip     = mgodb.NewSwitchMogoDB(config.Config.Mongodb.Database, db.Setting)
	log        = logger.GetLogger("[Device service]")
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
	// Query and can use for scale up
	pipeline := nil
	mgstream.Worker = make(map[string]DataRow)
	// stream func to handle event
	streamer := func(cs *mongo.ChangeStream) {
		for cs.Next(context.Background()) {

			changeDoc := new(model.Mgostream)
			if err := cs.Decode(changeDoc); err != nil {
				log.Info("err:", err)
				break
			}

			setting := new(model.Setting)
			changeDoc.ToStruct(setting)
			log.Info("App event tracking code:", setting)
		}
	}

	go dbtrip.ChangeStream(pipeline, streamer)

	return mgstream
}
