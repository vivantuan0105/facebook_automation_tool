package providers

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

type FacebookProvider struct {
	client *http.Client
}

type FBUploadResponse struct {
	PhotoID string `json:"fbid"`
}

var jsonTextValueRegex = regexp.MustCompile(`"(?:text|title|subtitle)"\s*:\s*"((?:\\.|[^"\\])*)"`)

func normalizeWhitespace(s string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(s)), " ")
}

func parseLocationFromLabel(label string) string {
	cleaned := normalizeWhitespace(label)
	if cleaned == "" {
		return ""
	}

	re := regexp.MustCompile(`(?i)^(sống tại|đến từ|lives in|from)\s*:?\s*(.+)$`)
	if m := re.FindStringSubmatch(cleaned); len(m) >= 3 {
		return normalizeWhitespace(m[2])
	}
	return ""
}

func isLocationSubtitle(subtitle string) bool {
	low := strings.ToLower(normalizeWhitespace(subtitle))
	if low == "" {
		return false
	}

	keywords := []string{
		"tỉnh/thành phố hiện tại",
		"quê quán",
		"sống tại",
		"đến từ",
		"current city",
		"hometown",
		"lives in",
		"from",
	}
	for _, k := range keywords {
		if strings.Contains(low, k) {
			return true
		}
	}
	return false
}

func extractTextCandidatesFromJSON(raw string) []string {
	matches := jsonTextValueRegex.FindAllStringSubmatch(raw, -1)
	if len(matches) == 0 {
		return nil
	}

	result := make([]string, 0, len(matches))
	for _, m := range matches {
		if len(m) < 2 || m[1] == "" {
			continue
		}
		decoded, err := strconv.Unquote(`"` + m[1] + `"`)
		if err != nil {
			decoded = m[1]
		}
		decoded = normalizeWhitespace(decoded)
		if decoded != "" {
			result = append(result, decoded)
		}
	}
	return result
}

func parseLocationFromLabelSafe(label string) string {
	cleaned := normalizeWhitespace(label)
	if cleaned == "" {
		return ""
	}

	prefixes := []string{
		"s\u1ed1ng t\u1ea1i",
		"\u0111\u1ebfn t\u1eeb",
		"song tai",
		"den tu",
		"lives in",
		"from",
	}
	for _, p := range prefixes {
		re := regexp.MustCompile(`(?i)^` + regexp.QuoteMeta(p) + `\s*:?\s*(.+)$`)
		if m := re.FindStringSubmatch(cleaned); len(m) >= 2 {
			return normalizeWhitespace(m[1])
		}
	}
	return ""
}

func isLocationSubtitleSafe(subtitle string) bool {
	low := strings.ToLower(normalizeWhitespace(subtitle))
	if low == "" {
		return false
	}

	keywords := []string{
		"t\u1ec9nh/th\u00e0nh ph\u1ed1 hi\u1ec7n t\u1ea1i",
		"qu\u00ea qu\u00e1n",
		"s\u1ed1ng t\u1ea1i",
		"\u0111\u1ebfn t\u1eeb",
		"tinh/thanh pho hien tai",
		"que quan",
		"song tai",
		"den tu",
		"current city",
		"hometown",
		"lives in",
		"from",
	}
	for _, k := range keywords {
		if strings.Contains(low, k) {
			return true
		}
	}
	return false
}

func isBirthdaySubtitleSafe(subtitle string) bool {
	low := strings.ToLower(normalizeWhitespace(subtitle))
	if low == "" {
		return false
	}

	keywords := []string{
		"ng\u00e0y sinh",
		"n\u0103m sinh",
		"sinh nh\u1eadt",
		"ngay sinh",
		"nam sinh",
		"sinh nhat",
		"birthday",
		"birth date",
		"date of birth",
	}
	for _, k := range keywords {
		if strings.Contains(low, k) {
			return true
		}
	}
	return false
}

func isGenderLabelSafe(text string) bool {
	low := strings.ToLower(normalizeWhitespace(text))
	if low == "" {
		return false
	}
	keywords := []string{
		"giới tính",
		"gioi tinh",
		"gender",
		"sex",
	}
	for _, k := range keywords {
		if strings.Contains(low, k) {
			return true
		}
	}
	return false
}

func normalizeGenderValue(raw string) string {
	low := strings.ToLower(normalizeWhitespace(raw))
	if low == "" {
		return ""
	}

	if strings.Contains(low, "không công khai") || strings.Contains(low, "khong cong khai") ||
		strings.Contains(low, "not public") || strings.Contains(low, "private") || strings.Contains(low, "only me") {
		return "Không công khai"
	}

	if regexp.MustCompile(`(?i)\bfemale\b`).MatchString(low) || strings.Contains(low, "nữ") || strings.Contains(low, " nu ") {
		return "Nữ (FEMALE)"
	}

	// Tránh bắt nhầm cụm "nam sinh".
	if regexp.MustCompile(`(?i)\bmale\b`).MatchString(low) {
		return "Nam (MALE)"
	}
	if strings.Contains(low, "nam sinh") || strings.Contains(low, "năm sinh") {
		return ""
	}
	if strings.Contains(low, " nam ") || strings.HasPrefix(low, "nam ") || strings.HasSuffix(low, " nam") || low == "nam" {
		return "Nam (MALE)"
	}

	return ""
}

