package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/epointpayment/customerprofilingengine-demo/models/faker"

	nats "github.com/nats-io/go-nats"
)

type Collector struct {
	nc *nats.Conn
	c  *nats.EncodedConn

	ApiURL string
}

func NewCollector() *Collector {
	c := &Collector{
		ApiURL: apiURL,
	}

	return c
}

func (c *Collector) Run() {
	done := make(chan bool)

	c.nc, _ = nats.Connect(nats.DefaultURL)
	c.c, _ = nats.NewEncodedConn(c.nc, nats.JSON_ENCODER)
	defer c.nc.Close()

	c.c.Subscribe("collector", func(customer *faker.Customer) {

		customerInfo := c.getCustomerInfo(int(customer.ID))
		customerTransactions := c.getCustomerTransactions(int(customer.ID), customer.Date, customer.Date.AddDate(0, 0, 1))

		payload := &faker.Payload{
			Info:         customerInfo,
			Transactions: customerTransactions,
		}

		c.c.Publish("recorder", payload)
	})

	<-done
}

func (c *Collector) getCustomerInfo(customerID int) faker.Customer {
	var customer faker.Customer

	// Make a get request
	res, err := http.Get(c.ApiURL + "/customer/" + strconv.Itoa(customerID) + "/info")
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &customer); err != nil {
		panic(err)
	}

	return customer
}

func (c *Collector) getCustomerTransactions(customerID int, startDate, endDate time.Time) []faker.Transaction {
	var transaction []faker.Transaction

	// Make a get request
	res, err := http.Get(c.ApiURL + "/customer/" + strconv.Itoa(customerID) + "/transactions?startDate=" + startDate.Format("20060102") + "&endDate=" + endDate.Format("20060102"))
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &transaction); err != nil {
		panic(err)
	}

	return transaction
}
