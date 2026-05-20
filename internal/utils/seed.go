package utils

import (
	"fmt"
	"time"

	"backend-assignment/internal/models"
	"backend-assignment/internal/storage"
)

func SeedProducts(count int) {

	storage.ProductMu.Lock()
	defer storage.ProductMu.Unlock()

	for i := 1; i <= count; i++ {

		productID := GenerateID()

		sku := fmt.Sprintf("SKU-%04d", i)

		product := &models.Product{
			ID:        productID,
			Name:      fmt.Sprintf("Product %d", i),
			SKU:       sku,
			CreatedAt: time.Now(),
		}

		storage.Products[productID] = product

		storage.ProductOrder = append(storage.ProductOrder, productID)

		storage.SKUIndex[sku] = productID

		// Generate sample media
		var imageURLs []string

		for j := 1; j <= 10; j++ {
			imageURLs = append(
				imageURLs,
				fmt.Sprintf(
					"https://cdn.example.com/products/%s/image-%d.jpg",
					sku,
					j,
				),
			)
		}

		storage.ProductMediaStore[productID] = &models.ProductMedia{
			ProductID: productID,
			ImageURLs: imageURLs,
			VideoURLs: []string{
				fmt.Sprintf(
					"https://cdn.example.com/products/%s/demo.mp4",
					sku,
				),
			},
		}
	}
}