package external

import "context"

// FileStorage defines the interface for the file storage service
type FileStorage interface {
	UploadFile(ctx context.Context, file *domain.File) error
	DownloadFile(ctx context.Context, id string) (*domain.File, error)
	DeleteFile(ctx context.Context, id string) error
	ListFiles(ctx context.Context, limit, offset int) ([]*domain.File, int64, error)
}
