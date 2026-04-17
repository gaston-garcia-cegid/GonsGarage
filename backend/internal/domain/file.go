package domain

import "time"

type File struct {
	ID        string
	Name      string
	Size      int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewFile(id, name string, size int64) *File {
	return &File{
		ID:        id,
		Name:      name,
		Size:      size,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (f *File) UpdateName(name string) {
	f.Name = name
	f.UpdatedAt = time.Now()
}

func (f *File) UpdateSize(size int64) {
	f.Size = size
	f.UpdatedAt = time.Now()
}
