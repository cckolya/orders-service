package service

import (
	"context"
	"encoding/json"
	"github.com/nats-io/nats.go"
	"github.com/patrickmn/go-cache"
	"github.com/rs/zerolog"
	"orders-service/internal/model"
	"time"
)

type MQBroker interface {
	GetConsumer() (msg chan *nats.Msg, err error)
}

type Repository interface {
	CreateValidMsg(orderUID string, data []byte) error
	CreateInvalidMsg(data []byte) error
	GetOrderByID(uid string) ([]byte, error)
	GetValidOrders(limit int) (map[string]interface{}, error)
}

type Order struct {
	cache *cache.Cache
	log   *zerolog.Logger
	repo  Repository
	mq    MQBroker
}

func NewOrder(cache *cache.Cache, log *zerolog.Logger, repo Repository, mq MQBroker) *Order {
	return &Order{cache: cache, log: log, repo: repo, mq: mq}
}

func (o *Order) ValidateCache() {
	orders, err := o.repo.GetValidOrders(100)
	if err != nil {
		o.log.Err(err).Send()
		return
	}
	for key, bytes := range orders {
		o.cache.Set(key, bytes, time.Minute)
		o.log.Debug().Msgf("Message with key=[%s] saved into cache", key)
	}
}

func (o *Order) ListenAndServeMQ(ctx context.Context) error {
	msgs, err := o.mq.GetConsumer()
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg := <-msgs:
			var order *model.Order
			err = json.Unmarshal(msg.Data, &order)
			if err != nil {
				o.log.Err(err).Send()
				err := o.repo.CreateInvalidMsg(msg.Data) // TODO: продумать куда сохранять если база не доступна
				if err != nil {
					o.log.Err(err).Send()
				} else {
					o.log.Debug().Msg("invalid message was successful save")
				}
				continue
			}

			err = o.repo.CreateValidMsg(order.OrderUID, msg.Data)
			if err != nil {
				o.log.Err(err).Send()
			}
			o.log.Debug().Msg("valid message was successful save")
			o.cache.Set(order.OrderUID, msg.Data, time.Minute)
		}
	}
}

func (o *Order) GetOrderById(uid string) ([]byte, error) {
	val, exist := o.cache.Get(uid)
	if exist {
		return val.([]byte), nil
	}
	order, err := o.repo.GetOrderByID(uid)
	if err != nil {
		return nil, err
	}
	o.cache.Set(uid, order, time.Minute)
	return order, nil
}
