package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	Username     string         `gorm:"size:100;unique" json:"username"`
	PasswordHash string         `gorm:"size:255;not null" json:"-"`
	Email        string         `gorm:"size:255;unique" json:"email"`
	Nickname     string         `gorm:"size:100" json:"nickname"`
	Role         string         `gorm:"size:20;default:user" json:"role"` // admin, user
	Avatar       string         `gorm:"size:500" json:"avatar"`
	Status       string         `gorm:"size:20;default:active" json:"status"` // active, inactive
}

type WhatsAppAccount struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	UserID       uint           `gorm:"not null;index" json:"user_id"`
	Phone        string         `gorm:"size:20;unique" json:"phone"`
	Nickname     string         `gorm:"size:100" json:"nickname"`
	SessionData  string         `gorm:"type:text" json:"session_data"` // Encrypted session
	DeviceName   string         `gorm:"size:100" json:"device_name"`
	Status       string         `gorm:"size:20;default:offline" json:"status"` // online, offline, connecting, disconnected
	LastSeen     *time.Time     `json:"last_seen"`
}

type Customer struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	AccountID  uint           `gorm:"not null;index" json:"account_id"`
	Phone      string         `gorm:"size:20;not null;index" json:"phone"`
	Name       string         `gorm:"size:100" json:"name"`
	Avatar     string         `gorm:"size:500" json:"avatar"`
	Country    string         `gorm:"size:50" json:"country"`
	Tags       string         `gorm:"type:json" json:"tags"` // JSON array
	Notes      string         `gorm:"type:text" json:"notes"`
	Source     string         `gorm:"size:50" json:"source"` // whatsapp, import, manual
	LastMsgAt  *time.Time     `json:"last_msg_at"`
}

type Message struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	AccountID   uint           `gorm:"not null;index" json:"account_id"`
	CustomerID  uint           `gorm:"not null;index" json:"customer_id"`
	Direction   string         `gorm:"size:10;not null" json:"direction"` // inbound, outbound
	Content     string         `gorm:"type:text" json:"content"`
	MessageType string         `gorm:"size:20;default:text" json:"message_type"` // text, image, video, audio, document
	MediaURL    string         `gorm:"size:500" json:"media_url"`
	MediaID     string         `gorm:"size:100" json:"media_id"`
	Status      string         `gorm:"size:20;default:sent" json:"status"` // sent, delivered, read, failed
	SentAt      *time.Time    `json:"sent_at"`
}

type MessageTemplate struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time     `json:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	UserID     uint           `gorm:"not null;index" json:"user_id"`
	Name       string         `gorm:"size:100;not null" json:"name"`
	Content    string         `gorm:"type:text;not null" json:"content"`
	Variables  string         `gorm:"type:json" json:"variables"` // JSON array of variable names
	Category   string         `gorm:"size:50" json:"category"`   // greeting, inquiry, followup, etc
	Status     string         `gorm:"size:20;default:active" json:"status"` // active, inactive
}

type ScheduledTask struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	UserID      uint           `gorm:"not null;index" json:"user_id"`
	AccountID   uint           `gorm:"not null;index" json:"account_id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	TemplateID  uint           `json:"template_id"`
	CustomerIDs string         `gorm:"type:json" json:"customer_ids"` // JSON array of customer IDs
	Message     string         `gorm:"type:text" json:"message"`      // Direct message content
	ScheduledAt *time.Time     `json:"scheduled_at"`
	RepeatType  string         `gorm:"size:20" json:"repeat_type"` // once, daily, weekly, monthly
	RepeatEnd   *time.Time    `json:"repeat_end"`
	Status      string         `gorm:"size:20;default:pending" json:"status"` // pending, running, completed, failed, cancelled
	LastRunAt   *time.Time    `json:"last_run_at"`
	NextRunAt   *time.Time    `json:"next_run_at"`
	RunCount    int           `gorm:"default:0" json:"run_count"`
	ErrorMsg    string        `gorm:"type:text" json:"error_msg"`
}

type AutoReply struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	UserID     uint           `gorm:"not null;index" json:"user_id"`
	AccountID  uint           `gorm:"not null;index" json:"account_id"`
	Keyword    string         `gorm:"size:100;not null" json:"keyword"`
	MatchType  string         `gorm:"size:20;default:exact" json:"match_type"` // exact, contains, regex
	Response   string         `gorm:"type:text;not null" json:"response"`
	Priority   int            `gorm:"default:0" json:"priority"`
	Status     string         `gorm:"size:20;default:active" json:"status"` // active, inactive
}

type Tag struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	UserID    uint           `gorm:"not null;index" json:"user_id"`
	Name      string         `gorm:"size:50;not null" json:"name"`
	Color     string         `gorm:"size:20;default:#2563EB" json:"color"`
	Count     int            `gorm:"default:0" json:"count"`
}
