package partneradapters

import (
	"context"

	"github.com/nlsnnn/berezhok/internal/adapters/rabbitmq"
)

const (
	// Notification templates
	TemplatePartnerApproved   = "application-approved-1"
	TemplatePartnerRejected   = "partner_rejected"
	TemplateEmployeeInvite    = "application-approved-1"
	TemplateLegalInfoVerified = "partner_legal_info_verified"
	TemplateLegalInfoRejected = "partner_legal_info_rejected"
)

type NotificationAdapter struct {
	publisher *rabbitmq.NotificationPublisher
}

func NewNotificationAdapter(publisher *rabbitmq.NotificationPublisher) *NotificationAdapter {
	return &NotificationAdapter{publisher: publisher}
}

// Send partner approval notification with credentials
func (a *NotificationAdapter) SendPartnerApprovalNotification(ctx context.Context, email, name, password string) error {
	return a.publisher.SendEmail(ctx, email, TemplatePartnerApproved, map[string]string{
		"name":     name,
		"password": password,
	})
}

// Send partner rejection notification with reason
func (a *NotificationAdapter) SendPartnerRejectionNotification(ctx context.Context, email, name, reason string) error {
	return a.publisher.SendEmail(ctx, email, TemplatePartnerRejected, map[string]string{
		"name":   name,
		"reason": reason,
	})
}

func (a *NotificationAdapter) SendEmployeeInviteNotification(ctx context.Context, email, name, password string) error {
	return a.publisher.SendEmail(ctx, email, TemplateEmployeeInvite, map[string]string{
		"name":     name,
		"password": password,
	})
}

func (a *NotificationAdapter) SendPartnerLegalInfoVerifiedNotification(ctx context.Context, email, name string) error {
	return a.publisher.SendEmail(ctx, email, TemplateLegalInfoVerified, map[string]string{
		"name": name,
	})
}

func (a *NotificationAdapter) SendPartnerLegalInfoRejectedNotification(ctx context.Context, email, name, reason string) error {
	return a.publisher.SendEmail(ctx, email, TemplateLegalInfoRejected, map[string]string{
		"name":   name,
		"reason": reason,
	})
}
