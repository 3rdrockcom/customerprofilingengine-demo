package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/epointpayment/customerprofilingengine-demo/models/faker"
	nats "github.com/nats-io/go-nats"
)

type Dispatcher struct {
	nc *nats.Conn
	c  *nats.EncodedConn

	StartTime time.Time
	StopTime  time.Time
	BatchSize int
	ApiURL    string
}

func NewDispatcher() *Dispatcher {
	d := &Dispatcher{
		StartTime: time.Date(2015, time.January, 1, 0, 0, 0, 0, time.UTC),
		StopTime:  time.Now().UTC(),
		BatchSize: 1000,
		ApiURL:    apiURL,
	}

	return d
}

func (d *Dispatcher) Run() {
	d.nc, _ = nats.Connect(nats.DefaultURL)
	d.c, _ = nats.NewEncodedConn(d.nc, nats.JSON_ENCODER)
	defer d.nc.Close()

	subject := "collector"

	start := 0
	limit := 100
	for {
		customers := d.makeRequest(start, limit)

		count := len(customers)
		if count == 0 {
			break
		}

		for i := 0; i < count; i++ {
			customer := customers[i]
			d.c.Publish(subject, customer)
			log.Println("Dispatching: Customer " + strconv.FormatUint(customers[i].ID, 10))
		}

		start = start + limit
	}
}

func (d *Dispatcher) makeRequest(start, limit int) []faker.Customer {
	var customers []faker.Customer

	// Make a get request
	res, err := http.Get(d.ApiURL + "/customers/list?start=" + strconv.Itoa(start) + "&limit=" + strconv.Itoa(limit))
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &customers); err != nil {
		panic(err)
	}

	return customers
}
