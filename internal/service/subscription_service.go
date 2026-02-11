package service

import (
	"context"
	"errors"
	"strconv"
	"subscriptions/internal/models"
	"subscriptions/internal/storage"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type SubscriptionService interface {
	CreateSubscription(ctx context.Context, subscription *models.Subscription) (*uint, error)
	GetSubscriptionByID(ctx context.Context, id uint) (*models.Subscription, error)
	UpdateSubscription(ctx context.Context, id uint, req *models.UpdateSubscriptionRequest) (*models.Subscription, error)
	DeleteSubscription(ctx context.Context, id uint) error
	ListSubscriptions(ctx context.Context, limit, offset string) ([]models.Subscription, error)
	CalculateSubscriptionsTotalCost(ctx context.Context, startDate, endDate models.MonthYear, userID *uuid.UUID, serviceName *string) (uint, error)
}

type subscriptionService struct {
	storage storage.SubscriptionStorage
	logger  *zap.Logger
}

func NewSubscriptionService(storage storage.SubscriptionStorage, logger *zap.Logger) SubscriptionService {
	return &subscriptionService{
		storage: storage,
		logger:  logger,
	}
}

func (s *subscriptionService) CreateSubscription(ctx context.Context, subscription *models.Subscription) (*uint, error) {
	// проверка дат (дата окончания не нулевая и идет после даты начала)
	if subscription.EndDate != nil {
		if subscription.EndDate.Before(subscription.StartDate) {
			s.logger.Error("subscription end date must be after start date")
			return nil, errors.New("subscription end date must be after start date")
		}
	}

	return s.storage.CreateSubscription(ctx, subscription)
}

func (s *subscriptionService) GetSubscriptionByID(ctx context.Context, id uint) (*models.Subscription, error) {
	return s.storage.GetSubscriptionByID(ctx, id)
}

func (s *subscriptionService) UpdateSubscription(ctx context.Context, id uint, req *models.UpdateSubscriptionRequest) (*models.Subscription, error) {
	// Запрашиваем подписку по id
	subscription, err := s.storage.GetSubscriptionByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Проверяем наличие полей для обновления
	if req.ServiceName != nil {
		subscription.ServiceName = *req.ServiceName
	}
	if req.Price != nil {
		subscription.Price = *req.Price
	}
	if req.StartDate != nil {
		subscription.StartDate = req.StartDate.ToTime()
	}

	if req.EndDate != nil {
		endDate := req.EndDate.ToEndTime()
		// Проверяем, что дата окончания позже даты начала
		if endDate.Before(subscription.StartDate) {
			s.logger.Error("subscription end date must be after start date")
			return nil, errors.New("subscription end date must be after start date")
		}
	}

	// Сохраняем изменения
	err = s.storage.UpdateSubscription(ctx, subscription)
	if err != nil {
		return nil, err
	}

	return subscription, nil
}

func (s *subscriptionService) DeleteSubscription(ctx context.Context, id uint) error {
	// Запрашиваем подписку по id
	_, err := s.storage.GetSubscriptionByID(ctx, id)
	if err != nil {
		return err
	}
	return s.storage.DeleteSubscription(ctx, id)
}

func (s *subscriptionService) ListSubscriptions(ctx context.Context, limit, offset string) ([]models.Subscription, error) {
	// смещение и лимит проверяем на числовой тип (и больше 0)
	offsetInt, err := strconv.Atoi(offset)
	if err != nil || offsetInt <= 0 {
		offsetInt = 0
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt <= 0 {
		limitInt = 100
	}

	s.logger.Info("Listing subscriptions")
	return s.storage.ListSubscriptions(ctx, uint(limitInt), uint(offsetInt))
}

func (s *subscriptionService) CalculateSubscriptionsTotalCost(ctx context.Context, startDate, endDate models.MonthYear, userID *uuid.UUID, serviceName *string) (uint, error) {
	startTime := startDate.ToTime()
	endTime := endDate.ToEndTime()

	// Проверяем, что даты корректны
	if endTime.Before(startTime) {
		s.logger.Error("subscription end date must be after start date")
		return 0, errors.New("subscription end date must be after start date")
	}

	return s.storage.CalculateSubscriptionsTotalCost(ctx, startTime, endTime, userID, serviceName)
}
