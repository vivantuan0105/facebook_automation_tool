package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"socialmanager/backend/models"
	"time"
)

var DataDir = "Data"

// Khởi tạo hệ thống thư mục
func InitFileSystem() error {
	return os.MkdirAll(DataDir, 0755)
}

func GetAccountPath(uid string) string {
	return filepath.Join(DataDir, uid)
}

func GetAccountInfoPath(uid string) string {
	return filepath.Join(GetAccountPath(uid), "Info.txt")
}

func GetAllAccounts() ([]models.FacebookAccount, error) {
	accounts := make([]models.FacebookAccount, 0)

	entries, err := os.ReadDir(DataDir)
	if err != nil {
		if os.IsNotExist(err) {
			return accounts, nil
		}
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			uid := entry.Name()
			infoPath := GetAccountInfoPath(uid)
			
			data, err := os.ReadFile(infoPath)
			if err == nil {
				var acc models.FacebookAccount
				if err := json.Unmarshal(data, &acc); err == nil {
					accounts = append(accounts, acc)
				}
			}
		}
	}
	return accounts, nil
}

func SaveAccount(acc models.FacebookAccount) error {
	if acc.UID == "" {
		return fmt.Errorf("UID không được để trống")
	}

	if acc.CreatedAt == "" {
		acc.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	}

	dirPath := GetAccountPath(acc.UID)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return fmt.Errorf("không thể tạo thư mục cho UID %s: %v", acc.UID, err)
	}

	// Tạo các file rỗng theo thiết kế để dễ mở rộng sau này
	friendsPath := filepath.Join(dirPath, "Friends.txt")
	if _, err := os.Stat(friendsPath); os.IsNotExist(err) {
		os.WriteFile(friendsPath, []byte(""), 0644)
	}

	groupsPath := filepath.Join(dirPath, "Groups.txt")
	if _, err := os.Stat(groupsPath); os.IsNotExist(err) {
		os.WriteFile(groupsPath, []byte(""), 0644)
	}

	// Lưu thông tin account (dưới dạng JSON để dễ lấy lại object) vào file Info.txt
	infoPath := GetAccountInfoPath(acc.UID)
	data, err := json.MarshalIndent(acc, "", "  ")
	if err != nil {
		return fmt.Errorf("không thể parse dữ liệu account: %v", err)
	}

	return os.WriteFile(infoPath, data, 0644)
}

func DeleteAccount(uid string) error {
	dirPath := GetAccountPath(uid)
	return os.RemoveAll(dirPath)
}
