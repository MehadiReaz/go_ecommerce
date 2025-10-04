package shipping

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

// ListAddresses retrieves all addresses for the user
func (h *Handler) ListAddresses(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int64)

	addresses, err := h.service.ListAddresses(userID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Addresses retrieved successfully", addresses)
}

// CreateAddress creates a new address
func (h *Handler) CreateAddress(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int64)

	var req CreateAddressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := utils.ValidateStruct(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	address, err := h.service.CreateAddress(userID, &req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusCreated, "Address created successfully", address)
}

// UpdateAddress updates an address
func (h *Handler) UpdateAddress(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int64)

	vars := mux.Vars(r)
	addressID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid address ID")
		return
	}

	var req UpdateAddressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	address, err := h.service.UpdateAddress(userID, addressID, &req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Address updated successfully", address)
}

// DeleteAddress deletes an address
func (h *Handler) DeleteAddress(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int64)

	vars := mux.Vars(r)
	addressID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid address ID")
		return
	}

	if err := h.service.DeleteAddress(userID, addressID); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Address deleted successfully", nil)
}
