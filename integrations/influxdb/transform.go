package influxdb

import (
	"time"
	influx "github.com/influxdata/influxdb/client/v2"
	"github.com/alivinco/fimpgo"
)

// DefaultTransform - transforms IotMsg into InfluxDb datapoint
func DefaultTransform(context *MsgContext, topic string, iotMsg *fimpgo.FimpMessage, domain string) (*influx.Point, error) {
	tags := map[string]string{
		"topic":  topic,
		"domain": domain,
		"mtype":    iotMsg.Type,
		"serv": iotMsg.Service,
	}
	var fields map[string]interface{}
	var vInt int64
	var err error
	switch iotMsg.Type {
	case "evt.sensor.report","evt.meter.report":
		val ,err := iotMsg.GetFloatValue()
		if err != nil {
			fields = map[string]interface{}{
				"value": val,
				"unit":  iotMsg.Properties["unit"],
			}
		}

	case "evt.binary.report","evt.presence.report","evt.open.report":
		val ,err := iotMsg.GetBoolValue()
		if err != nil {
			fields = map[string]interface{}{
				"value": val,
			}
		}

	case "evt.lvl.report":
		vInt, err = iotMsg.GetIntValue()
		if err != nil {
			return nil, err
		}
		fields = map[string]interface{}{
			"value": vInt,
		}

	}
	if fields != nil {
		point, _ := influx.NewPoint(iotMsg.Service, tags, fields, time.Now())
		return point, nil
	}

	return nil, nil

}


