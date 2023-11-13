package nats

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
)

var (
	timeout = 5 * time.Second
	buffer  = 128 * 1024 * 1024
)

type Conn struct {
	rwc *nats.Conn
}

func Connect(s *Server) (*Conn, error) {
	opt := nats.Options{
		InProcessServer:  s.s,
		AllowReconnect:   true,
		MaxReconnect:     -1,
		ReconnectWait:    timeout,
		Timeout:          timeout,
		ReconnectBufSize: buffer,
	}

	rwc, err := opt.Connect()
	if err != nil {
		return nil, fmt.Errorf("parse conn options: %w", err)
	}

	return &Conn{rwc: rwc}, nil
}

func (c *Conn) JetStream(opts ...nats.JSOpt) (nats.JetStreamContext, error) {
	return c.rwc.JetStream(opts...)
}

func (c *Conn) Close() error {
	if c.rwc != nil {
		c.rwc.Close()
	}
	return nil
}

func Encode(value any) (b []byte, err error) {
	var buf bytes.Buffer
	err = gob.NewEncoder(&buf).Encode(value)
	if err != nil {
		return nil, nil
	}
	return buf.Bytes(), nil
}

func Decode(p []byte, data any) (err error) {
	var buf = bytes.NewReader(p)
	err = gob.NewDecoder(buf).Decode(data)
	return
}
