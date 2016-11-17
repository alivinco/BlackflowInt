## Blackflow Integrator 

### The application is container for blackflow integrations . 

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
---------------             --------------          ----------------        --------------           ------------
|             |             |            |          |              |        |   Batch    |           |          | 
| MQTT broker | ---MSG----> |  Filters   | --MSG--> |  Tag msg     |--MSG-->| aggregator |---BATCH-->| InfluxDB |
|             |             |            |          |              |        |            |           |          |
---------------             --------------          ----------------        --------------           ------------

````

**Process configuration**

````                      

                          -----------------------------------
                          | ProcessConfig,Selectors,Filters | 
                          -----------------------------------  
                               |
                              \|/ 
----------------         ------------------------          -------------
| MQTT broker  |-------->| Process A instance 1 |--------->| InfluxDB  | 
|              |         |                      |          |           |
----------------         ------------------------          -------------

                          -----------------------------------
                          | ProcessConfig,Selectors,Filters | 
                          -----------------------------------
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
}
```

Selector structure :

``` 
type Selector struct {
	Topic string
}
```

Tags : 
 TODO 

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
}

```

Application config :

```


```
