package main

import (
	"flag"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/alivinco/blackflowint/integrations/influxdb"
	"github.com/alivinco/blackflowint/models"
	"github.com/labstack/echo"
	"github.com/spf13/viper"
)

var conf models.MainConfig

// SetupLog configures default logger
func SetupLog() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, ForceColors: true})
	log.SetLevel(log.DebugLevel)
}

func ConfigureApp() {
	var configFileDir string
	viper.SetDefault("StorageLocation", "./")
	viper.SetDefault("AdminRestAPIBindAddres", ":8091")
	flag.StringVar(&configFileDir, "c", "./", "Config file")

	flag.Parse()
	viper.SetConfigName("blackflowint")
	viper.AddConfigPath(configFileDir)
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	conf = models.MainConfig{}

	err = viper.Unmarshal(&conf)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
	log.Infof("Starting app using config parameters : %+v", conf)

	// v := viper.New()

}

func main() {
	SetupLog()
	ConfigureApp()
	SetupLog()
	adminRestHandler := echo.New()

	influxdb.Boot(&conf, adminRestHandler)

	adminRestHandler.Logger.Fatal(adminRestHandler.Start(conf.AdminRestAPIBindAddres))

}
