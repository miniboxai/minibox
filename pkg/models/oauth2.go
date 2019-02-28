package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

type OAuth2Client struct {
	ID          string
	Secret      string
	RedirectUri string
	Name        string
	Summary     string
	UserID      sql.NullInt64
	User        User
	// EndpointID  int
	// Endpoint    OAuth2Endpoint `gorm:"foreignkey:EndpointID"`
	UserData  json.RawMessage
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

type OAuth2Authorize struct {
	Model
	Client      OAuth2Client
	Code        string
	ExpiresIn   int32
	Scope       string
	RedirectUri string
	State       string
	// CreatedAt           time.Time
	UserData            json.RawMessage
	CodeChallenge       string
	CodeChallengeMethod string
}

type OAuth2Token struct {
	ModelByInt
	Client   OAuth2Client `gorm:"foreignkey:ClientID"`
	ClientID string
	// Previous access data, for refresh token
	AccessData    *OAuth2Token
	AccessTokenID int
	Token         string
	Refresh       string
	ExpiresIn     int32
	Scope         string
	RedirectUri   string
	UserID        int
	User          User
	// CreatedAt   time.Time
	UserData json.RawMessage
}

type OAuth2Endpoint struct {
	Model
	AuthURL  string
	TokenURL string
}
