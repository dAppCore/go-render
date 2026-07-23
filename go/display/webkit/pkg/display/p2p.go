package display

import (
	"context"

	core "dappco.re/go"
	"dappco.re/go/render/display/webkit/pkg/p2p"
)

func (s *Service) attachP2PBridge() {
	router, ok := core.ServiceFor[*p2p.Service](s.Core(), "p2p")
	if !ok || router == nil {
		return
	}
	if s.p2pBridgeCancel != nil {
		s.p2pBridgeCancel()
		s.p2pBridgeCancel = nil
	}
	ctx, cancel := context.WithCancel(context.Background())
	if err := router.Subscribe(ctx, "display", func(envelope p2p.Envelope) {
		if s.events == nil {
			return
		}
		s.events.Emit(Event{
			Type: EventCustomEvent,
			Data: map[string]any{
				"source":    "p2p",
				"topic":     envelope.Topic,
				"route":     envelope.Route,
				"sender_id": envelope.SenderID,
				"payload":   envelope.Payload,
			},
		})
	}); err != nil {
		cancel()
		if s.app != nil {
			s.app.Logger().Info("p2p bridge subscribe failed", "err", err)
		}
		return
	}
	s.p2pBridgeCancel = cancel
}
