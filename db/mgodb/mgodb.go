package mgodb

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"

	"github.com/labstack/gommon/log"
	"github.com/tientp-floware/mgodb-stream/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	//log               = zap.NewExample()
	host              string
	authdb            string
	dbname            string
	usrname           string
	usrpass           string
	fullconnectstring string
)

//Init mongo client
func init() {
	authdb = "admin"
	usrname = config.Config.Mongodb.User
	usrpass = config.Config.Mongodb.Pass
	dbname = config.Config.Mongodb.Database
	host = config.Config.Mongodb.Host
	fullconnectstring = fmt.Sprintf(`mongodb://%s:%s@%s/%s?retryWrites=true&w=majority&authSource=admin`, usrname, usrpass, host, dbname)
	connection()
}

// connection for mongodb official
func connection() {

	opt := options.Client()
	opt.ApplyURI(fullconnectstring)
	opt.SetTLSConfig(&tls.Config{})
	opt.SetMaxPoolSize(8)
	opt.SetMinPoolSize(3)
	opt.SetReadPreference(readpref.Nearest())
	//opt.SetDirect(true)
	err := opt.Validate()
	// Connect to MongoDB
	MgoClient, err = mongo.Connect(context.Background(), opt)
	if err != nil {
		log.Fatal(err)
	}
	// Check the connection
	err = MgoClient.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Driver offical - connected to MongoDB!")
	MgoDatabase = MgoClient.Database(dbname)
	log.Info("Using DB:", dbname)

}

// NewMogoDB returns a new layer database instance
func NewMogoDB(collection string) *MogoDB {
	return &MogoDB{c: MgoDatabase.Collection(collection)}
}

// NewSwitchMogoDB returns a new layer database instance
func NewSwitchMogoDB(dbname, collection string) *MogoDB {
	return &MogoDB{c: MgoSwitchDB(dbname).Collection(collection)}
}

var (
	// MgoClient connection client
	MgoClient *mongo.Client
	// MgoDatabase would you change db
	MgoDatabase *mongo.Database
)

// This package using offical mongodb driver for golang.
// Current have been released. We should switch to this driver.
// It's optimize and support full features for newest version mongodb.

// MogoDB hold layer struct
type MogoDB struct {
	c       *mongo.Collection
	filter  interface{}
	limit   int64
	project interface{}
	skip    int64
	sort    interface{}
}

// var result []bson.M
// filter = bson.D{{Key: "username", Value: "baba"}}
// project = bson.D{{Key: "_id", Value: 0}, {Key: "username", Value: 1}, {Key: "password", Value: 0}}
// c.Find(filter).Project(project).Skip(10).Limit(3).Decode(&result);

// MgoSwitchDB if you want to switch database
func MgoSwitchDB(dbName string) *mongo.Database {
	return MgoClient.Database(dbName)
}

// From you want to change collection
func (mg *MogoDB) From(coll string) *MogoDB {
	mg.c = MgoDatabase.Collection(coll)
	return mg
}

// MgoClient get client
func (mg *MogoDB) MgoClient() *mongo.Client {
	return MgoClient
}

// Find find with filter
func (mg *MogoDB) Find(filter interface{}) *MogoDB {
	mg.filter = filter
	return mg
}

// Limit sets sorting
func (mg *MogoDB) Limit(limit int64) *MogoDB {
	mg.limit = limit
	return mg
}

// Project sets sorting
func (mg *MogoDB) Project(project interface{}) *MogoDB {
	mg.project = project
	return mg
}

// Skip sets sorting
func (mg *MogoDB) Skip(skip int64) *MogoDB {
	mg.skip = skip
	return mg
}

// Sort sets sorting
func (mg *MogoDB) Sort(sort interface{}) *MogoDB {
	mg.sort = sort
	return mg
}

// Decode returns all docs
func (mg *MogoDB) Decode(result interface{}) error {
	opts := options.Find()
	if mg.limit > 0 {
		opts.SetLimit(mg.limit)
	}
	if mg.project != nil {
		opts.SetProjection(mg.project)
	}
	if mg.skip > 0 {
		opts.SetSkip(mg.skip)
	}
	if mg.sort != nil {
		opts.SetSort(mg.sort)
	}
	cur, err := mg.c.Find(nil, mg.filter, opts)
	if err != nil {
		return err
	}

	var docs []bson.M
	for cur.Next(nil) {
		var doc bson.M
		cur.Decode(&doc)
		docs = append(docs, doc)
	}
	b, _ := json.Marshal(docs)
	json.Unmarshal(b, result)
	return nil
}

// FindOne find with filter
func (mg *MogoDB) FindOne(filter interface{}) *mongo.SingleResult {
	return mg.c.FindOne(context.TODO(), filter)
}

// InsertOne find with filter
func (mg *MogoDB) InsertOne(doc interface{}) (*mongo.InsertOneResult, error) {
	return mg.c.InsertOne(context.TODO(), doc)
}

// InsertMany insert one document
func (mg *MogoDB) InsertMany(docs []interface{}) (*mongo.InsertManyResult, error) {
	return mg.c.InsertMany(context.TODO(), docs)
}

// UpdateOne update with filter
func (mg *MogoDB) UpdateOne(filter, updateDoc interface{}) (*mongo.UpdateResult, error) {
	return mg.c.UpdateOne(context.TODO(), filter, updateDoc)
}

// UpdateMany with filter
func (mg *MogoDB) UpdateMany(updateDocs, filter bson.M) (*mongo.UpdateResult, error) {
	return mg.c.UpdateMany(context.TODO(), filter, updateDocs)
}

// DeleteOne with filter
func (mg *MogoDB) DeleteOne(filter bson.M) (*mongo.DeleteResult, error) {
	return mg.c.DeleteOne(context.TODO(), filter)
}

// DeleteMany with filter
func (mg *MogoDB) DeleteMany(filter bson.M) (*mongo.DeleteResult, error) {
	return mg.c.DeleteMany(context.TODO(), filter)
}

// ChangeStream implement
func (mg *MogoDB) ChangeStream(pipeline interface{}, s func(*mongo.ChangeStream)) {
	ctx := context.Background()
	cur, err := mg.c.Watch(ctx, pipeline,
		options.ChangeStream().SetFullDocument(options.UpdateLookup))
	if err != nil {
		// Handle err
		return
	}
	defer cur.Close(ctx)
	//Handling change stream in a cycle
	for {
		select {
		case <-ctx.Done():
			err := cur.Close(ctx)
			if err != nil {
				log.Info("change stream err:", err)
				break
			}
		default:
			s(cur)
		}
	}
}
