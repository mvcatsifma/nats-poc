package main

import (
	"fmt"
	"github.com/nats-io/go-nats"
	"github.com/satori/go.uuid"
	"log"
	"math/rand"
	"sync"
	"time"
)

var machines = []string{"machine-1", "machine-2", "machine-3"}

func main() {
	var wg = sync.WaitGroup{}

	for _, machine := range machines {
		wg.Add(1)
		go func(machine string) {
			defer wg.Done()
			connectionName := nats.Name(machine)
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

		MachineLoop:
			for {
				time.Sleep(time.Second * time.Duration(rand.Int63n(10)))
				product := &Product{
					UUID:      uuid.Must(uuid.NewV4()),
					Produced:  time.Now(),
					Ok:        rand.Float32() < 0.9,
					MachineId: machine,
				}

				if err := ec.Publish("products", product); err != nil {
					break MachineLoop
				} else {
					fmt.Printf("\nSend: %v: %v", machine, product.UUID)
				}

				// Make sure the message goes through before we close
				nc.Flush()
			}
		}(machine)
	}

	wg.Wait()
}

type Product struct {
	UUID      uuid.UUID
	MachineId string
	Ok        bool
	Produced  time.Time
}
