package app

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"regexp"
	"socialmanager/backend/models"
	"socialmanager/backend/providers"
	"socialmanager/backend/store"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) SelectPhotoDialog() (string, error) {
	options := runtime.OpenDialogOptions{
		Title: "Chọn file ảnh",
		Filters: []runtime.FileFilter{
			{DisplayName: "Image Files (*.jpg, *.jpeg, *.png)", Pattern: "*.jpg;*.jpeg;*.png"},
		},
	}
	return runtime.OpenFileDialog(a.ctx, options)
}

func (a *App) GetOverview() models.DashboardStats {
	store.DB.Mu.RLock()
	defer store.DB.Mu.RUnlock()

	var pending, running, success, failed int
	for _, t := range store.DB.Tasks {
		switch t.Status {
		case "Pending":
			pending++
		case "Running", "WaitingConfirmation":
			running++
		case "Success":
			success++
		case "Failed":
			failed++
		}
	}

	return models.DashboardStats{
		TotalTasks: len(store.DB.Tasks),
		Pending:    pending,
		Running:    running,
		Success:    success,
		Failed:     failed,
	}
}

func (a *App) GetReactionTasks() []models.ReactionTask {
	store.DB.Mu.RLock()
	defer store.DB.Mu.RUnlock()
	return store.DB.Tasks
}

func (a *App) GetExecutions() []models.TaskExecution {
	store.DB.Mu.RLock()
	defer store.DB.Mu.RUnlock()
	return store.DB.Executions
}

func (a *App) GetLogs() []models.ActivityLog {
	store.DB.Mu.RLock()
	defer store.DB.Mu.RUnlock()
	return store.DB.Logs
}

func (a *App) GetSettings() models.AppSettings {
	store.DB.Mu.RLock()
	defer store.DB.Mu.RUnlock()
	return store.DB.Settings
}

func (a *App) ValidateReactionTask(input models.ReactionTask) string {
	if strings.TrimSpace(input.Cookie) == "" {
		return "Vui lòng nhập Cookie"
	}
	if input.TaskType != "Đăng bài viết" {
		if input.TargetMode == "post_url" && strings.TrimSpace(input.PostURL) == "" {
			return "Vui lòng nhập URL của post"
		}
		if input.TargetMode == "post_id" && strings.TrimSpace(input.PostID) == "" {
			return "Vui lòng nhập ID của post"
		}
	}
	if input.TaskType == "Like bài viết" && input.ReactionType == "" {
		return "Vui lòng chọn loại reaction"
	}
	if input.TaskType == "Comment bài viết" && strings.TrimSpace(input.Message) == "" {
		return "Vui lòng nhập nội dung bình luận"
	}
	if input.TaskType == "Đăng bài viết" && strings.TrimSpace(input.Message) == "" && len(input.PhotoPaths) == 0 {
		return "Vui lòng nhập nội dung bài viết hoặc đính kèm ít nhất 1 ảnh"
	}
	return ""
}

func (a *App) CreateReactionTask(input models.ReactionTask) models.ReactionTask {
	store.DB.Mu.Lock()
	defer store.DB.Mu.Unlock()

	input.ID = store.DB.NextTaskID()
	input.Status = "Pending"
	input.CreatedAt = time.Now().Format(time.RFC3339)
	input.UpdatedAt = time.Now().Format(time.RFC3339)

	// mask cookie
	cookieLen := len(input.Cookie)
	if cookieLen > 20 {
		input.CookieMasked = input.Cookie[:10] + "..." + input.Cookie[cookieLen-10:]
	} else {
		input.CookieMasked = "***"
	}

	store.DB.Tasks = append([]models.ReactionTask{input}, store.DB.Tasks...)
	store.DB.AddLog("Task Center", "CreateTask", "Success", fmt.Sprintf("Đã thêm task '%v' vào queue.", input.TaskType))

	return input
}

