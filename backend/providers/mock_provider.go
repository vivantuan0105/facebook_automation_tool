package providers

import (
	"context"
	"fmt"
	"time"
)

type MockProvider struct{}

func NewMockProvider() *MockProvider {
	return &MockProvider{}
}

func (m *MockProvider) Connect(ctx context.Context, config map[string]interface{}) error {
	// Simulate connection flow via official supported flow abstraction
	time.Sleep(500 * time.Millisecond)
	return nil
}

func (m *MockProvider) ReactToPost(ctx context.Context, postURL string, reactionType string) error {
	// Mock the executor
	time.Sleep(1 * time.Second)
	fmt.Printf("[MockExecutor] Analyzed compliant endpoints. Sent %s reaction to %s\n", reactionType, postURL)
	return nil
}

func (m *MockProvider) CreatePost(ctx context.Context, content string, files []string) error {
	time.Sleep(1 * time.Second)
	fmt.Println("[MockExecutor] Created post.")
	return nil
}
