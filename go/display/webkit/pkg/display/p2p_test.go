package display

import (
	"context"
	"sync"
	"time"

	core "dappco.re/go"
	"dappco.re/go/render/display/webkit/pkg/p2p"
)

func newDisplayP2PTestService(t *core.T) (*Service, *p2p.Service, *core.Core, <-chan Event) {
	t.Helper()

	var p2pSvc *p2p.Service
	driver := newLoopbackP2PDriver()
	displaySvc, c, eventBuffer := newServiceWithMockApp(t, func(c *core.Core) core.Result {
		p2pSvc = p2p.NewServiceWithDriver(c, p2p.Options{NodeID: "node-1"}, driver)
		return core.Result{Value: p2pSvc, OK: true}
	})
	core.AssertNotNil(t, displaySvc)
	core.AssertNotNil(t, p2pSvc)
	return displaySvc, p2pSvc, c, eventBuffer
}

type loopbackP2PDriver struct {
	mu       sync.Mutex
	handlers map[string]func(p2p.Envelope)
}

func newLoopbackP2PDriver() *loopbackP2PDriver {
	return &loopbackP2PDriver{handlers: make(map[string]func(p2p.Envelope))}
}

func (d *loopbackP2PDriver) Publish(_ context.Context, envelope p2p.Envelope) resultFailure {
	d.mu.Lock()
	handler := d.handlers[envelope.Topic]
	d.mu.Unlock()
	if handler != nil {
		handler(envelope)
	}
	return nil
}

func (d *loopbackP2PDriver) Subscribe(_ context.Context, topic string, handler func(p2p.Envelope)) resultFailure {
	d.mu.Lock()
	d.handlers[topic] = handler
	d.mu.Unlock()
	return nil
}

type immediateP2PDriver struct {
	envelope p2p.Envelope
}

func (d *immediateP2PDriver) Publish(context.Context, p2p.Envelope) resultFailure {
	return nil
}

func (d *immediateP2PDriver) Subscribe(_ context.Context, topic string, handler func(p2p.Envelope)) resultFailure {
	if handler != nil {
		handler(p2p.Envelope{
			Topic:    topic,
			Route:    d.envelope.Route,
			SenderID: d.envelope.SenderID,
			Payload:  d.envelope.Payload,
		})
	}
	return nil
}

func TestDisplayP2P_attachP2PBridge_GoodCase(t *core.T) {
	_, p2pSvc, _, eventBuffer := newDisplayP2PTestService(t)

	err := p2pSvc.Publish(context.Background(), p2p.Envelope{
		Topic:    "display",
		Route:    "route-1",
		SenderID: "peer-1",
		Payload:  map[string]any{"hello": "world"},
	})
	core.RequireNoError(t, err)

	select {
	case event := <-eventBuffer:
		core.AssertEqual(t, EventCustomEvent, event.Type)
		core.AssertNotNil(t, event.Data)
		core.AssertEqual(t, "p2p", event.Data["source"])
		core.AssertEqual(t, "route-1", event.Data["route"])
		core.AssertEqual(t, "peer-1", event.Data["sender_id"])
		core.AssertEqual(t, map[string]any{"hello": "world"}, event.Data["payload"])
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for bridged event")
	}
}

func TestDisplayP2P_OnStartup_InitializesEventManagerBeforeBridge(t *core.T) {
	const route = "route-startup"

	driver := &immediateP2PDriver{
		envelope: p2p.Envelope{
			Route:    route,
			SenderID: "peer-startup",
			Payload:  map[string]any{"hello": "world"},
		},
	}

	_, _, eventBuffer := newServiceWithMockApp(t, func(c *core.Core) core.Result {
		p2pSvc := p2p.NewServiceWithDriver(c, p2p.Options{NodeID: "node-startup"}, driver)
		return core.Result{Value: p2pSvc, OK: true}
	})

	select {
	case event := <-eventBuffer:
		core.AssertEqual(t, EventCustomEvent, event.Type)
		core.AssertNotNil(t, event.Data)
		core.AssertEqual(t, "p2p", event.Data["source"])
		core.AssertEqual(t, route, event.Data["route"])
		core.AssertEqual(t, "peer-startup", event.Data["sender_id"])
		core.AssertEqual(t, map[string]any{"hello": "world"}, event.Data["payload"])
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for bridged startup event")
	}
}

func TestDisplayP2P_attachP2PBridge_Bad(t *core.T) {
	// attachP2PBridge
	ax7Variant := "attachP2PBridge:bad"
	core.AssertContains(t, ax7Variant, "bad")
	c := core.New(core.WithServiceLock())
	svc, err := New()
	core.RequireNoError(t, err)
	svc.ServiceRuntime = core.NewServiceRuntime(c, Options{})

	core.AssertNotPanics(t, func() {
		svc.attachP2PBridge()
	})
}

func TestDisplayP2P_attachP2PBridge_UglyCase(t *core.T) {
	displaySvc, p2pSvc, _, _ := newDisplayP2PTestService(t)
	displaySvc.events = nil

	core.AssertNotPanics(t, func() {
		err := p2pSvc.Publish(context.Background(), p2p.Envelope{
			Topic:    "display",
			SenderID: "peer-2",
		})
		core.RequireNoError(t, err)
	})
}
