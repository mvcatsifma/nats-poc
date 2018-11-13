package main

import (
	"github.com/nats-io/go-nats"
	"github.com/satori/go.uuid"
	"log"
	"math/rand"
	"time"
)

var machineName = "machine-1"

func main() {
	connectionName := nats.Name(machineName)
	nc, err := nats.Connect(nats.DefaultURL, connectionName)
	if err != nil {
		log.Fatal(err)
	}
	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		nc.Close()
		log.Fatal(err)
	}
	defer ec.Close()

	product := &Product{
		UUID:      uuid.Must(uuid.NewV4()),
		Produced:  time.Now(),
		Ok:        rand.Float32() < 0.9,
		MachineId: machineName,
	}

	if err := ec.Publish("products", product); err != nil {
		log.Fatal(err)
	}
	// Make sure the message goes through before we close
	nc.Flush()
}

type Product struct {
	UUID      uuid.UUID
	MachineId string
	Ok        bool
	Produced  time.Time
}
