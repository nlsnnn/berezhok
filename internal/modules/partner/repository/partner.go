package repository

import (
	"context"
	"errors"
	"math/big"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/nlsnnn/berezhok/internal/adapters/postgresql/sqlc"
	"github.com/nlsnnn/berezhok/internal/lib/pgconverter"
	"github.com/nlsnnn/berezhok/internal/modules/partner/domain"
	partnerErrors "github.com/nlsnnn/berezhok/internal/modules/partner/errors"
)

type PartnerRepo struct {
	q *sqlc.Queries
}

func NewPartnerRepo(q *sqlc.Queries) *PartnerRepo {
	return &PartnerRepo{q: q}
}

func (r *PartnerRepo) FindByID(ctx context.Context, id string) (domain.Partner, error) {
	uid := uuid.MustParse(id)
	p, err := r.q.FindPartnerByID(ctx, uid)
	if err != nil {
		return domain.Partner{}, err
	}
	return partnerToDomain(p), nil
}

func (r *PartnerRepo) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	return r.q.CheckEmailExists(ctx, email)
}

func (r *PartnerRepo) List(ctx context.Context) ([]domain.Partner, error) {
	rows, err := r.q.ListPartners(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]domain.Partner, len(rows))
	for i, p := range rows {
		result[i] = partnerToDomain(p)
	}
	return result, nil
}

func (r *PartnerRepo) Create(ctx context.Context, name string) (domain.Partner, error) {
	p, err := r.q.CreatePartner(ctx, sqlc.CreatePartnerParams{
		BrandName:      name,
		Status:         string(domain.PartnerStatusPendingDocuments),
		CommissionRate: pgtype.Numeric{Int: big.NewInt(20), Exp: -2, Valid: true},
	})
	if err != nil {
		return domain.Partner{}, err
	}
	return partnerToDomain(p), nil
}

func (r *PartnerRepo) GetProfile(ctx context.Context, employeeID string) (domain.PartnerProfile, error) {
	uid := uuid.MustParse(employeeID)
	row, err := r.q.GetPartnerProfile(ctx, uid)
	if err != nil {
		return domain.PartnerProfile{}, err
	}

	commission, err := getCommissionFromRow(row)
	if err != nil {
		return domain.PartnerProfile{}, err
	}

	profile := domain.PartnerProfile{
		Partner: domain.Partner{
			ID:         row.PartnerID.String(),
			BrandName:  row.BrandName,
			Status:     domain.PartnerStatus(row.PartnerStatus),
			Commission: commission,
			CreatedAt:  row.PartnerCreatedAt,
		},
		Employee: domain.Employee{
			ID:                 row.EmployeeID.String(),
			Email:              row.Email,
			Name:               row.EmployeeName.String,
			Role:               domain.EmployeeRole(row.Role),
			MustChangePassword: row.MustChangePassword.Bool,
			CreatedAt:          row.EmployeeCreatedAt,
		},
	}

	if row.LocationID.Valid {
		profile.Location = &domain.LocationSummary{
			ID:        row.LocationID.String(),
			Name:      row.LocationName.String,
			Address:   row.LocationAddress.String,
			Status:    domain.LocationStatusActive,
			CreatedAt: row.LocationCreatedAt.Time,
		}
	}

	// Get all locations for the partner
	locationRows, err := r.q.FindLocationsByPartnerID(ctx, row.PartnerID)
	if err != nil {
		return domain.PartnerProfile{}, err
	}

	locations := make([]domain.LocationSummary, len(locationRows))
	for i, loc := range locationRows {
		locations[i] = domain.LocationSummary{
			ID:        loc.ID.String(),
			Name:      loc.Name,
			Address:   loc.Address,
			Status:    domain.LocationStatus(loc.Status),
			CreatedAt: time.Now(), // TODO: Add CreatedAt to the SQL query and use it here
		}
	}
	profile.Locations = locations

	return profile, nil
}

