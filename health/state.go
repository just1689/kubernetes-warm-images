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
			return status{
				OK:  false,
				Msg: name,
			}
		}
	}
	return status{OK: true}
}

func (h *health) run() {
	go func() {
		for {
			select {
			case next := <-h.NotifyOK:
				h.Lock()
				h.statuses[next] = true
				h.Unlock()
			case next := <-h.NotifyFail:
				h.Lock()
				h.statuses[next] = false
				h.Unlock()
			}
		}
	}()
}
