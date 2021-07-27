package health

import "sync"

var GlobalHealth = newHealthServer()

type status struct {
	OK  bool   `json:"bool"`
	Msg string `json:"msg"`
}

func newHealthServer() *health {
	result := &health{
		statuses:   make(map[string]bool),
		NotifyOK:   make(chan string, 256),
		NotifyFail: make(chan string, 256),
	}
	result.run()
	return result
}

type health struct {
	sync.Mutex
	statuses   map[string]bool
	NotifyOK   chan string
	NotifyFail chan string
}

func (h *health) IsSystemOK() status {
	h.Lock()
	defer h.Unlock()
	for name, ok := range h.statuses {
		if !ok {
			return status{OK: false, Msg: name}
		}
	}
	return status{OK: true}
}

func (h *health) run() {
	go func() {
		for {
			select {
			case name := <-h.NotifyOK:
				h.changeStatus(name, true)
			case name := <-h.NotifyFail:
				h.changeStatus(name, false)
			}
		}
	}()
}

func (h *health) changeStatus(name string, val bool) {
	h.Lock()
	defer h.Unlock()
	h.statuses[name] = val
}
