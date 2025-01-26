package events

import "sync"

var eventPool = sync.Pool{
	New: func() interface{} {
		return New()
	},
}

func GetEvent() *Event {
	return eventPool.Get().(*Event)
}

func PutEvent(e *Event) {
	e.Reset()
	eventPool.Put(e)
}
