package hub

import (
	"sync"

	"github.com/google/uuid"
	"github.com/melfish/br-api/internal/models"
)

type key struct {
	recipientType models.RecipientType
	recipientID   uuid.UUID
}

type Hub struct {
	mu   sync.Mutex
	subs map[key][]chan *models.Notification
}

func New() *Hub {
	return &Hub{subs: make(map[key][]chan *models.Notification)}
}

func (h *Hub) Subscribe(recipientType models.RecipientType, recipientID uuid.UUID) chan *models.Notification {
	ch := make(chan *models.Notification, 8)
	k := key{recipientType, recipientID}
	h.mu.Lock()
	h.subs[k] = append(h.subs[k], ch)
	h.mu.Unlock()
	return ch
}

func (h *Hub) Unsubscribe(recipientType models.RecipientType, recipientID uuid.UUID, ch chan *models.Notification) {
	k := key{recipientType, recipientID}
	h.mu.Lock()
	defer h.mu.Unlock()
	chans := h.subs[k]
	for i, c := range chans {
		if c == ch {
			h.subs[k] = append(chans[:i], chans[i+1:]...)
			break
		}
	}
}

func (h *Hub) Publish(n *models.Notification) {
	k := key{n.RecipientType, n.RecipientID}
	h.mu.Lock()
	chans := make([]chan *models.Notification, len(h.subs[k]))
	copy(chans, h.subs[k])
	h.mu.Unlock()
	for _, ch := range chans {
		select {
		case ch <- n:
		default:
		}
	}
}
