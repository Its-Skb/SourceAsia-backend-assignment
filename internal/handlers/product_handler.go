package handlers

import (
	"net/http"
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