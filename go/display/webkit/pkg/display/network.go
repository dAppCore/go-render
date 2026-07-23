package display

import (
	"html"
	"net"
	"sort"
	"time"

	core "dappco.re/go"
	"dappco.re/go/render/display/webkit/pkg/p2p"
)

type NetworkInterfaceState struct {
	Name         string   `json:"name"`
	Index        int      `json:"index"`
	MTU          int      `json:"mtu"`
	HardwareAddr string   `json:"hardware_addr,omitempty"`
	Flags        []string `json:"flags,omitempty"`
	Addresses    []string `json:"addresses,omitempty"`
	Up           bool     `json:"up"`
	Loopback     bool     `json:"loopback"`
}

type NetworkPeerState struct {
	ID        string    `json:"id"`
	Topic     string    `json:"topic"`
	Connected bool      `json:"connected"`
	SeenAt    time.Time `json:"seen_at"`
}

type NetworkState struct {
	Hostname   string                  `json:"hostname"`
	Interfaces []NetworkInterfaceState `json:"interfaces"`
	Peers      []NetworkPeerState      `json:"peers,omitempty"`
	ObservedAt time.Time               `json:"observed_at"`
}

func (s *Service) networkState() NetworkState {
	state := NetworkState{
		Hostname:   hostname(),
		Interfaces: make([]NetworkInterfaceState, 0),
		ObservedAt: time.Now().UTC(),
	}

	interfaces, err := net.Interfaces()
	if err != nil {
		return state
	}

	sort.Slice(interfaces, func(i, j int) bool {
		return core.Lower(interfaces[i].Name) < core.Lower(interfaces[j].Name)
	})

	for _, iface := range interfaces {
		addresses := make([]string, 0)
		if addrs, err := iface.Addrs(); err == nil {
			for _, addr := range addrs {
				addresses = append(addresses, addr.String())
			}
			sort.Strings(addresses)
		}

		state.Interfaces = append(state.Interfaces, NetworkInterfaceState{
			Name:         iface.Name,
			Index:        iface.Index,
			MTU:          iface.MTU,
			HardwareAddr: iface.HardwareAddr.String(),
			Flags:        interfaceFlags(iface.Flags),
			Addresses:    addresses,
			Up:           iface.Flags&net.FlagUp != 0,
			Loopback:     iface.Flags&net.FlagLoopback != 0,
		})
	}

	state.Peers = s.p2pPeers()
	return state
}

type peerLister interface {
	Peers() []p2p.Peer
}

func (s *Service) p2pPeers() []NetworkPeerState {
	if s == nil || s.Core() == nil {
		return nil
	}

	for _, serviceName := range []string{"p2p", "network"} {
		serviceResult := s.Core().Service(serviceName)
		if !serviceResult.OK || serviceResult.Value == nil {
			continue
		}
		lister, ok := serviceResult.Value.(peerLister)
		if !ok {
			continue
		}

		peers := lister.Peers()
		if len(peers) == 0 {
			continue
		}

		peerStates := make([]NetworkPeerState, 0, len(peers))
		for _, peer := range peers {
			peerStates = append(peerStates, NetworkPeerState{
				ID:        peer.ID,
				Topic:     peer.Topic,
				Connected: peer.Connected,
				SeenAt:    peer.SeenAt,
			})
		}
		sort.Slice(peerStates, func(i, j int) bool {
			if peerStates[i].SeenAt.Equal(peerStates[j].SeenAt) {
				return core.Lower(peerStates[i].ID) < core.Lower(peerStates[j].ID)
			}
			return peerStates[i].SeenAt.After(peerStates[j].SeenAt)
		})
		return peerStates
	}

	return nil
}

func (s *Service) renderNetworkPage(state NetworkState) string {
	builder := core.NewBuilder()
	builder.WriteString("<!doctype html><html><head><meta charset=\"utf-8\"><title>core://network</title><style>")
	builder.WriteString("body{font:14px/1.5 ui-sans-serif,system-ui,-apple-system,BlinkMacSystemFont,\"Segoe UI\",sans-serif;background:#07111f;color:#e2e8f0;margin:0}")
	builder.WriteString("header{padding:20px;border-bottom:1px solid #203047;background:linear-gradient(180deg,#0f172a,#07111f)}")
	builder.WriteString("main{padding:20px;display:grid;gap:16px}section{background:#0b1220;border:1px solid #203047;border-radius:16px;padding:16px}")
	builder.WriteString("ul{list-style:none;padding:0;margin:0;display:grid;gap:12px}.iface{padding:12px;border:1px solid #203047;border-radius:12px;background:#020617}")
	builder.WriteString(".meta{color:#94a3b8}.name{font-weight:700;color:#7dd3fc}code{background:#111827;border-radius:8px;padding:2px 6px}")
	builder.WriteString("</style></head><body><header><strong>core://network</strong><div class=\"meta\">")
	builder.WriteString(html.EscapeString(state.Hostname))
	builder.WriteString(" · ")
	builder.WriteString(html.EscapeString(state.ObservedAt.Format(time.RFC3339)))
	builder.WriteString("</div></header><main><section><ul>")

	if len(state.Interfaces) == 0 {
		builder.WriteString("<li class=\"meta\">No network interfaces were detected.</li>")
	} else {
		for _, iface := range state.Interfaces {
			builder.WriteString("<li class=\"iface\"><div class=\"name\">")
			builder.WriteString(html.EscapeString(iface.Name))
			builder.WriteString("</div><div class=\"meta\">Index ")
			builder.WriteString(core.Sprintf("%d", iface.Index))
			builder.WriteString(" · MTU ")
			builder.WriteString(core.Sprintf("%d", iface.MTU))
			builder.WriteString(" · ")
			if iface.Up {
				builder.WriteString("up")
			} else {
				builder.WriteString("down")
			}
			if iface.Loopback {
				builder.WriteString(" · loopback")
			}
			builder.WriteString("</div>")
			if len(iface.Addresses) > 0 {
				builder.WriteString("<pre>")
				builder.WriteString(html.EscapeString(core.Join("\n", iface.Addresses...)))
				builder.WriteString("</pre>")
			}
			builder.WriteString("</li>")
		}
	}

	if len(state.Peers) > 0 {
		builder.WriteString("</ul></section><section><div class=\"meta\">Registered peers</div><ul>")
		for _, peer := range state.Peers {
			builder.WriteString("<li class=\"iface\"><div class=\"name\">")
			builder.WriteString(html.EscapeString(peer.ID))
			builder.WriteString("</div><div class=\"meta\">")
			builder.WriteString(html.EscapeString(peer.Topic))
			builder.WriteString(" · ")
			if peer.Connected {
				builder.WriteString("connected")
			} else {
				builder.WriteString("disconnected")
			}
			if !peer.SeenAt.IsZero() {
				builder.WriteString(" · ")
				builder.WriteString(html.EscapeString(peer.SeenAt.Format(time.RFC3339)))
			}
			builder.WriteString("</div></li>")
		}
		builder.WriteString("</ul></section>")
	} else {
		builder.WriteString("</ul></section>")
	}
	builder.WriteString("</main></body></html>")
	return builder.String()
}

