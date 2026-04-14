package store

import (
	"socialmanager/backend/models"
	"sync"
	"time"
)

type Store struct {
	Mu          sync.RWMutex
	Tasks       []models.ReactionTask
	Executions  []models.TaskExecution
	Logs        []models.ActivityLog
	Settings    models.AppSettings
	TaskIDSeq   uint
	ExecIDSeq   uint
	LogIDSeq    uint
}

var DB *Store

func InitStore() {
	DB = &Store{
		Tasks:      make([]models.ReactionTask, 0),
		Executions: make([]models.TaskExecution, 0),
		Logs:       make([]models.ActivityLog, 0),
		Settings: models.AppSettings{
			AppName:              "Social Desktop Manager",
			Theme:                "light",
			DefaultExecutionMode: "mock",
			MaskCookieByDefault:  true,
			PersistToJson:        false,
			GraphqlLikeDocId:     "", // Default is empty so the user configures it
		},
		TaskIDSeq: 1,
		ExecIDSeq: 1,
		LogIDSeq:  1,
	}

	// Khởi tạo thư mục Data
	InitFileSystem()
	// Không khởi tạo dữ liệu ảo (Mock logs) nữa để UI sạch sẽ cho môi trường thật
}

func (s *Store) NextTaskID() uint {
	s.TaskIDSeq++
	return s.TaskIDSeq - 1
}

func (s *Store) NextExecID() uint {
	s.ExecIDSeq++
	return s.ExecIDSeq - 1
}

func (s *Store) AddLog(module, action, status, details string) {
	s.LogIDSeq++
	s.Logs = append([]models.ActivityLog{{
		ID:      s.LogIDSeq - 1,
		Time:    time.Now().Format(time.RFC3339),
		Module:  module,
		Action:  action,
		Status:  status,
		Details: details,
	}}, s.Logs...)
}
