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

- MQTT event stream dump into InfluxDB



### MQTT event stream dump into InfluxDB

**Message processing pipeline** : 
````


                        ------------------
                        |                |
       -----SUB---------|   Selectors    |
       |                | (subsriptions) |
       |                ------------------
      \|/ 
---------------             ------------                    -------------                      ---------                 ---------------       
|             |             |          |                    |           |                      |       |                 |   DataPoint |        
| MQTT broker | ---MSG----> |  Filters | --Context+IotMsg-->| Transform |--Context+DataPoint-->| Write | ---DataPoint--> |     Batch   |
|             |             |          |                    |           |                      |       |                 |             |       
---------------             ------------                    -------------                      ---------                 ---------------       
																															  |	
                           ------------------------Batch size trigger----------------------------------------------------------
                           |
                          \|/
                      --------------           ------------
   ---------          |   Write    |           |          |
   | Timer |--------->|   Batch    |---BATCH-->| InfluxDB |
   ---------          |            |           |          |
                      --------------           ------------





````

**Process configuration**

````                      

                          ------------------------------------------------
                          | ProcessConfig,Selectors,Filters,Measurements | 
                          ------------------------------------------------
                               |
                              \|/ 
----------------         ------------------------          -------------
| MQTT broker  |-------->| Process A instance 1 |--------->| InfluxDB  | 
|              |         |                      |          |           |
----------------         ------------------------          -------------

                          ------------------------------------------------
                          | ProcessConfig,Selectors,Filters,Measurements | 
                          ------------------------------------------------
                               |
                              \|/ 
----------------         ------------------------          -------------
| MQTT broker  |-------->| Process A instance 2 |--------->| InfluxDB  | 
|              |         |                      |          |           |
----------------         ------------------------          -------------

````
Configuration data model :

````

   --------------------------                                               
   |     ProcessConfig      |
   |________________________|
     |                 |
    /|\               /|\
-------------    ------------ 
| Selectors |    |  Filters | 
-------------    ------------ 
                    |    /|\
                    |     |
                    -------

````  

Filter structure : 

```
type Filter struct {
	ID          int
	Name        string
	Topic       string
	Domain      string
	MsgType     string
	MsgClass    string
	MsgSubClass string
	// If true then returns everythin except matching value
	Negation bool
	// Boolean operation between 2 filters , supported values : and , or
	LinkedFilterBooleanOperation string
	LinkedFilterID               int
	IsAtomic                     bool
	// Optional field , all tags defined here will be converted into influxDb tags
	Tags map[string]string
	// If set , then the value will overrride default measurement name defined in transformation
	MeasurementName string
	// definies if filter is temporary and can be stored in memory or persisted to disk
	InMemory bool
}
```

Selector structure :

``` 
type Selector struct {
	ID    IDt
	Topic string
	// definies if filter is temporary and can be stored in memory or persisted to disk
	InMemory bool
}
```

Tags : 
 TODO 

Measurement structure : 

```
// Measurement stores measurement specific configs like retention policy
type Measurement struct {
	Name string
	// Normally should be composed of Name + RetentionPolicy , for instance : sensor_1d
	RetentionPolicyName string
	// Have to be in format 1d , 1w , 1h .
	RetentionPolicyDuration string
}

```

Process config :

```
type ProcessConfig struct {
	ID                 string
	MqttBrokerAddr     string
	MqttClientID       string
	MqttBrokerUsername string
	MqttBrokerPassword string
	InfluxAddr         string
	InfluxUsername     string
	InfluxPassword     string
	InfluxDB           string
	// DataPoints are saved in batches .
	// Batch is sent to DB once it reaches BatchMaxSize or SaveInterval .
	// depends on what comes first .
	BatchMaxSize int
	// Interval in miliseconds
	SaveInterval time.Duration
	Filters      []Filter
	Selectors    []Selector
	Measurements []Measurement
}

```

Message context :
```
type MsgContext struct {
	FilterID        IDt
	MeasurementName string
}

```

