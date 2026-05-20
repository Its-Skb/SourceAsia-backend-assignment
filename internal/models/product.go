package models

import "time"

type Product struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	SKU       string    `json:"sku"`
	CreatedAt time.Time `json:"created_at"`
}

type ProductMedia struct {
	ProductID string   `json:"product_id"`
	ImageURLs []string `json:"image_urls"`
	VideoURLs []string `json:"video_urls"`
}