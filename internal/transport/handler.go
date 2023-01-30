package transport

import (
	"github.com/gorilla/mux"
	"net/http"
)

type OrderTransport interface {
	GetOrderByID(http.ResponseWriter, *http.Request)
}

type Handler struct {
	OrderTransport
}

func NewHandler(orderTransport OrderTransport) *Handler {
	return &Handler{OrderTransport: orderTransport}
}

func (h *Handler) Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/order/{id}", h.GetOrderByID).Methods(http.MethodGet)

	return router
}
