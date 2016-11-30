## Blackflow Integrator 

### The application is container for blackflow integrations . 

#### BlackflowInt application config :

The application can be configured either using -c flag with config file location or using ENV variables instead . 

-c should point to folder where config file resides , config file should have a name blackflowint.toml

Start application using config file : 
```
blackflowint -c ./
``` 

Supported ENV variables : 
- ZM_LOGLEVEL="info"
- ZM_STORAGELOCATION="/var/lib/blackflowint"
- ZM_ADMINRESTAPIBINDADDRES=":5015"


```
type MainConfig struct {
	StorageLocation        string
	AdminRestAPIBindAddres string
	LogPath                string
	LogPath                string
}

```



#### Integrations : 

- [MQTT event stream dump into InfluxDB](Integrations/influxdb/README.md) 

#### Docker 

- start container : docker run --name blackflowint -d -p 5016:5015 --link influxdb alivinco/blackflowint