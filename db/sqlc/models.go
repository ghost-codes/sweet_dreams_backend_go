// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0

package db

import (
	"time"
)

type Admin struct {
	ID             int64     `json:"id"`
	Username       string    `json:"username"`
	FullName       string    `json:"full_name"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"hashed_password"`
	IsSuper        bool      `json:"is_super"`
	CreatedAt      time.Time `json:"created_at"`
}

type Approval struct {
	ID            int64     `json:"id"`
	RequestID     int64     `json:"request_id"`
	AssignedNurse int64     `json:"assigned_nurse"`
	ApprovedBy    int64     `json:"approved_by"`
	Status        string    `json:"status"`
	Notes         *string   `json:"notes"`
	UpdatedAt     time.Time `json:"updated_at"`
	CreatedAt     time.Time `json:"created_at"`
}

type Nurse struct {
	ID             int64     `json:"id"`
	FullName       string    `json:"full_name"`
	Email          string    `json:"email"`
	Contact        string    `json:"contact"`
	ProfilePicture *string   `json:"profile_picture"`
	CreatedAt      time.Time `json:"created_at"`
}

type Request struct {
	ID            int64       `json:"id"`
	UserID        *int64      `json:"user_id"`
	Type          *string     `json:"type"`
	PreferedNurse *int64      `json:"prefered_nurse"`
	StartDate     time.Time   `json:"start_date"`
	EndDate       time.Time   `json:"end_date"`
	Location      interface{} `json:"location"`
	CreatedAt     time.Time   `json:"created_at"`
}

type User struct {
	ID                int64      `json:"id"`
	Username          string     `json:"username"`
	FirstName         string     `json:"first_name"`
	LastName          string     `json:"last_name"`
	Email             string     `json:"email"`
	HashedPassword    *string    `json:"hashed_password"`
	AvatarUrl         *string    `json:"avatar_url"`
	Contact           *string    `json:"contact"`
	SecurityKey       string     `json:"security_key"`
	PasswordChangedAt time.Time  `json:"password_changed_at"`
	VerifiedAt        *time.Time `json:"verified_at"`
	CreatedAt         time.Time  `json:"created_at"`
	TwitterSocial     bool       `json:"twitter_social"`
	GoogleSocial      bool       `json:"google_social"`
	AppleSocial       bool       `json:"apple_social"`
}

type VerifyEmail struct {
	ID        int64     `json:"id"`
	Username  *string   `json:"username"`
	Email     string    `json:"email"`
	SecretKey string    `json:"secret_key"`
	IsUsed    bool      `json:"is_used"`
	ExpiredAt time.Time `json:"expired_at"`
	CreatedAt time.Time `json:"created_at"`
}
