package nats

import (
	"context"
	"fmt"
	"log"

	"github.com/cenkalti/backoff"
	"github.com/nats-io/nats-server/v2/server"
	"github.com/rs/xid"
)

var cluster = "nats-route://lhr:4248"

type Server struct {
	s *server.Server
}

func NewServer(ctx context.Context /* ,options */) (s *Server, err error) {
	s = &Server{}
	s.s, err = server.NewServer(&server.Options{
		Port:      4222,
		HTTPPort:  8222,
		JetStream: true,
		StoreDir:  "/data/nats",
		// make simpler
		ServerName: xid.New().String(),
		Cluster: server.ClusterOpts{
			Name: "nats-cluster",
			Port: 4248,
		},
		Routes:    server.RoutesFromStr(cluster),
		RoutesStr: cluster,
	})
	if err != nil {
		return nil, fmt.Errorf("parse server options: %w", err)
	}
	s.s.Start()
	return
}

func (s *Server) Close() error {
	if s.s != nil {
		s.s.Shutdown()
	}
	return nil
}

func (n *Server) WaitForServer() {
	b := backoff.NewExponentialBackOff()

	for {
		d := b.NextBackOff()
		ready := n.s.ReadyForConnections(d)
		if ready {
			break
		}

		log.Printf("NATS server not ready, waited %s, retrying...", d)
	}
}
