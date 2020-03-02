package mysql

import (
	"fmt"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/labstack/gommon/log"
	"github.com/tientp-floware/mgodb-stream/config"
)

var (
	once   sync.Once
	dbPoll *gorm.DB
)

func init() {
	gorm.NowFunc = func() time.Time {
		return time.Now().UTC().Truncate(1000 * time.Nanosecond)
	}
	/* db, err := gorm.Open("mysql", fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s",
	config.Config.Db.Host, config.Config.Db.User, config.Config.Db.Name, config.Config.Db.Password)) */
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local", config.Config.Db.User,
		config.Config.Db.Password, config.Config.Db.Host, config.Config.Db.Name))
	if err != nil {
		fmt.Println(err)
		log.Panicf("Connect not success to postgres database at host:%s with user:%s and db:%s",
			config.Config.Db.Host, config.Config.Db.User, config.Config.Db.Name)
	}

	db.DB().SetMaxIdleConns(20)

	db.DB().SetMaxOpenConns(100)

	db.DB().Ping()

	if config.Config.Db.DebugMode {
		db.LogMode(config.Config.Db.DebugMode)
	}
	dbPoll = db
	log.Info("init connected db")
}

func new() *gorm.DB {
	gorm.NowFunc = func() time.Time {
		return time.Now().UTC().Truncate(1000 * time.Nanosecond)
	}
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True", config.Config.Db.User,
		config.Config.Db.Password, config.Config.Db.Host, config.Config.Db.Name))
	if err != nil {
		log.Infof("Connect not success to postgres database at host:%s with user:%s and db:%s",
			config.Config.Db.Host, config.Config.Db.User, config.Config.Db.Name)
	}
	// db.SingularTable(true)
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	// db.DB().SetConnMaxLifetime(time.Nanosecond)
	// db.DB().SetConnMaxLifetime(3 * time.Second)
	db.DB().Ping()
	// db.DB().Exec("SET timezone TO 'Asia/Ho_Chi_Minh';")
	if config.Config.Db.DebugMode {
		db.LogMode(true)
	}
	return db
}

func GetDB() *gorm.DB {
	once.Do(func() {
		dbPoll = new()
	})
	return dbPoll
}

// IntervalTime return range time
func IntervalTime(n int64) (b, e string) {
	nowTime := time.Now()
	afterTime := nowTime.Add(time.Duration(n) * time.Minute)
	return afterTime.Format(`2006-01-02T15:04:05.999`), nowTime.Format(`2006-01-02T15:04:05.999`)
}
