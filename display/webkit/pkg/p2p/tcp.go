package p2p

import (
	"bufio"
	"context"
	"net"
	"sync"

	core "dappco.re/go"
)

type TCPOptions struct {
	ListenAddr string
	PeerAddrs  []string
	NodeID     string
}

type TCPDriver struct {
	options       TCPOptions
	mu            sync.RWMutex
	listener      net.Listener
	subscriptions map[string][]*subscription
}

func NewTCPDriver(options TCPOptions) *TCPDriver {
	return &TCPDriver{
		options: TCPOptions{
			ListenAddr: core.Trim(options.ListenAddr),
			PeerAddrs:  append([]string(nil), options.PeerAddrs...),
			NodeID:     core.Trim(options.NodeID),
		},
		subscriptions: make(map[string][]*subscription),
	}
}

func (d *TCPDriver) ListenAddr() string {
	d.mu.RLock()
	defer d.mu.RUnlock()
	if d.listener != nil {
		return d.listener.Addr().String()
	}
	return d.options.ListenAddr
}

func (d *TCPDriver) Subscribe(ctx context.Context, topic string, handler func(Envelope)) resultFailure {
	if ctx == nil {
		ctx = context.Background()
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	topic = core.Trim(topic)
	if topic == "" {
		return core.NewError("topic is required")
	}
	if handler == nil {
		return core.NewError("handler is required")
	}
	if err := d.ensureListener(); err != nil {
		return err
	}

	sub := &subscription{
		handler: func(envelope Envelope) {
			select {
			case <-ctx.Done():
				return
			default:
				handler(envelope)
			}
		},
	}

	d.mu.Lock()
	d.subscriptions[topic] = append(d.subscriptions[topic], sub)
	d.mu.Unlock()

	context.AfterFunc(ctx, func() {
		d.removeSubscription(topic, sub)
	})

	return nil
}

func (d *TCPDriver) Publish(ctx context.Context, envelope Envelope) resultFailure {
	if core.Trim(envelope.Topic) == "" {
		return core.NewError("topic is required")
	}
	if core.Trim(envelope.SenderID) == "" {
		envelope.SenderID = d.options.NodeID
	}
	d.dispatch(envelope)
	payload, err := jsonMarshal(envelope)
	if err != nil {
		return err
	}
	var publishErr error
	for _, peer := range d.options.PeerAddrs {
		peer = core.Trim(peer)
		if peer == "" {
			continue
		}
		conn, err := (&net.Dialer{}).DialContext(ctx, "tcp", peer)
		if err != nil {
			publishErr = core.ErrorJoin(publishErr, err)
			continue
		}
		if _, err := conn.Write(append(payload, '\n')); err != nil {
			publishErr = core.ErrorJoin(publishErr, err)
			if closeErr := conn.Close(); closeErr != nil {
				publishErr = core.ErrorJoin(publishErr, closeErr)
			}
			continue
		}
		if closeErr := conn.Close(); closeErr != nil {
			publishErr = core.ErrorJoin(publishErr, closeErr)
		}
	}
	return publishErr
}

func (d *TCPDriver) Close() resultFailure {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.listener == nil {
		d.subscriptions = make(map[string][]*subscription)
		return nil
	}
	err := d.listener.Close()
	d.listener = nil
	d.subscriptions = make(map[string][]*subscription)
	return err
}

func (d *TCPDriver) ensureListener() resultFailure {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.listener != nil || core.Trim(d.options.ListenAddr) == "" {
		return nil
	}
	listener, err := net.Listen("tcp", d.options.ListenAddr)
	if err != nil {
		return err
	}
	d.listener = listener
	go d.acceptLoop(listener)
	return nil
}

func (d *TCPDriver) acceptLoop(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		go d.readConn(conn)
	}
}

func (d *TCPDriver) readConn(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		var envelope Envelope
		if err := jsonUnmarshal(scanner.Bytes(), &envelope); err != nil {
			continue
		}
		d.dispatch(envelope)
	}
}

func (d *TCPDriver) dispatch(envelope Envelope) {
	d.mu.RLock()
	handlers := append([]*subscription{}, d.subscriptions[envelope.Topic]...)
	handlers = append(handlers, d.subscriptions["*"]...)
	d.mu.RUnlock()
	for _, sub := range handlers {
		if sub == nil {
			continue
		}
		sub.handler(envelope)
	}
}

func (d *TCPDriver) removeSubscription(topic string, target *subscription) {
	d.mu.Lock()
	defer d.mu.Unlock()

	subs := d.subscriptions[topic]
	for i, sub := range subs {
		if sub == target {
			copy(subs[i:], subs[i+1:])
			subs[len(subs)-1] = nil
			subs = subs[:len(subs)-1]
			break
		}
	}
	if len(subs) == 0 {
		delete(d.subscriptions, topic)
		return
	}
	d.subscriptions[topic] = subs
}

type subscription struct {
	handler func(Envelope)
}
