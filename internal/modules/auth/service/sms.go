package service

import "context"

type smsService struct {
	storage smsStorage
	sender  smsSender
}

type smsSender interface {
	SendCode(ctx context.Context, phone, code string) error
}

type smsStorage interface {
	Save(ctx context.Context, phone, code string) error
	Validate(ctx context.Context, phone, code string) (bool, error)
}

func NewSMSService(storage smsStorage, sender smsSender) *smsService {
	return &smsService{storage: storage, sender: sender}
}

func (s *smsService) SendCode(ctx context.Context, phone, code string) error {
	if err := s.storage.Save(ctx, phone, code); err != nil {
		return err
	}

	if err := s.sender.SendCode(ctx, phone, code); err != nil {
		return err
	}

	return nil
}

func (s *smsService) ValidateCode(ctx context.Context, phone, code string) (bool, error) {
	return s.storage.Validate(ctx, phone, code)
}
