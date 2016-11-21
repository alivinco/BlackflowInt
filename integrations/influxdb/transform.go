package influxdb

import (
	"time"

	iotmsg "github.com/alivinco/iotmsglibgo"
	influx "github.com/influxdata/influxdb/client/v2"
)

// DefaultTransform - transforms IotMsg into InfluxDb datapoint
func DefaultTransform(context *MsgContext, topic string, iotMsg *iotmsg.IotMsg, domain string) (*influx.Point, error) {
	var dpName string
	tags := map[string]string{
		"topic":  topic,
		"domain": domain,
		"cls":    iotMsg.Class,
		"subcls": iotMsg.SubClass,
		"type":   MapIotMsgType(iotMsg.Type),
	}
	var fields map[string]interface{}
	var vInt int
	var err error
	switch iotMsg.Class {
	case "sensor":
		fields = map[string]interface{}{
			"value": iotMsg.GetDefaultFloat(),
			"unit":  iotMsg.Default.Unit,
		}
		dpName = "sensor"
		break
	case "binary":
		fields = map[string]interface{}{
			"value": iotMsg.GetDefaultBool(),
		}
		dpName = "binary"
		break
	case "level":
		vInt, err = iotMsg.GetDefaultInt()
		if err != nil {
			return nil, err
		}
		fields = map[string]interface{}{
			"value": vInt,
		}
		dpName = "level"
		break
	default:
		fields = map[string]interface{}{
			"value": iotMsg.GetDefaultStr(),
		}
		dpName = "default"
	}
	if fields != nil {
		if context.MeasurementName == "" {
			context.MeasurementName = dpName
		} else {
			dpName = context.MeasurementName
		}
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
