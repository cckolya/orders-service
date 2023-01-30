package main

import (
	"context"
	"encoding/json"
	"log"
	mrand "math/rand"
	"orders-service/internal/delivery"
	"orders-service/internal/model"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"
)

func main() {
	mq, err := delivery.NewNats()
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.TODO())

	sing := make(chan os.Signal, 1)
	signal.Notify(sing, os.Interrupt)
	go func() {
		<-sing
		cancel()
	}()

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		validMessageProducer(ctx, mq, time.Second*2)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		invalidMessageProducer(ctx, mq, time.Second*5)
		wg.Done()
	}()
	wg.Wait()

}

func validMessageProducer(ctx context.Context, mq *delivery.Nats, delay time.Duration) {
	ticker := time.NewTicker(delay)
	for {
		select {
		case <-ticker.C:
			data, err := json.Marshal(generateOrder())
			if err != nil {
				log.Println(err)
				continue
			}
			err = mq.PutMessage(data)
			if err != nil {
				log.Println(err)
			}
		case <-ctx.Done():
			log.Println("validMessageProducer is stopped")
			return
		}
	}
}

func invalidMessageProducer(ctx context.Context, mq *delivery.Nats, delay time.Duration) {
	ticker := time.NewTicker(delay)
	for {
		select {
		case <-ticker.C:
			data := generateRandomData()
			err := mq.PutMessage(data)
			if err != nil {
				log.Println(err)
			}
		case <-ctx.Done():
			log.Println("invalidMessageProducer is stopped")
			return
		}
	}
}

func generateRandomData() []byte {
	return []byte(randomStr(128))
}

func generateOrder() *model.Order {
	return &model.Order{
		OrderUID:    randomStr(10),
		TrackNumber: randomStr(5),
		Entry:       randomStr(6),
		Delivery: model.Delivery{
			Name:    randomStr(10),
			Phone:   randomStr(7),
			Zip:     randomStr(4),
			City:    randomStr(7),
			Address: randomStr(11),
			Region:  randomStr(8),
			Email:   randomStr(10),
		},
		Payment: model.Payment{
			Transaction:  randomStr(10),
			RequestID:    randomStr(10),
			Currency:     randomStr(10),
			Provider:     randomStr(10),
			Amount:       mrand.Intn(1000),
			PaymentDt:    mrand.Intn(1000),
			Bank:         randomStr(10),
			DeliveryCost: mrand.Intn(1000),
			GoodsTotal:   mrand.Intn(1000),
			CustomFee:    mrand.Intn(1000),
		},
		Items: []model.Item{
			{
				ChrtID:      mrand.Intn(1000),
				TrackNumber: randomStr(10),
				Price:       mrand.Intn(1000),
				Rid:         randomStr(10),
				Name:        randomStr(10),
				Sale:        mrand.Intn(1000),
				Size:        randomStr(10),
				TotalPrice:  mrand.Intn(1000),
				NmID:        mrand.Intn(1000),
				Brand:       randomStr(10),
				Status:      mrand.Intn(1000),
			},
		},
		Locale:            randomStr(10),
		InternalSignature: randomStr(10),
		CustomerID:        randomStr(10),
		DeliveryService:   randomStr(10),
		Shardkey:          randomStr(10),
		SmID:              mrand.Intn(1000),
		DateCreated:       time.Now(),
		OofShard:          randomStr(10),
	}
}

func randomStr(n int) string {
	str := "abcdefghijklmnopqrstuvwxyz"
	bld := strings.Builder{}
	for i := 0; i < n; i++ {
		mrand.Seed(time.Now().UnixNano())
		bld.WriteByte(str[mrand.Intn(len(str))])
	}
	return bld.String()
}
