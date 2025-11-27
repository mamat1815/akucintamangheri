package services

import (
	"campus-lost-and-found/internal/models"
	"campus-lost-and-found/internal/repository"

	"github.com/google/uuid"
)

type NotificationService struct {
	Repo *repository.NotificationRepository
}

func NewNotificationService(repo *repository.NotificationRepository) *NotificationService {
	return &NotificationService{Repo: repo}
}

func (s *NotificationService) CreateNotification(userID uuid.UUID, title, body, refType string, refID uuid.UUID) error {
	notification := &models.Notification{
		UserID:  userID,
		Title:   title,
		Body:    body,
		RefType: refType,
		RefID:   refID,
	}
	return s.Repo.Create(notification)
}

func (s *NotificationService) GetUserNotifications(userID string) ([]models.Notification, error) {
	return s.Repo.FindByUserID(userID)
}

func (s *NotificationService) MarkAsRead(id string) error {
	return s.Repo.MarkAsRead(id)
}
