package main

import (
	"context"
	"errors"
	"github.com/patrickmn/go-cache"
	"github.com/rs/zerolog"
	"net"
	"net/http"
	"orders-service/config"
	"orders-service/internal/delivery"
	"orders-service/internal/repository"
	"orders-service/internal/service"
	"orders-service/internal/transport"
	"os"
	"os/signal"
	"sync"
	"time"
)

func main() {
	log := zerolog.New(os.Stdout).Level(zerolog.DebugLevel).With().Caller().Timestamp().Logger()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	log.Info().Msg("Config uploaded")

	pg, err := repository.NewPg(&log, cfg.Postgres)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	log.Info().Msg("Postgres uploaded")

	mq, err := delivery.NewNats()
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	log.Info().Msg("Nats-streaming uploaded")

	cash := cache.New(time.Minute, time.Minute)

	orderService := service.NewOrder(cash, &log, pg, mq)

	orderService.ValidateCache()

	ctx, cancel := context.WithCancel(context.TODO())
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		err = orderService.ListenAndServeMQ(ctx)
		if err != nil {
			log.Fatal().Err(err).Send()
		} else {
			log.Info().Msg("Consumer is successful stopped")
		}
	}()
	log.Info().Msg("Consumer listening started")

	orderTransport := transport.NewOrder(&log, orderService)

	handler := transport.NewHandler(orderTransport)

	wg.Add(1)
	listener, err := net.Listen("tcp", cfg.Handler.Url)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	go func() {
		defer wg.Done()
		err = http.Serve(listener, handler.Router())
		if err != nil && !errors.Is(err, net.ErrClosed) {
			log.Fatal().Err(err).Send()
		} else {
			log.Info().Msg("HTTP is successful stopped")
		}
	}()
	log.Info().Msg("HTTP listening started")

	sign := make(chan os.Signal, 1)
	signal.Notify(sign, os.Interrupt)

	<-sign
	cancel()
	err = listener.Close()
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	wg.Wait()
	log.Info().Msg("Program successful stopped")
}
