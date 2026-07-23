package p2p

import (
	"context"
	"sync"
	"time"
)

type Envelope struct {
	Topic      string         `json:"topic"`
	Route      string         `json:"route"`
	SenderID   string         `json:"sender_id"`
	Payload    map[string]any `json:"payload"`
	ReceivedAt time.Time      `json:"received_at"`
}

type Peer struct {
	ID        string    `json:"id"`
	Topic     string    `json:"topic"`
	Connected bool      `json:"connected"`
	SeenAt    time.Time `json:"seen_at"`
}

type Driver interface {
	Publish(context.Context, Envelope) resultFailure
	Subscribe(context.Context, string, func(Envelope)) resultFailure
}

type Router struct {
	driver Driver
	mu     sync.RWMutex
	peers  map[string]Peer
}

func New(driver Driver) *Router {
	return &Router{
		driver: driver,
		peers:  make(map[string]Peer),
	}
}

func (r *Router) Subscribe(ctx context.Context, topic string, handler func(Envelope)) resultFailure {
	if r.driver == nil {
		return nil
	}
	return r.driver.Subscribe(ctx, topic, func(envelope Envelope) {
		r.mu.Lock()
		r.peers[envelope.SenderID] = Peer{
			ID:        envelope.SenderID,
			Topic:     topic,
			Connected: true,
			SeenAt:    time.Now(),
		}
		r.mu.Unlock()
		handler(envelope)
	})
}

func (r *Router) Publish(ctx context.Context, envelope Envelope) resultFailure {
	if r.driver == nil {
		return nil
	}
	envelope.ReceivedAt = time.Now()
	return r.driver.Publish(ctx, envelope)
}

func (r *Router) Peers() []Peer {
	r.mu.RLock()
	defer r.mu.RUnlock()
	peers := make([]Peer, 0, len(r.peers))
	for _, peer := range r.peers {
		peers = append(peers, peer)
	}
	return peers
}
