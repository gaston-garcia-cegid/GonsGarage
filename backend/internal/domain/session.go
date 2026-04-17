package domain

import "time"

type Session struct {
	ID        string
	UserID    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewSession(id, userID string) *Session {
	return &Session{
		ID:        id,
		UserID:    userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (s *Session) UpdateUserID(userID string) {
	s.UserID = userID
	s.UpdatedAt = time.Now()
}
