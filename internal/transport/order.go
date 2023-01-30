package transport

import (
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"net/http"
)

type OrderService interface {
	GetOrderById(uid string) ([]byte, error)
}

type Order struct {
	log     *zerolog.Logger
	service OrderService
}

func NewOrder(log *zerolog.Logger, service OrderService) *Order {
	return &Order{log: log, service: service}
}

func (o *Order) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	o.log.Info().Msg("Start GetOrderByID")
	defer o.log.Info().Msg("Finish GetOrderByID")

	vars := mux.Vars(r)
	idStr := vars["id"]

	order, err := o.service.GetOrderById(idStr)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Oops, order not found"))
		o.log.Err(err).Send()
		return
	}
	o.log.Info().Msgf("Order by id=[%s] was found", idStr)

	w.Write(order)
	w.WriteHeader(http.StatusOK)
}
