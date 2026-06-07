package models

import (
	"time"
)

type Admin struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Username string `gorm:"type:varchar(100);uniqueIndex;not null" json:"username"`
	Password string `gorm:"type:varchar(255);not null" json:"-"`
}

type Poll struct {
	ID          uint         `gorm:"primaryKey" json:"id"`
	Title       string       `gorm:"type:varchar(255);not null" json:"title"`
	Description string       `gorm:"type:text" json:"description"`
	Slug        string       `gorm:"type:varchar(100);uniqueIndex;not null" json:"slug"`
	IsActive    bool         `gorm:"default:true" json:"is_active"`
	ExpiresAt   time.Time    `json:"expires_at"`
	CreatedAt   time.Time    `json:"created_at"`
	Options     []Option     `gorm:"foreignKey:PollID;constraint:OnDelete:CASCADE" json:"options"`
	Tokens      []VoterToken `gorm:"foreignKey:PollID;constraint:OnDelete:CASCADE" json:"-"`
}

type Option struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	PollID    uint   `json:"poll_id"`
	Value     string `gorm:"type:varchar(255);not null" json:"value"`
	VoteCount int    `gorm:"default:0" json:"vote_count"`
}

type VoterToken struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	PollID      uint   `gorm:"not null" json:"poll_id"`
	TokenString string `gorm:"type:varchar(50);uniqueIndex;not null" json:"token_string"`
	IsUsed      bool   `gorm:"default:false" json:"is_used"`
}

type Vote struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	PollID    uint      `gorm:"not null" json:"poll_id"`
	OptionID  uint      `gorm:"not null" json:"option_id"`
	TokenID   uint      `gorm:"not null;uniqueIndex" json:"token_id"`
	CreatedAt time.Time `json:"created_at"`
}