func (r *PartnerRepo) GetDashboard(ctx context.Context, employeeID string) (domain.PartnerDashboard, error) {
	profile, err := r.GetProfile(ctx, employeeID)
	if err != nil {
		return domain.PartnerDashboard{}, err
	}

	locations := make([]domain.DashboardLocation, len(profile.Locations))
	for i, loc := range profile.Locations {
		count, err := r.q.CountActiveBoxesByLocationID(ctx, uuid.MustParse(loc.ID))
		if err != nil {
			return domain.PartnerDashboard{}, err
		}

		locations[i] = domain.DashboardLocation{
			ID:               loc.ID,
			Name:             loc.Name,
			Address:          loc.Address,
			Status:           loc.Status,
			ActiveBoxesCount: int(count),
		}
	}

	now := time.Now()
	startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	todayStats, err := r.GetStats(ctx, employeeID, domain.StatsFilter{
		Period:           "today",
		DateFrom:         startOfToday,
		DateTo:           startOfToday,
		TopLocationsSort: "revenue_desc",
		TopBoxesSort:     "revenue_desc",
		OrdersSort:       "created_at_desc",
		Limit:            5,
	})
	if err != nil {
		return domain.PartnerDashboard{}, err
	}

	weekStats, err := r.GetStats(ctx, employeeID, domain.StatsFilter{
		Period:           "last_7_days",
		DateFrom:         startOfToday.AddDate(0, 0, -6),
		DateTo:           startOfToday,
		TopLocationsSort: "revenue_desc",
		TopBoxesSort:     "revenue_desc",
		OrdersSort:       "created_at_desc",
		Limit:            5,
	})
	if err != nil {
		return domain.PartnerDashboard{}, err
	}

	nextPayoutDate := nextPartnerPayoutDate()

	legalRow, legalErr := r.q.GetPartnerLegalInfoStatusByPartnerID(ctx, uuid.MustParse(profile.Partner.ID))
	if legalErr != nil && !errors.Is(legalErr, pgx.ErrNoRows) {
		return domain.PartnerDashboard{}, legalErr
	}

	hasLegalInfo := legalErr == nil
	legalInfoStatus := ""
	if hasLegalInfo {
		legalInfoStatus = legalRow.String
	}

	partnerID := uuid.MustParse(profile.Partner.ID)
	_, payoutErr := r.q.GetPayoutDestination(ctx, partnerID)
	hasPayoutDestination := payoutErr == nil

	return domain.PartnerDashboard{
		Partner:              profile.Partner,
		Employee:             profile.Employee,
		Locations:            locations,
		HasLegalInfo:         hasLegalInfo,
		LegalInfoStatus:      legalInfoStatus,
		HasPayoutDestination: hasPayoutDestination,
		Today: domain.DashboardTodayStats{
			PendingConfirmation: findStatusCount(todayStats.StatusBreakdown, "paid"),
			Confirmed:           findStatusCount(todayStats.StatusBreakdown, "confirmed"),
			PickedUp:            findStatusCount(todayStats.StatusBreakdown, "picked_up"),
			Completed:           todayStats.Summary.OrdersCompleted,
		},
		Week: domain.DashboardWeekStats{
			OrdersCompleted: weekStats.Summary.OrdersCompleted,
			GrossRevenue:    weekStats.Summary.GrossRevenue,
			NetRevenue:      weekStats.Summary.NetRevenue,
			AvgRating:       weekStats.Summary.AvgRating,
		},
		Finance: domain.DashboardFinance{
			BalancePending: weekStats.Summary.NetRevenue,
			NextPayoutDate: &nextPayoutDate,
		},
	}, nil
}

