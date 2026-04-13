package scheduler

import (
	"fmt"
	"time"

	"socialmanager/backend/database"
	"socialmanager/backend/models"
)

type Scheduler struct {
	ticker *time.Ticker
	quit   chan struct{}
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		quit: make(chan struct{}),
	}
}

func (s *Scheduler) Start() {
	s.ticker = time.NewTicker(1 * time.Minute)
	go func() {
		for {
			select {
			case <-s.ticker.C:
				s.runTasks()
			case <-s.quit:
				s.ticker.Stop()
				return
			}
		}
	}()
}

func (s *Scheduler) Stop() {
	close(s.quit)
}

func (s *Scheduler) runTasks() {
	// Simple mock task runner
	var tasks []models.Task
	database.DB.Where("enabled = ?", true).Find(&tasks)

	for _, t := range tasks {
		// Mock logic: randomly decide it ran
		now := time.Now()
		t.LastRun = &now
		database.DB.Save(&t)

		log := models.Log{
			Time:    now,
			Module:  "Automation",
			Action:  fmt.Sprintf("Task executed: %s", t.Name),
			Status:  "Success",
			Details: "Mock task ran safely.",
		}
		database.DB.Create(&log)
	}
}
