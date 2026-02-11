package storage

import (
	"context"
	"subscriptions/internal/models"
	"time"

	"github.com/google/uuid"
)

type SubscriptionStorage interface {
	CreateSubscription(ctx context.Context, subscription *models.Subscription) (*uint, error)
	GetSubscriptionByID(ctx context.Context, id uint) (*models.Subscription, error)
	UpdateSubscription(ctx context.Context, subscription *models.Subscription) error
	DeleteSubscription(ctx context.Context, id uint) error
	ListSubscriptions(ctx context.Context, limit, offset uint) ([]models.Subscription, error)
	CalculateSubscriptionsTotalCost(ctx context.Context, startDate, endDate time.Time, userID *uuid.UUID, serviceName *string) (uint, error)
}