func (r *PartnerRepo) GetStats(ctx context.Context, employeeID string, filter domain.StatsFilter) (domain.PartnerStats, error) {
	profile, err := r.GetProfile(ctx, employeeID)
	if err != nil {
		return domain.PartnerStats{}, err
	}

	partnerID := uuid.MustParse(profile.Partner.ID)
	locationID, err := parseOptionalUUID(filter.LocationID)
	if err != nil {
		return domain.PartnerStats{}, partnerErrors.ErrInvalidStatsDateRange
	}

	summaryRow, err := r.q.GetPartnerStatsSummary(ctx, sqlc.GetPartnerStatsSummaryParams{
		PartnerID:  partnerID,
		StartDate:  filter.DateFrom,
		EndDate:    filter.DateTo,
		LocationID: locationID,
		Status:     filter.Status,
	})
	if err != nil {
		return domain.PartnerStats{}, err
	}

	timelineRows, err := r.q.GetPartnerStatsTimeline(ctx, sqlc.GetPartnerStatsTimelineParams{
		PartnerID:  partnerID,
		StartDate:  filter.DateFrom,
		EndDate:    filter.DateTo,
		LocationID: locationID,
		Status:     filter.Status,
	})
	if err != nil {
		return domain.PartnerStats{}, err
	}

	statusRows, err := r.q.GetPartnerStatsStatusBreakdown(ctx, sqlc.GetPartnerStatsStatusBreakdownParams{
		PartnerID:  partnerID,
		StartDate:  filter.DateFrom,
		EndDate:    filter.DateTo,
		LocationID: locationID,
		Status:     filter.Status,
	})
	if err != nil {
		return domain.PartnerStats{}, err
	}

	topLocationRows, err := r.q.GetPartnerStatsTopLocations(ctx, sqlc.GetPartnerStatsTopLocationsParams{
		PartnerID:  partnerID,
		StartDate:  filter.DateFrom,
		EndDate:    filter.DateTo,
		LocationID: locationID,
		Status:     filter.Status,
		Sort:       filter.TopLocationsSort,
	})
	if err != nil {
		return domain.PartnerStats{}, err
	}

	topBoxRows, err := r.q.GetPartnerStatsTopBoxes(ctx, sqlc.GetPartnerStatsTopBoxesParams{
		PartnerID:  partnerID,
		StartDate:  filter.DateFrom,
		EndDate:    filter.DateTo,
		LocationID: locationID,
		Status:     filter.Status,
		Sort:       filter.TopBoxesSort,
	})
	if err != nil {
		return domain.PartnerStats{}, err
	}

	orderRows, err := r.q.ListPartnerStatsOrders(ctx, sqlc.ListPartnerStatsOrdersParams{
		PartnerID:  partnerID,
		StartDate:  filter.DateFrom,
		EndDate:    filter.DateTo,
		LocationID: locationID,
		Status:     filter.Status,
		Sort:       filter.OrdersSort,
		PageLimit:  int32(filter.Limit),
		PageOffset: int32(filter.Offset),
	})
	if err != nil {
		return domain.PartnerStats{}, err
	}

	totalOrders, err := r.q.CountPartnerStatsOrders(ctx, sqlc.CountPartnerStatsOrdersParams{
		PartnerID:  partnerID,
		StartDate:  filter.DateFrom,
		EndDate:    filter.DateTo,
		LocationID: locationID,
		Status:     filter.Status,
	})
	if err != nil {
		return domain.PartnerStats{}, err
	}

	timeline := make([]domain.PartnerStatsTimelinePoint, len(timelineRows))
	for i, row := range timelineRows {
		timeline[i] = domain.PartnerStatsTimelinePoint{
			Date:            row.Date,
			OrdersTotal:     int(row.OrdersTotal),
			OrdersCompleted: int(row.OrdersCompleted),
			GrossRevenue:    int(row.GrossRevenue),
			NetRevenue:      int(row.NetRevenue),
		}
	}

	statusBreakdown := make([]domain.PartnerStatsStatusBreakdownItem, len(statusRows))
	for i, row := range statusRows {
		statusBreakdown[i] = domain.PartnerStatsStatusBreakdownItem{
			Status: row.Status,
			Count:  int(row.Count),
			Share:  row.Share,
		}
	}

	topLocations := make([]domain.PartnerStatsLocation, len(topLocationRows))
	for i, row := range topLocationRows {
		topLocations[i] = domain.PartnerStatsLocation{
			LocationID:      row.LocationID.String(),
			Name:            row.Name,
			Address:         row.Address,
			OrdersTotal:     int(row.OrdersTotal),
			OrdersCompleted: int(row.OrdersCompleted),
			GrossRevenue:    int(row.GrossRevenue),
			NetRevenue:      int(row.NetRevenue),
			AvgRating:       row.AvgRating,
		}
	}

	topBoxes := make([]domain.PartnerStatsBox, len(topBoxRows))
	for i, row := range topBoxRows {
		topBoxes[i] = domain.PartnerStatsBox{
			BoxID:           row.BoxID.String(),
			Name:            row.Name,
			ImageURL:        row.ImageUrl,
			LocationName:    row.LocationName,
			OrdersTotal:     int(row.OrdersTotal),
			OrdersCompleted: int(row.OrdersCompleted),
			GrossRevenue:    int(row.GrossRevenue),
			NetRevenue:      int(row.NetRevenue),
		}
	}

	orders := make([]domain.PartnerStatsOrder, len(orderRows))
	for i, row := range orderRows {
		statusValue := string(row.Status)
		orders[i] = domain.PartnerStatsOrder{
			ID:              row.ID.String(),
			Status:          statusValue,
			PickupCode:      row.PickupCode,
			Amount:          pgconverter.NumericToDecimalOrZero(row.Amount).InexactFloat64(),
			BoxName:         row.BoxName,
			BoxImageURL:     row.BoxImageUrl,
			CustomerPhone:   row.CustomerPhone,
			CustomerName:    row.CustomerName,
			LocationID:      row.LocationID.String(),
			LocationName:    row.LocationName,
			LocationAddress: row.LocationAddress,
			PickupTimeStart: row.PickupTimeStart,
			PickupTimeEnd:   row.PickupTimeEnd,
			CreatedAt:       row.CreatedAt,
			CanPickup:       statusValue == "confirmed",
		}
	}

	total := int(totalOrders)
	return domain.PartnerStats{
		Summary: domain.PartnerStatsSummary{
			OrdersTotal:               int(summaryRow.OrdersTotal),
			OrdersCompleted:           int(summaryRow.OrdersCompleted),
			OrdersCancelled:           int(summaryRow.OrdersCancelled),
			OrdersPendingConfirmation: int(summaryRow.OrdersPendingConfirmation),
			GrossRevenue:              int(summaryRow.GrossRevenue),
			NetRevenue:                int(summaryRow.NetRevenue),
			AvgOrderValue:             summaryRow.AvgOrderValue,
			AvgRating:                 summaryRow.AvgRating,
			ReviewsCount:              int(summaryRow.ReviewsCount),
		},
		Timeline:        timeline,
		StatusBreakdown: statusBreakdown,
		TopLocations:    topLocations,
		TopBoxes:        topBoxes,
		Orders:          orders,
		Meta: domain.PartnerStatsMeta{
			Period:           filter.Period,
			DateFrom:         filter.DateFrom,
			DateTo:           filter.DateTo,
			LocationID:       filter.LocationID,
			Status:           filter.Status,
			TopLocationsSort: filter.TopLocationsSort,
			TopBoxesSort:     filter.TopBoxesSort,
			OrdersSort:       filter.OrdersSort,
			Pagination: domain.PartnerStatsPagination{
				Total:   total,
				Limit:   filter.Limit,
				Offset:  filter.Offset,
				HasMore: filter.Offset+len(orders) < total,
			},
		},
	}, nil
}

