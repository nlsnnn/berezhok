package handlers

import (
	"log/slog"

	"github.com/nlsnnn/berezhok/internal/shared/validator"
)

type partnerHandler struct {
	partService partnerSvc
	log         *slog.Logger
	validator   *validator.Validator
}

func NewPartnerHandler(partService partnerSvc, log *slog.Logger) partnerHandler {
	return partnerHandler{
		partService: partService,
		log:         log,
		validator:   validator.New(),
	}
}
