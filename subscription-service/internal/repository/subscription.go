package repository

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
	"subscription-service/internal/model"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, sub *model.Subscription) error
	GetByID(ctx context.Context, id string) (*model.Subscription, error)
	Update(ctx context.Context, id string, req *model.UpdateSubscriptionRequest) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, userID *string, serviceName *string) ([]*model.Subscription, error)
	GetSummary(ctx context.Context, req *model.SummaryRequest) (*model.SubscriptionSummary, error)
}

type subscriptionRepo struct {
	db *sqlx.DB
}

func NewSubscriptionRepository(db *sqlx.DB) SubscriptionRepository {
	return &subscriptionRepo{db: db}
}

func (r *subscriptionRepo) Create(ctx context.Context, sub *model.Subscription) error {
	query := `
		INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	log.Printf("Creating subscription for user %s, service: %s", sub.UserID, sub.ServiceName)

	return r.db.QueryRowContext(ctx, query,
		sub.ServiceName,
		sub.Price,
		sub.UserID,
		sub.StartDate,
		sub.EndDate,
	).Scan(&sub.ID, &sub.CreatedAt, &sub.UpdatedAt)
}

func (r *subscriptionRepo) GetByID(ctx context.Context, id string) (*model.Subscription, error) {
	query := `SELECT * FROM subscriptions WHERE id = $1`

	var sub model.Subscription
	log.Printf("Fetching subscription with ID: %s", id)

	err := r.db.GetContext(ctx, &sub, query, id)
	if err != nil {
		return nil, err
	}

	return &sub, nil
}

func (r *subscriptionRepo) Update(ctx context.Context, id string, req *model.UpdateSubscriptionRequest) error {
	var setClauses []string
	var args []interface{}
	argPos := 1

	if req.ServiceName != nil {
		setClauses = append(setClauses, fmt.Sprintf("service_name = $%d", argPos))
		args = append(args, *req.ServiceName)
		argPos++
	}

	if req.Price != nil {
		setClauses = append(setClauses, fmt.Sprintf("price = $%d", argPos))
		args = append(args, *req.Price)
		argPos++
	}

	if req.StartDate != nil {
		setClauses = append(setClauses, fmt.Sprintf("start_date = $%d", argPos))
		args = append(args, *req.StartDate)
		argPos++
	}

	if req.EndDate != nil {
		setClauses = append(setClauses, fmt.Sprintf("end_date = $%d", argPos))
		args = append(args, *req.EndDate)
		argPos++
	}

	if len(setClauses) == 0 {
		return fmt.Errorf("no fields to update")
	}

	setClauses = append(setClauses, "updated_at = CURRENT_TIMESTAMP")
	args = append(args, id)

	query := fmt.Sprintf("UPDATE subscriptions SET %s WHERE id = $%d",
		strings.Join(setClauses, ", "), argPos)

	log.Printf("Updating subscription with ID: %s", id)

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("subscription not found")
	}

	return nil
}

func (r *subscriptionRepo) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM subscriptions WHERE id = $1`

	log.Printf("Deleting subscription with ID: %s", id)

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("subscription not found")
	}

	return nil
}

func (r *subscriptionRepo) List(ctx context.Context, userID *string, serviceName *string) ([]*model.Subscription, error) {
	query := `SELECT * FROM subscriptions WHERE 1=1`
	var args []interface{}
	argPos := 1

	if userID != nil {
		query += fmt.Sprintf(" AND user_id = $%d", argPos)
		args = append(args, *userID)
		argPos++
	}

	if serviceName != nil {
		query += fmt.Sprintf(" AND service_name = $%d", argPos)
		args = append(args, *serviceName)
		argPos++
	}

	query += " ORDER BY created_at DESC"

	log.Printf("Listing subscriptions, userID: %v, serviceName: %v", userID, serviceName)

	var subscriptions []*model.Subscription
	err := r.db.SelectContext(ctx, &subscriptions, query, args...)
	if err != nil {
		return nil, err
	}

	return subscriptions, nil
}

func (r *subscriptionRepo) GetSummary(ctx context.Context, req *model.SummaryRequest) (*model.SubscriptionSummary, error) {
	query := `
		SELECT COALESCE(SUM(price), 0) as total_cost, COUNT(*) as count
		FROM subscriptions 
		WHERE start_date <= $1 
		AND (end_date IS NULL OR end_date >= $2)
	`

	var args []interface{}
	args = append(args, req.EndPeriod, req.StartPeriod)
	argPos := 3

	if req.UserID != nil {
		query += fmt.Sprintf(" AND user_id = $%d", argPos)
		args = append(args, *req.UserID)
		argPos++
	}

	if req.ServiceName != nil {
		query += fmt.Sprintf(" AND service_name = $%d", argPos)
		args = append(args, *req.ServiceName)
		argPos++
	}

	log.Printf("Calculating summary for period %s to %s, userID: %v, serviceName: %v",
		req.StartPeriod, req.EndPeriod, req.UserID, req.ServiceName)

	var summary model.SubscriptionSummary
	err := r.db.GetContext(ctx, &summary, query, args...)
	if err != nil {
		return nil, err
	}

	return &summary, nil
}
