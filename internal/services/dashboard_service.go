package services

import (
	"fmt"
	"math"
	"time"
	"github.com/maxintelinno/sport-hub-profile/internal/models"
	"github.com/maxintelinno/sport-hub-profile/internal/repositories"
)

type DashboardService interface {
	GetDashboard(userID string) (*models.DashboardResponse, error)
}

type dashboardService struct {
	userRepo      repositories.UserRepository
	dashboardRepo repositories.DashboardRepository
}

func NewDashboardService(userRepo repositories.UserRepository, dashboardRepo repositories.DashboardRepository) DashboardService {
	return &dashboardService{userRepo: userRepo, dashboardRepo: dashboardRepo}
}

func (s *dashboardService) GetDashboard(userID string) (*models.DashboardResponse, error) {
	// 1. Get Owner Info
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	avatarInitial := ""
	if len(user.Fullname) > 0 {
		avatarInitial = string([]rune(user.Fullname)[0])
	}

	// 2. Get Data from Repository
	sub, _ := s.dashboardRepo.GetActiveSubscription(userID)
	summaryStats, _ := s.dashboardRepo.GetSummaryStats(userID)
	trend, _ := s.dashboardRepo.GetRevenueTrend(userID)
	todayBookings, _ := s.dashboardRepo.GetTodayBookingCount(userID)

	if summaryStats == nil {
		summaryStats = &models.DashboardSummary{}
	}

	response := &models.DashboardResponse{
		Owner: models.DashboardOwner{
			ID:            user.ID,
			Fullname:      user.Fullname,
			Phone:         user.Phone,
			AvatarInitial: avatarInitial,
		},
		Summary:        *summaryStats,
		RevenueTrend7d: trend,
		Alerts:         []models.DashboardAlert{},
		NextActions:    []models.DashboardAction{},
	}

	// 3. Plan Info & Business Rules
	if sub != nil {
		trialDaysLeft := 0
		if sub.Status == "trial" && sub.TrialEndAt != nil {
			trialDaysLeft = int(math.Max(0, time.Until(*sub.TrialEndAt).Hours()/24))
		}

		response.Plan = models.DashboardPlan{
			Code:            sub.Plan.Code,
			Name:            sub.Plan.Name,
			Status:          sub.Status,
			TrialDaysLeft:   trialDaysLeft,
			IsTrial:         sub.Status == "trial",
			PriceAfterTrial: sub.Plan.Price,
		}

		// Next Action: Trial warning
		if sub.Status == "trial" && trialDaysLeft <= 3 {
			response.NextActions = append(response.NextActions, models.DashboardAction{
				Title:      fmt.Sprintf("ลอง PRO ฟรี (เหลือ %d วันสุดท้าย)", trialDaysLeft),
				Message:    "ยกเลิกได้ทุกเมื่อ • ไม่มีการผูกมัด • ไม่คิดเงินทันที",
				ActionText: "เริ่มใช้ฟรีเลย",
				ActionType: "start_trial",
			})
		}

		// Next Action: Create field suggestion
		if sub.Plan.MaxFields != nil {
			maxFields := int64(*sub.Plan.MaxFields)
			if summaryStats.FieldCount < maxFields {
				response.NextActions = append(response.NextActions, models.DashboardAction{
					Title:      "เพิ่มสนามใหม่ เพื่อรับลูกค้าเพิ่ม",
					Message:    fmt.Sprintf("คุณสามารถเพิ่มสนามได้อีก %d สนาม", maxFields-summaryStats.FieldCount),
					ActionText: "เพิ่มสนาม",
					ActionType: "create_field",
				})
			}
		} else {
			// Unlimited fields plan
			response.NextActions = append(response.NextActions, models.DashboardAction{
				Title:      "เพิ่มสนามใหม่ เพื่อรับลูกค้าเพิ่ม",
				Message:    "คุณสามารถเพิ่มสนามได้ไม่จำกัดจำนวน",
				ActionText: "เพิ่มสนาม",
				ActionType: "create_field",
			})
		}

		// Upsell: If on free plan with 1 field
		if sub.Plan.Code == "free" && summaryStats.FieldCount >= 1 {
			response.Upsell = &models.DashboardUpsell{
				Title:      "เพิ่มรายได้ +27,000/เดือน",
				Subtitle:   "จ่ายแค่ ฿999/เดือน • กำไรเพิ่ม ~฿26K+",
				ActionText: "ดูรายละเอียด",
				ActionType: "open_upgrade",
			}
		}
	} else {
		// No subscription (Fallback to default free plan view)
		response.Plan = models.DashboardPlan{
			Code:   "none",
			Name:   "No Plan",
			Status: "inactive",
		}
	}

	// Alert: No bookings today
	if todayBookings == 0 {
		response.Alerts = append(response.Alerts, models.DashboardAlert{
			Type:       "warning",
			Title:      "วันนี้ยังไม่มีการจอง",
			Message:    "เปิดโปรเพื่อเพิ่มลูกค้าทันที",
			ActionText: "สร้างโปรโมชัน",
			ActionType: "open_promotion",
		})
	}

	return response, nil
}
