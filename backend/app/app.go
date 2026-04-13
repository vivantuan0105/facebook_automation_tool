package app

import (
	"context"
	"errors"

	"socialmanager/backend/database"
	"socialmanager/backend/models"
	"socialmanager/backend/scheduler"
)

type App struct {
	ctx   context.Context
	sched *scheduler.Scheduler
}

func NewApp() *App {
	return &App{
		sched: scheduler.NewScheduler(),
	}
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	a.sched.Start()
}

// Stats
func (a *App) GetDashboardStats() models.DashboardStats {
	var totalPosts, totalTasks, success, errorsCount int64
	database.DB.Model(&models.Post{}).Count(&totalPosts)
	database.DB.Model(&models.Task{}).Count(&totalTasks)
	database.DB.Model(&models.Log{}).Where("status = ?", "Success").Count(&success)
	database.DB.Model(&models.Log{}).Where("status = ?", "Error").Count(&errorsCount)

	var settings models.Settings
	database.DB.First(&settings)

	return models.DashboardStats{
		TotalPosts:       totalPosts,
		TotalTasks:       totalTasks,
		SuccessfulRuns:   success,
		RecentErrors:     errorsCount,
		ConnectionStatus: settings.ConnectionStatus,
		ConnectedAccount: settings.ConnectedAccount,
	}
}

// --- POSTS ---
func (a *App) GetPosts() []models.Post {
	var posts []models.Post
	database.DB.Order("created_at desc").Find(&posts)
	return posts
}

func (a *App) CreatePost(post models.Post) (models.Post, error) {
	result := database.DB.Create(&post)
	return post, result.Error
}

func (a *App) UpdatePost(post models.Post) (models.Post, error) {
	result := database.DB.Save(&post)
	return post, result.Error
}

func (a *App) DeletePost(id uint) error {
	return database.DB.Delete(&models.Post{}, id).Error
}

// --- TASKS ---
func (a *App) GetTasks() []models.Task {
	var tasks []models.Task
	database.DB.Order("created_at desc").Find(&tasks)
	return tasks
}

func (a *App) CreateTask(task models.Task) (models.Task, error) {
	result := database.DB.Create(&task)
	return task, result.Error
}

func (a *App) UpdateTask(task models.Task) (models.Task, error) {
	result := database.DB.Save(&task)
	return task, result.Error
}

func (a *App) DeleteTask(id uint) error {
	return database.DB.Delete(&models.Task{}, id).Error
}

func (a *App) RunTaskOnce(id uint) error {
	var task models.Task
	if err := database.DB.First(&task, id).Error; err != nil {
		return err
	}
	
	database.DB.Create(&models.Log{
		Time:    database.DB.NowFunc(),
		Module:  "Automation",
		Action:  "Manual Run: " + task.Name,
		Status:  "Success",
		Details: "Triggered manually by user.",
	})
	return nil
}

// --- LOGS ---
func (a *App) GetLogs() []models.Log {
	var logs []models.Log
	database.DB.Order("time desc").Limit(100).Find(&logs)
	return logs
}

// --- SETTINGS ---
func (a *App) GetSettings() (models.Settings, error) {
	var settings models.Settings
	err := database.DB.First(&settings).Error
	if err != nil && settings.ID == 0 {
	    return settings, errors.New("settings not found")
	}
	return settings, nil
}

func (a *App) UpdateSettings(s models.Settings) error {
	return database.DB.Save(&s).Error
}

func (a *App) DisconnectAccount() error {
	return database.DB.Model(&models.Settings{}).Where("1=1").Updates(map[string]interface{}{
		"connection_status": "disconnected",
		"connected_account": "",
	}).Error
}

func (a *App) ConnectAccount(token string) error {
	// Mock connection
	if token == "" {
		return errors.New("Invalid token")
	}
	return database.DB.Model(&models.Settings{}).Where("1=1").Updates(map[string]interface{}{
		"connection_status": "connected",
		"connected_account": "new_page_demo",
	}).Error
}
