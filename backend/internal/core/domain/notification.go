package domain

import "time"

type Notification struct {
	ID        string
	Message   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewNotification(id, message string) *Notification {
	return &Notification{
		ID:        id,
		Message:   message,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (n *Notification) UpdateMessage(message string) {
	n.Message = message
	n.UpdatedAt = time.Now()
}
