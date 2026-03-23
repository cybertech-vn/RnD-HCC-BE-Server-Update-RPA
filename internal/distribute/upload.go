package distribute

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/SaraNguyen999/setup-server/internal/db/minio"
	"github.com/SaraNguyen999/setup-server/internal/db/model"
	"github.com/SaraNguyen999/setup-server/internal/software"
	"github.com/SaraNguyen999/setup-server/pkg/server"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func UploadVersion(c *fiber.Ctx) error {
	appID := c.Locals("app_id").(string)
	version := c.FormValue("version")

	if appID == "" || version == "" {
		return server.ResponseReturn(c, true, "app_id and version required", nil, 422)
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return server.ResponseReturn(c, true, "file required", nil, 422)
	}

	// check version tồn tại
	if v, _ := software.GetVersion(appID, version); v != nil {
		return server.ResponseReturn(c, true, "version already exists", nil, 409)
	}

	// mở stream từ request
	file, err := fileHeader.Open()
	if err != nil {
		return server.ResponseReturn(c, true, "failed to open file", nil, 500)
	}
	defer file.Close()

	// tạo hash
	hash := sha256.New()

	// tee reader: vừa upload vừa tính checksum
	reader := io.TeeReader(file, hash)

	objectName := fmt.Sprintf("%s/%s", appID, filepath.Base(fileHeader.Filename))

	// upload lên MinIO
	err = minio.MinioDB.UploadStream("", objectName, reader, fileHeader.Size, fileHeader.Header.Get("Content-Type"))
	if err != nil {
		fmt.Println("MinIO upload error:", err)
		return server.ResponseReturn(c, true, "upload to minio failed", nil, 500)
	}

	// lấy checksum
	checksum := hex.EncodeToString(hash.Sum(nil))

	meta := &software.SoftwareVersion{
		SoftwareVersion: model.SoftwareVersion{
			ID:        uuid.New().String(),
			AppID:     appID,
			Version:   version,
			FileName:  filepath.Base(fileHeader.Filename),
			Checksum:  checksum,
			Size:      fileHeader.Size,
			CreatedAt: time.Now().Unix(),
		},
	}

	// lưu metadata
	if err := software.SaveMetadata(meta); err != nil {
		return server.ResponseReturn(c, true, "failed to save metadata", nil, 500)
	}

	return server.ResponseReturn(c, false, "file uploaded successfully", meta, 200)
}
