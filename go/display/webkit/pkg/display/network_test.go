package display

import (
	core "dappco.re/go"
	"net"
	"time"
)

func TestNetwork_InterfaceFlags_GoodCase(t *core.T) {
	flags := interfaceFlags(net.FlagUp | net.FlagLoopback | net.FlagRunning)

	core.AssertEqual(t, []string{"up", "loopback", "running"}, flags)
	core.AssertNotEmpty(t, core.Sprintf("%T", flags))
}

func TestNetwork_InterfaceFlags_BadCase(t *core.T) {
	core.AssertEmpty(t, interfaceFlags(0))
	observedType := core.Sprintf("%T", interfaceFlags(0))
	core.AssertNotEmpty(t, observedType)
}

func TestNetwork_InterfaceFlags_UglyCase(t *core.T) {
	core.AssertEmpty(t, interfaceFlags(net.Flags(1<<30)))
	observedType := core.Sprintf("%T", interfaceFlags(net.Flags(1<<30)))
	core.AssertNotEmpty(t, observedType)
}

func TestNetwork_RenderNetworkPage_GoodCase(t *core.T) {
	svc := &Service{}
	state := NetworkState{
		Hostname:   "core-host",
		ObservedAt: time.Unix(1_700_000_000, 0).UTC(),
		Interfaces: []NetworkInterfaceState{
			{
				Name:      "en0",
				Index:     2,
				MTU:       1500,
				Addresses: []string{"192.168.0.10/24", "fe80::1/64"},
				Up:        true,
			},
		},
		Peers: []NetworkPeerState{
			{ID: "peer-1", Topic: "timeline", Connected: true, SeenAt: time.Unix(1_700_000_100, 0).UTC()},
		},
	}

	body := svc.renderNetworkPage(state)

	core.AssertContains(t, body, "core://network")
	core.AssertContains(t, body, "core-host")
	core.AssertContains(t, body, "en0")
	core.AssertContains(t, body, "192.168.0.10/24")
	core.AssertContains(t, body, "Registered peers")
	core.AssertContains(t, body, "peer-1")
}

func TestNetwork_RenderNetworkPage_BadCase(t *core.T) {
	svc := &Service{}

	body := svc.renderNetworkPage(NetworkState{
		Hostname:   "<host>",
		ObservedAt: time.Unix(1, 0).UTC(),
	})

	core.AssertContains(t, body, "No network interfaces were detected.")
	core.AssertContains(t, body, "&lt;host&gt;")
}

func TestNetwork_RenderNetworkPage_UglyCase(t *core.T) {
	svc := &Service{}

	body := svc.renderNetworkPage(NetworkState{
		Hostname:   repeatString("x", 128),
		ObservedAt: time.Unix(1, 0).UTC(),
		Interfaces: []NetworkInterfaceState{
			{Name: "\"quoted\"", Index: 99, MTU: 9, Addresses: []string{"<addr>"}, Up: false, Loopback: true},
		},
	})

	core.AssertContains(t, body, "&#34;quoted&#34;")
	core.AssertContains(t, body, "&lt;addr&gt;")
	core.AssertContains(t, body, "loopback")
}

func TestNetwork_RenderNetworkInterfacePage_GoodCase(t *core.T) {
	svc := &Service{}
	state := NetworkState{
		Hostname:   "core-host",
		ObservedAt: time.Unix(1_700_000_000, 0).UTC(),
		Peers: []NetworkPeerState{
			{ID: "peer-1", Topic: "timeline", Connected: true, SeenAt: time.Unix(1_700_000_100, 0).UTC()},
		},
	}
	iface := NetworkInterfaceState{
		Name:      "en0",
		Index:     2,
		MTU:       1500,
		Addresses: []string{"192.168.0.10/24"},
		Flags:     []string{"up", "running"},
		Up:        true,
	}

	body := svc.renderNetworkInterfacePage(state, iface)

	core.AssertContains(t, body, "core://network/en0")
	core.AssertContains(t, body, "core-host")
	core.AssertContains(t, body, "192.168.0.10/24")
	core.AssertContains(t, body, "Flags: up, running")
	core.AssertContains(t, body, "Registered peers")
	core.AssertContains(t, body, "peer-1")
}

func TestNetwork_RenderNetworkInterfacePage_BadCase(t *core.T) {
	svc := &Service{}
	state := NetworkState{Hostname: "core-host", ObservedAt: time.Unix(1, 0).UTC()}
	iface := NetworkInterfaceState{Name: "en0", Index: 2, MTU: 1500, Up: false}

	body := svc.renderNetworkInterfacePage(state, iface)

	core.AssertContains(t, body, "Flags: none")
	core.AssertNotContains(t, body, "Registered peers")
}

func TestNetwork_RenderNetworkInterfacePage_UglyCase(t *core.T) {
	svc := &Service{}
	state := NetworkState{Hostname: "<host>", ObservedAt: time.Unix(1, 0).UTC()}
	iface := NetworkInterfaceState{
		Name:      "\"quoted\"",
		Index:     99,
		MTU:       9,
		Addresses: []string{"<addr>"},
		Up:        false,
		Loopback:  true,
	}

	body := svc.renderNetworkInterfacePage(state, iface)

	core.AssertContains(t, body, "&#34;quoted&#34;")
	core.AssertContains(t, body, "&lt;addr&gt;")
	core.AssertContains(t, body, "&lt;host&gt;")
	core.AssertContains(t, body, "loopback")
}
