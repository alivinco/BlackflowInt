## Blackflow Integrator 

### The application is container for blackflow integrations . 

#### Integrations : 

- MQTT event stream dump into InfluxDB



#### MQTT event stream dump into InfluxDB

**Event processing pipeline** : 
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

