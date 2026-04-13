package models

import "time"

type Post struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Summary   string    `json:"summary"`
	Status    string    `json:"status"` // Draft, Scheduled, Published
	Notes     string    `json:"notes"`
	PostAt    time.Time `json:"post_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Task struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Type        string     `json:"type"` // Reminder, Sync, Check
	Schedule    string     `json:"schedule"` // e.g. "0 9 * * *"
	Enabled     bool       `json:"enabled"`
	LastRun     *time.Time `json:"last_run"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type Log struct {
	ID      uint      `json:"id" gorm:"primaryKey"`
	Time    time.Time `json:"time"`
	Module  string    `json:"module"`
	Action  string    `json:"action"`
	Status  string    `json:"status"` // Success, Warning, Error, Info
	Details string    `json:"details"`
}

type Settings struct {
	ID               uint   `json:"id" gorm:"primaryKey"`
	Theme            string `json:"theme"`
	AppName          string `json:"app_name"`
	DatabasePath     string `json:"database_path"`
	AutoStart        bool   `json:"auto_start"`
	LogLevel         string `json:"log_level"`
	Timezone         string `json:"timezone"`
	ConnectionStatus string `json:"connection_status"`
	ConnectedAccount string `json:"connected_account"`
}

// DashboardStats for stats cards
type DashboardStats struct {
	TotalPosts       int64 `json:"total_posts"`
	TotalTasks       int64 `json:"total_tasks"`
	SuccessfulRuns   int64 `json:"successful_runs"`
	RecentErrors     int64 `json:"recent_errors"`
	ConnectionStatus string `json:"connection_status"`
	ConnectedAccount string `json:"connected_account"`
}