func (r *PartnerRepo) UpdateEmployeePassword(ctx context.Context, employeeID, newHash string) error {
	uid := uuid.MustParse(employeeID)
	return r.q.UpdatePartnerEmployeePassword(ctx, sqlc.UpdatePartnerEmployeePasswordParams{
		ID:                 uid,
		PasswordHash:       newHash,
		MustChangePassword: pgtype.Bool{Valid: true, Bool: false},
	})
}

func (r *PartnerRepo) UpdateEmployeeName(ctx context.Context, employeeID, name string) error {
	uid := uuid.MustParse(employeeID)
	return r.q.UpdatePartnerEmployeeName(ctx, sqlc.UpdatePartnerEmployeeNameParams{
		ID:   uid,
		Name: pgtype.Text{String: name, Valid: true},
	})
}

func (r *PartnerRepo) UpsertLegalInfo(ctx context.Context, info domain.LegalInfo) error {
	partnerID := uuid.MustParse(info.PartnerID)

	_, err := r.q.UpsertPartnerLegalInfo(ctx, sqlc.UpsertPartnerLegalInfoParams{
		PartnerID:    partnerID,
		Inn:          info.Inn,
		Ogrn:         pgtype.Text{String: info.Ogrn, Valid: info.Ogrn != ""},
		Kpp:          pgtype.Text{String: info.Kpp, Valid: info.Kpp != ""},
		LegalAddress: info.LegalAddress,
		LegalName:    info.LegalName,
	})
	return err
}

