package storage

import (
	"sync"

	"backend-assignment/internal/models"
)

var (

	// Product metadata storage
	Products = make(map[string]*models.Product)

	// Media storage separated for performance optimization
	ProductMediaStore = make(map[string]*models.ProductMedia)

	// SKU uniqueness index
	SKUIndex = make(map[string]string)

	// Stable insertion order for pagination
	ProductOrder []string

	// Concurrency safety
	ProductMu sync.RWMutex
)