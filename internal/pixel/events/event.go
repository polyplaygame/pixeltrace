package events

import (
	"encoding/json"
	"sync"
	"time"
)

const (
	IPGeoExtraKey = "ip_geo"
)

type Event struct {
	mrw       *sync.RWMutex   `json:"-"`         // extra字段的读写锁
	App       string          `json:"app"`       // 事件所属应用
	IP        string          `json:"ip"`        // 事件IP
	Timestamp int64           `json:"timestamp"` // 事件时间戳,毫秒
	Data      json.RawMessage `json:"data"`      // 事件数据
	Extra     map[string]any  `json:"extra"`     // 事件额外数据
}

func New() *Event {
	return &Event{
		mrw:       new(sync.RWMutex),
		Extra:     make(map[string]any),
		Timestamp: time.Now().UnixMilli(),
	}
}

// DeepCopy 深拷贝
func (e *Event) DeepCopy() *Event {
	clone := &Event{
		mrw:       new(sync.RWMutex),
		App:       e.App,
		IP:        e.IP,
		Timestamp: e.Timestamp,
		Extra:     make(map[string]any),
	}
	if e.Data != nil {
		clone.Data = make(json.RawMessage, len(e.Data))
		copy(clone.Data, e.Data)
	}
	e.mrw.RLock()
	defer e.mrw.RUnlock()
	for k, v := range e.Extra {
		clone.Extra[k] = v
	}
	return clone
}

func (e *Event) Reset() {
	e.mrw = nil
	e.App = ""
	e.IP = ""
	e.Timestamp = 0
	e.Data = nil
	e.Extra = nil
}

func (e *Event) SetApp(app string) *Event {
	e.App = app
	return e
}

func (e *Event) SetIP(ip string) *Event {
	e.IP = ip
	return e
}

func (e *Event) SetTimestamp(timestamp int64) *Event {
	e.Timestamp = timestamp
	return e
}

func (e *Event) SetData(data json.RawMessage) *Event {
	e.Data = data
	return e
}

func (e *Event) SetExtra(key string, value any) *Event {
	e.mrw.Lock()
	defer e.mrw.Unlock()
	e.Extra[key] = value
	return e
}
