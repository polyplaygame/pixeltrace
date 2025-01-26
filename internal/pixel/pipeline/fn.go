package pipeline

import (
	"pixeltrace/internal/geoip"
	"pixeltrace/internal/pixel/events"

	jsoniter "github.com/json-iterator/go"

	"github.com/cloudwego/hertz/pkg/common/hlog"
)

// addIPGeoExtra 添加IP地理信息
var addIPGeoExtra = func(e *events.Event) *events.Event {
	if e.IP == "" {
		return e
	}
	result, err := geoip.Lookup(e.IP, "en")
	if err != nil {
		hlog.Errorf("[stream]ip lookup error: %v", err)
		return e
	}
	e.SetExtra(events.IPGeoExtraKey, result)
	return e
}

// unmarshalEventData 反序列化事件数据
// ! why? bi那边反馈直接[]byte 会有乱码，所以这里做了一次json序列化
var unmarshalEventData = func(e *events.Event) *events.Event {
	if e.Data == nil {
		return e
	}
	var m map[string]any
	if err := jsoniter.Unmarshal(e.Data, &m); err != nil {
		hlog.Errorf("[stream]failed to unmarshal event data: %v", err)
		return e
	}
	data, err := jsoniter.Marshal(m)
	if err != nil {
		hlog.Errorf("[stream]failed to marshal event data: %v", err)
		return e
	}
	e.Data = data
	return e
}