func extractGenderFromCandidates(candidates []string) string {
	if len(candidates) == 0 {
		return ""
	}

	// 1) Nếu text tự chứa giá trị rõ ràng.
	for _, c := range candidates {
		if g := normalizeGenderValue(c); g != "" && !isGenderLabelSafe(c) {
			return g
		}
	}

	// 2) Tìm quanh nhãn "Giới tính".
	neighborOffsets := []int{-2, -1, 1, 2}
	for i, c := range candidates {
		text := normalizeWhitespace(c)
		if !isGenderLabelSafe(text) {
			continue
		}

		if g := normalizeGenderValue(text); g != "" {
			return g
		}
		for _, off := range neighborOffsets {
			j := i + off
			if j < 0 || j >= len(candidates) {
				continue
			}
			if g := normalizeGenderValue(candidates[j]); g != "" {
				return g
			}
		}
	}
	return ""
}

func parseBirthdayFromLabelSafe(label string) string {
	cleaned := normalizeWhitespace(label)
	if cleaned == "" {
		return ""
	}

	low := strings.ToLower(cleaned)
	keywords := []string{
		"ng\u00e0y sinh",
		"n\u0103m sinh",
		"sinh nh\u1eadt",
		"ngay sinh",
		"nam sinh",
		"sinh nhat",
		"birthday",
		"birth date",
		"date of birth",
		"born on",
	}
	hasKeyword := false
	for _, k := range keywords {
		if strings.Contains(low, k) {
			hasKeyword = true
			break
		}
	}
	if !hasKeyword {
		return ""
	}

	hiddenKeywords := []string{
		"kh\u00f4ng c\u00f4ng khai",
		"khong cong khai",
		"ch\u1ec9 m\u00ecnh t\u00f4i",
		"chi minh toi",
		"not public",
		"private",
		"only me",
	}
	for _, k := range hiddenKeywords {
		if strings.Contains(low, k) {
			return "Khong cong khai"
		}
	}

	noiseKeywords := []string{
		"b\u1ea1n ch\u1ec9 c\u00f3 th\u1ec3 ch\u1ec9nh s\u1eeda",
		"ban chi co the chinh sua",
		"s\u1ed1 l\u1ea7n nh\u1ea5t \u0111\u1ecbnh",
		"so lan nhat dinh",
		"t\u00ecm hi\u1ec3u th\u00eam",
		"tim hieu them",
		"learn more",
		"can be edited",
		"policy",
	}
	for _, k := range noiseKeywords {
		if strings.Contains(low, k) {
			return ""
		}
	}

	value := cleaned
	for _, k := range keywords {
		re := regexp.MustCompile(`(?i)\b` + regexp.QuoteMeta(k) + `\b`)
		value = re.ReplaceAllString(value, " ")
	}
	value = strings.Trim(value, ":-|.,;•· ")
	value = normalizeWhitespace(value)
	if value == "" {
		return ""
	}
	if !looksLikeBirthdayValue(value) {
		return ""
	}
	return value
}

func looksLikeBirthdayValue(value string) bool {
	low := strings.ToLower(normalizeWhitespace(value))
	if low == "" {
		return false
	}

	monthKeywords := []string{
		"thang", "th\u00e1ng", "month",
		"january", "february", "march", "april", "may", "june",
		"july", "august", "september", "october", "november", "december",
		"jan", "feb", "mar", "apr", "jun", "jul", "aug", "sep", "oct", "nov", "dec",
	}
	for _, m := range monthKeywords {
		if strings.Contains(low, m) {
			return true
		}
	}

	numericDatePatterns := []string{
		`^\d{1,2}[\/\-.]\d{1,2}([\/\-.]\d{2,4})?$`,
		`^\d{4}[\/\-.]\d{1,2}[\/\-.]\d{1,2}$`,
	}
	for _, p := range numericDatePatterns {
		if regexp.MustCompile(p).MatchString(low) {
			return true
		}
	}

	dayMonthWords := regexp.MustCompile(`\b\d{1,2}\b.*\b(tháng|thang|january|february|march|april|may|june|july|august|september|october|november|december|jan|feb|mar|apr|jun|jul|aug|sep|oct|nov|dec)\b`)
	return dayMonthWords.MatchString(low)
}