func (s *Service) renderNetworkInterfacePage(state NetworkState, iface NetworkInterfaceState) string {
	builder := core.NewBuilder()
	builder.WriteString("<!doctype html><html><head><meta charset=\"utf-8\"><title>core://network/")
	builder.WriteString(html.EscapeString(iface.Name))
	builder.WriteString("</title><style>")
	builder.WriteString("body{font:14px/1.5 ui-sans-serif,system-ui,-apple-system,BlinkMacSystemFont,\"Segoe UI\",sans-serif;background:#07111f;color:#e2e8f0;margin:0}")
	builder.WriteString("header{padding:20px;border-bottom:1px solid #203047;background:linear-gradient(180deg,#0f172a,#07111f)}")
	builder.WriteString("main{padding:20px;display:grid;gap:16px}section{background:#0b1220;border:1px solid #203047;border-radius:16px;padding:16px}")
	builder.WriteString("pre{margin:0;white-space:pre-wrap;word-break:break-word}code{background:#111827;border-radius:8px;padding:2px 6px}")
	builder.WriteString(".meta{color:#94a3b8}.name{font-weight:700;color:#7dd3fc}")
	builder.WriteString("</style></head><body><header><strong>core://network/")
	builder.WriteString(html.EscapeString(iface.Name))
	builder.WriteString("</strong><div class=\"meta\">")
	builder.WriteString(html.EscapeString(state.Hostname))
	builder.WriteString("</div></header><main><section><div class=\"name\">")
	builder.WriteString(html.EscapeString(iface.Name))
	builder.WriteString("</div><div class=\"meta\">Index ")
	builder.WriteString(core.Sprintf("%d", iface.Index))
	builder.WriteString(" · MTU ")
	builder.WriteString(core.Sprintf("%d", iface.MTU))
	builder.WriteString(" · ")
	if iface.Up {
		builder.WriteString("up")
	} else {
		builder.WriteString("down")
	}
	if iface.Loopback {
		builder.WriteString(" · loopback")
	}
	builder.WriteString("</div><pre>")
	builder.WriteString(html.EscapeString(core.Join("\n", iface.Addresses...)))
	builder.WriteString("</pre><div class=\"meta\">Flags: ")
	if len(iface.Flags) == 0 {
		builder.WriteString("none")
	} else {
		builder.WriteString(html.EscapeString(core.Join(", ", iface.Flags...)))
	}
	builder.WriteString("</div></section>")
	if len(state.Peers) > 0 {
		builder.WriteString("<section><div class=\"meta\">Registered peers</div><ul>")
		for _, peer := range state.Peers {
			builder.WriteString("<li class=\"iface\"><div class=\"name\">")
			builder.WriteString(html.EscapeString(peer.ID))
			builder.WriteString("</div><div class=\"meta\">")
			builder.WriteString(html.EscapeString(peer.Topic))
			builder.WriteString(" · ")
			if peer.Connected {
				builder.WriteString("connected")
			} else {
				builder.WriteString("disconnected")
			}
			if !peer.SeenAt.IsZero() {
				builder.WriteString(" · ")
				builder.WriteString(html.EscapeString(peer.SeenAt.Format(time.RFC3339)))
			}
			builder.WriteString("</div></li>")
		}
		builder.WriteString("</ul></section>")
	}
	builder.WriteString("</main></body></html>")
	return builder.String()
}

func interfaceFlags(flags net.Flags) []string {
	values := make([]string, 0, 4)
	if flags&net.FlagUp != 0 {
		values = append(values, "up")
	}
	if flags&net.FlagBroadcast != 0 {
		values = append(values, "broadcast")
	}
	if flags&net.FlagLoopback != 0 {
		values = append(values, "loopback")
	}
	if flags&net.FlagPointToPoint != 0 {
		values = append(values, "point-to-point")
	}
	if flags&net.FlagMulticast != 0 {
		values = append(values, "multicast")
	}
	if flags&net.FlagRunning != 0 {
		values = append(values, "running")
	}
	return values
}

func hostname() string {
	name, err := coreHostname()
	if err != nil || core.Trim(name) == "" {
		return "localhost"
	}
	return name
}
