package application

import (
	"context"
	"fmt"
	"mime/multipart"
	"strings"
	"time"

	"github.com/boilerplate/internal/config"
	"github.com/boilerplate/internal/modules/upload/domain"
	"github.com/boilerplate/internal/shared/app_errors"
	"github.com/boilerplate/pkg/logger"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

const maxFileSize = 5 * 1024 * 1024 // 5MB

type UploadService struct {
	log logger.Logger
	cld *cloudinary.Cloudinary
	cfg config.Cloudinary
}

func NewUploadService(log logger.Logger, cfg config.Cloudinary) *UploadService {
	cld, err := cloudinary.NewFromParams(
		cfg.CloudName,
		cfg.ApiKey,
		cfg.Secret,
	)
	if err != nil {
		panic(fmt.Sprintf("failed to init cloudinary: %v", err))
	}
	return &UploadService{log: log, cld: cld}
}

func (s *UploadService) UploadImage(ctx context.Context, file *multipart.FileHeader, uploadType domain.UploadType) (string, error) {
	if !uploadType.IsValid() {
		s.log.Error("invalid upload type", "type", uploadType.String())
		return "", app_errors.InvalidInput()
	}

	if file.Size > maxFileSize {
		s.log.Warn("file too large", "size", file.Size, "max", maxFileSize)
		return "", app_errors.ValidationError(fmt.Sprintf("file size exceeds %dMB limit", maxFileSize/1024/1024))
	}

	src, err := file.Open()
	if err != nil {
		s.log.Error("failed to open file", "filename", file.Filename, "error", err)
		return "", app_errors.InvalidInput().WithCause(err)
	}
	defer src.Close()

	// Validate mime type
	mimeType, err := s.validateMimeType(src)
	if err != nil {
		return "", err
	}
	if !strings.HasPrefix(mimeType, "image/") {
		return "", app_errors.ValidationError("file must be an image")
	}

	// Upload to Cloudinary
	publicID := fmt.Sprintf("%s/%d", uploadType.String(), time.Now().UnixMilli())
	resp, err := s.cld.Upload.Upload(ctx, src, uploader.UploadParams{
		PublicID: publicID,
		Folder:   "portfolio/" + uploadType.String(),
	})
	if err != nil {
		s.log.Error("cloudinary upload failed", "error", err)
		return "", app_errors.InternalError("failed to upload file").WithCause(err)
	}

	s.log.Info("file uploaded to cloudinary", "url", resp.SecureURL)
	return resp.SecureURL, nil
}

func (s *UploadService) validateMimeType(file multipart.File) (string, error) {
	header := make([]byte, 512)
	n, err := file.Read(header)
	if err != nil {
		return "", app_errors.InvalidInput().WithCause(err)
	}
	file.Seek(0, 0)
	return detectMimeType(header[:n]), nil
}

func detectMimeType(data []byte) string {
	if len(data) >= 4 {
		if data[0] == 0xFF && data[1] == 0xD8 && data[2] == 0xFF {
			return "image/jpeg"
		}
		if data[0] == 0x89 && data[1] == 0x50 && data[2] == 0x4E && data[3] == 0x47 {
			return "image/png"
		}
		if data[0] == 0x47 && data[1] == 0x49 && data[2] == 0x46 {
			return "image/gif"
		}
		if data[0] == 0x52 && data[1] == 0x49 && data[2] == 0x46 && data[3] == 0x46 {
			return "image/webp"
		}
	}
	return "application/octet-stream"
}
