package manager

import (
	"errors"
	"net/http"
	"wb_test_1/internal/models/order"
	"wb_test_1/pkg/api_helper"

	"github.com/go-chi/chi"
)

type Servicer interface {
	GetById(string) (*order.Order, error)
	GetIdsList() []string
}

type Handler struct {
	service Servicer
}

func NewHandler(service Servicer) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(router *chi.Mux) {
	router.Get("/order/{id}", h.GetOrderById)
	router.Get("/ordersIdsList", h.GetOrdersIdsList)
}

func (h *Handler) GetOrderById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		api_helper.ErrorJSON(w, errors.New("id parameter is not defined"))
		return
	}

	order, err := h.service.GetById(id)
	if err != nil {
		api_helper.ErrorJSON(w, err)
		return
	}

	var payload api_helper.JsonResponse
	payload.Error = false
	payload.Message = "Succeed!"
	payload.Data = order

	api_helper.WriteJSON(w, http.StatusAccepted, payload)
}

func (h *Handler) GetOrdersIdsList(w http.ResponseWriter, r *http.Request) {
	orders := h.service.GetIdsList()

	var payload api_helper.JsonResponse
	payload.Error = false
	payload.Message = "Succeed!"
	payload.Data = orders

	api_helper.WriteJSON(w, http.StatusAccepted, payload)
}
