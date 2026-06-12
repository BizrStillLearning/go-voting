package models

import (
	"time"

	"gorm.io/gorm"
)

type Poll struct {
	gorm.Model
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Slug        string       `json:"slug" gorm:"uniqueIndex"`
	IsActive    bool         `json:"is_active" gorm:"default:true"`
	ExpiresAt   time.Time    `json:"expires_at"`
	Positions   []Position   `json:"positions"`
	Tokens      []VoterToken `json:"tokens"`
}

type Position struct {
	gorm.Model
	PollID  uint     `json:"poll_id"`
	Title   string   `json:"title"`
	Options []Option `json:"options"`
}

type Option struct {
	gorm.Model
	PositionID uint   `json:"position_id"`
	Value      string `json:"value"`
	PhotoURL   string `json:"photo_url"`
	Vision     string `json:"vision"`
	Mission    string `json:"mission"`
	VoteCount  int    `json:"vote_count" gorm:"default:0"`
}

type VoterToken struct {
	gorm.Model
	PollID      uint   `json:"poll_id"`
	TokenString string `json:"TokenString" gorm:"uniqueIndex"`
	IsUsed      bool   `json:"is_used" gorm:"default:false"`
}

type Admin struct {
	gorm.Model
	Username string `json:"username" gorm:"unique"`
	Password string `json:"password"`
}

type Vote struct {
	gorm.Model
	PollID     uint   `json:"poll_id"`
	PositionID uint   `json:"position_id"`
	OptionID   uint   `json:"option_id"`
	TokenID    uint   `json:"token_id"`
	Token      string `json:"token"`
}
