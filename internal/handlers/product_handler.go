package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"backend-assignment/internal/models"
	"backend-assignment/internal/storage"
	"backend-assignment/internal/utils"
	"backend-assignment/internal/validators"

	"github.com/gin-gonic/gin"
)

func CreateProduct(c *gin.Context) {

	var req models.CreateProductRequest

	// Validate JSON body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "invalid JSON body",
		})
		return
	}

	// Trim input values
	req.Name = strings.TrimSpace(req.Name)
	req.SKU = strings.TrimSpace(req.SKU)

	// Validate required fields
	if req.Name == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "name is required",
		})
		return
	}

	if req.SKU == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "sku is required",
		})
		return
	}

	// Validate image URLs
	if err := validators.ValidateURLs(req.ImageURLs); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Validate video URLs
	if err := validators.ValidateURLs(req.VideoURLs); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Concurrency-safe write operation
	storage.ProductMu.Lock()
	defer storage.ProductMu.Unlock()

	// Check duplicate SKU
	if _, exists := storage.SKUIndex[req.SKU]; exists {
		c.JSON(http.StatusConflict, models.ErrorResponse{
			Success: false,
			Error:   "SKU already exists",
		})
		return
	}

	// Generate product ID
	productID := utils.GenerateID()

	// Create product metadata
	product := &models.Product{
		ID:        productID,
		Name:      req.Name,
		SKU:       req.SKU,
		CreatedAt: time.Now(),
	}

	// Store product
	storage.Products[productID] = product

	// Store media separately
	storage.ProductMediaStore[productID] = &models.ProductMedia{
		ProductID: productID,
		ImageURLs: req.ImageURLs,
		VideoURLs: req.VideoURLs,
	}

	// Update SKU index
	storage.SKUIndex[req.SKU] = productID

	// Response
	c.JSON(http.StatusCreated, models.SuccessResponse{
		Success: true,
		Message: "product created successfully",
		Data: gin.H{
			"id":          product.ID,
			"name":        product.Name,
			"sku":         product.SKU,
			"image_urls":  req.ImageURLs,
			"video_urls":  req.VideoURLs,
			"created_at":  product.CreatedAt,
		},
	})
}

func GetProducts(c *gin.Context) {

	// Default pagination values
	limit := 20
	offset := 0

	// Parse limit
	if limitQuery := c.Query("limit"); limitQuery != "" {
		parsedLimit, err := strconv.Atoi(limitQuery)

		if err != nil || parsedLimit <= 0 {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Success: false,
				Error:   "invalid limit parameter",
			})
			return
		}

		// Max limit protection
		if parsedLimit > 100 {
			parsedLimit = 100
		}

		limit = parsedLimit
	}

	// Parse offset
	if offsetQuery := c.Query("offset"); offsetQuery != "" {
		parsedOffset, err := strconv.Atoi(offsetQuery)

		if err != nil || parsedOffset < 0 {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Success: false,
				Error:   "invalid offset parameter",
			})
			return
		}

		offset = parsedOffset
	}

	// Read lock for concurrency-safe reads
	storage.ProductMu.RLock()
	defer storage.ProductMu.RUnlock()

	// Convert map to slice
	productIDs := make([]string, 0, len(storage.Products))

	for productID := range storage.Products {
		productIDs = append(productIDs, productID)
	}

	// Pagination bounds
	start := offset

	if start > len(productIDs) {
		start = len(productIDs)
	}

	end := start + limit

	if end > len(productIDs) {
		end = len(productIDs)
	}

	paginatedIDs := productIDs[start:end]

	// Build lightweight response
	var products []models.ProductListItem

	for _, productID := range paginatedIDs {

		product := storage.Products[productID]
		media := storage.ProductMediaStore[productID]

		item := models.ProductListItem{
			ID:         product.ID,
			Name:       product.Name,
			SKU:        product.SKU,
			ImageCount: len(media.ImageURLs),
			VideoCount: len(media.VideoURLs),
			CreatedAt:  product.CreatedAt,
		}

		// Optional thumbnail URL
		if len(media.ImageURLs) > 0 {
			item.ThumbnailURL = media.ImageURLs[0]
		}

		products = append(products, item)
	}

	// Response
	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "products fetched successfully",
		Data: gin.H{
			"products": products,
			"pagination": gin.H{
				"limit":  limit,
				"offset": offset,
				"count":  len(products),
				"total":  len(storage.Products),
			},
		},
	})
}

func GetProductByID(c *gin.Context) {

	// Get product ID from URL
	productID := c.Param("id")

	// Concurrency-safe read
	storage.ProductMu.RLock()
	defer storage.ProductMu.RUnlock()

	// Check if product exists
	product, exists := storage.Products[productID]

	if !exists {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Success: false,
			Error:   "product not found",
		})
		return
	}

	// Get media data
	media := storage.ProductMediaStore[productID]

	// Build detail response
	response := models.ProductDetailResponse{
		ID:         product.ID,
		Name:       product.Name,
		SKU:        product.SKU,
		ImageURLs:  media.ImageURLs,
		VideoURLs:  media.VideoURLs,
		CreatedAt:  product.CreatedAt,
	}

	// Return response
	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "product fetched successfully",
		Data:     response,
	})
}