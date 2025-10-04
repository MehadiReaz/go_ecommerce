package product

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

// List retrieves all products with filtering
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	filter := &ProductFilter{}

	// Parse query parameters
	if categoryID := r.URL.Query().Get("category_id"); categoryID != "" {
		if id, err := strconv.ParseInt(categoryID, 10, 64); err == nil {
			filter.CategoryID = id
		}
	}

	if minPrice := r.URL.Query().Get("min_price"); minPrice != "" {
		if price, err := strconv.ParseFloat(minPrice, 64); err == nil {
			filter.MinPrice = price
		}
	}

	if maxPrice := r.URL.Query().Get("max_price"); maxPrice != "" {
		if price, err := strconv.ParseFloat(maxPrice, 64); err == nil {
			filter.MaxPrice = price
		}
	}

	if featured := r.URL.Query().Get("is_featured"); featured != "" {
		if val, err := strconv.ParseBool(featured); err == nil {
			filter.IsFeatured = &val
		}
	}

	if limit := r.URL.Query().Get("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil {
			filter.Limit = l
		}
	}

	if offset := r.URL.Query().Get("offset"); offset != "" {
		if o, err := strconv.Atoi(offset); err == nil {
			filter.Offset = o
		}
	}

	products, err := h.service.List(filter)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Products retrieved successfully", products)
}

// GetByID retrieves a product by ID
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	product, err := h.service.GetByID(id)
	if err != nil {
		utils.ErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Product retrieved successfully", product)
}

// Create creates a new product (admin only)
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := utils.ValidateStruct(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	product, err := h.service.Create(&req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusCreated, "Product created successfully", product)
}

// Update updates a product (admin only)
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	var req UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	product, err := h.service.Update(id, &req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Product updated successfully", product)
}

// Delete deletes a product (admin only)
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	if err := h.service.Delete(id); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Product deleted successfully", nil)
}

// Search searches for products
func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	searchTerm := r.URL.Query().Get("q")
	if searchTerm == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "Search term is required")
		return
	}

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

	products, err := h.service.Search(searchTerm, limit, offset)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Search results", products)
}
