package models

import (
	"github.com/lib/pq"
	"time"
)

type Order struct {
	Listings            pq.StringArray `json:"listings" gorm:"type:text[]" validate:"required"`
	TransactionRef      string         `json:"transaction_ref"`
	Amount              float64        `json:"amount" validate:"required number"`
	Status              string         `json:"status" validate:"required"`
	CompletedAt         time.Time      `json:"completed_at"`
	PurchasedBy         string         `json:"purchased_by" validate:"required"`
	OrderId             string         `json:"order_id" validate:"required"`
	TotalEmissionOffset float32        `json:"total_mission_offset" validate:"required"`
}
