package sms

import (
	"context"
	"log"
)

type ConsoleSender struct{}

func NewConsoleSender() *ConsoleSender {
	return &ConsoleSender{}
}

func (s *ConsoleSender) SendCode(ctx context.Context, phone, code string) error {
	log.Printf("[SMS STUB] Sending code %s to phone %s", code, phone)
	return nil
}
