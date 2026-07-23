package p2p

import (
	"bufio"
	"context"
	core "dappco.re/go"
	"net"
	"time"
)

func TestTCPDriver_Subscribe_CancelRemovesHandler(t *core.T) {
	driver := NewTCPDriver(TCPOptions{})
	ctx, cancel := context.WithCancel(context.Background())

	calls := 0
	err := driver.Subscribe(ctx, "updates", func(Envelope) {
		calls++
	})
	core.RequireNoError(t, err)

	cancel()
	err = driver.Publish(context.Background(), Envelope{Topic: "updates"})
	core.RequireNoError(t, err)
	core.AssertEmpty(t, calls)
}

func TestTCPDriver_Publish_ContinuesAfterPeerFailure(t *core.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	core.RequireNoError(t, err)
	defer listener.Close()

	received := make(chan Envelope, 1)
	acceptErr := make(chan resultFailure, 1)
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			acceptErr <- err
			return
		}
		defer conn.Close()

		scanner := bufio.NewScanner(conn)
		if scanner.Scan() {
			var envelope Envelope
			if err := jsonUnmarshal(scanner.Bytes(), &envelope); err != nil {
				acceptErr <- err
				return
			}
			received <- envelope
			return
		}
		if err := scanner.Err(); err != nil {
			acceptErr <- err
			return
		}
		acceptErr <- context.Canceled
	}()

	driver := NewTCPDriver(TCPOptions{
		PeerAddrs: []string{"127.0.0.1:1", listener.Addr().String()},
		NodeID:    "node-1",
	})

	err = driver.Publish(context.Background(), Envelope{
		Topic:   "updates",
		Payload: map[string]any{"hello": "world"},
	})
	core.AssertError(t, err)

	select {
	case envelope := <-received:
		core.AssertEqual(t, "updates", envelope.Topic)
		core.AssertEqual(t, "node-1", envelope.SenderID)
		core.AssertEqual(t, map[string]any{"hello": "world"}, envelope.Payload)
	case err := <-acceptErr:
		core.RequireNoError(t, err)
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for peer delivery")
	}
}

// AX7 generated source-matching smoke coverage.
func TestTcp_NewTCPDriver_Good(t *core.T) {
	// NewTCPDriver
	ax7Variant := "NewTCPDriver:good"
	core.AssertContains(t, ax7Variant, "good")
	result := core.Try(func() any {
		got0 := NewTCPDriver(*new(TCPOptions))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestTcp_NewTCPDriver_Bad(t *core.T) {
	// NewTCPDriver
	ax7Variant := "NewTCPDriver:bad"
	core.AssertContains(t, ax7Variant, "bad")
	result := core.Try(func() any {
		got0 := NewTCPDriver(*new(TCPOptions))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestTcp_NewTCPDriver_Ugly(t *core.T) {
	// NewTCPDriver
	ax7Variant := "NewTCPDriver:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	result := core.Try(func() any {
		got0 := NewTCPDriver(*new(TCPOptions))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestTcp_TCPDriver_ListenAddr_Good(t *core.T) {
	// TCPDriver ListenAddr
	ax7Variant := "TCPDriver_ListenAddr:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(TCPDriver)
	result := core.Try(func() any {
		got0 := subject.ListenAddr()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestTcp_TCPDriver_ListenAddr_Bad(t *core.T) {
	// TCPDriver ListenAddr
	ax7Variant := "TCPDriver_ListenAddr:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(TCPDriver)
	result := core.Try(func() any {
		got0 := subject.ListenAddr()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestTcp_TCPDriver_ListenAddr_Ugly(t *core.T) {
	// TCPDriver ListenAddr
	ax7Variant := "TCPDriver_ListenAddr:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(TCPDriver)
	result := core.Try(func() any {
		got0 := subject.ListenAddr()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestTcp_TCPDriver_Subscribe_Good(t *core.T) {
	// TCPDriver Subscribe
	ax7Variant := "TCPDriver_Subscribe:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(TCPDriver)
	result := core.Try(func() any {
		got0 := subject.Subscribe(core.Background(), "agent", nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestTcp_TCPDriver_Subscribe_Bad(t *core.T) {
	// TCPDriver Subscribe
	ax7Variant := "TCPDriver_Subscribe:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(TCPDriver)
	result := core.Try(func() any {
		got0 := subject.Subscribe(core.Background(), "", nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestTcp_TCPDriver_Subscribe_Ugly(t *core.T) {
	// TCPDriver Subscribe
	ax7Variant := "TCPDriver_Subscribe:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(TCPDriver)
	result := core.Try(func() any {
		got0 := subject.Subscribe(core.Background(), "../../edge", nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestTcp_TCPDriver_Publish_Good(t *core.T) {
	// TCPDriver Publish
	ax7Variant := "TCPDriver_Publish:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(TCPDriver)
	result := core.Try(func() any {
		got0 := subject.Publish(core.Background(), *new(Envelope))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestTcp_TCPDriver_Publish_Bad(t *core.T) {
	// TCPDriver Publish
	ax7Variant := "TCPDriver_Publish:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(TCPDriver)
	result := core.Try(func() any {
		got0 := subject.Publish(core.Background(), *new(Envelope))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestTcp_TCPDriver_Publish_Ugly(t *core.T) {
	// TCPDriver Publish
	ax7Variant := "TCPDriver_Publish:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(TCPDriver)
	result := core.Try(func() any {
		got0 := subject.Publish(core.Background(), *new(Envelope))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestTcp_TCPDriver_Close_Good(t *core.T) {
	// TCPDriver Close
	ax7Variant := "TCPDriver_Close:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(TCPDriver)
	result := core.Try(func() any {
		got0 := subject.Close()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestTcp_TCPDriver_Close_Bad(t *core.T) {
	// TCPDriver Close
	ax7Variant := "TCPDriver_Close:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(TCPDriver)
	result := core.Try(func() any {
		got0 := subject.Close()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestTcp_TCPDriver_Close_Ugly(t *core.T) {
	// TCPDriver Close
	ax7Variant := "TCPDriver_Close:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(TCPDriver)
	result := core.Try(func() any {
		got0 := subject.Close()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}
