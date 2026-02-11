package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type MonthYear string

const MonthYearFormat = "01-2006"

func (my MonthYear) ToTime() time.Time {
	t, _ := ParseMonthYear(string(my))
	return t
}

// ParseMonthYear парсит строку в формате MM-YYYY
func ParseMonthYear(s string) (time.Time, error) {
	t, err := time.Parse(MonthYearFormat, s)

	if err != nil {
		fmt.Println("Error parsing time:", err)
		return time.Time{}, err
	}

	return t, nil
}

// FormatMonthYear форматирует time.Time в MM-YYYY
func FormatMonthYear(t time.Time) MonthYear {
	return MonthYear(t.Format(MonthYearFormat))
}

// IsValidMonthYear проверяет валидность строки MM-YYYY
func IsValidMonthYear(s string) bool {
	_, err := ParseMonthYear(s)
	return err == nil
}

// MarshalJSON для MonthYearString
func (my MonthYear) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(my))
}

// UnmarshalJSON для MonthYearString
func (my *MonthYear) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	// Проверяем формат
	if !IsValidMonthYear(s) {
		return fmt.Errorf("invalid month-year format: %s", s)
	}

	*my = MonthYear(s)
	return nil
}

func (my MonthYear) ToEndTime() time.Time {
	endTime, _ := ParseMonthYear(string(my))

	// Дата конца - последний день последнего месяца
	endDateTime := time.Date(
		endTime.Year(),
		endTime.Month()+1, // Следующий месяц
		0,                 // Последний день предыдущего месяца
		23, 59, 59, 999999999,
		time.UTC,
	)

	return endDateTime
}

var (
	ErrSubscriptionNotFound = errors.New("subscription not found")
)

type Subscription struct {
	ID          uint       `json:"id" db:"id"`
	UserID      uuid.UUID  `json:"user_id" db:"user_id"`
	ServiceName string     `json:"service_name" db:"service_name"`
	Price       uint       `json:"price" db:"price"`
	StartDate   time.Time  `json:"start_date" db:"start_date"`
	EndDate     *time.Time `json:"end_date,omitempty" db:"end_date"`
}

type SubscriptionResponse struct {
	ID          uint      `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	ServiceName string    `json:"service_name"`
	Price       uint      `json:"price"`
	StartDate   string    `json:"start_date"`
	EndDate     *string   `json:"end_date,omitempty"`
}

type CreateSubscriptionRequest struct {
	UserID      uuid.UUID  `json:"user_id" binding:"required"`
	ServiceName string     `json:"service_name" binding:"required,min=1,max=255"`
	Price       uint       `json:"price" binding:"required,numeric,gt=0"`
	StartDate   MonthYear  `json:"start_date" binding:"required"`
	EndDate     *MonthYear `json:"end_date,omitempty"`
}

type UpdateSubscriptionRequest struct {
	ServiceName *string    `json:"service_name,omitempty"`
	Price       *uint      `json:"price,omitempty"`
	StartDate   *MonthYear `json:"start_date,omitempty"`
	EndDate     *MonthYear `json:"end_date,omitempty"`
}

type TotalCostRequest struct {
	UserID      *uuid.UUID `json:"user_id,omitempty"`
	StartDate   MonthYear  `json:"start_date" binding:"required"`
	EndDate     MonthYear  `json:"end_date" binding:"required"`
	ServiceName *string    `json:"service_name,omitempty"`
}

type TotalCostResponse struct {
	TotalCost uint `json:"total_cost"`
	Period    struct {
		StartDate MonthYear `json:"start_date"`
		EndDate   MonthYear `json:"end_date"`
	} `json:"period"`
	Filters struct {
		UserID      *uuid.UUID `json:"user_id,omitempty"`
		ServiceName *string    `json:"service_name,omitempty"`
	} `json:"filters"`
}

func (subscription Subscription) GetResponse() SubscriptionResponse {
	subscriptionResponse := &SubscriptionResponse{
		ID:          subscription.ID,
		UserID:      subscription.UserID,
		ServiceName: subscription.ServiceName,
		Price:       subscription.Price,
		StartDate:   subscription.StartDate.Format(MonthYearFormat),
		EndDate:     nil,
	}

	if subscription.EndDate != nil {
		endDate := subscription.EndDate.Format(MonthYearFormat)
		subscriptionResponse.EndDate = &endDate
	}
	return *subscriptionResponse
}
