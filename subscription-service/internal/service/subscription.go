package service

import (
	"context"
	"log"

	"subscription-service/internal/model"
	"subscription-service/internal/repository"
)

type SubscriptionService interface {
	CreateSubscription(ctx context.Context, req *model.CreateSubscriptionRequest) (*model.Subscription, error)
	GetSubscription(ctx context.Context, id string) (*model.Subscription, error)
	UpdateSubscription(ctx context.Context, id string, req *model.UpdateSubscriptionRequest) error
	DeleteSubscription(ctx context.Context, id string) error
	ListSubscriptions(ctx context.Context, userID *string, serviceName *string) ([]*model.Subscription, error)
	GetSummary(ctx context.Context, req *model.SummaryRequest) (*model.SubscriptionSummary, error)
}

type subscriptionService struct {
	repo repository.SubscriptionRepository
}

func NewSubscriptionService(repo repository.SubscriptionRepository) SubscriptionService {
	return &subscriptionService{repo: repo}
}

func (s *subscriptionService) CreateSubscription(ctx context.Context, req *model.CreateSubscriptionRequest) (*model.Subscription, error) {
	log.Printf("Creating subscription for user %s", req.UserID)

	subscription := &model.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	}

	if err := s.repo.Create(ctx, subscription); err != nil {
		log.Printf("Error creating subscription: %v", err)
		return nil, err
	}

	log.Printf("Subscription created successfully with ID: %s", subscription.ID)
	return subscription, nil
}

func (s *subscriptionService) GetSubscription(ctx context.Context, id string) (*model.Subscription, error) {
	log.Printf("Getting subscription with ID: %s", id)

	subscription, err := s.repo.GetByID(ctx, id)
	if err != nil {
		log.Printf("Error getting subscription %s: %v", id, err)
		return nil, err
	}

	return subscription, nil
}

func (s *subscriptionService) UpdateSubscription(ctx context.Context, id string, req *model.UpdateSubscriptionRequest) error {
	log.Printf("Updating subscription with ID: %s", id)

	if err := s.repo.Update(ctx, id, req); err != nil {
		log.Printf("Error updating subscription %s: %v", id, err)
		return err
	}

	log.Printf("Subscription %s updated successfully", id)
	return nil
}

func (s *subscriptionService) DeleteSubscription(ctx context.Context, id string) error {
	log.Printf("Deleting subscription with ID: %s", id)

	if err := s.repo.Delete(ctx, id); err != nil {
		log.Printf("Error deleting subscription %s: %v", id, err)
		return err
	}

	log.Printf("Subscription %s deleted successfully", id)
	return nil
}

func (s *subscriptionService) ListSubscriptions(ctx context.Context, userID *string, serviceName *string) ([]*model.Subscription, error) {
	log.Printf("Listing subscriptions, filters - userID: %v, serviceName: %v", userID, serviceName)

	subscriptions, err := s.repo.List(ctx, userID, serviceName)
	if err != nil {
		log.Printf("Error listing subscriptions: %v", err)
		return nil, err
	}

	log.Printf("Found %d subscriptions", len(subscriptions))
	return subscriptions, nil
}

func (s *subscriptionService) GetSummary(ctx context.Context, req *model.SummaryRequest) (*model.SubscriptionSummary, error) {
	log.Printf("Getting summary for period %s to %s", req.StartPeriod, req.EndPeriod)

	summary, err := s.repo.GetSummary(ctx, req)
	if err != nil {
		log.Printf("Error getting summary: %v", err)
		return nil, err
	}

	log.Printf("Summary calculated: total cost %d, count %d", summary.TotalCost, summary.Count)
	return summary, nil
}
