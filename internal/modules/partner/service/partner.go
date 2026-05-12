package service

import (
	"context"
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"

	"github.com/nlsnnn/berezhok/internal/modules/partner/domain"
	"github.com/nlsnnn/berezhok/internal/modules/partner/errors"
	"github.com/nlsnnn/berezhok/internal/shared/auth"
	"github.com/nlsnnn/berezhok/internal/shared/authz"
)

type ChangePasswordInput struct {
	UserID          string
	CurrentPassword string
	NewPassword     string
}

type AddLegalInfoInput struct {
	PartnerID    string
	Inn          string
	Ogrn         string
	Kpp          string
	LegalAddress string
	LegalName    string
}

const (
	defaultStatsPeriod      = "last_7_days"
	defaultTopLocationsSort = "revenue_desc"
	defaultTopBoxesSort     = "revenue_desc"
	defaultOrdersSort       = "created_at_desc"
	defaultStatsLimit       = 20
	maxStatsLimit           = 100
	statsPeriodToday        = "today"
	statsPeriodLast7Days    = "last_7_days"
	statsPeriodLast30Days   = "last_30_days"
	statsPeriodThisMonth    = "this_month"
	statsPeriodLastMonth    = "last_month"
)

type partService struct {
	repo    partnerRepo
	empRepo employeeRepoForPartner
}

type partnerRepo interface {
	FindByID(ctx context.Context, id string) (domain.Partner, error)
	List(ctx context.Context) ([]domain.Partner, error)
	Create(ctx context.Context, name string) (domain.Partner, error)
	CheckEmailExists(ctx context.Context, email string) (bool, error)
	GetProfile(ctx context.Context, employeeID string) (domain.PartnerProfile, error)
	GetDashboard(ctx context.Context, employeeID string) (domain.PartnerDashboard, error)
	GetStats(ctx context.Context, employeeID string, filter domain.StatsFilter) (domain.PartnerStats, error)
	UpdateEmployeePassword(ctx context.Context, employeeID, newHash string) error
	UpsertLegalInfo(ctx context.Context, info domain.LegalInfo) error
	UpdateStatus(ctx context.Context, partnerID string, status domain.PartnerStatus) error
}

type employeeRepoForPartner interface {
	FindByID(ctx context.Context, id string) (domain.Employee, error)
}

func NewPartnerService(repo partnerRepo, empRepo employeeRepoForPartner) *partService {
	return &partService{repo: repo, empRepo: empRepo}
}

func (s *partService) List(ctx context.Context) ([]domain.Partner, error) {
	return s.repo.List(ctx)
}

func (s *partService) FindByID(ctx context.Context, id string) (domain.Partner, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *partService) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	emailExists, err := s.repo.CheckEmailExists(ctx, email)
	if err != nil {
		return false, err
	}

	return emailExists, nil
}

func (s *partService) Create(ctx context.Context, name string) (domain.Partner, error) {
	return s.repo.Create(ctx, name)
}

func (s *partService) ChangePassword(ctx context.Context, input ChangePasswordInput) error {
	employee, err := s.empRepo.FindByID(ctx, input.UserID)
	if err != nil {
		return err
	}

	if !auth.Compare(employee.PasswordHash, input.CurrentPassword) {
		return errors.ErrInvalidCredentials
	}

	if input.CurrentPassword == input.NewPassword {
		return errors.ErrPasswordUnchanged
	}

	newHash, err := auth.Hash(input.NewPassword)
	if err != nil {
		return err
	}

	return s.repo.UpdateEmployeePassword(ctx, input.UserID, newHash)
}

func (s *partService) Profile(ctx context.Context, userID string) (domain.PartnerProfile, error) {
	return s.repo.GetProfile(ctx, userID)
}

func (s *partService) Dashboard(ctx context.Context, userID string) (domain.PartnerDashboard, error) {
	return s.repo.GetDashboard(ctx, userID)
}

func (s *partService) Stats(ctx context.Context, actor authz.PartnerActor, filter domain.StatsFilter) (domain.PartnerStats, error) {
	resolvedFilter, err := NormalizeStatsFilter(filter)
	if err != nil {
		return domain.PartnerStats{}, err
	}

	if actor.Role != authz.RoleOwner {
		if actor.LocationID == nil {
			return domain.PartnerStats{}, authz.ErrLocationScopeDenied
		}
		resolvedFilter.LocationID = actor.LocationID.String()
	}

	return s.repo.GetStats(ctx, actor.EmployeeID.String(), resolvedFilter)
}

func (s *partService) AddLegalInfo(ctx context.Context, input AddLegalInfoInput) error {
	// TODO: move validation to handler (?)
	if !isDigitsOnly(input.Inn) || (len(input.Inn) != 10 && len(input.Inn) != 12) {
		return errors.ErrInvalidINN
	}

	if input.Ogrn != "" && (!isDigitsOnly(input.Ogrn) || (len(input.Ogrn) != 13 && len(input.Ogrn) != 15)) {
		return errors.ErrInvalidOGRN
	}

	if input.Kpp != "" && (!isDigitsOnly(input.Kpp) || len(input.Kpp) != 9) {
		return errors.ErrInvalidKPP
	}

	partner, err := s.repo.FindByID(ctx, input.PartnerID)
	if err != nil {
		return err
	}

	if partner.Status != domain.PartnerStatusPendingDocuments && partner.Status != domain.PartnerStatusActive {
		return errors.ErrPartnerStatusInvalid
	}

	err = s.repo.UpsertLegalInfo(ctx, domain.LegalInfo{
		PartnerID:    input.PartnerID,
		Inn:          input.Inn,
		Ogrn:         input.Ogrn,
		Kpp:          input.Kpp,
		LegalAddress: input.LegalAddress,
		LegalName:    input.LegalName,
	})
	if err != nil {
		return err
	}

	return s.repo.UpdateStatus(ctx, input.PartnerID, domain.PartnerStatusActive)
}

