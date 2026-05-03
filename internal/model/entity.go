package model

import (
	"time"
)

type Poll struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"type:varchar(255);not null" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	Slug        string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"slug"`
	CreatedAt   time.Time `json:"created_at"`
	Options     []Option  `gorm:"foreignKey:PollID" json:"options"`
}

type Option struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	PollID    uint   `json:"poll_id"`
	Value     string `gorm:"type:varchar(255);not null" json:"value"`
	VoteCount int    `gorm:"default:0" json:"vote_count"`
}

type Vote struct {
	ID         uint      `gorm:"primaryKey"`
	PollID     uint      `gorm:"not null"`
	OptionID   uint      `gorm:"not null"`
	Identifier string    `gorm:"type:varchar(255);not null"` // IP Address
	VotedAt    time.Time `gorm:"autoCreateTime"`
}
