package flow

import "time"

type CredentialsProvider interface {
	Username() string
	Password() string
	TwoFactorCode() string
}

type TokenStorage interface {
	Token() string
	IsValid() bool
	SetToken(token string)
}

type MemoryTokenStorage struct {
	token      string
	validUntil time.Time
}

func (t *MemoryTokenStorage) Token() string {
	return t.token
}

func (t *MemoryTokenStorage) IsValid() bool {
	return t.token != "" && t.validUntil.After(time.Now())
}

func (t *MemoryTokenStorage) SetToken(token string) {
	t.token = token
	t.validUntil = time.Now().Add(2 * time.Hour)
}
