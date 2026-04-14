package providers

import "context"

// SocialProvider abstraction for executing social media operations
type SocialProvider interface {
	Connect(ctx context.Context, config map[string]interface{}) error
	ReactToPost(ctx context.Context, postURL string, reactionType string) error
	CreatePost(ctx context.Context, content string, files []string) error
}
