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
	done := make(chan bool)

	d.nc, _ = nats.Connect(nats.DefaultURL)
	d.c, _ = nats.NewEncodedConn(d.nc, nats.JSON_ENCODER)
	defer d.nc.Close()

	currentDate := d.StartTime

	for {
		start := 0
		limit := d.BatchSize

		for {
			customers := d.getCustomersList(start, limit)

			count := len(customers)
			if count == 0 {
				break
			}

			for i := 0; i < count; i++ {
				customer := customers[i]

				customer.Date = currentDate

				d.c.Publish("collector", customer)
				log.Println("Dispatching: Customer " + strconv.FormatUint(customers[i].ID, 10))
			}

			start = start + limit
		}

		currentDate = currentDate.AddDate(0, 0, 1)
		if currentDate.After(d.StopTime) {
			break
		}
	}

	<-done
}

func (d *Dispatcher) getCustomersList(start, limit int) []faker.Customer {
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
