package influxdb

import (
	"time"

	iotmsg "github.com/alivinco/iotmsglibgo"
	influx "github.com/influxdata/influxdb/client/v2"
)

// DefaultTransform - transforms IotMsg into InfluxDb datapoint
func DefaultTransform(topic string, iotMsg *iotmsg.IotMsg, domain string) (*influx.Point, error) {
	var dpName string
	tags := map[string]string{
		"topic":  topic,
		"domain": domain,
		"cls":    iotMsg.Class,
		"subcls": iotMsg.SubClass,
		"type":   MapIotMsgType(iotMsg.Type),
	}
	var fields map[string]interface{}
	switch iotMsg.Class {
	case "sensor":
		fields = map[string]interface{}{
			"value": iotMsg.GetDefaultFloat(),
			"unit":  iotMsg.Default.Unit,
		}
		dpName = "sensors"
    case "binary":
		fields = map[string]interface{}{
			"value": iotMsg.GetDefaultBool(),
		}
		dpName = "binary"    
	default:
		fields = nil
	}
	if fields != nil {
		point, _ := influx.NewPoint(dpName, tags, fields, time.Now())
		return point, nil
	}
	return nil, nil

}

// MapIotMsgType maps int to string
func MapIotMsgType(msgType iotmsg.IotMsgType) string {
	typeMap := map[iotmsg.IotMsgType]string{
		iotmsg.MsgTypeCmd: "cmd",
		iotmsg.MsgTypeEvt: "evt",
		iotmsg.MsgTypeGet: "get",
	}
	return typeMap[msgType]
}
