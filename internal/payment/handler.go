package payment

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

// CreatePayment creates a new payment
func (h *Handler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int64)

	var req CreatePaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := utils.ValidateStruct(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	payment, err := h.service.CreatePayment(userID, &req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusCreated, "Payment initiated", payment)
}

// GetPayment retrieves a payment
func (h *Handler) GetPayment(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int64)

	vars := mux.Vars(r)
	paymentID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid payment ID")
		return
	}

	payment, err := h.service.GetPayment(paymentID, userID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Payment retrieved successfully", payment)
}

// StripeWebhook handles Stripe webhooks
func (h *Handler) StripeWebhook(w http.ResponseWriter, r *http.Request) {
	var payload PaymentWebhookPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid webhook payload")
		return
	}

	if err := h.service.ProcessWebhook("stripe", &payload); err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"success"}`))
}

// BkashWebhook handles bKash webhooks
func (h *Handler) BkashWebhook(w http.ResponseWriter, r *http.Request) {
	var payload PaymentWebhookPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid webhook payload")
		return
	}

	if err := h.service.ProcessWebhook("bkash", &payload); err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"success"}`))
}
