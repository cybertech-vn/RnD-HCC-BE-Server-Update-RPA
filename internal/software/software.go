package software

import (
	"github.com/SaraNguyen999/setup-server/internal/db/postgres"
	"github.com/SaraNguyen999/setup-server/pkg/server"
	"github.com/gofiber/fiber/v2"
)

func CheckUpdate(c *fiber.Ctx) error {

	req := new(CheckUpdateRequest)

	if err := c.BodyParser(req); err != nil {
		return server.ResponseReturn(c, true, "invalid request body", nil, 400)
	}

	latest, err := GetLatestVersion(req.AppID)
	if err != nil {
		return server.ResponseReturn(c, true, "failed to get latest version", nil, 500)
	}

	if latest == nil {
		return server.ResponseReturn(c, true, "no versions available", nil, 205)
	}

	if latest.Version != req.Version {
		return server.ResponseReturn(c, true, "update available", fiber.Map{
			"version":  latest.Version,
			"checksum": latest.Checksum,
			"size":     latest.Size,
		}, 200)
	}

	return server.ResponseReturn(c, true, "no update available", nil, 204)
}

func GetVersion(appID string, version string) (*SoftwareVersion, error) {

	if version == "latest" {
		return GetLatestVersion(appID)
	}

	var v SoftwareVersion

	err := postgres.DB.
		Where("app_id = ? AND version = ?", appID, version).
		First(&v).Error

	if err != nil {
		return nil, err
	}

	return &v, nil
}

func GetLatestVersion(appID string) (*SoftwareVersion, error) {

	var v SoftwareVersion

	err := postgres.DB.
		Where("app_id = ?", appID).
		Order("created_at DESC").
		First(&v).Error

	if err != nil {
		return nil, err
	}

	return &v, nil
}

func GetLatestAfter(version string) (*SoftwareVersion, error) {

	var v SoftwareVersion

	err := postgres.DB.
		Where("version > ?", version).
		Order("version DESC").
		First(&v).Error

	if err != nil {
		return nil, err
	}

	return &v, nil
}

func SaveMetadata(meta *SoftwareVersion) error {
	return postgres.DB.Create(meta).Error
}
