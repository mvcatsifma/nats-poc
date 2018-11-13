package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/nats-io/go-nats"
	"github.com/satori/go.uuid"
	"log"
	"sync"
	"time"
	// implicitely required db drivers
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var database *gorm.DB

func init() {
	database, _ = gorm.Open(
		"postgres",
		"host=database port=5432 user=postgres dbname=postgres password=postgres sslmode=disable",
	)
	if database.Error != nil {
		panic("failed to connect database")
	}
	database.AutoMigrate(&Product{})
}

func main() {
	defer database.Close()

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
		fmt.Printf("\nRecieve: %v: %v", p.MachineId, p.UUID)
		database.Save(p)
		database.Exec("INSERT INTO product_stream values(?, ?)", p.UUID, p.Ok)
	}); err != nil {
		log.Fatal()
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
