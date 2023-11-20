package handler

import (
	mClient "RIP/internal/app/s3/minio"
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/minio/minio-go"
)

func (h *Handler) SaveImage(ctx context.Context, file multipart.File, header *multipart.FileHeader) (string, error) {
	objectName := uuid.New().String() + filepath.Ext(header.Filename)
	if _, err := h.Minio.PutObject(mClient.BucketName, objectName, file, header.Size, minio.PutObjectOptions{
		ContentType: header.Header.Get("Content-Type"),
	}); err != nil {
		return "", err
	}

	return fmt.Sprintf("http://%s/%s/%s", mClient.MinioHost, mClient.BucketName, objectName), nil
}

func (h *Handler) DeleteImage(objectName string) error {
	return h.Minio.RemoveObject(mClient.BucketName, objectName)
}
