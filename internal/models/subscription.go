package models

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID          uint       `json:"id" db:"id"`
	UserID      uuid.UUID  `json:"user_id" db:"user_id"`
	ServiceName string     `json:"service_name" db:"service_name"`
	Price       uint       `json:"price" db:"price"`
	StartDate   time.Time  `json:"start_date" db:"start_date"`
	EndDate     *time.Time `json:"end_date,omitempty" db:"end_date"`
}

type CreateSubscriptionRequest struct {
	UserID      uuid.UUID  `json:"user_id" binding:"required"`
	ServiceName string     `json:"service_name" binding:"required,min=1,max=255"`
	Price       uint       `json:"price" binding:"required,numeric,gt=0"`
	StartDate   time.Time  `json:"start_date" binding:"required"`
	EndDate     *time.Time `json:"end_date,omitempty"`
}

type UpdateSubscriptionRequest struct {
	ServiceName *string    `json:"service_name,omitempty"`
	Price       *uint      `json:"price,omitempty"`
	StartDate   *time.Time `json:"start_date,omitempty"`
	EndDate     *time.Time `json:"end_date,omitempty"`
}

type TotalCostRequest struct {
	UserID      *uuid.UUID `json:"user_id,omitempty"`
	StartDate   time.Time  `json:"start_date" binding:"required"`
	EndDate     time.Time  `json:"end_date" binding:"required"`
	ServiceName *string    `json:"service_name,omitempty"`
}

type TotalCostResponse struct {
	TotalCost uint `json:"total_cost"`
	Period    struct {
		StartDate time.Time `json:"start_date"`
		EndDate   time.Time `json:"end_date"`
	} `json:"period"`
	Filters struct {
		UserID      *uuid.UUID `json:"user_id,omitempty"`
		ServiceName *string    `json:"service_name,omitempty"`
	} `json:"filters"`
}
