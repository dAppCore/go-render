package p2p

import (
	"context"

	core "dappco.re/go"
)

// NewService / NewServiceWithDriver wire a router over the driver and expose
// the configured node identity through State.
//
//	svc := NewServiceWithDriver(c, Options{NodeID: "node-1"}, &fakeDriver{})
func TestServiceBehaviour_NewService_Good(t *core.T) {
	c := core.New(core.WithServiceLock())
	svc := NewService(c, Options{ListenAddr: "127.0.0.1:0", NodeID: "node-1"})
	core.AssertNotNil(t, svc)

	state := svc.State()
	core.AssertEqual(t, "node-1", state.NodeID)
}

// Publish routes an envelope through the driver and Peers reflects subscribers.
func TestServiceBehaviour_PublishSubscribe_Good(t *core.T) {
	driver := &fakeDriver{}
	c := core.New(core.WithServiceLock())
	svc := NewServiceWithDriver(c, Options{NodeID: "node-1"}, driver)

	var got Envelope
	err := svc.Subscribe(context.Background(), "topic", func(e Envelope) { got = e })
	core.AssertNil(t, err)
	core.AssertEqual(t, "peer-1", got.SenderID)

	err = svc.Publish(context.Background(), Envelope{Topic: "topic", SenderID: "node-1"})
	core.AssertNil(t, err)
	core.AssertLen(t, driver.published, 1)

	// Subscribing recorded a peer; Peers surfaces it.
	core.AssertGreaterOrEqual(t, len(svc.Peers()), 1)
}

// Publish surfaces a driver failure.
func TestServiceBehaviour_Publish_Bad(t *core.T) {
	driver := &fakeDriver{publishErr: core.NewError("publish failed")}
	c := core.New(core.WithServiceLock())
	svc := NewServiceWithDriver(c, Options{}, driver)

	err := svc.Publish(context.Background(), Envelope{Topic: "t"})
	core.AssertNotNil(t, err)
	core.AssertContains(t, err.Error(), "publish failed")
}

// OnStartup registers the p2p.publish and p2p.state actions; OnShutdown is a
// safe no-op when the driver has no Close method.
func TestServiceBehaviour_StartupActions_Good(t *core.T) {
	driver := &fakeDriver{}
	c := core.New(core.WithServiceLock())
	svc := NewServiceWithDriver(c, Options{NodeID: "node-1"}, driver)

	core.RequireTrue(t, svc.OnStartup(context.Background()).OK)

	// p2p.publish routes through the driver.
	publish := c.Action("p2p.publish").Run(context.Background(), core.NewOptions(
		core.Option{Key: "topic", Value: "demo"},
		core.Option{Key: "payload", Value: map[string]any{"k": "v"}},
	))
	core.AssertTrue(t, publish.OK)
	core.AssertLen(t, driver.published, 1)

	// p2p.state returns the current State.
	state := c.Action("p2p.state").Run(context.Background(), core.NewOptions())
	core.AssertTrue(t, state.OK)
	value, ok := state.Value.(State)
	core.AssertTrue(t, ok)
	core.AssertEqual(t, "node-1", value.NodeID)

	core.AssertTrue(t, svc.OnShutdown(context.Background()).OK)
}

// OptionsFromEnv parses peer lists and trims empties.
//
//	CORE_P2P_PEERS="a, ,b" -> PeerAddrs ["a","b"]
func TestServiceBehaviour_OptionsFromEnv_Good(t *core.T) {
	t.Setenv("CORE_P2P_ADDR", "127.0.0.1:9000")
	t.Setenv("CORE_P2P_PEERS", "10.0.0.1:9000, ,10.0.0.2:9000")
	t.Setenv("CORE_P2P_NODE_ID", "env-node")

	opts := OptionsFromEnv()
	core.AssertEqual(t, "127.0.0.1:9000", opts.ListenAddr)
	core.AssertEqual(t, "env-node", opts.NodeID)
	core.AssertLen(t, opts.PeerAddrs, 2)
}

// OptionsFromEnv yields empty fields when nothing is set.
func TestServiceBehaviour_OptionsFromEnv_Ugly(t *core.T) {
	t.Setenv("CORE_P2P_ADDR", "")
	t.Setenv("CORE_P2P_PEERS", "")
	t.Setenv("CORE_P2P_NODE_ID", "")

	opts := OptionsFromEnv()
	core.AssertEqual(t, "", opts.ListenAddr)
	core.AssertEmpty(t, opts.PeerAddrs)
}

// mapValue returns a map directly, normalises a struct via JSON, and returns
// nil when the key is absent.
func TestServiceBehaviour_mapValue(t *core.T) {
	direct := mapValue(core.NewOptions(
		core.Option{Key: "payload", Value: map[string]any{"a": 1}},
	), "payload")
	core.AssertEqual(t, 1, direct["a"])

	core.AssertNil(t, mapValue(core.NewOptions(), "payload"))

	type wrapped struct {
		Name string `json:"name"`
	}
	normalised := mapValue(core.NewOptions(
		core.Option{Key: "payload", Value: wrapped{Name: "x"}},
	), "payload")
	core.AssertEqual(t, "x", normalised["name"])
}

// coalesce returns the first non-blank value, or "" when all are blank.
func TestServiceBehaviour_coalesce(t *core.T) {
	core.AssertEqual(t, "b", coalesce("", "  ", "b"))
	core.AssertEqual(t, "", coalesce("", "   "))
}
