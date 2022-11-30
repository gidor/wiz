package events

import (
	"sync"
)

// Listener is the type of a Listener, it's a func which receives any,optional, arguments from the caller/emmiter
type Listener func(string, ...interface{}) error

type dispatcher struct {
	events map[string][]Listener
	mutex  sync.RWMutex
}

func (d *dispatcher) clearAll() {
	for k := range d.events {
		delete(d.events, k)
	}

}
func (d *dispatcher) clear(event string) {
	delete(d.events, event)
}

func (d *dispatcher) addEvent(events ...string) {

	d.mutex.Lock()
	defer d.mutex.Unlock()

	if d.events == nil {
		d.events = make(map[string][]Listener, 5)
	}

	for _, event := range events {
		l, ok := d.events[event]
		if !ok {
			d.events[event] = make([]Listener, 5)
		} else {
			if l == nil {
				d.events[event] = make([]Listener, 5)
			}
		}

	}

}

func (d *dispatcher) addListener(event string, listener ...Listener) {
	if len(listener) == 0 {
		return
	}

	d.mutex.Lock()
	defer d.mutex.Unlock()

	if d.events == nil {
		d.events = make(map[string][]Listener, 5)
	}

	listeners := d.events[event]

	if listeners == nil {
		listeners = make([]Listener, 5)
	}
	d.events[event] = append(listeners, listener...)
}

func (d *dispatcher) trigger(event string, data ...interface{}) {

	d.mutex.RLock()
	defer d.mutex.RUnlock()

	if d.events != nil {
		listeners := d.events[event]
		for i := range listeners {
			l := listeners[i]
			if l != nil {
				l(event, data...)
			}
		}
	}
}

var (
	__disp__ dispatcher
)

func AddListener(event string, listener ...Listener) {
	__disp__.addListener(event, listener...)
}

func On(event string, listener ...Listener) {
	__disp__.addListener(event, listener...)
}

func Trigger(event string, data ...interface{}) {
	__disp__.trigger(event, data...)
}

func ClearAll() {
	__disp__.clearAll()
}
func Clear(event string) {
	__disp__.clear(event)
}
