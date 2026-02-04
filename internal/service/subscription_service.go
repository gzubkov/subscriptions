package service

import (
	"context"
	"errors"
	"strconv"
	"subscriptions/internal/models"
	"subscriptions/internal/storage"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type SubscriptionService interface {
	CreateSubscription(ctx context.Context, subscription *models.Subscription) error
	GetSubscriptionByID(ctx context.Context, id uint) (*models.Subscription, error)
	UpdateSubscription(ctx context.Context, id uint, req *models.UpdateSubscriptionRequest) (*models.Subscription, error)
	DeleteSubscription(ctx context.Context, id uint) error
	ListSubscriptions(ctx context.Context, limit, offset string) ([]models.Subscription, error)
	CalculateSubscriptionsTotalCost(ctx context.Context, startDate, endDate time.Time, userID *uuid.UUID, serviceName *string) (uint, error)
}

type subscriptionService struct {
	repo   storage.SubscriptionStorage
	logger *zap.Logger
}

func NewSubscriptionService(repo storage.SubscriptionStorage, logger *zap.Logger) SubscriptionService {
	return &subscriptionService{
		repo:   repo,
		logger: logger,
	}
}

func (s *subscriptionService) CreateSubscription(ctx context.Context, subscription *models.Subscription) error {
	// проверка дат (дата окончания не нулевая и идет после даты начала)
	if subscription.EndDate != nil && subscription.EndDate.Before(subscription.StartDate) {
		s.logger.Error("subscription end date must be after start date")
		return errors.New("subscription end date must be after start date")
	}

	return s.repo.CreateSubscription(ctx, subscription)
}

func (s *subscriptionService) GetSubscriptionByID(ctx context.Context, id uint) (*models.Subscription, error) {
	return s.repo.GetSubscriptionByID(ctx, id)
}

func (s *subscriptionService) UpdateSubscription(ctx context.Context, id uint, req *models.UpdateSubscriptionRequest) (*models.Subscription, error) {
	// Запрашиваем подписку по id
	subscription, err := s.repo.GetSubscriptionByID(ctx, id)
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
	if req.EndDate != nil {
		// Проверяем, что дата окончания позже даты начала
		if req.EndDate.Before(subscription.StartDate) {
			s.logger.Error("subscription end date must be after start date")
			return nil, errors.New("subscription end date must be after start date")
		}
		subscription.EndDate = req.EndDate
	}

	// Сохраняем изменения
	err = s.repo.UpdateSubscription(ctx, subscription)
	if err != nil {
		return nil, err
	}

	return subscription, nil
}

func (s *subscriptionService) DeleteSubscription(ctx context.Context, id uint) error {
	// Запрашиваем подписку по id
	_, err := s.repo.GetSubscriptionByID(ctx, id)
	if err != nil {
		return err
	}
	return s.repo.DeleteSubscription(ctx, id)
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

	return s.repo.ListSubscriptions(ctx, uint(limitInt), uint(offsetInt))
}

func (s *subscriptionService) CalculateSubscriptionsTotalCost(ctx context.Context, startDate, endDate time.Time, userID *uuid.UUID, serviceName *string) (uint, error) {
	// Проверяем, что даты корректны
	if endDate.Before(startDate) {
		s.logger.Error("subscription end date must be after start date")
		return 0, errors.New("subscription end date must be after start date")
	}

	return s.repo.CalculateSubscriptionsTotalCost(ctx, startDate, endDate, userID, serviceName)
}
