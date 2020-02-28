package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Schema struct {
	Db struct {
		Host      string `mapstructure:"host"`
		User      string `mapstructure:"user"`
		Name      string `mapstructure:"name_db"`
		Password  string `mapstructure:"password"`
		DebugMode bool   `mapstructure:"debug"`
	} `mapstructure:"mysql"`

	Encryption struct {
		OIDKey    string `mapstructure:"oid_key"`
		JWTSecret string `mapstructure:"jwt_secret"`
		JWTExp    int    `mapstructure:"jwt_exp"`
		JWTPol    string `mapstructure:"jwt_pol"`
	} `mapstructure:"encryption"`

	Profiler struct {
		Prometheus    bool   `mapstructure:"prometheus"`
		StatsdAddress string `mapstructure:"statsd_address"`
		Service       string `mapstructure:"service"`
	} `mapstructure:"profiler"`

	Partner struct {
		URLDevices   string `mapstructure:"url_devices"`
		URLMetricGPS string `mapstructure:"url_metric_gps"`
	} `mapstructure:"partner"`

	TimeZone struct {
		ZoneST  string `mapstructure:"time_zone_str"`
		ZoneNum int    `mapstructure:"time_zone_num"`
		Layout  string `mapstructure:"layout"`
	} `mapstructure:"time_zone"`

	Mongodb struct {
		Host        string `mapstructure:"host"`
		Port        int    `mapstructure:"port"`
		Database    string `mapstructure:"database"`
		DatabaseIOT string `mapstructure:"database_oit"`
		User        string `mapstructure:"user"`
		Pass        string `mapstructure:"pass"`
	} `mapstructure:"mongodb"`
}

var Config Schema

func init() {
	config := viper.New()
	config.SetConfigName("config")
	config.AddConfigPath(".")          // Look for config in current directory
	config.AddConfigPath("config/")    // Optionally look for config in the working directory.
	config.AddConfigPath("../config/") // Look for config needed for tests.
	config.AddConfigPath("../")        // Look for config needed for tests.

	config.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	config.AutomaticEnv()

	err := config.ReadInConfig() // Find and read the config file
	if err != nil {              // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	err = config.Unmarshal(&Config)
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}
