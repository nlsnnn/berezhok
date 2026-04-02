package domain

import "time"

type PartnerDashboard struct {
	Partner   Partner
	Employee  Employee
	Locations []DashboardLocation
	Today     DashboardTodayStats
	Week      DashboardWeekStats
	Finance   DashboardFinance
}

type DashboardLocation struct {
	ID               string
	Name             string
	Address          string
	Status           LocationStatus
	ActiveBoxesCount int
}

type DashboardTodayStats struct {
	PendingConfirmation int
	Confirmed           int
	PickedUp            int
	Completed           int
}

type DashboardWeekStats struct {
	OrdersCompleted int
	GrossRevenue    int
	NetRevenue      int
	AvgRating       float64
}

type DashboardFinance struct {
	BalancePending int
	NextPayoutDate *time.Time
}
