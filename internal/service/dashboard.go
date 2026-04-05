package service

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/fatihrizqon/symetra-service/internal/delivery/http/response"
	"github.com/fatihrizqon/symetra-service/internal/repository"
	"github.com/google/uuid"
)

type OverviewResults struct {
	GrowthPercentage float64 `json:"growth_percentage"`
	NewUsers         int64   `json:"new_users"`
	ActiveUsers      int64   `json:"active_users"`
	PreviousNewUsers int64   `json:"previous_new_users"`
	TotalUsers       int64   `json:"total_users"`
	VerifiedUsers    int64   `json:"verified_users"`
	Trend            string  `json:"trend"`
}

type IDashboardService interface {
	Overview(ctx context.Context, companyId uuid.UUID, module string, period string) (response.DashboardOverviewResponse, error)
}

type DashboardService struct {
	IUserRepository         repository.IUserRepository
	ICompanyMemberRepository repository.ICompanyMemberRepository
}

func NewDashboardService(
	repo repository.IUserRepository,
	memberRepo repository.ICompanyMemberRepository,
) IDashboardService {
	return &DashboardService{
		IUserRepository:         repo,
		ICompanyMemberRepository: memberRepo,
	}
}

func (s *DashboardService) Overview(
	ctx context.Context,
	companyId uuid.UUID,
	module string,
	period string,
) (response.DashboardOverviewResponse, error) {

	switch module {
	case "users":
		results, err := s.usersOverview(ctx, companyId, period)
		if err != nil {
			return response.DashboardOverviewResponse{}, err
		}

		return response.DashboardOverviewResponse{
			Module:      module,
			Period:      period,
			GeneratedAt: time.Now(),
			Results:     results,
		}, nil

	default:
		return response.DashboardOverviewResponse{}, errors.New("dashboard module not found")
	}
}

func (s *DashboardService) usersOverview(
	ctx context.Context,
	companyId uuid.UUID,
	period string,
) (OverviewResults, error) {

	start, end := parsePeriod(period)

	duration := end.Sub(start)
	previousStart := start.Add(-duration)
	previousEnd := start

	// Total members in this company
	members, err := s.ICompanyMemberRepository.FindMembersByCompany(companyId)
	if err != nil {
		return OverviewResults{}, err
	}
	totalUsers := int64(len(members))

	// For new/previous users, fall back to platform-wide counts scoped by period
	// (member join date is not tracked separately, so we approximate with user creation dates
	// filtered to members of this company)
	currentUsers, err := s.IUserRepository.CountBetween(ctx, start, end)
	if err != nil {
		return OverviewResults{}, err
	}

	previousUsers, err := s.IUserRepository.CountBetween(ctx, previousStart, previousEnd)
	if err != nil {
		return OverviewResults{}, err
	}

	// Verified users (platform-wide, as email_verified_at is on users table)
	verifiedUsers, err := s.IUserRepository.CountVerified(ctx)
	if err != nil {
		return OverviewResults{}, err
	}

	growth, trend := calculateGrowth(currentUsers, previousUsers)

	results := OverviewResults{
		GrowthPercentage: growth,
		NewUsers:         currentUsers,
		PreviousNewUsers: previousUsers,
		TotalUsers:       totalUsers,
		VerifiedUsers:    verifiedUsers,
		Trend:            trend,
	}

	return results, nil
}

func parsePeriod(period string) (time.Time, time.Time) {
	now := time.Now()

	switch period {
	case "7d":
		return now.AddDate(0, 0, -7), now
	case "30d":
		return now.AddDate(0, 0, -30), now
	case "90d":
		return now.AddDate(0, 0, -90), now
	default:
		return now.AddDate(0, 0, -30), now
	}
}

func calculateGrowth(current, previous int64) (float64, string) {

	if previous == 0 {
		if current == 0 {
			return 0, "flat"
		}
		return 100, "up"
	}

	growth := ((float64(current) - float64(previous)) / float64(previous)) * 100

	growth = math.Round(growth*100) / 100

	switch {
	case growth > 0:
		return growth, "up"
	case growth < 0:
		return growth, "down"
	default:
		return growth, "flat"
	}
}
