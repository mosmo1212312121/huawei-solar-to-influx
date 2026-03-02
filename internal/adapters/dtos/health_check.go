package dtos

type HealthCheckResponse struct {
	Status   Status `json:"status"`
	Database Status `json:"database"`
}
type Status string

const (
	UP   Status = "up"
	DOWN Status = "down"
)
