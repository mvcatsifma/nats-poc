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
var productTypes = []string{"pt1", "pt2", "pt3"}
var status = []string{"RUNNING", "STOPPED", "SERVICING", "RESETTING"}

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

			go func() {
			ProductLoop:
				for {
					time.Sleep(time.Second * time.Duration(rand.Int63n(2)))
					product := &Product{
						UUID:        uuid.Must(uuid.NewV4()),
						MachineId:   machine,
						ProductType: productTypes[rand.Intn(len(productTypes))],
						Ok:          rand.Float32() < 0.9,
						CreatedAt:   time.Now(),
					}
					if err := ec.Publish("products", product); err != nil {
						fmt.Println(err)
						break ProductLoop
					} else {
						fmt.Printf("\nSend: %+v", product)
					}
					nc.Flush()
				}
			}()

		StatusLoop:
			for {
				time.Sleep(time.Second * time.Duration(rand.Int63n(10)))
				status := &Status{
					Status:    status[rand.Intn(len(status))],
					MachineId: machine,
					UpdatedAt: time.Now(),
				}
				if err := ec.Publish("status", status); err != nil {
					fmt.Println(err)
					break StatusLoop
				} else {
					fmt.Printf("\nSend: %+v", status)
				}
				nc.Flush()
			}
		}(machine)
	}

	wg.Wait()
}

type Product struct {
	UUID        uuid.UUID
	MachineId   string
	ProductType string
	Ok          bool
	CreatedAt   time.Time
}

type Status struct {
	MachineId string
	Status    string
	UpdatedAt time.Time
}
