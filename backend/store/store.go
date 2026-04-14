package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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
	
	DB.LoadSettings()
}

func (s *Store) SaveSettings() error {
	data, err := json.MarshalIndent(s.Settings, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(DataDir, "settings.json"), data, 0644)
}

func (s *Store) LoadSettings() {
	data, err := os.ReadFile(filepath.Join(DataDir, "settings.json"))
	if err == nil {
		var settings models.AppSettings
		if err := json.Unmarshal(data, &settings); err == nil {
			s.Settings = settings
		}
	}
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
	
	logEntry := fmt.Sprintf("[%s] %s | %s | %s | %s\n", time.Now().Format(time.RFC3339), module, action, status, details)
	f, _ := os.OpenFile("debug_logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if f != nil {
		f.WriteString(logEntry)
		f.Close()
	}

	s.Logs = append([]models.ActivityLog{{
		ID:      s.LogIDSeq - 1,
		Time:    time.Now().Format(time.RFC3339),
		Module:  module,
		Action:  action,
		Status:  status,
		Details: details,
	}}, s.Logs...)
}
