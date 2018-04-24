package main

import (
	"github.com/epointpayment/customerprofilingengine-demo/models"
	"github.com/epointpayment/customerprofilingengine-demo/models/faker"

	nats "github.com/nats-io/go-nats"
)

type Analyzer struct {
	nc *nats.Conn
	c  *nats.EncodedConn

	Duration int //days
}

func NewAnalyzer() *Analyzer {
	a := &Analyzer{
		Duration: 30,
	}

	return a
}

func (a *Analyzer) Run() {
	done := make(chan bool)

	a.nc, _ = nats.Connect(nats.DefaultURL)
	a.c, _ = nats.NewEncodedConn(a.nc, nats.JSON_ENCODER)
	defer a.nc.Close()

	a.c.Subscribe("analyzer", func(analyze faker.Payload) {

		id := analyze.Info.ID
		t1 := analyze.Date.AddDate(0, 0, -a.Duration)
		t2 := analyze.Date

		transactions := []models.Transactions{}

		query := db.Where("customers_id = ?", id).
			Where("date_time >= ?", t1).
			Where("date_time < ?", t2)
		err := query.Find(&transactions).Error
		if err != nil {
			panic(err)
		}

		analyzer := &models.Analyzer{
			CustomersID: uint(id),
			Date:        t2,
		}

		for i := range transactions {
			analyzer.SumCredits += transactions[i].Credit
			if transactions[i].Credit != 0 {
				analyzer.NumCredits++
			}

			analyzer.SumDebits += transactions[i].Debit
			if transactions[i].Debit != 0 {
				analyzer.NumDebits++
			}

			analyzer.NumTransactions++
		}

		if analyzer.NumCredits > 0 {
			analyzer.AvgCredits = analyzer.SumCredits / float64(analyzer.NumCredits)
		}

		if analyzer.NumDebits > 0 {
			analyzer.AvgDebits = analyzer.SumDebits / float64(analyzer.NumDebits)
		}

		db.Create(analyzer)

	})

	<-done
}
