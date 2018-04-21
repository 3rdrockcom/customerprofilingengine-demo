package main

import (
	"github.com/epointpayment/customerprofilingengine-demo/models/faker"

	"github.com/epointpayment/customerprofilingengine-demo/models"

	nats "github.com/nats-io/go-nats"
)

type Recorder struct {
	nc *nats.Conn
	c  *nats.EncodedConn

	ApiURL string
}

func NewRecorder() *Recorder {
	r := &Recorder{}

	return r
}

func (r *Recorder) Run() {
	done := make(chan bool)

	r.nc, _ = nats.Connect(nats.DefaultURL)
	r.c, _ = nats.NewEncodedConn(r.nc, nats.JSON_ENCODER)
	defer r.nc.Close()

	r.c.Subscribe("recorder", func(payload *faker.Payload) {
		customer := &models.Customers{}
		errNoRecords := db.First(&customer, uint(payload.Info.ID)).Error

		tx := db.Begin()
		if tx.Error != nil {
			panic(tx.Error)
		}

		if errNoRecords != nil {
			customer = &models.Customers{
				ID:        uint(payload.Info.ID),
				Email:     payload.Info.Email,
				Gender:    payload.Info.Gender,
				FirstName: payload.Info.FirstName,
				LastName:  payload.Info.LastName,
			}
			tx.Create(customer)
		} else {
			customer.Email = payload.Info.Email
			customer.Gender = payload.Info.Gender
			customer.FirstName = payload.Info.FirstName
			customer.LastName = payload.Info.LastName
			tx.Save(customer)
		}

		entries := payload.Transactions
		for j := range entries {
			entry := &models.Transactions{
				ID:          entries[j].ID,
				CustomersID: entries[j].CustomersID,
				DateTime:    entries[j].DateTime,
				Description: entries[j].Description,
				Credit:      entries[j].Credit,
				Debit:       entries[j].Debit,
				Balance:     entries[j].Balance,
			}

			if err := tx.Create(entry).Error; err != nil {
				tx.Rollback()
				panic(err)
			}
		}

		if err := tx.Commit().Error; err != nil {
			panic(err)
		}
	})

	<-done
}