func (a *App) RunReactionTaskNow(taskID uint, mode string) error {
	var task *models.ReactionTask
	store.DB.Mu.Lock()
	for i := range store.DB.Tasks {
		if store.DB.Tasks[i].ID == taskID {
			task = &store.DB.Tasks[i]
			break
		}
	}
	if task == nil {
		store.DB.Mu.Unlock()
		return errors.New("không tìm thấy task")
	}

	if task.Status == "Running" || task.Status == "WaitingConfirmation" {
		store.DB.Mu.Unlock()
		return errors.New("task đang chạy")
	}

	task.Status = "Running"
	task.UpdatedAt = time.Now().Format(time.RFC3339)

	exec := models.TaskExecution{
		ID:        store.DB.NextExecID(),
		TaskID:    taskID,
		Mode:      mode,
		Status:    "Running",
		StartedAt: time.Now().Format(time.RFC3339),
	}
	store.DB.Executions = append([]models.TaskExecution{exec}, store.DB.Executions...)
	store.DB.AddLog("Queue", "RunTask", "Running", fmt.Sprintf("Bắt đầu chạy task %d ở mode %v", taskID, mode))
	store.DB.Mu.Unlock()

	// async execution
	go func(id uint, eID uint, tMode string, targetURL string, targetID string) {
		time.Sleep(1500 * time.Millisecond) // Simulate some work
		store.DB.Mu.Lock()
		defer store.DB.Mu.Unlock()

		// locate structs
		var t *models.ReactionTask
		var e *models.TaskExecution
		for i := range store.DB.Tasks {
			if store.DB.Tasks[i].ID == id {
				t = &store.DB.Tasks[i]
				break
			}
		}
		for i := range store.DB.Executions {
			if store.DB.Executions[i].ID == eID {
				e = &store.DB.Executions[i]
				break
			}
		}

		if t == nil || e == nil {
			return
		}

			// Lấy docId từ Cài Đặt (Hệ thống) thay vì từ Account
			var docId string

			if t.TaskType == "Comment bài viết" {
				docId = store.DB.Settings.GraphqlCommentDocId
			} else if t.TaskType == "Đăng bài viết" {
				docId = store.DB.Settings.GraphqlPostDocId
			} else {
				docId = store.DB.Settings.GraphqlLikeDocId
			}

			// Gọi FacebookProvider thật
			fbProvider := providers.NewFacebookProvider()
			store.DB.AddLog("Queue", "Execution", "Info", "Đang phân tích trang và chuẩn bị payload...")
			
			var resultLog string
			var err error

			if t.TaskType == "Comment bài viết" {
				resultLog, err = fbProvider.CommentToPost(t.Cookie, targetURL, t.Message, docId)
			} else if t.TaskType == "Đăng bài viết" {
				photoIDs := []string{}
				for _, path := range t.PhotoPaths {
					if path != "" {
						store.DB.AddLog("FacebookProvider", "UploadPhoto", "Info", fmt.Sprintf("Đang tải ảnh từ %s lên Facebook...", path))
						photoID, upErr := fbProvider.UploadPhoto(t.Cookie, path)
						if upErr != nil {
							err = fmt.Errorf("Lỗi tải ảnh %s: %v", path, upErr)
							break
						}
						photoIDs = append(photoIDs, photoID)
					}
				}
				
				if err == nil {
					resultLog, err = fbProvider.PostToFacebook(t.Cookie, t.Message, photoIDs, docId)
				}
			} else {
				resultLog, err = fbProvider.ReactToPost(t.Cookie, targetURL, t.ReactionType, docId)
			}
			
			if err != nil {
				t.Status = "Failed"
				t.UpdatedAt = time.Now().Format(time.RFC3339)
				e.Status = "Failed"
				e.FinishedAt = time.Now().Format(time.RFC3339)
				e.ErrorMessage = err.Error()

				store.DB.AddLog("FacebookProvider", "Execute", "Error", fmt.Sprintf("Task %d thất bại: %v", id, err.Error()))
			} else {
				t.Status = "Success"
				t.UpdatedAt = time.Now().Format(time.RFC3339)
				e.Status = "Success"
				e.FinishedAt = time.Now().Format(time.RFC3339)
				e.ResultSummary = resultLog

				store.DB.AddLog("FacebookProvider", "Execute", "Success", fmt.Sprintf("Task %d chạy thật thành công", id))
			}

	}(taskID, exec.ID, mode, task.PostURL, task.PostID)

	return nil
}

func (a *App) PauseReactionTask(taskID uint) error {
	store.DB.Mu.Lock()
	defer store.DB.Mu.Unlock()
	for i := range store.DB.Tasks {
		if store.DB.Tasks[i].ID == taskID {
			store.DB.Tasks[i].Status = "Paused"
			store.DB.Tasks[i].UpdatedAt = time.Now().Format(time.RFC3339)
			store.DB.AddLog("Queue", "PauseTask", "Info", fmt.Sprintf("Đã tạm dừng task %d", taskID))
			return nil
		}
	}
	return errors.New("không tìm thấy task")
}

