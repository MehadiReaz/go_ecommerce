package order

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"ecommerce_project/pkg/utils"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// List retrieves user's orders
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int64)

	limit := 20
	offset := 0

	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	if o := r.URL.Query().Get("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil {
			offset = parsed
		}
	}

	orders, err := h.service.List(userID, limit, offset)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Orders retrieved successfully", orders)
}

// Create creates a new order
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int64)

	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := utils.ValidateStruct(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	order, err := h.service.Create(userID, &req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusCreated, "Order created successfully", order)
}

// GetByID retrieves an order by ID
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int64)

	vars := mux.Vars(r)
	orderID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid order ID")
		return
	}

	order, err := h.service.GetByID(orderID, userID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Order retrieved successfully", order)
}

// Cancel cancels an order
func (h *Handler) Cancel(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int64)

	vars := mux.Vars(r)
	orderID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid order ID")
		return
	}

	if err := h.service.Cancel(orderID, userID); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Order cancelled successfully", nil)
}
