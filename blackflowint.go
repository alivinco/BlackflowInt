package main

import (
	"flag"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/alivinco/blackflowint/integrations/influxdb"
	"github.com/alivinco/blackflowint/models"
	"github.com/labstack/echo"
	"github.com/spf13/viper"

	"github.com/alivinco/blackflowint/integrations/restmqttproxy"
)

var conf models.MainConfig

// SetupLog configures default logger
func SetupLog() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, ForceColors: true})
}

// ConfigureApp configures the application .
func ConfigureApp() {
	var configFileDir string
	conf = models.MainConfig{}
	viper.AutomaticEnv()
	// Is used as default broker config
	viper.SetEnvPrefix("zm")
	viper.SetDefault("StorageLocation", "./")
	viper.SetDefault("AdminRestAPIBindAddres", ":5015")
	viper.SetDefault("LogLevel", "info")

	// if config file is set then use config otherwise , otherwise use ENV var
	flag.StringVar(&configFileDir, "c", "", "Config file")
	flag.Parse()
	if configFileDir != "" {
		viper.SetConfigName("blackflowint")
		viper.AddConfigPath(configFileDir)
		err := viper.ReadInConfig() // Find and read the config file
		if err != nil {             // Handle errors reading the config file
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}
		err = viper.Unmarshal(&conf)
		if err != nil {
			log.Fatalf("unable to decode into struct, %v", err)
		}
	} else {
		conf.AdminRestAPIBindAddres = viper.GetString("AdminRestAPIBindAddres")
		conf.StorageLocation = viper.GetString("StorageLocation")
		conf.LogLevel = viper.GetString("LogLevel")
	}

	log.Infof("Starting app using config parameters : %+v", conf)

	logLevel, err := log.ParseLevel(conf.LogLevel)
	if err != nil {
		logLevel = log.InfoLevel
	}

	log.SetLevel(logLevel)
}

func main() {
	SetupLog()
	ConfigureApp()
	SetupLog()
	adminRestHandler := echo.New()

	influxdb.Boot(&conf, adminRestHandler)
	restmqttproxy.Boot(&conf,adminRestHandler)

	adminRestHandler.Logger.Fatal(adminRestHandler.Start(conf.AdminRestAPIBindAddres))

}
