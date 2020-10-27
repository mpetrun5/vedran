// Copyright (C) 2017 Michał Matczuk
// Use of this source code is governed by an AGPL-style
// license that can be found in the LICENSE file.

package server

import (
	"context"
	"fmt"
	"github.com/NodeFactoryIo/vedran/pkg/http-tunnel"
	"golang.org/x/net/http2"
	"net"
	"net/http"
	"sync"
	"time"
)

type connPair struct {
	conn       net.Conn
	clientConn *http2.ClientConn
}

type connPool struct {
	t     *http2.Transport
	conns map[string]connPair // key is host:port
	free  func(identifier string)
	mu    sync.RWMutex
}

func newConnPool(t *http2.Transport, f func(identifier string)) *connPool {
	return &connPool{
		t:     t,
		free:  f,
		conns: make(map[string]connPair),
	}
}

func (p *connPool) URL(identifier string) string {
	return fmt.Sprint("https://", identifier)
}

func (p *connPool) GetClientConn(req *http.Request, addr string) (*http2.ClientConn, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if cp, ok := p.conns[addr]; ok && cp.clientConn.CanTakeNewRequest() {
		return cp.clientConn, nil
	}
	return nil, errClientNotConnected
}

func (p *connPool) MarkDead(c *http2.ClientConn) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for addr, cp := range p.conns {
		if cp.clientConn == c {
			p.close(cp, addr)
			return
		}
	}
}

func (p *connPool) AddConn(conn net.Conn, identifier string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	addr := identifier

	if cp, ok := p.conns[addr]; ok {
		if err := p.ping(cp); err != nil {
			p.close(cp, addr)
		} else {
			return errClientAlreadyConnected
		}
	}

	c, err := p.t.NewClientConn(conn)
	if err != nil {
		return err
	}
	p.conns[addr] = connPair{
		conn:       conn,
		clientConn: c,
	}

	return nil
}

func (p *connPool) DeleteConn(identifier string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	addr := identifier

	if cp, ok := p.conns[addr]; ok {
		p.close(cp, addr)
	}
}

func (p *connPool) Ping(identifier string) (time.Duration, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	addr := identifier

	if cp, ok := p.conns[addr]; ok {
		start := time.Now()
		err := p.ping(cp)
		return time.Since(start), err
	}

	return 0, errClientNotConnected
}

func (p *connPool) ping(cp connPair) error {
	ctx, cancel := context.WithTimeout(context.Background(), tunnel.DefaultPingTimeout)
	defer cancel()

	return cp.clientConn.Ping(ctx)
}

func (p *connPool) close(cp connPair, addr string) {
	cp.conn.Close()
	delete(p.conns, addr)
	if p.free != nil {
		p.free(addr)
	}
}