func looksLikeLocationValue(value string) bool {
	low := strings.ToLower(normalizeWhitespace(value))
	if low == "" {
		return false
	}

	invalid := []string{
		"khong cong khai",
		"không công khai",
		"sinh nhat",
		"sinh nhật",
		"ngay sinh",
		"ngày sinh",
		"nam sinh",
		"năm sinh",
		"chinh sua",
		"chỉnh sửa",
	}
	for _, x := range invalid {
		if strings.Contains(low, x) {
			return false
		}
	}

	if strings.Contains(low, "thành phố") || strings.Contains(low, "thanh pho") ||
		strings.Contains(low, "tỉnh") || strings.Contains(low, "tinh") ||
		strings.Contains(low, "city") || strings.Contains(low, "hometown") {
		return true
	}

	// Cho phép tên địa danh ngắn nếu không chứa số/câu hệ thống.
	if len([]rune(low)) >= 3 && len([]rune(low)) <= 60 && !regexp.MustCompile(`\d{2,}`).MatchString(low) {
		return true
	}
	return false
}

func extractBirthYearFromBirthday(birthday string) string {
	b := strings.TrimSpace(strings.ToLower(birthday))
	if b == "" {
		return ""
	}
	invalid := []string{
		"khong cong khai",
		"không công khai",
		"not public",
		"private",
		"only me",
	}
	for _, x := range invalid {
		if strings.Contains(b, x) {
			return ""
		}
	}

	years := regexp.MustCompile(`\b(19\d{2}|20\d{2})\b`).FindAllString(birthday, -1)
	if len(years) == 0 {
		return ""
	}
	return years[len(years)-1]
}

