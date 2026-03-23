package software

import (
	"github.com/SaraNguyen999/setup-server/internal/db/model"
)

type SoftwareVersion struct {
	model.SoftwareVersion
}

type CheckUpdateRequest struct {
	AppID   string `json:"app_id"`
	Version string `json:"version"`
}
