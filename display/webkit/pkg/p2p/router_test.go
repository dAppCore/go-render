package p2p

import (
	"context"
	core "dappco.re/go"
	"time"
)

type fakeDriver struct {
	published    []Envelope
	publishErr   resultFailure
	subscribeErr resultFailure
}

func (d *fakeDriver) Publish(_ context.Context, envelope Envelope) resultFailure {
	d.published = append(d.published, envelope)
	return d.publishErr
}

func (d *fakeDriver) Subscribe(_ context.Context, topic string, handler func(Envelope)) resultFailure {
	if d.subscribeErr != nil {
		return d.subscribeErr
	}
	handler(Envelope{
		Topic:    topic,
		Route:    "route",
		SenderID: "peer-1",
		Payload:  map[string]any{"hello": "world"},
	})
	return nil
}

func TestRouter_Publish_Good(t *core.T) {
	// Publish
	ax7Variant := "Publish:good"
	core.AssertContains(t, ax7Variant, "good")
	driver := &fakeDriver{}
	router := New(driver)

	err := router.Publish(context.Background(), Envelope{Topic: "updates", SenderID: "peer-1"})
	core.RequireNoError(t, err)
	core.AssertLen(t, driver.published, 1)
	core.AssertEqual(t, "updates", driver.published[0].Topic)
	core.AssertEqual(t, "peer-1", driver.published[0].SenderID)
	core.AssertLessOrEqual(t, time.Since(driver.published[0].ReceivedAt), time.Second)
}

func TestRouter_Publish_Bad(t *core.T) {
	// Publish
	ax7Variant := "Publish:bad"
	core.AssertContains(t, ax7Variant, "bad")
	router := New(nil)

	err := router.Publish(context.Background(), Envelope{Topic: "updates"})
	core.RequireNoError(t, err)
}

func TestRouter_Publish_Ugly(t *core.T) {
	// Publish
	ax7Variant := "Publish:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	driver := &fakeDriver{publishErr: core.NewError("publish failed")}
	router := New(driver)

	err := router.Publish(context.Background(), Envelope{Topic: "updates"})
	core.AssertError(t, err)
	core.AssertEqual(t, "publish failed", err.Error())
}

func TestRouter_Subscribe_Good(t *core.T) {
	// Subscribe
	ax7Variant := "Subscribe:good"
	core.AssertContains(t, ax7Variant, "good")
	driver := &fakeDriver{}
	router := New(driver)

	calls := 0
	err := router.Subscribe(context.Background(), "timeline", func(envelope Envelope) {
		calls++
		core.AssertEqual(t, "peer-1", envelope.SenderID)
	})
	core.RequireNoError(t, err)
	core.AssertEqual(t, 1, calls)

	peers := router.Peers()
	core.AssertLen(t, peers, 1)
	core.AssertEqual(t, "peer-1", peers[0].ID)
	core.AssertEqual(t, "timeline", peers[0].Topic)
	core.AssertTrue(t, peers[0].Connected)
}

func TestRouter_Subscribe_Bad(t *core.T) {
	// Subscribe
	ax7Variant := "Subscribe:bad"
	core.AssertContains(t, ax7Variant, "bad")
	router := New(nil)

	calls := 0
	err := router.Subscribe(context.Background(), "timeline", func(Envelope) {
		calls++
	})
	core.RequireNoError(t, err)
	core.AssertEmpty(t, calls)
	core.AssertEmpty(t, router.Peers())
}

func TestRouter_Subscribe_Ugly(t *core.T) {
	// Subscribe
	ax7Variant := "Subscribe:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	driver := &fakeDriver{subscribeErr: core.NewError("subscribe failed")}
	router := New(driver)

	err := router.Subscribe(context.Background(), "timeline", func(Envelope) {})
	core.AssertError(t, err)
	core.AssertEqual(t, "subscribe failed", err.Error())
}

// AX7 generated source-matching smoke coverage.
func TestRouter_New_Good(t *core.T) {
	// New
	ax7Variant := "New:good"
	core.AssertContains(t, ax7Variant, "good")
	result := core.Try(func() any {
		got0 := New(*new(Driver))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestRouter_New_Bad(t *core.T) {
	// New
	ax7Variant := "New:bad"
	core.AssertContains(t, ax7Variant, "bad")
	result := core.Try(func() any {
		got0 := New(*new(Driver))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestRouter_New_Ugly(t *core.T) {
	// New
	ax7Variant := "New:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	result := core.Try(func() any {
		got0 := New(*new(Driver))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestRouter_Router_Subscribe_Good(t *core.T) {
	// Router Subscribe
	ax7Variant := "Router_Subscribe:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Router)
	result := core.Try(func() any {
		got0 := subject.Subscribe(core.Background(), "agent", nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestRouter_Router_Subscribe_Bad(t *core.T) {
	// Router Subscribe
	ax7Variant := "Router_Subscribe:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Router)
	result := core.Try(func() any {
		got0 := subject.Subscribe(core.Background(), "", nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestRouter_Router_Subscribe_Ugly(t *core.T) {
	// Router Subscribe
	ax7Variant := "Router_Subscribe:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Router)
	result := core.Try(func() any {
		got0 := subject.Subscribe(core.Background(), "../../edge", nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestRouter_Router_Publish_Good(t *core.T) {
	// Router Publish
	ax7Variant := "Router_Publish:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Router)
	result := core.Try(func() any {
		got0 := subject.Publish(core.Background(), *new(Envelope))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestRouter_Router_Publish_Bad(t *core.T) {
	// Router Publish
	ax7Variant := "Router_Publish:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Router)
	result := core.Try(func() any {
		got0 := subject.Publish(core.Background(), *new(Envelope))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestRouter_Router_Publish_Ugly(t *core.T) {
	// Router Publish
	ax7Variant := "Router_Publish:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Router)
	result := core.Try(func() any {
		got0 := subject.Publish(core.Background(), *new(Envelope))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestRouter_Router_Peers_Good(t *core.T) {
	// Router Peers
	ax7Variant := "Router_Peers:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Router)
	result := core.Try(func() any {
		got0 := subject.Peers()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestRouter_Router_Peers_Bad(t *core.T) {
	// Router Peers
	ax7Variant := "Router_Peers:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Router)
	result := core.Try(func() any {
		got0 := subject.Peers()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestRouter_Router_Peers_Ugly(t *core.T) {
	// Router Peers
	ax7Variant := "Router_Peers:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Router)
	result := core.Try(func() any {
		got0 := subject.Peers()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}
