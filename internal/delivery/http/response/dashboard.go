package response

import "time"

type DashboardOverviewResponse struct {
	Module      string      `json:"module"`
	Period      string      `json:"period"`
	GeneratedAt time.Time   `json:"generated_at"`
	Results     interface{} `json:"results"`
}
