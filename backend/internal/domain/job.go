package domain

import "time"

type Job struct {
	ID        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewJob(id, name string) *Job {
	return &Job{
		ID:        id,
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (j *Job) UpdateName(name string) {
	j.Name = name
	j.UpdatedAt = time.Now()
}
