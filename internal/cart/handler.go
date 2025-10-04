package cart

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

// Get retrieves the user's cart
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int64)
	
	cart, err := h.service.Get(userID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	utils.SuccessResponse(w, http.StatusOK, "Cart retrieved successfully", cart)
}

// AddItem adds an item to the cart
func (h *Handler) AddItem(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int64)
	
	var req AddItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	
	if err := utils.ValidateStruct(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	
	if err := h.service.AddItem(userID, &req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	
	utils.SuccessResponse(w, http.StatusOK, "Item added to cart", nil)
}

// UpdateItem updates a cart item
func (h *Handler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int64)
	
	vars := mux.Vars(r)
	itemID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid item ID")
		return
	}
	
	var req UpdateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	
	if err := utils.ValidateStruct(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	
	if err := h.service.UpdateItem(userID, itemID, &req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	
	utils.SuccessResponse(w, http.StatusOK, "Cart item updated", nil)
}

// RemoveItem removes an item from the cart
func (h *Handler) RemoveItem(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int64)
	
	vars := mux.Vars(r)
	itemID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid item ID")
		return
	}
	
	if err := h.service.RemoveItem(userID, itemID); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	
	utils.SuccessResponse(w, http.StatusOK, "Item removed from cart", nil)
}

// Clear clears all items from the cart
func (h *Handler) Clear(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int64)
	
	if err := h.service.Clear(userID); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	
	utils.SuccessResponse(w, http.StatusOK, "Cart cleared successfully", nil)
}
