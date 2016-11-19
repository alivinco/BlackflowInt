package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/alivinco/blackflowint/integrations/influxdb"
	"github.com/alivinco/blackflowint/models"
	"github.com/labstack/echo"
)

// SetupLog configures default logger
func SetupLog() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	log.SetLevel(log.DebugLevel)
}

func main() {

	conf := models.MainConfig{
		StorageLocation:        "./",
		AdminRestAPIBindAddres: ":8091",
	}
	SetupLog()
	adminRestHandler := echo.New()

	influxdb.Boot(&conf, adminRestHandler)

	adminRestHandler.Logger.Fatal(adminRestHandler.Start(conf.AdminRestAPIBindAddres))

}