func extractBirthYearFromCandidates(candidates []string, birthday string) string {
	if len(candidates) == 0 {
		return ""
	}

	findYear := func(text string) string {
		years := regexp.MustCompile(`\b(19\d{2}|20\d{2})\b`).FindAllString(text, -1)
		if len(years) == 0 {
			return ""
		}
		return years[len(years)-1]
	}

	birthdayNorm := strings.ToLower(normalizeWhitespace(birthday))
	neighborOffsets := []int{-4, -3, -2, -1, 1, 2, 3, 4}

	// 1) Ưu tiên quanh nhãn "năm sinh"
	for i, c := range candidates {
		text := strings.ToLower(normalizeWhitespace(c))
		if text == "" {
			continue
		}
		if strings.Contains(text, "năm sinh") || strings.Contains(text, "nam sinh") || strings.Contains(text, "birth year") {
			if y := findYear(c); y != "" {
				return y
			}
			for _, off := range neighborOffsets {
				j := i + off
				if j < 0 || j >= len(candidates) {
					continue
				}
				if y := findYear(candidates[j]); y != "" {
					return y
				}
			}
		}
	}

	// 2) Quanh cụm sinh nhật/ngày sinh hoặc text ngày sinh đã parse được
	for i, c := range candidates {
		text := strings.ToLower(normalizeWhitespace(c))
		if text == "" {
			continue
		}
		isBirthdayAnchor := isBirthdaySubtitleSafe(text)
		if !isBirthdayAnchor && birthdayNorm != "" && text == birthdayNorm {
			isBirthdayAnchor = true
		}
		if !isBirthdayAnchor && birthdayNorm != "" && strings.Contains(text, birthdayNorm) {
			isBirthdayAnchor = true
		}
		if !isBirthdayAnchor {
			continue
		}

		if y := findYear(c); y != "" {
			return y
		}
		for _, off := range neighborOffsets {
			j := i + off
			if j < 0 || j >= len(candidates) {
				continue
			}
			if y := findYear(candidates[j]); y != "" {
				return y
			}
		}
	}

	return ""
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
	if actorId == "0" || actorId == "" {
		fmt.Println("[DEBUG-ERROR] Không trích xuất được USER_ID, có thể Cookie đã chết, bị Checkpoint, hoặc định dạng trống/sai.")
		return "", errors.New("Cookie đã chết (Trang trả về dưới quyền Khách/USER_ID=0) hoặc bị checkpoint!")
	}

	// 1. Lấy Numeric ID của bài viết
	numericID := ""
	idRegex := regexp.MustCompile(`"top_level_post_id":"(\d+)"`)
	if m := idRegex.FindStringSubmatch(htmlText); len(m) > 1 {
		numericID = m[1]
	} else {
		idRegex = regexp.MustCompile(`"ent_id":"(\d+)"`)
		if m := idRegex.FindStringSubmatch(htmlText); len(m) > 1 {
			numericID = m[1]
		} else {
			// Cào tham số từ URL
			parsedUrl, err := url.Parse(postURL)
			if err == nil {
				fallbackId := parsedUrl.Query().Get("story_fbid")
				if fallbackId != "" && !strings.HasPrefix(fallbackId, "pfbid") {
					numericID = fallbackId
				}
			}
		}
	}

	feedbackIDBase64 := ""

	// 2. THUẬT TOÁN TÌM KIẾM THÔNG MINH (Smart Match)
	// Tránh dùng Regex tĩnh ngay từ đầu vì thường bắt nhầm ID của Composer.
	if numericID != "" {
		allFeedbacks := regexp.MustCompile(`(ZmVlZGJhY2s6[a-zA-Z0-9+_/=]+)`).FindAllStringSubmatch(htmlText, -1)
		for _, match := range allFeedbacks {
			if len(match) > 1 {
				decodedBytes, err := base64.StdEncoding.DecodeString(match[1])
				if err == nil {
					decodedStr := string(decodedBytes)
					if strings.HasPrefix(decodedStr, "feedback:") && strings.Contains(decodedStr, numericID) {
						feedbackIDBase64 = match[1]
						fmt.Printf("[DEBUG] Smart Match Found: %s -> %s\n", match[1], decodedStr)
						break
					}
				}
			}
		}
	}

	// 3. Nếu chưa tìm thấy, dùng Regex bóc tách tĩnh
	if feedbackIDBase64 == "" {
		patterns := []string{
			`"feedback"\s*:\s*\{\s*"id"\s*:\s*"(ZmVlZGJhY2s6[^"]+)"`,
			`"feedback_id"\s*:\s*"(ZmVlZGJhY2s6[^"]+)"`,
			`"target_feedback"\s*:\s*\{\s*"id"\s*:\s*"(ZmVlZGJhY2s6[^"]+)"`,
			`"node"\s*:\s*\{\s*"id"\s*:\s*"(ZmVlZGJhY2s6[^"]+)"\s*,\s*"__isFeedback"`,
		}
		for _, p := range patterns {
			rx := regexp.MustCompile(p)
			if m := rx.FindStringSubmatch(htmlText); len(m) > 1 {
				feedbackIDBase64 = m[1]
				break
			}
		}
	}

	// 4. Nếu vẫn trống và có Numeric ID, ghép cơ bản
	if feedbackIDBase64 == "" && numericID != "" {
		feedbackRaw := "feedback:" + numericID
		feedbackIDBase64 = base64.StdEncoding.EncodeToString([]byte(feedbackRaw))
	}

	if feedbackIDBase64 == "" {
		return "", errors.New("Không bóc tách được Feedback ID hoặc Numeric ID Bài viết từ mã nguồn html.")
	}

	// 3. Mapping Reaction ID
	reactionToFB := map[string]string{
		"Like":  "1635855486666999",
		"Love":  "1678524932434102",
		"Care":  "613557422527858",
		"Haha":  "115940658764963",
		"Wow":   "478547315650144",
		"Sad":   "908563459236466",
		"Angry": "444813342392137",
	}

	fbReactionType := reactionToFB[reactionType]
	if fbReactionType == "" {
		fbReactionType = "1635855486666999" // Mặc định là Like
	}

	// doc_id được lấy từ tham số truyền vào, không hardcode nữa

	// 4. Bắn Request GraphQL
	variables := fmt.Sprintf(`{"input":{"feedback_id":"%s","feedback_reaction_id":"%s","feedback_source":"OBJECT","is_tracking_encrypted":true,"tracking":[],"session_id":"%s","actor_id":"%s","client_mutation_id":"1"},"useDefaultActor":false,"scale":1.5}`,
		feedbackIDBase64, fbReactionType, generatePseudoUUID(), actorId)

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
	if strings.Contains(fbBody, `"viewer_feedback_reaction_info":{"`) {
		return fmt.Sprintf("Thành công! ID mã hóa: %s. Trạng thái HTTP: 200", feedbackIDBase64), nil
	}

	// Xử lý lỗi tinh vi hơn
	if strings.Contains(fbBody, "was not found") && strings.Contains(fbBody, "The GraphQL document") {
		return "", errors.New("GraphQL Doc_ID Like đã hết hạn (Not Found). Vui lòng cập nhật Doc_ID mới trong Cài Đặt!")
	}

	if strings.Contains(fbBody, `"severity":"CRITICAL"`) {
		summary := ""
		if m := regexp.MustCompile(`"summary":"([^"]+)"`).FindStringSubmatch(fbBody); len(m) > 1 {
			summary = m[1]
		}
		description := ""
		if m := regexp.MustCompile(`"description":"([^"]+)"`).FindStringSubmatch(fbBody); len(m) > 1 {
			description = m[1]
		}

		if summary != "" {
			return "", fmt.Errorf("Bị chặn/Từ chối: %s (Chi tiết: %s)", summary, description)
		}
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

	// 1. Lấy Numeric ID của bài viết
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

	feedbackIDBase64 := ""

	if numericID != "" {
		allFeedbacks := regexp.MustCompile(`(ZmVlZGJhY2s6[a-zA-Z0-9+_/=]+)`).FindAllStringSubmatch(htmlText, -1)
		for _, match := range allFeedbacks {
			if len(match) > 1 {
				decodedBytes, err := base64.StdEncoding.DecodeString(match[1])
				if err == nil {
					decodedStr := string(decodedBytes)
					if strings.HasPrefix(decodedStr, "feedback:") && strings.Contains(decodedStr, numericID) {
						feedbackIDBase64 = match[1]
						break
					}
				}
			}
		}
	}

	if feedbackIDBase64 == "" {
		patterns := []string{
			`"feedback"\s*:\s*\{\s*"id"\s*:\s*"(ZmVlZGJhY2s6[^"]+)"`,
			`"feedback_id"\s*:\s*"(ZmVlZGJhY2s6[^"]+)"`,
			`"target_feedback"\s*:\s*\{\s*"id"\s*:\s*"(ZmVlZGJhY2s6[^"]+)"`,
			`"node"\s*:\s*\{\s*"id"\s*:\s*"(ZmVlZGJhY2s6[^"]+)"\s*,\s*"__isFeedback"`,
		}
		for _, p := range patterns {
			rx := regexp.MustCompile(p)
			if m := rx.FindStringSubmatch(htmlText); len(m) > 1 {
				feedbackIDBase64 = m[1]
				break
			}
		}
	}

	if feedbackIDBase64 == "" && numericID != "" {
		feedbackRaw := "feedback:" + numericID
		feedbackIDBase64 = base64.StdEncoding.EncodeToString([]byte(feedbackRaw))
	}

	if feedbackIDBase64 == "" {
		return "", errors.New("Không bóc tách được Feedback ID hoặc Numeric ID Bài viết từ mã nguồn html.")
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
			"actor_id":                uid,
			"message":                 map[string]string{"text": message},
			"attachments":             attachments,
			"source":                  "WWW",
			"composer_entry_point":    "inline_composer",
			"composer_source_surface": "timeline",
			"idempotence_token":       idempotencyToken,
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

// ---------------------------------------------------------
// TÍNH NĂNG QUÉT THÔNG TIN PROFILE VÀ LƯU INFO.TXT
// ---------------------------------------------------------

func (f *FacebookProvider) ScanAccountInfo(cookie string, uid string, docId string) (map[string]string, error) {
	// Lấy trang cá nhân section About để có nhiều thông tin nhất
	targetPages := []string{
		"https://www.facebook.com/profile.php?id=" + uid + "&sk=directory_personal_details",
		"https://www.facebook.com/profile.php?id=" + uid + "&sk=about",
	}

	htmlParts := make([]string, 0, len(targetPages))
	for _, pageURL := range targetPages {
		req, err := http.NewRequest("GET", pageURL, nil)
		if err != nil {
			continue
		}
		req.Header.Set("Cookie", cookie)
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
		req.Header.Set("Accept", "text/html,application/xhtml+xml")
		req.Header.Set("Accept-Language", "vi-VN,vi;q=0.9,en-US;q=0.8,en;q=0.7")

		resp, err := f.client.Do(req)
		if err != nil {
			continue
		}
		bodyBytes, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if len(bodyBytes) > 0 {
			htmlParts = append(htmlParts, string(bodyBytes))
		}
	}

	if len(htmlParts) == 0 {
		return nil, errors.New("không tải được trang profile để quét thông tin")
	}

	htmlText := strings.Join(htmlParts, "\n")

	// Lưới Vét Regex: Trích xuất tất cả các khối JSON trong thẻ script
	jsonFragments := []string{}
	scripts := regexp.MustCompile(`<script type="application/json"[^>]*>(.*?)</script>`).FindAllStringSubmatch(htmlText, -1)
	for _, m := range scripts {
		if len(m) > 1 {
			jsonFragments = append(jsonFragments, m[1])
		}
	}

	// Bổ sung lưới phụ vét sâu hơn từ các biến JS (nếu có)
	relayScripts := regexp.MustCompile(`"data"\s*:\s*(\{.*?\})\s*,\s*"errors"`).FindAllStringSubmatch(htmlText, -1)
	for _, m := range relayScripts {
		if len(m) > 1 {
			jsonFragments = append(jsonFragments, `{"data":` + m[1] + `}`)
		}
	}

	resultMap := map[string]string{
		"gender":    "",
		"birthday":  "",
		"birthYear": "",
		"location":  "",
		"friends":   "",
		"followers": "",
		"name":      "",
	}

	var details string

	for _, rawJson := range jsonFragments {
		// 1. Tìm Tên (Name) - Có thể nằm ở nhiều chỗ khác nhau tùy section
		namePaths := []string{
			"data.node.name",
			"data.user.name",
			"data.viewer.actor.name",
			"__bbox.result.data.node.name",
			"__bbox.result.data.user.name",
			"__bbox.result.data.name",
			"data.name",
		}
		for _, path := range namePaths {
			if n := gjson.Get(rawJson, path).String(); n != "" {
				resultMap["name"] = n
				break
			}
		}

		// 2. Tìm Giới tính (Gender)
		genderPaths := []string{
			"data.node.gender",
			"data.user.gender",
			"data.viewer.actor.gender",
			"__bbox.result.data.node.gender",
			"__bbox.result.data.user.gender",
			"__bbox.result.data.gender",
		}
		for _, path := range genderPaths {
			if g := gjson.Get(rawJson, path).String(); g != "" {
				if ng := normalizeGenderValue(g); ng != "" {
					resultMap["gender"] = ng
				} else {
					resultMap["gender"] = g
				}
				break
			}
		}

		// 3. Tìm Followers / Friends
		fCount := gjson.Get(rawJson, "data.user.profile_header_actions.follower_count.count").String()
		if fCount == "" {
			fCount = gjson.Get(rawJson, "data.node.profile_header_actions.follower_count.count").String()
		}
		if fCount != "" && resultMap["followers"] == "" {
			resultMap["followers"] = fCount
		}

		// 4. Quét Context Items (Học vấn, Nơi ở, v.v.)
		// Comet layout thường để ở timeline_context_item_sections hoặc profile_about_all_sections
		contextPaths := []string{
			"data.node.timeline_context_item_sections.0.items",
			"data.user.timeline_context_item_sections.0.items",
			"__bbox.result.data.node.timeline_context_item_sections.0.items",
			"__bbox.result.data.user.timeline_context_item_sections.0.items",
			"data.node.profile_about_all_sections.edges",
			"data.user.profile_about_all_sections.edges",
			"__bbox.result.data.node.profile_about_all_sections.edges",
			"__bbox.result.data.user.profile_about_all_sections.edges",
		}

		for _, cp := range contextPaths {
			items := gjson.Get(rawJson, cp)
			if items.Exists() {
				items.ForEach(func(key, value gjson.Result) bool {
					// Thử lấy text từ các cấu trúc phức tạp của Facebok
					label := value.Get("renderer.context_item.title.text").String()
					if label == "" {
						label = value.Get("node.title.text").String() // fallback cho about sections
					}
					subtitle := value.Get("renderer.context_item.subtitle.text").String()
					if subtitle == "" {
						subtitle = value.Get("node.subtitle.text").String()
					}

					if label != "" {
						details += fmt.Sprintf("- %s\n", label)
						
						lowLabel := strings.ToLower(label)
						if isGenderLabelSafe(subtitle) {
							if g := normalizeGenderValue(label); g != "" {
								resultMap["gender"] = g
							}
						} else if isGenderLabelSafe(label) {
							if g := normalizeGenderValue(subtitle); g != "" {
								resultMap["gender"] = g
							}
						}
						if loc := parseLocationFromLabelSafe(label); loc != "" {
							resultMap["location"] = loc
						} else if resultMap["location"] == "" && isLocationSubtitleSafe(subtitle) {
							// Nhiều profile trả title là thành phố, subtitle mới cho biết đây là vị trí
								resultMap["location"] = normalizeWhitespace(label)
							}
							if resultMap["birthday"] == "" {
								if b := parseBirthdayFromLabelSafe(label); b != "" {
									resultMap["birthday"] = b
								} else if isBirthdaySubtitleSafe(subtitle) {
									resultMap["birthday"] = normalizeWhitespace(label)
								} else if b := parseBirthdayFromLabelSafe(subtitle); b != "" {
									resultMap["birthday"] = b
								} else if isBirthdaySubtitleSafe(label) && subtitle != "" {
									resultMap["birthday"] = normalizeWhitespace(subtitle)
								}
							}
						if strings.Contains(lowLabel, "sống tại") || strings.Contains(lowLabel, "đến từ") {
							resultMap["location"] = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(label, "Sống tại", ""), "Đến từ", ""))
						}
						if strings.Contains(lowLabel, "người theo dõi") {
							parts := strings.Fields(label)
							if len(parts) > 0 {
								resultMap["followers"] = parts[0]
							}
						}
					}
					return true
				})
				}
			}

			if resultMap["location"] == "" {
				for _, candidate := range extractTextCandidatesFromJSON(rawJson) {
					if loc := parseLocationFromLabelSafe(candidate); loc != "" {
						resultMap["location"] = loc
						break
					}
				}
			}
			if resultMap["birthday"] == "" {
				for _, candidate := range extractTextCandidatesFromJSON(rawJson) {
					if b := parseBirthdayFromLabelSafe(candidate); b != "" {
						resultMap["birthday"] = b
						break
					}
				}
			}
			if g := extractGenderFromCandidates(extractTextCandidatesFromJSON(rawJson)); g != "" {
				resultMap["gender"] = g
			}
		}

	// FALLBACK Cố định bằng Regex nếu JSON Parser thất bại toàn tập
	if resultMap["name"] == "" {
		re := regexp.MustCompile(`"NAME":"([^"]+)"`)
		if m := re.FindStringSubmatch(htmlText); len(m) > 1 {
			resultMap["name"] = m[1]
		}
	}
	if resultMap["gender"] == "" {
		if strings.Contains(htmlText, `"gender":"MALE"`) {
			resultMap["gender"] = "Nam (MALE)"
		} else if strings.Contains(htmlText, `"gender":"FEMALE"`) {
			resultMap["gender"] = "Nữ (FEMALE)"
		}
	}

	// Xử lý các trường trống
	if g := extractGenderFromCandidates(extractTextCandidatesFromJSON(htmlText)); g != "" {
		resultMap["gender"] = g
	}
	if resultMap["birthday"] == "" || resultMap["location"] == "" {
		candidates := extractTextCandidatesFromJSON(htmlText)
		neighborOffsets := []int{-3, -2, -1, 1, 2, 3}
		for i, c := range candidates {
			text := normalizeWhitespace(c)
			if text == "" {
				continue
			}
			if resultMap["birthday"] == "" {
				if b := parseBirthdayFromLabelSafe(text); b != "" {
					resultMap["birthday"] = b
				} else if isBirthdaySubtitleSafe(text) {
					for _, off := range neighborOffsets {
						j := i + off
						if j < 0 || j >= len(candidates) {
							continue
						}
						near := normalizeWhitespace(candidates[j])
						if looksLikeBirthdayValue(near) {
							resultMap["birthday"] = near
							break
						}
					}
				}
			}
			if resultMap["location"] == "" {
				if loc := parseLocationFromLabelSafe(text); loc != "" {
					resultMap["location"] = loc
				} else if isLocationSubtitleSafe(text) {
					for _, off := range neighborOffsets {
						j := i + off
						if j < 0 || j >= len(candidates) {
							continue
						}
						near := normalizeWhitespace(candidates[j])
						if looksLikeLocationValue(near) {
							resultMap["location"] = near
							break
						}
					}
				}
			}
			if resultMap["birthday"] != "" && resultMap["location"] != "" {
				break
			}
		}
	}
	if resultMap["location"] == "" {
		resultMap["location"] = "Không công khai"
	}
	if resultMap["gender"] == "" {
		resultMap["gender"] = "Không công khai"
	}
	if resultMap["followers"] == "" {
		resultMap["followers"] = "0"
	}
	if resultMap["birthday"] == "" {
		resultMap["birthday"] = "Khong cong khai"
	}
	if resultMap["birthYear"] == "" {
		resultMap["birthYear"] = extractBirthYearFromBirthday(resultMap["birthday"])
	}
	if resultMap["birthYear"] == "" {
		resultMap["birthYear"] = extractBirthYearFromCandidates(extractTextCandidatesFromJSON(htmlText), resultMap["birthday"])
	}
	if resultMap["friends"] == "" {
		resultMap["friends"] = "0"
	}

	bio := "" // Thường bị ẩn

	// Tạo Data Thư mục
	path := filepath.Join("Data", uid)
	os.MkdirAll(path, os.ModePerm)

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	content := fmt.Sprintf(`--- THÔNG TIN TÀI KHOẢN ---
Quét lúc: %s
UID: %s

[CƠ BẢN]
Họ tên: %s
Giới tính: %s
Tiểu sử: %s

[VỊ TRÍ & MỐI QUAN HỆ]
Nơi ở hiện tại / Quê quán: %s

[CHI TIẾT TRÊN TRANG CÁ NHÂN]
%s
--------------------------`, timestamp, uid, resultMap["name"], resultMap["gender"], bio, resultMap["location"], details)

	// File path
	filePath := filepath.Join(path, "Profile_Scan.txt")
	os.WriteFile(filePath, []byte(content), 0644)

	return resultMap, nil
}

// ---------------------------------------------------------
// TÍNH NĂNG QUÉT BẠN BÈ VÀ LƯU FRIENDS.TXT
// ---------------------------------------------------------

func (f *FacebookProvider) ScanAccountFriends(cookie string, uid string, docId string) (int, error) {
	// Sử dụng doc_id chuyên biệt cho lấy danh sách bạn bè như hướng dẫn
	docIdToUse := "26206414195674994"

	fbDtsg := ""
	lsd := ""
	jazoest := ""
	cursor := ""
	hasNextPage := true
	friendList := []string{}
	friendMap := make(map[string]bool)

	// Lấy fb_dtsg, lsd, jazoest lần đầu tiên từ trang cá nhân
	req, err := http.NewRequest("GET", "https://www.facebook.com/profile.php?id="+uid, nil)
	if err == nil {
		req.Header.Set("Cookie", cookie)
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
		req.Header.Set("Accept", "text/html,application/xhtml+xml")
		resp, err := f.client.Do(req)
		if err == nil {
			bodyBytes, _ := io.ReadAll(resp.Body)
			htmlText := string(bodyBytes)
			
			dtsgRegex1 := regexp.MustCompile(`\["DTSGInitialData",\[\],{"token":"([^"]+)"}`)
			if m := dtsgRegex1.FindStringSubmatch(htmlText); len(m) > 1 {
				fbDtsg = m[1]
			} else {
				dtsgRegex2 := regexp.MustCompile(`"DTSGInitialData",.*?"token":"([^"]+)"`)
				if m := dtsgRegex2.FindStringSubmatch(htmlText); len(m) > 1 {
					fbDtsg = m[1]
				}
			}

			lsdRegex := regexp.MustCompile(`"LSD",\[\],{"token":"([^"]+)"}`)
			if m := lsdRegex.FindStringSubmatch(htmlText); len(m) > 1 {
				lsd = m[1]
			} else {
				lsdRegex = regexp.MustCompile(`name="lsd" value="([^"]+)"`)
				if m := lsdRegex.FindStringSubmatch(htmlText); len(m) > 1 {
					lsd = m[1]
				}
			}

			jazoestRegex := regexp.MustCompile(`name="jazoest" value="(\d+)"`)
			if m := jazoestRegex.FindStringSubmatch(htmlText); len(m) > 1 {
				jazoest = m[1]
			} else {
				jazoestRegex = regexp.MustCompile(`"jazoest":"(\d+)"`)
				if m := jazoestRegex.FindStringSubmatch(htmlText); len(m) > 1 {
					jazoest = m[1]
				}
			}

			os.WriteFile("DEBUG_PROFILE_DUMP.html", bodyBytes, 0644)
			resp.Body.Close()
		}
	}

	if fbDtsg == "" {
		return 0, errors.New("Không lấy được fb_dtsg từ trang cá nhân.")
	}

	// Bắt đầu vòng lặp lấy bạn bè
	for hasNextPage {
		// Tạo biến variable JSON sạch gọn, không dư dấu space
		var variables string
		if cursor == "" {
			variables = `{"count":20,"cursor":null,"name":""}`
		} else {
			variables = fmt.Sprintf(`{"count":20,"cursor":"%s","name":""}`, cursor)
		}

		payload := url.Values{}
		payload.Set("av", uid)
		payload.Set("__user", uid)
		payload.Set("__a", "1")
		payload.Set("fb_dtsg", fbDtsg)
		if lsd != "" {
			payload.Set("lsd", lsd)
		}
		if jazoest != "" {
			payload.Set("jazoest", jazoest)
		}
		payload.Set("doc_id", docIdToUse)
		payload.Set("variables", variables)
		payload.Set("fb_api_req_friendly_name", "FriendingCometFriendsListPaginationQuery")

		graphReq, _ := http.NewRequest("POST", "https://www.facebook.com/api/graphql/", strings.NewReader(payload.Encode()))
		graphReq.Header.Set("Cookie", cookie)
		graphReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		graphReq.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
		graphReq.Header.Set("X-FB-Friendly-Name", "FriendingCometFriendsListPaginationQuery")

		graphResp, _ := f.client.Do(graphReq)
		if graphResp == nil {
			break
		}

		bodyBytes, _ := io.ReadAll(graphResp.Body)
		graphResp.Body.Close()
		jsonStr := string(bodyBytes)

		// Thêm dòng ghi file debug ra ngoài để theo dõi
		os.WriteFile("DEBUG_FRIENDS_RESPONSE.json", bodyBytes, 0644)
		os.WriteFile("DEBUG_FRIENDS_PAYLOAD.txt", []byte(payload.Encode() + "\nfb_dtsg: " + fbDtsg), 0644)

		graphqlErrors := gjson.Get(jsonStr, "errors")
		if graphqlErrors.Exists() && graphqlErrors.IsArray() {
			errMsg := graphqlErrors.Get("0.message").String()
			return len(friendList), fmt.Errorf("GraphQL Error: %s", errMsg)
		}

		// Đường dẫn dữ liệu theo đúng cấu trúc FriendingCometFriendsListPaginationQuery
		edgesPath := "data.viewer.all_friends.edges"
		pageInfoPath := "data.viewer.all_friends.page_info"

		edges := gjson.Get(jsonStr, edgesPath)
		pageInfo := gjson.Get(jsonStr, pageInfoPath)

		// Dự phòng nếu FB trả theo username path (đôi khi đổi từ viewer sang user)
		if !edges.Exists() {
			edges = gjson.Get(jsonStr, "data.user.all_friends.edges")
			pageInfo = gjson.Get(jsonStr, "data.user.all_friends.page_info")
		}

		if edges.Exists() && edges.IsArray() {
			edges.ForEach(func(key, value gjson.Result) bool {
				fName := value.Get("node.name").String()
				if fName == "" {
					fName = value.Get("node.title.text").String() // Fallback
				}
				
				fId := value.Get("node.id").String()
				fUrl := value.Get("node.url").String()
				
				displayStr := ""
				if fId != "" {
					displayStr = fmt.Sprintf("%s | %s", fId, fName)
				} else if fUrl != "" {
					displayStr = fmt.Sprintf("%s | %s", fName, fUrl)
				}
				
				if displayStr != "" && fName != "" && !friendMap[displayStr] {
					friendMap[displayStr] = true
					friendList = append(friendList, displayStr)
				}
				return true
			})
		}

		if pageInfo.Exists() {
			hasNextPage = pageInfo.Get("has_next_page").Bool()
			if hasNextPage {
				cursor = pageInfo.Get("end_cursor").String()
				time.Sleep(2 * time.Second) // Nghỉ 2 giây để chống block
			}
		} else {
			hasNextPage = false
		}
	}

	// Ghi danh sách ra file
	path := filepath.Join("Data", uid)
	os.MkdirAll(path, os.ModePerm)

	content := fmt.Sprintf("Tổng cộng: %d bạn bè\n==============================\n", len(friendList))
	for _, f := range friendList {
		content += f + "\n"
	}
	os.WriteFile(filepath.Join(path, "Friends.txt"), []byte(content), 0644)

	return len(friendList), nil
}
