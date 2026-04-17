package models

// LoginResult là kết quả trả về sau khi đăng nhập bằng user/pass
type LoginResult struct {
	UID         string `json:"uid"`         // c_user UID của Facebook
	CookieFull  string `json:"cookieFull"`  // Chuỗi cookie đầy đủ (dùng thêm vào DB)
	Status      string `json:"status"`      // Live | Failed | Checkpoint | WrongPassword
	RawResponse string `json:"rawResponse"` // Snippet response để debug (chỉ khi lỗi)
}

// ImportResult là kết quả trả về sau khi import tài khoản hàng loạt
type ImportResult struct {
	TotalProcessed int      `json:"totalProcessed"`
	SuccessCount   int      `json:"successCount"`
	FailedCount    int      `json:"failedCount"`
	Errors         []string `json:"errors"`
}

type FacebookAccount struct {
	UID       string `json:"uid"`
	Name             string `json:"name"` // Tên gợi nhớ hoặc tên Facebook
	Cookie              string `json:"cookie"`
	Status              string `json:"status"` // Live, Die, Checkpoint, Unchecked
	Gender              string `json:"gender"`
	Birthday            string `json:"birthday"`
	BirthYear           string `json:"birthYear"`
	Location            string `json:"location"`
	Friends             string `json:"friends"`
	Posts               string `json:"posts"`
	Followers           string `json:"followers"`
	CreatedAt           string `json:"createdAt"`
}

type ReactionTask struct {
	ID            uint   `json:"id"`
	Cookie        string `json:"cookie"`
	CookieMasked  string `json:"cookieMasked"`
	TaskType      string `json:"taskType"`
	TargetMode    string `json:"targetMode"`
	PostID        string `json:"postId"`
	PostURL       string `json:"postUrl"`
	ReactionType  string `json:"reactionType"`
	Message       string   `json:"message"`
	PhotoPaths    []string `json:"photoPaths"`
	FeedbackID    string   `json:"feedbackId"`
	FeedbackReactionID string `json:"feedbackReactionId"`
	ActorID            string `json:"actorId"`
	FeedbackSource     string `json:"feedbackSource"`
	SessionID          string `json:"sessionId"`
	ClientMutationID   string `json:"clientMutationId"`
	Notes              string `json:"notes"`
	ScheduleAt         string `json:"scheduleAt"`
	QuantityLimit int    `json:"quantityLimit"`
	Status        string `json:"status"` // Pending, Running, Success, Failed, Paused, WaitingConfirmation
	CreatedAt     string `json:"createdAt"`
	UpdatedAt     string `json:"updatedAt"`
}

type TaskExecution struct {
	ID            uint   `json:"id"`
	TaskID        uint   `json:"taskId"`
	Mode          string `json:"mode"`
	Status        string `json:"status"`
	StartedAt     string `json:"startedAt"`
	FinishedAt    string `json:"finishedAt"`
	ErrorMessage  string `json:"errorMessage"`
	ResultSummary string `json:"resultSummary"`
}

type ActivityLog struct {
	ID      uint   `json:"id"`
	Time    string `json:"time"`
	Module  string `json:"module"`
	Action  string `json:"action"`
	Status  string `json:"status"`
	Details string `json:"details"`
}

type AppSettings struct {
	AppName              string `json:"appName"`
	Theme                string `json:"theme"`
	DefaultExecutionMode string `json:"defaultExecutionMode"`
	MaskCookieByDefault  bool   `json:"maskCookieByDefault"`
	PersistToJson        bool   `json:"persistToJson"`
	GraphqlLoginDocId    string `json:"graphqlLoginDocId"`
	GraphqlLikeDocId     string `json:"graphqlLikeDocId"`
	GraphqlCommentDocId  string `json:"graphqlCommentDocId"`
	GraphqlPostDocId     string `json:"graphqlPostDocId"`
	GraphqlScanPostDocId string `json:"graphqlScanPostDocId"`
	GraphqlProfileDocId  string `json:"graphqlProfileDocId"`
	GraphqlFriendsDocId  string `json:"graphqlFriendsDocId"`
}

type DashboardStats struct {
	TotalTasks int `json:"totalTasks"`
	Pending    int `json:"pending"`
	Running    int `json:"running"`
	Success    int `json:"success"`
	Failed     int `json:"failed"`
}
