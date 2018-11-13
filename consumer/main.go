package main

import (
	"fmt"
	"github.com/nats-io/go-nats"
	"github.com/satori/go.uuid"
	"log"
	"sync"
	"time"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		nc.Close()
		log.Fatal(err)
	}
	defer ec.Close()

	wg := sync.WaitGroup{}
	wg.Add(1)

	if _, err := ec.Subscribe("products", func(p *Product) {
		fmt.Printf("%+v", p)
		wg.Done()
	}); err != nil {
		log.Fatal(err)
	}

	wg.Wait()

	nc.Close()
}

type Product struct {
	UUID      uuid.UUID
	MachineId string
	Ok        bool
	Produced  time.Time
}
