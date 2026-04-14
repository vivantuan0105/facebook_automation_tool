package providers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"bytes"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

type FacebookProvider struct {
	client *http.Client
}

type FBUploadResponse struct {
	PhotoID string `json:"fbid"`
}

func NewFacebookProvider() *FacebookProvider {
	return &FacebookProvider{
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

func (f *FacebookProvider) ReactToPost(cookie string, postURL string, reactionType string, docId string) (string, error) {
	if docId == "" {
		return "", errors.New("CHƯA CẤU HÌNH DOC_ID LIKE! Vui lòng vào Cài đặt (Settings) và điền dãy số GraphQL React Doc_ID lấy từ trình duyệt của bạn.")
	}

	// 1. Phân tích Link và cào thông số
	req, err := http.NewRequest("GET", postURL, nil)
	if err != nil {
		return "", fmt.Errorf("Lỗi tạo request GET: %v", err)
	}
	req.Header.Set("Cookie", cookie)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml")
	req.Header.Set("Accept-Language", "vi-VN,vi;q=0.9,en-US;q=0.8,en;q=0.7")

	resp, err := f.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Lỗi kết nối Facebook: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	htmlText := string(bodyBytes)

	// Lấy fb_dtsg
	fbDtsg := ""
	dtsgRegex1 := regexp.MustCompile(`\["DTSGInitialData",\[\],{"token":"([^"]+)"}`)
	if m := dtsgRegex1.FindStringSubmatch(htmlText); len(m) > 1 {
		fbDtsg = m[1]
	} else {
		dtsgRegex2 := regexp.MustCompile(`"DTSGInitialData",.*?"token":"([^"]+)"`)
		if m := dtsgRegex2.FindStringSubmatch(htmlText); len(m) > 1 {
			fbDtsg = m[1]
		}
	}

	if fbDtsg == "" {
		return "", errors.New("Không thể cào được token fb_dtsg (Cookie có thể đã chết hoặc IP bị giới hạn)")
	}

	// Lấy jazoest
	jazoest := ""
	jazoestRegex := regexp.MustCompile(`name="jazoest" value="(\d+)"`)
	if m := jazoestRegex.FindStringSubmatch(htmlText); len(m) > 1 {
		jazoest = m[1]
	} else {
		jazoestRegex = regexp.MustCompile(`"jazoest":"(\d+)"`)
		if m := jazoestRegex.FindStringSubmatch(htmlText); len(m) > 1 {
			jazoest = m[1]
		}
	}

	// Lấy lsd
	lsd := ""
	lsdRegex := regexp.MustCompile(`"LSD",\[\],{"token":"([^"]+)"}`)
	if m := lsdRegex.FindStringSubmatch(htmlText); len(m) > 1 {
		lsd = m[1]
	} else {
		lsdRegex = regexp.MustCompile(`name="lsd" value="([^"]+)"`)
		if m := lsdRegex.FindStringSubmatch(htmlText); len(m) > 1 {
			lsd = m[1]
		}
	}

	// Lấy Actor ID
	actorId := "0"
	actorRegex := regexp.MustCompile(`"USER_ID":"(\d+)"`)
	if m := actorRegex.FindStringSubmatch(htmlText); len(m) > 1 {
		actorId = m[1]
	}

	// Cố gắng tìm trực tiếp chuỗi Base64 (Thường bắt đầu bằng ZmVl)
	feedbackIDBase64 := ""
	encodedRegex1 := regexp.MustCompile(`"feedback":{"id":"(ZmV[^"]+)"}`)
	if m := encodedRegex1.FindStringSubmatch(htmlText); len(m) > 1 {
		feedbackIDBase64 = m[1]
	} else {
		encodedRegex2 := regexp.MustCompile(`"feedback_id":"(ZmV[^"]+)"`)
		if m := encodedRegex2.FindStringSubmatch(htmlText); len(m) > 1 {
			feedbackIDBase64 = m[1]
		}
	}

	// Lấy Numeric ID của bài viết nếuchưa có chuỗi mã hóa
	if feedbackIDBase64 == "" {
		numericID := ""
		idRegex := regexp.MustCompile(`"top_level_post_id":"(\d+)"`)
		if m := idRegex.FindStringSubmatch(htmlText); len(m) > 1 {
			numericID = m[1]
		} else {
			idRegex = regexp.MustCompile(`"ent_id":"(\d+)"`)
			if m := idRegex.FindStringSubmatch(htmlText); len(m) > 1 {
				numericID = m[1]
			} else {
				// Cào tham số từ URL nhưng cẩn thận với pfbid
				parsedUrl, err := url.Parse(postURL)
				if err == nil {
					fallbackId := parsedUrl.Query().Get("story_fbid")
					if fallbackId != "" && !strings.HasPrefix(fallbackId, "pfbid") {
						numericID = fallbackId
					}
				}
			}
		}

		if numericID == "" {
			return "", errors.New("Không bóc tách được Feedback ID hoặc Numeric ID Bài viết từ mã nguồn html. Có thể Cookie đã hết hạn hoặc định dạng page bị đổi.")
		}

		// 2. Chuyển đổi ID sang định dạng Base64
		feedbackRaw := "feedback:" + numericID
		feedbackIDBase64 = base64.StdEncoding.EncodeToString([]byte(feedbackRaw))
	}

	// 3. Mapping Reaction ID
	reactionID := "1635855486666999" // ID mặc định cho LIKE (Gợi ý của hệ thống)

	// doc_id được lấy từ tham số truyền vào, không hardcode nữa

	// 4. Bắn Request GraphQL
	variables := fmt.Sprintf(`{"input":{"feedback_id":"%s","feedback_reaction_id":"%s","feedback_source":"OBJECT","is_tracking_encrypted":true,"tracking":[],"session_id":"%s","actor_id":"%s","client_mutation_id":"1"},"useDefaultActor":false,"scale":1.5}`,
		feedbackIDBase64, reactionID, generatePseudoUUID(), actorId)

	data := url.Values{}
	data.Set("fb_dtsg", fbDtsg)
	if jazoest != "" {
		data.Set("jazoest", jazoest)
	}
	if lsd != "" {
		data.Set("lsd", lsd)
	}
	
	// Facebook strict payload parameters
	data.Set("__user", actorId)
	data.Set("__a", "1")
	data.Set("__req", "1")
	data.Set("__comet_req", "15")

	data.Set("variables", variables)
	data.Set("doc_id", docId)

	postReq, _ := http.NewRequest("POST", "https://www.facebook.com/api/graphql/", strings.NewReader(data.Encode()))
	postReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	postReq.Header.Set("Cookie", cookie)
	postReq.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")

	res, err := f.client.Do(postReq)
	if err != nil {
		return "", fmt.Errorf("Network lỗi khi bắn POST: %v", err)
	}
	defer res.Body.Close()

	fbBodyBytes, _ := io.ReadAll(res.Body)
	fbBody := string(fbBodyBytes)
	
	// Lưu lại response log để AI đọc
	os.WriteFile("debug_facebook_response.json", fbBodyBytes, 0644)

	if res.StatusCode != 200 {
		return "", fmt.Errorf("GraphQL HTTP %d: %s", res.StatusCode, fbBody)
	}

	// Kiểm tra xem có bằng chứng thành công trong data không. Nếu có thì ưu tiên trả về thành công dù GraphQL có kèm theo warning.
	if strings.Contains(fbBody, `"viewer_feedback_reaction_info"`) {
		return fmt.Sprintf("Thành công! ID mã hóa: %s. Trạng thái HTTP: 200", feedbackIDBase64), nil
	}

	// Xử lý lỗi tinh vi hơn
	if strings.Contains(fbBody, "was not found") && strings.Contains(fbBody, "The GraphQL document") {
		return "", errors.New("GraphQL Doc_ID Like đã hết hạn (Not Found). Vui lòng cập nhật Doc_ID mới trong Cài Đặt!")
	}
	
	if strings.Contains(fbBody, `"severity":"CRITICAL"`) {
		return "", fmt.Errorf("GraphQL thực thi lỗi nghiêm trọng: %s", strings.TrimPrefix(fbBody, "for (;;);"))
	}

	if strings.Contains(fbBody, `"errors":[{`) {
		return "", fmt.Errorf("GraphQL trả về lỗi nhưng không có dữ liệu thành công: %s", strings.TrimPrefix(fbBody, "for (;;);"))
	}

	return fmt.Sprintf("Thành công (Không có errors nhưng thiếu viewer_feedback)! ID: %s. Trạng thái HTTP: 200", feedbackIDBase64), nil
}

func (f *FacebookProvider) CommentToPost(cookie string, postURL string, message string, docId string) (string, error) {
	if docId == "" {
		return "", errors.New("CHƯA CẤU HÌNH DOC_ID COMMENT! Vui lòng vào Cài đặt (Settings) và điền dãy số GraphQL Comment Doc_ID lấy từ trình duyệt của bạn.")
	}

	// 1. Phân tích Link và cào thông số giống Like
	req, err := http.NewRequest("GET", postURL, nil)
	if err != nil {
		return "", fmt.Errorf("Lỗi tạo request GET: %v", err)
	}
	req.Header.Set("Cookie", cookie)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml")
	req.Header.Set("Accept-Language", "vi-VN,vi;q=0.9,en-US;q=0.8,en;q=0.7")

	resp, err := f.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Lỗi kết nối Facebook: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	htmlText := string(bodyBytes)

	// Lấy fb_dtsg
	fbDtsg := ""
	dtsgRegex1 := regexp.MustCompile(`\["DTSGInitialData",\[\],{"token":"([^"]+)"}`)
	if m := dtsgRegex1.FindStringSubmatch(htmlText); len(m) > 1 {
		fbDtsg = m[1]
	} else {
		dtsgRegex2 := regexp.MustCompile(`"DTSGInitialData",.*?"token":"([^"]+)"`)
		if m := dtsgRegex2.FindStringSubmatch(htmlText); len(m) > 1 {
			fbDtsg = m[1]
		}
	}

	if fbDtsg == "" {
		return "", errors.New("Không thể cào được token fb_dtsg")
	}

	// Lấy jazoest
	jazoest := ""
	jazoestRegex := regexp.MustCompile(`name="jazoest" value="(\d+)"`)
	if m := jazoestRegex.FindStringSubmatch(htmlText); len(m) > 1 {
		jazoest = m[1]
	} else {
		jazoestRegex = regexp.MustCompile(`"jazoest":"(\d+)"`)
		if m := jazoestRegex.FindStringSubmatch(htmlText); len(m) > 1 {
			jazoest = m[1]
		}
	}

	// Lấy lsd
	lsd := ""
	lsdRegex := regexp.MustCompile(`"LSD",\[\],{"token":"([^"]+)"}`)
	if m := lsdRegex.FindStringSubmatch(htmlText); len(m) > 1 {
		lsd = m[1]
	} else {
		lsdRegex = regexp.MustCompile(`name="lsd" value="([^"]+)"`)
		if m := lsdRegex.FindStringSubmatch(htmlText); len(m) > 1 {
			lsd = m[1]
		}
	}

	// Lấy Actor ID
	actorId := "0"
	actorRegex := regexp.MustCompile(`"USER_ID":"(\d+)"`)
	if m := actorRegex.FindStringSubmatch(htmlText); len(m) > 1 {
		actorId = m[1]
	}

	// Cố gắng tìm trực tiếp chuỗi Base64
	feedbackIDBase64 := ""
	encodedRegex1 := regexp.MustCompile(`"feedback":{"id":"(ZmV[^"]+)"}`)
	if m := encodedRegex1.FindStringSubmatch(htmlText); len(m) > 1 {
		feedbackIDBase64 = m[1]
	} else {
		encodedRegex2 := regexp.MustCompile(`"feedback_id":"(ZmV[^"]+)"`)
		if m := encodedRegex2.FindStringSubmatch(htmlText); len(m) > 1 {
			feedbackIDBase64 = m[1]
		}
	}

	// Lấy Numeric ID của bài viết
	if feedbackIDBase64 == "" {
		numericID := ""
		idRegex := regexp.MustCompile(`"top_level_post_id":"(\d+)"`)
		if m := idRegex.FindStringSubmatch(htmlText); len(m) > 1 {
			numericID = m[1]
		} else {
			idRegex = regexp.MustCompile(`"ent_id":"(\d+)"`)
			if m := idRegex.FindStringSubmatch(htmlText); len(m) > 1 {
				numericID = m[1]
			} else {
				parsedUrl, err := url.Parse(postURL)
				if err == nil {
					fallbackId := parsedUrl.Query().Get("story_fbid")
					if fallbackId != "" && !strings.HasPrefix(fallbackId, "pfbid") {
						numericID = fallbackId
					}
				}
			}
		}

		if numericID == "" {
			return "", errors.New("Không bóc tách được Feedback ID hoặc Numeric ID Bài viết từ mã nguồn html.")
		}

		feedbackRaw := "feedback:" + numericID
		feedbackIDBase64 = base64.StdEncoding.EncodeToString([]byte(feedbackRaw))
	}

	// 4. Bắn Request GraphQL Comment
	variablesMap := map[string]interface{}{
		"input": map[string]interface{}{
			"feedback_id":        feedbackIDBase64,
			"message":            map[string]interface{}{"ranges": []interface{}{}, "text": message},
			"actor_id":           actorId,
			"client_mutation_id": "1",
		},
	}
	varsJSON, _ := json.Marshal(variablesMap)

	data := url.Values{}
	data.Set("fb_dtsg", fbDtsg)
	if jazoest != "" {
		data.Set("jazoest", jazoest)
	}
	if lsd != "" {
		data.Set("lsd", lsd)
	}
	
	data.Set("__user", actorId)
	data.Set("__a", "1")
	data.Set("__req", "1")
	data.Set("variables", string(varsJSON))
	data.Set("doc_id", docId)

	postReq, _ := http.NewRequest("POST", "https://www.facebook.com/api/graphql/", strings.NewReader(data.Encode()))
	postReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	postReq.Header.Set("Cookie", cookie)
	postReq.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")

	res, err := f.client.Do(postReq)
	if err != nil {
		return "", fmt.Errorf("Network lỗi khi bắn POST: %v", err)
	}
	defer res.Body.Close()

	fbBodyBytes, _ := io.ReadAll(res.Body)
	fbBody := string(fbBodyBytes)
	
	os.WriteFile("debug_facebook_response.json", fbBodyBytes, 0644)

	if res.StatusCode != 200 {
		return "", fmt.Errorf("GraphQL HTTP %d: %s", res.StatusCode, fbBody)
	}
	
	// Kiểm tra xem có bằng chứng thành công trong data không.
	// Sử dụng `"comment_create"` hoặc `"comment"` để nhận diện data hợp lệ.
	if strings.Contains(fbBody, `"comment_create"`) || strings.Contains(fbBody, `"comment"`) {
		return fmt.Sprintf("Comment thành công! Trạng thái HTTP: 200"), nil
	}

	if strings.Contains(fbBody, "was not found") && strings.Contains(fbBody, "The GraphQL document") {
		return "", errors.New("GraphQL Doc_ID Comment đã hết hạn (Not Found). Vui lòng cập nhật Doc_ID mới trong Cài Đặt!")
	}

	if strings.Contains(fbBody, `"severity":"CRITICAL"`) {
		return "", fmt.Errorf("GraphQL thực thi lỗi nghiêm trọng: %s", strings.TrimPrefix(fbBody, "for (;;);"))
	}

	if strings.Contains(fbBody, `"errors":[{`) {
		return "", fmt.Errorf("GraphQL trả về lỗi nhưng không có dữ liệu thành công: %s", strings.TrimPrefix(fbBody, "for (;;);"))
	}

	return fmt.Sprintf("Comment có thể thành công nhưng không có marker! Trạng thái HTTP: 200"), nil
}

func createPostPayload(uid, message string, photoIDs []string) map[string]interface{} {
	attachments := []interface{}{}
	for _, pid := range photoIDs {
		if pid != "" {
			attachments = append(attachments, map[string]interface{}{
				"photo": map[string]string{"id": pid},
			})
		}
	}

	idempotencyToken := generatePseudoUUID() + "_FEED"
	
	return map[string]interface{}{
		"input": map[string]interface{}{
			"actor_id":                 uid,
			"message":                  map[string]string{"text": message},
			"attachments":              attachments,
			"source":                   "WWW",
			"composer_entry_point":     "inline_composer",
			"composer_source_surface":  "timeline",
			"idempotence_token":        idempotencyToken,
			"audience": map[string]interface{}{
				"privacy": map[string]interface{}{
					"base_state": "EVERYONE",
				},
			},
			"client_mutation_id": "1",
		},
	}
}

func (f *FacebookProvider) PostToFacebook(cookie string, message string, photoIDs []string, docId string) (string, error) {
	if docId == "" {
		return "", errors.New("CHƯA CẤU HÌNH DOC_ID ĐĂNG BÀI! Vui lòng vào Cài đặt (Settings) và điền dãy số GraphQL Post Doc_ID.")
	}

	// Lấy token từ trang chủ
	req, err := http.NewRequest("GET", "https://www.facebook.com/", nil)
	if err != nil {
		return "", fmt.Errorf("Lỗi tạo request GET: %v", err)
	}
	req.Header.Set("Cookie", cookie)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml")

	resp, err := f.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Lỗi kết nối Facebook: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	htmlText := string(bodyBytes)

	// Lấy fb_dtsg
	fbDtsg := ""
	dtsgRegex1 := regexp.MustCompile(`\["DTSGInitialData",\[\],{"token":"([^"]+)"}`)
	if m := dtsgRegex1.FindStringSubmatch(htmlText); len(m) > 1 {
		fbDtsg = m[1]
	} else {
		dtsgRegex2 := regexp.MustCompile(`"DTSGInitialData",.*?"token":"([^"]+)"`)
		if m := dtsgRegex2.FindStringSubmatch(htmlText); len(m) > 1 {
			fbDtsg = m[1]
		}
	}

	if fbDtsg == "" {
		return "", errors.New("Không thể cào được token fb_dtsg từ trang chủ")
	}

	actorId := "0"
	actorRegex := regexp.MustCompile(`"USER_ID":"(\d+)"`)
	if m := actorRegex.FindStringSubmatch(htmlText); len(m) > 1 {
		actorId = m[1]
	}

	// Encode variables
	variables := createPostPayload(actorId, message, photoIDs)
	varsJSON, _ := json.Marshal(variables)

	data := url.Values{}
	data.Set("fb_dtsg", fbDtsg)
	data.Set("__user", actorId)
	data.Set("__a", "1")
	data.Set("__req", "1")
	data.Set("variables", string(varsJSON))
	data.Set("doc_id", docId)

	postReq, _ := http.NewRequest("POST", "https://www.facebook.com/api/graphql/", strings.NewReader(data.Encode()))
	postReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	postReq.Header.Set("Cookie", cookie)
	postReq.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")

	res, err := f.client.Do(postReq)
	if err != nil {
		return "", fmt.Errorf("Network lỗi khi bắn POST: %v", err)
	}
	defer res.Body.Close()

	fbBodyBytes, _ := io.ReadAll(res.Body)
	fbBody := string(fbBodyBytes)
	
	os.WriteFile("debug_facebook_response.json", fbBodyBytes, 0644)

	if res.StatusCode != 200 {
		return "", fmt.Errorf("GraphQL HTTP %d: %s", res.StatusCode, fbBody)
	}

	// Đăng bài thành công nếu response có story_create marker
	if strings.Contains(fbBody, `"story_create"`) {
		return fmt.Sprintf("Đăng bài thành công! Trạng thái HTTP: 200"), nil
	}

	if strings.Contains(fbBody, "was not found") && strings.Contains(fbBody, "The GraphQL document") {
		return "", errors.New("GraphQL Doc_ID Đăng bài đã hết hạn (Not Found). Vui lòng cập nhật Doc_ID mới trong Cài Đặt!")
	}

	if strings.Contains(fbBody, `"severity":"CRITICAL"`) {
		return "", fmt.Errorf("GraphQL thực thi lỗi nghiêm trọng: %s", strings.TrimPrefix(fbBody, "for (;;);"))
	}

	if strings.Contains(fbBody, `"errors":[{`) {
		return "", fmt.Errorf("GraphQL trả về lỗi nhưng không có dữ liệu thành công: %s", strings.TrimPrefix(fbBody, "for (;;);"))
	}

	return fmt.Sprintf("Đăng bài có thể thành công (không có error array)! Trạng thái HTTP: 200"), nil
}

func generatePseudoUUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "07eaf94f-ce2c-4840-8011-9151a61e8118"
	}
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func (f *FacebookProvider) UploadPhoto(cookie, filePath string) (string, error) {
	// Bước 1: Lấy fb_dtsg và uid
	reqGet, err := http.NewRequest("GET", "https://www.facebook.com/", nil)
	if err != nil {
		return "", fmt.Errorf("lỗi tạo request GET: %v", err)
	}
	reqGet.Header.Set("Cookie", cookie)
	reqGet.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
	reqGet.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml")

	respGet, err := f.client.Do(reqGet)
	if err != nil {
		return "", fmt.Errorf("lỗi kết nối Facebook: %v", err)
	}
	bodyBytes, _ := io.ReadAll(respGet.Body)
	respGet.Body.Close()
	htmlText := string(bodyBytes)

	fbDtsg := ""
	dtsgRegex1 := regexp.MustCompile(`\["DTSGInitialData",\[\],{"token":"([^"]+)"}`)
	if m := dtsgRegex1.FindStringSubmatch(htmlText); len(m) > 1 {
		fbDtsg = m[1]
	} else {
		dtsgRegex2 := regexp.MustCompile(`"DTSGInitialData",.*?"token":"([^"]+)"`)
		if m := dtsgRegex2.FindStringSubmatch(htmlText); len(m) > 1 {
			fbDtsg = m[1]
		}
	}

	uid := ""
	uidRegex := regexp.MustCompile(`"USER_ID":"(\d+)"`)
	if m := uidRegex.FindStringSubmatch(htmlText); len(m) > 1 {
		uid = m[1]
	} else if cRegex := regexp.MustCompile(`c_user=(\d+)`); len(cRegex.FindStringSubmatch(cookie)) > 1 {
		uid = cRegex.FindStringSubmatch(cookie)[1]
	}

	if fbDtsg == "" {
		return "", errors.New("không tìm thấy fb_dtsg, cookie có thể đã chết")
	}

	uploadURL := fmt.Sprintf("https://upload.facebook.com/ajax/react_composer/attachments/photo/upload?__a=1&__req=1&__user=%s", uid)

	// 2. Mở file ảnh từ ổ đĩa
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("không thể mở file: %v", err)
	}
	defer file.Close()

	// 3. Tạo body dạng multipart
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	writer.WriteField("fb_dtsg", fbDtsg)
	writer.WriteField("profile_id", uid)
	writer.WriteField("source", "8")

	part, err := writer.CreateFormFile("f", "image.jpg")
	if err != nil {
		return "", fmt.Errorf("lỗi tạo form data: %v", err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return "", fmt.Errorf("lỗi copy dữ liệu file: %v", err)
	}

	err = writer.Close()
	if err != nil {
		return "", fmt.Errorf("lỗi đóng multipart writer: %v", err)
	}

	// 3. Khởi tạo request
	req, err := http.NewRequest("POST", uploadURL, body)
	if err != nil {
		return "", fmt.Errorf("lỗi tạo request HTTP: %v", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Cookie", cookie)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
	req.Header.Set("Origin", "https://www.facebook.com")
	req.Header.Set("Referer", "https://www.facebook.com/")
	req.Header.Set("Accept", "*/*")

	// 4. Gửi request
	resp, err := f.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("lỗi kết nối upload: %v", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	// Trim JSON hijacking prevention prefix
	bodyStr := string(respBody)
	if strings.HasPrefix(bodyStr, "for (;;);") {
		bodyStr = strings.TrimPrefix(bodyStr, "for (;;);")
		respBody = []byte(bodyStr)
	}

	var fbResp struct {
		PhotoID string `json:"photo_id"`
	}
	json.Unmarshal(respBody, &fbResp)
	
	photoID := fbResp.PhotoID

	if photoID == "" {
		var payloadResp struct {
			Payload struct {
				Fbid    string `json:"fbid"`
				PhotoID string `json:"photoID"`
			} `json:"payload"`
		}
		json.Unmarshal(respBody, &payloadResp)
		if payloadResp.Payload.PhotoID != "" {
			photoID = payloadResp.Payload.PhotoID
		} else if payloadResp.Payload.Fbid != "" {
			photoID = payloadResp.Payload.Fbid
		}
	}

	if photoID == "" {
		// Fallback dynamic parse just in case it's numeric
		var dynamic map[string]interface{}
		if err := json.Unmarshal(respBody, &dynamic); err == nil {
			if pInfo, ok := dynamic["payload"].(map[string]interface{}); ok {
				if id, ok := pInfo["photoID"]; ok && fmt.Sprintf("%v", id) != "" {
					photoID = fmt.Sprintf("%v", id)
				} else if id, ok := pInfo["fbid"]; ok && fmt.Sprintf("%v", id) != "" {
					photoID = fmt.Sprintf("%v", id)
				}
			}
		}
	}

	if photoID == "" {
		// Truncate the response body to avoid 4MB log lines.
		errorResp := string(respBody)
		if len(errorResp) > 1000 {
			errorResp = errorResp[:1000] + "...(truncated)"
		}
		return "", fmt.Errorf("Không lấy được Photo ID. FB phản hồi: %s", errorResp)
	}

	return photoID, nil
}
