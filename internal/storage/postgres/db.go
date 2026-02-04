package postgres

import (
	"context"
	"fmt"
	"subscriptions/internal/models"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type SubscriptionStorage struct {
	db     *sqlx.DB
	logger *zap.Logger
}

func NewSubscriptionStorage(db *sqlx.DB, logger *zap.Logger) *SubscriptionStorage {
	return &SubscriptionStorage{
		db:     db,
		logger: logger,
	}
}

func (r *SubscriptionStorage) CreateSubscription(ctx context.Context, subscription *models.Subscription) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO subscriptions (user_id, service_name, price, start_date, end_date) VALUES ($1, $2, $3, $4, $5)`,
		subscription.UserID,
		subscription.ServiceName,
		subscription.Price,
		subscription.StartDate,
		subscription.EndDate,
	)

	if err != nil {
		r.logger.Error("Failed to create new subscription", zap.Error(err))
		return err
	}

	return nil
}

func (r *SubscriptionStorage) GetSubscriptionByID(ctx context.Context, id uint) (*models.Subscription, error) {
	var subscription models.Subscription

	query := `SELECT * FROM subscriptions WHERE id = $1`
	err := r.db.GetContext(ctx, &subscription, query, id)
	if err != nil {
		return nil, err
	}

	return &subscription, nil
}

func (r *SubscriptionStorage) UpdateSubscription(ctx context.Context, subscription *models.Subscription) error {
	query := `
  UPDATE subscriptions 
  SET service_name = $2, price = $3, start_date = $4, end_date = $5
  WHERE id = $1
 `

	_, err := r.db.ExecContext(
		ctx,
		query,
		subscription.ID,
		subscription.ServiceName,
		subscription.Price,
		subscription.StartDate,
		subscription.EndDate,
	)

	if err != nil {
		r.logger.Error("Failed to update subscription", zap.Error(err))
		return err
	}

	return nil
}

func (r *SubscriptionStorage) DeleteSubscription(ctx context.Context, id uint) error {
	query := `DELETE FROM subscriptions WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		r.logger.Error("Failed to delete subscription", zap.Error(err))
		return err
	}

	return nil
}

// чисто формально в описании задачи не указано, что надо выполнить пагинацию, но можно легко добавить, что я и сделал
func (r *SubscriptionStorage) ListSubscriptions(ctx context.Context, limit, offset uint) ([]models.Subscription, error) {
	var subscriptions []models.Subscription

	query := `SELECT * FROM subscriptions ORDER BY id DESC LIMIT $1 OFFSET $2`
	err := r.db.SelectContext(ctx, &subscriptions, query, limit, offset)
	if err != nil {
		return nil, err
	}

	return subscriptions, nil
}

func (r *SubscriptionStorage) CalculateSubscriptionsTotalCost(
	ctx context.Context,
	startDate, endDate time.Time,
	userID *uuid.UUID,
	serviceName *string,
) (uint, error) {
	query := `
  SELECT COALESCE(SUM(price), 0) as total_cost
  FROM subscriptions
  WHERE start_date <= $2
   AND (end_date >= $1 OR end_date IS NULL)
 `

	args := []interface{}{startDate, endDate}
	argIndex := 3

	if userID != nil {
		query += fmt.Sprintf(" AND user_id = $%d", argIndex)
		args = append(args, *userID)
		argIndex++
	}

	if serviceName != nil {
		query += fmt.Sprintf(" AND service_name = $%d", argIndex)
		args = append(args, *serviceName)
	}

	var totalCost uint
	err := r.db.GetContext(ctx, &totalCost, query, args...)
	if err != nil {
		r.logger.Error("Failed to calculate total cost", zap.Error(err))
		return 0, err
	}

	return totalCost, nil
}
