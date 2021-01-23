package event

import (
	"sync"

	"github.com/Mahamed-Belkheir/sunduq"
)

type handlerIndex struct {
	channelMap map[int]chan sunduq.Message
	mut        *sync.RWMutex
}

func (h handlerIndex) send(envelope sunduq.Envelope) {
	h.mut.RLock()
	defer h.mut.RUnlock()
	handler, ok := h.channelMap[envelope.ID]
	if !ok {
		//log it
		return
	}
	handler <- envelope.Message
}

func (h handlerIndex) register(id int, channel chan sunduq.Message) {
	h.mut.Lock()
	defer h.mut.Unlock()
	h.channelMap[id] = channel
}

func (h handlerIndex) unregister(id int) {
	h.mut.Lock()
	defer h.mut.Unlock()
	delete(h.channelMap, id)
}
