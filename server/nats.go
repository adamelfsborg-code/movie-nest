package server

import (
	"log"

	"github.com/nats-io/nats.go"
)

type Nats struct {
	host string
}

func ConnectNats(n *Nats) (*nats.Conn, error) {
	nc, err := nats.Connect(n.host)
	if err != nil {
		log.Fatal(err)
	}

	return nc, nil
}