func (r *PartnerRepo) UpdateStatus(ctx context.Context, partnerID string, status domain.PartnerStatus) error {
	uid := uuid.MustParse(partnerID)
	return r.q.UpdatePartnerStatusByID(ctx, sqlc.UpdatePartnerStatusByIDParams{
		ID:     uid,
		Status: string(status),
	})
}

func partnerToDomain(p sqlc.Partner) domain.Partner {
	commission, err := getCommission(p)
	if err != nil {
		commission = domain.Commission{}
	}

	return domain.Partner{
		ID:         p.ID.String(),
		BrandName:  p.BrandName,
		LogoURL:    p.LogoUrl.String,
		Status:     domain.PartnerStatus(p.Status),
		Commission: commission,
		CreatedAt:  p.CreatedAt,
	}
}

func getCommissionFromRow(row sqlc.GetPartnerProfileRow) (domain.Commission, error) {
	commissionRate, _ := numericToFloat64(row.CommissionRate.(pgtype.Numeric))

	var promoUntil *time.Time
	if row.PromoCommissionUntil.Valid {
		promoUntil = &row.PromoCommissionUntil.Time
	}

	commission, err := domain.NewCommission(commissionRate, promoUntil)
	if err != nil {
		return domain.Commission{}, err
	}

	return commission, nil
}

func getCommission(p sqlc.Partner) (domain.Commission, error) {
	commissionRate, _ := p.CommissionRate.Int.Float64()

	var promoUntil *time.Time
	if p.PromoCommissionUntil.Valid {
		promoUntil = &p.PromoCommissionUntil.Time
	}

	commission, err := domain.NewCommission(commissionRate, promoUntil)
	if err != nil {
		return domain.Commission{}, err
	}

	return commission, nil
}

func numericToFloat64(n pgtype.Numeric) (float64, error) {
	if !n.Valid || n.Int == nil {
		return 0, nil
	}

	f := new(big.Float).SetInt(n.Int)
	if n.Exp > 0 {
		mul := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(n.Exp)), nil))
		f.Mul(f, mul)
	} else if n.Exp < 0 {
		div := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(-n.Exp)), nil))
		f.Quo(f, div)
	}

	result, _ := f.Float64()
	return result, nil
}

func parseOptionalUUID(value string) (pgtype.UUID, error) {
	if value == "" {
		return pgtype.UUID{}, nil
	}

	parsed, err := uuid.Parse(value)
	if err != nil {
		return pgtype.UUID{}, err
	}

	return pgtype.UUID{Bytes: parsed, Valid: true}, nil
}

func findStatusCount(items []domain.PartnerStatsStatusBreakdownItem, status string) int {
	for _, item := range items {
		if item.Status == status {
			return item.Count
		}
	}

	return 0
}

func nextPartnerPayoutDate() time.Time {
	now := time.Now()
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}

	daysUntilNextMonday := 8 - weekday
	nextMonday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, daysUntilNextMonday)
	return nextMonday
}