func (a *App) RetryReactionTask(taskID uint) error {
	store.DB.Mu.Lock()
	defer store.DB.Mu.Unlock()
	for i := range store.DB.Tasks {
		if store.DB.Tasks[i].ID == taskID {
			store.DB.Tasks[i].Status = "Pending"
			store.DB.Tasks[i].UpdatedAt = time.Now().Format(time.RFC3339)
			store.DB.AddLog("Queue", "RetryTask", "Info", fmt.Sprintf("Đưa task %d về Pending", taskID))
			return nil
		}
	}
	return errors.New("không tìm thấy task")
}

func (a *App) RemoveReactionTask(taskID uint) error {
	store.DB.Mu.Lock()
	defer store.DB.Mu.Unlock()
	idx := -1
	for i := range store.DB.Tasks {
		if store.DB.Tasks[i].ID == taskID {
			idx = i
			break
		}
	}
	if idx != -1 {
		store.DB.Tasks = append(store.DB.Tasks[:idx], store.DB.Tasks[idx+1:]...)
		store.DB.AddLog("Queue", "RemoveTask", "Info", fmt.Sprintf("Xóa bỏ task %d", taskID))
		return nil
	}
	return errors.New("không tìm thấy task")
}

func (a *App) MarkManualCompleted(taskID uint) error {
	store.DB.Mu.Lock()
	defer store.DB.Mu.Unlock()
	for i := range store.DB.Tasks {
		if store.DB.Tasks[i].ID == taskID {
			if store.DB.Tasks[i].Status == "WaitingConfirmation" {
				store.DB.Tasks[i].Status = "Success"
				store.DB.Tasks[i].UpdatedAt = time.Now().Format(time.RFC3339)

				// find execution
				for j := range store.DB.Executions {
					if store.DB.Executions[j].TaskID == taskID && store.DB.Executions[j].Status == "WaitingConfirmation" {
						store.DB.Executions[j].Status = "Success"
						store.DB.Executions[j].FinishedAt = time.Now().Format(time.RFC3339)
						store.DB.Executions[j].ResultSummary = "Manually confirmed by user"
					}
				}

				store.DB.AddLog("ManualAssist", "ConfirmTask", "Success", fmt.Sprintf("Task %d đã được xác nhận thủ công bởi người dùng", taskID))
				return nil
			}
		}
	}
	return errors.New("Không thể xác nhận task này")
}

func (a *App) UpdateSettings(s models.AppSettings) error {
	store.DB.Mu.Lock()
	defer store.DB.Mu.Unlock()
	store.DB.Settings = s
	
	err := store.DB.SaveSettings()
	if err != nil {
		store.DB.AddLog("Settings", "Update", "Error", fmt.Sprintf("Lỗi khi lưu cấu hình: %v", err))
		return err
	}
	
	store.DB.AddLog("Settings", "Update", "Success", "Cấu hình được lưu thành công vào file settings.json.")
	return nil
}

// ---- QUẢN LÝ TÀI KHOẢN (ACCOUNTS) ----

func (a *App) GetAllAccounts() ([]models.FacebookAccount, error) {
	return store.GetAllAccounts()
}

func (a *App) DeleteAccount(uid string) error {
	err := store.DeleteAccount(uid)
	if err == nil {
		store.DB.AddLog("Accounts", "Delete", "Success", fmt.Sprintf("Đã xóa tài khoản UID %s", uid))
	}
	return err
}

func (a *App) AddAccount(name string, cookie string) (models.FacebookAccount, error) {
	// Lấy UID từ Cookie (c_user=...)
	re := regexp.MustCompile(`c_user=(\d+)`)
	matches := re.FindStringSubmatch(cookie)
	uid := ""
	if len(matches) > 1 {
		uid = matches[1]
	}

	if uid == "" {
		return models.FacebookAccount{}, errors.New("Cookie không hợp lệ. Không tìm thấy c_user (UID) trong Cookie")
	}

	acc := models.FacebookAccount{
		UID:                 uid,
		Name:                name,
		Cookie:              cookie,
		Status:              "Live", // Mặc định khi vừa thêm
	}

	err := store.SaveAccount(acc)
	if err == nil {
		store.DB.AddLog("Accounts", "Add", "Success", fmt.Sprintf("Đã thêm tài khoản %s (UID: %s)", name, uid))
	}
	return acc, err
}