func (s *partService) CanActivateBoxes(ctx context.Context, partnerID uuid.UUID) (bool, error) {
	partner, err := s.repo.FindByID(ctx, partnerID.String())
	if err != nil {
		return false, err
	}

	return partner.Status != domain.PartnerStatusPendingDocuments, nil
}

func isDigitsOnly(value string) bool {
	for _, r := range value {
		if !unicode.IsDigit(r) {
			return false
		}
	}

	return true
}

func NormalizeStatsFilter(filter domain.StatsFilter) (domain.StatsFilter, error) {
	resolved := filter
	resolved.Period = strings.TrimSpace(resolved.Period)
	resolved.LocationID = strings.TrimSpace(resolved.LocationID)
	resolved.Status = strings.TrimSpace(resolved.Status)
	resolved.TopLocationsSort = strings.TrimSpace(resolved.TopLocationsSort)
	resolved.TopBoxesSort = strings.TrimSpace(resolved.TopBoxesSort)
	resolved.OrdersSort = strings.TrimSpace(resolved.OrdersSort)

	if resolved.Limit <= 0 {
		resolved.Limit = defaultStatsLimit
	}
	if resolved.Limit > maxStatsLimit {
		resolved.Limit = maxStatsLimit
	}
	if resolved.Offset < 0 {
		resolved.Offset = 0
	}

	if resolved.TopLocationsSort == "" {
		resolved.TopLocationsSort = defaultTopLocationsSort
	}
	if !isAllowedValue(resolved.TopLocationsSort, defaultTopLocationsSort, "orders_desc", "rating_desc", "name_asc") {
		return domain.StatsFilter{}, errors.ErrInvalidStatsSort
	}

	if resolved.TopBoxesSort == "" {
		resolved.TopBoxesSort = defaultTopBoxesSort
	}
	if !isAllowedValue(resolved.TopBoxesSort, defaultTopBoxesSort, "orders_desc", "name_asc") {
		return domain.StatsFilter{}, errors.ErrInvalidStatsSort
	}

	if resolved.OrdersSort == "" {
		resolved.OrdersSort = defaultOrdersSort
	}
	if !isAllowedValue(resolved.OrdersSort, defaultOrdersSort, "created_at_asc", "amount_desc", "amount_asc", "pickup_time_desc", "pickup_time_asc") {
		return domain.StatsFilter{}, errors.ErrInvalidStatsSort
	}

	if resolved.LocationID != "" {
		if _, err := uuid.Parse(resolved.LocationID); err != nil {
			return domain.StatsFilter{}, fmt.Errorf("%w: invalid location_id", errors.ErrInvalidStatsDateRange)
		}
	}

	if resolved.Period == "" {
		if resolved.DateFrom.IsZero() || resolved.DateTo.IsZero() {
			resolved.Period = defaultStatsPeriod
		}
	}

	now := time.Now()
	switch resolved.Period {
	case statsPeriodToday:
		resolved.DateFrom = startOfDay(now)
		resolved.DateTo = startOfDay(now)
	case statsPeriodLast7Days:
		today := startOfDay(now)
		resolved.DateFrom = today.AddDate(0, 0, -6)
		resolved.DateTo = today
	case statsPeriodLast30Days:
		today := startOfDay(now)
		resolved.DateFrom = today.AddDate(0, 0, -29)
		resolved.DateTo = today
	case statsPeriodThisMonth:
		today := startOfDay(now)
		resolved.DateFrom = time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, today.Location())
		resolved.DateTo = today
	case statsPeriodLastMonth:
		today := startOfDay(now)
		firstDayCurrentMonth := time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, today.Location())
		lastMonthDate := firstDayCurrentMonth.AddDate(0, -1, 0)
		resolved.DateFrom = time.Date(lastMonthDate.Year(), lastMonthDate.Month(), 1, 0, 0, 0, 0, today.Location())
		resolved.DateTo = firstDayCurrentMonth.AddDate(0, 0, -1)
	case "":
		if resolved.DateFrom.IsZero() || resolved.DateTo.IsZero() {
			return domain.StatsFilter{}, errors.ErrInvalidStatsDateRange
		}
		resolved.Period = "custom"
	default:
		return domain.StatsFilter{}, errors.ErrInvalidStatsPeriod
	}

	if resolved.Period == "custom" {
		if resolved.DateFrom.IsZero() || resolved.DateTo.IsZero() {
			return domain.StatsFilter{}, errors.ErrInvalidStatsDateRange
		}
	}

	resolved.DateFrom = startOfDay(resolved.DateFrom)
	resolved.DateTo = startOfDay(resolved.DateTo)
	if resolved.DateFrom.After(resolved.DateTo) {
		return domain.StatsFilter{}, errors.ErrInvalidStatsDateRange
	}

	return resolved, nil
}

func startOfDay(value time.Time) time.Time {
	return time.Date(value.Year(), value.Month(), value.Day(), 0, 0, 0, 0, value.Location())
}

func isAllowedValue(value string, allowed ...string) bool {
	for _, item := range allowed {
		if value == item {
			return true
		}
	}

	return false
}
