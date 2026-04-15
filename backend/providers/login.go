package providers

// login.go — Module đăng nhập bằng User/Pass cho Facebook
//
// Triển khai 3 thành phần kỹ thuật:
//   1. encryptPassword  — Mã hóa mật khẩu dùng NaCl/XSalsa20 + định dạng #PWD_BROWSER:5
//   2. InitializeSession — GET www.facebook.com để bóc tách lsd / jazoest / fb_dtsg
//   3. SubmitLogin       — POST GraphQL doc_id=25845944448416411 để lấy session

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	mrand "math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/tidwall/gjson"
	"golang.org/x/crypto/nacl/box"
	"socialmanager/backend/models"
)

// ─────────────────────────────────────────────
// 1. MODULE MÃ HÓA (Cryptography)
// ─────────────────────────────────────────────

func encryptPassword(password string, serverPubKeyB64 string, keyID int, timestamp int64) (string, error) {
	serverPubBytes, err := base64.StdEncoding.DecodeString(serverPubKeyB64)
	if err != nil {
		serverPubBytes, err = base64.URLEncoding.DecodeString(serverPubKeyB64)
		if err != nil {
			return "", fmt.Errorf("giải mã publicKey thất bại: %w", err)
		}
	}
	if len(serverPubBytes) != 32 {
		return "", fmt.Errorf("publicKey phải đúng 32 bytes, nhận %d", len(serverPubBytes))
	}

	var serverPub [32]byte
	copy(serverPub[:], serverPubBytes)

	ephemeralPub, ephemeralPriv, err := box.GenerateKey(rand.Reader)
	if err != nil {
		return "", fmt.Errorf("sinh ephemeral key lỗi: %w", err)
	}

	var nonce [24]byte
	if _, err = io.ReadFull(rand.Reader, nonce[:]); err != nil {
		return "", fmt.Errorf("sinh nonce lỗi: %w", err)
	}

	tsStr := fmt.Sprintf("%d", timestamp)
	plaintext := []byte(tsStr + password)
	encrypted := box.Seal(nil, plaintext, &nonce, &serverPub, ephemeralPriv)

	var payload []byte
	payload = append(payload, 1) // version
	keyIDBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(keyIDBytes, uint16(keyID))
	payload = append(payload, keyIDBytes...)
	payload = append(payload, nonce[:]...)
	payload = append(payload, ephemeralPub[:]...)
	payload = append(payload, encrypted...)

	encodedPayload := base64.StdEncoding.EncodeToString(payload)
	result := fmt.Sprintf("#PWD_BROWSER:5:%d:%s", timestamp, encodedPayload)
	return result, nil
}

// ─────────────────────────────────────────────
// 2. MODULE THU THẬP SESSION (Scraping)
// ─────────────────────────────────────────────

type SessionData struct {
	LSD       string
	Jazoest   string
	FbDtsg    string
	PublicKey string
	KeyID     int
	Jar       http.CookieJar
	Client    *http.Client
}

func InitializeSession() (*SessionData, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("tạo cookie jar lỗi: %w", err)
	}

	client := &http.Client{
		Jar:     jar,
		Timeout: 30 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 5 {
				return errors.New("quá nhiều redirect")
			}
			for key, vals := range via[0].Header {
				req.Header[key] = vals
			}
			return nil
		},
	}

	req, err := http.NewRequest("GET", "https://www.facebook.com/", nil)
	if err != nil {
		return nil, fmt.Errorf("tạo GET request lỗi: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
	req.Header.Set("Sec-Ch-Ua", `"Chromium";v="124", "Google Chrome";v="124", "Not-A.Brand";v="99"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Windows"`)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "vi-VN,vi;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GET facebook.com lỗi: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("đọc response body lỗi: %w", err)
	}
	html := string(bodyBytes)

	sd := &SessionData{
		Jar:    jar,
		Client: client,
	}

	if m := regexp.MustCompile(`name="lsd"\s+value="([^"]+)"`).FindStringSubmatch(html); len(m) > 1 {
		sd.LSD = m[1]
	} else if m := regexp.MustCompile(`"LSD",\[\],\{"token":"([^"]+)"\}`).FindStringSubmatch(html); len(m) > 1 {
		sd.LSD = m[1]
	}

	if m := regexp.MustCompile(`name="jazoest"\s+value="([^"]+)"`).FindStringSubmatch(html); len(m) > 1 {
		sd.Jazoest = m[1]
	} else if m := regexp.MustCompile(`jazoest=(\d+)`).FindStringSubmatch(html); len(m) > 1 {
		sd.Jazoest = m[1]
	}

	pubKeyPattern := regexp.MustCompile(`"publicKey"\s*:\s*"([^"]+)"`)
	keyIDPattern := regexp.MustCompile(`"keyId"\s*:\s*(\d+)`)

	if m := pubKeyPattern.FindStringSubmatch(html); len(m) > 1 {
		sd.PublicKey = m[1]
	}
	if m := keyIDPattern.FindStringSubmatch(html); len(m) > 1 {
		fmt.Sscanf(m[1], "%d", &sd.KeyID)
	}

	if sd.PublicKey == "" {
		sd.PublicKey = "sMKRELYqYaPhonCc8UQ73oilhJaV3TzWCrJk/Kv+1fU="
		sd.KeyID = 7
	}

	if sd.LSD == "" {
		os.WriteFile("debug_login_init.html", bodyBytes, 0644)
		return nil, errors.New("không bóc tách được token LSD từ fb. Đã lưu debug_login_init.html")
	}

	return sd, nil
}

// ─────────────────────────────────────────────
// 3. MODULE THỰC THI
// ─────────────────────────────────────────────

func SubmitLogin(sd *SessionData, identifier string, encPassword string, _docId string) (string, string, error) {
	data := url.Values{}
	data.Set("lsd", sd.LSD)
	data.Set("jazoest", sd.Jazoest)
	data.Set("email", identifier)
	data.Set("encpass", encPassword)
	data.Set("login_source", "comet_header_shortwave")
	data.Set("next", "")

	req, err := http.NewRequest("POST", "https://www.facebook.com/login/device-based/regular/login/", strings.NewReader(data.Encode()))
	if err != nil {
		return "", "", fmt.Errorf("tạo request thất bại: %v", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
	req.Header.Set("Sec-Ch-Ua", `"Chromium";v="124", "Google Chrome";v="124", "Not-A.Brand";v="99"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Windows"`)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "vi-VN,vi;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Origin", "https://www.facebook.com")
	req.Header.Set("Referer", "https://www.facebook.com/")

	resp, err := sd.Client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("gửi POST thất bại: %v", err)
	}
	defer resp.Body.Close()

	finalURL := resp.Request.URL.String()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("đọc response body thất bại: %v", err)
	}

	return string(bodyBytes), finalURL, nil
}

// SubmitGraphQLLogin gửi GraphQL mutation cho Auth
func SubmitGraphQLLogin(sd *SessionData, identifier string, encPassword string, docId string) (string, string, error) {
	variables := fmt.Sprintf(
		`{"input":{"credentials_type":"password","error_detail_type":"button_with_disabled","source":"login","password":"%s","login_attempt_count":0,"identifier":"%s","lsd":"%s","jazoest":"%s"},"scale":1.5,"is_two_factor_auth_in_app_browser":false,"is_two_factor_auth_in_app_browser_webview":false,"is_two_factor_auth_in_app_browser_webview_bypass":false}`,
		escapeJSON(encPassword), escapeJSON(identifier), escapeJSON(sd.LSD), escapeJSON(sd.Jazoest),
	)

	data := url.Values{}
	data.Set("doc_id", docId)
	data.Set("variables", variables)
	data.Set("fb_dtsg", sd.FbDtsg)
	data.Set("lsd", sd.LSD)
	data.Set("jazoest", sd.Jazoest)

	req, err := http.NewRequest("POST", "https://www.facebook.com/api/graphql/", strings.NewReader(data.Encode()))
	if err != nil {
		return "", "", fmt.Errorf("tạo GraphQL request thất bại: %v", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
	req.Header.Set("Sec-Ch-Ua", `"Chromium";v="124", "Google Chrome";v="124", "Not-A.Brand";v="99"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Windows"`)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Origin", "https://www.facebook.com")
	req.Header.Set("Referer", "https://www.facebook.com/")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")

	resp, err := sd.Client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("gửi GraphQL POST thất bại: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	return string(bodyBytes), resp.Request.URL.String(), nil
}

// escapeJSON escape chuỗi để an toàn trong JSON string value
func escapeJSON(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	s = strings.ReplaceAll(s, "\n", `\n`)
	s = strings.ReplaceAll(s, "\r", `\r`)
	s = strings.ReplaceAll(s, "\t", `\t`)
	return s
}

// ─────────────────────────────────────────────
// 3.5 MODULE ĐĂNG NHẬP ANDROID APP API (BYPASS)
// ─────────────────────────────────────────────

func SubmitAndroidLogin(identifier string, password string) ([]*http.Cookie, string, error) {
	apiKey := "882a8490361da98702bf97a021ddc14d"
	apiSecret := "62f8ce9f74b12f84c123cc23437a4a32"

	machineId := fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", mrand.Int31(), mrand.Int31(), mrand.Int31(), mrand.Int31(), mrand.Int63())
	deviceId := fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", mrand.Int31(), mrand.Int31(), mrand.Int31(), mrand.Int31(), mrand.Int63())

	data := url.Values{}
	data.Set("api_key", apiKey)
	data.Set("credentials_type", "password")
	data.Set("email", identifier)
	data.Set("format", "JSON")
	data.Set("generate_machine_id", "1")
	data.Set("generate_session_cookies", "1")
	data.Set("locale", "en_US")
	data.Set("method", "auth.login")
	data.Set("password", password)
	data.Set("return_multiple_errors", "true")
	data.Set("device_id", deviceId)
	data.Set("machine_id", machineId)

	// Sort keys for signature
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	// Bubble sort keys or simple sort
	for i := 0; i < len(keys)-1; i++ {
		for j := i + 1; j < len(keys); j++ {
			if keys[i] > keys[j] {
				keys[i], keys[j] = keys[j], keys[i]
			}
		}
	}

	sigData := ""
	for _, k := range keys {
		sigData += fmt.Sprintf("%s=%s", k, data.Get(k))
	}
	sigData += apiSecret

	hash := md5.Sum([]byte(sigData))
	sig := fmt.Sprintf("%x", hash)
	data.Set("sig", sig)

	req, err := http.NewRequest("POST", "https://b-api.facebook.com/method/auth.login", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, "", fmt.Errorf("không thể tạo request Android App: %w", err)
	}

	req.Header.Set("User-Agent", "[FBAN/FB4A;FBAV/417.0.0.33.65;FBBV/480086274;FBDM/{density=3.0,width=1080,height=2132};FBLC/en_US;FBRV/481977799;FBCR/Viettel;FBMF/Xiaomi;FBBD/xiaomi;FBPN/com.facebook.katana;FBDV/M2007J20CG;FBSV/11;FBOP/1;FBCA/arm64-v8a:;]")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("lỗi kết nối API: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	return nil, string(bodyBytes), nil
}


// ─────────────────────────────────────────────
// 4. ORCHESTRATOR — LoginWithPassword
// ─────────────────────────────────────────────

// LoginWithPassword thực hiện quy trình đăng nhập
func LoginWithPassword(identifier string, password string, docId string) (models.LoginResult, error) {

	// 1. NGHIỆP VỤ B-API ANDROID CHÍNH GỐC (Bypass WEB Shadow-ban bằng UID)
	_, rawResponse, err := SubmitAndroidLogin(identifier, password)
	if err == nil && rawResponse != "" {
		j := gjson.Parse(rawResponse)
		
		if j.Get("session_cookies").Exists() {
			var uid, cookieFull string
			var cookieStrings []string
			
			for _, c := range j.Get("session_cookies").Array() {
				name := c.Get("name").String()
				val := c.Get("value").String()
				cookieStrings = append(cookieStrings, name+"="+val)
				if name == "c_user" {
					uid = val
				}
			}
			cookieFull = strings.Join(cookieStrings, "; ")
			
			if uid != "" {
				return models.LoginResult{
					UID:        uid,
					CookieFull: cookieFull,
					Status:     "Live",
				}, nil
			}
		}

		errCode := j.Get("error_code").Int()
		if errCode != 0 {
			msg := j.Get("error_msg").String()
			lowBody := strings.ToLower(msg)
			
			// 405 : Checkpoint / 2FA Verify
			if errCode == 405 {
				if strings.Contains(lowBody, "two-factor") || strings.Contains(lowBody, "approvals_code") || strings.Contains(lowBody, "2fa") || strings.Contains(lowBody, "verify their account") {
					return models.LoginResult{Status: "2FARequired", RawResponse: msg}, errors.New("tài khoản yêu cầu mã bảo mật 2 lớp (2FA). Vui lòng thêm 2FA vào hệ thống")
				}
				return models.LoginResult{Status: "Checkpoint", RawResponse: msg}, errors.New("tài khoản dính Checkpoint (cần xác minh danh tính trên trình duyệt thực)")
			}
			// 401/400: Sai mật khẩu
			if errCode == 401 || errCode == 400 || strings.Contains(lowBody, "invalid username") {
				return models.LoginResult{Status: "WrongPassword", RawResponse: msg}, errors.New("sai tài khoản / mật khẩu, hoặc tài khoản đã bị đổi pass")
			}
			
			// Lỗi khác
			os.WriteFile("debug_login_response_bapi.json", []byte(rawResponse), 0644)
			return models.LoginResult{Status: "Failed", RawResponse: msg}, fmt.Errorf("API từ chối truy cập. Lỗi: %s", msg)
		}
	}

	// 2. FALLBACK WEB GRAPHQL TRUYỀN THỐNG (Khởi tạo phiên nếu B-API thất bại vì kết nối)
	sd, err := InitializeSession()
	if err != nil {
		return models.LoginResult{}, fmt.Errorf("khởi tạo session thất bại: %w", err)
	}

	timestamp := time.Now().Unix()
	encPwd, err := encryptPassword(password, sd.PublicKey, sd.KeyID, timestamp)
	if err != nil {
		return models.LoginResult{}, fmt.Errorf("mã hóa mật khẩu thất bại: %w", err)
	}

	rawBody, finalUrl, err := SubmitGraphQLLogin(sd, identifier, encPwd, docId)
	if err != nil {
		rawBody, finalUrl, err = SubmitLogin(sd, identifier, encPwd, docId)
		if err != nil {
			return models.LoginResult{}, fmt.Errorf("gửi login request thất bại: %w", err)
		}
	}

	fbURL, _ := url.Parse("https://www.facebook.com")
	cookies := sd.Jar.Cookies(fbURL)

	uid := ""
	cookieParts := make([]string, 0, len(cookies))
	for _, c := range cookies {
		cookieParts = append(cookieParts, c.Name+"="+c.Value)
		if c.Name == "c_user" {
			uid = c.Value
		}
	}

	cookieFull := strings.Join(cookieParts, "; ")
	if uid != "" {
		return models.LoginResult{
			UID:        uid,
			CookieFull: cookieFull,
			Status:     "Live",
		}, nil
	}

	lowBody := strings.ToLower(rawBody)
	if strings.Contains(lowBody, "sai tài khoản") ||
		strings.Contains(lowBody, "incorrect") ||
		strings.Contains(lowBody, "you entered is incorrect") ||
		strings.Contains(lowBody, "không kết nối với tài khoản nào") ||
		strings.Contains(lowBody, "isn't connected to an account") ||
		strings.Contains(lowBody, "không tìm thấy tài khoản") {
		return models.LoginResult{Status: "WrongPassword"}, errors.New("sai tài khoản / mật khẩu, hoặc hệ thống chặn vì dùng UID. Hãy thử UID qua App API.")
	}

	if strings.Contains(strings.ToLower(finalUrl), "checkpoint") || strings.Contains(lowBody, "bạn đã thử quá nhiều lần") {
		return models.LoginResult{Status: "Checkpoint"}, errors.New("tài khoản bị checkpoint — cần xác minh bảo mật trên trình duyệt")
	}

	if strings.Contains(strings.ToLower(finalUrl), "two_step_verification") || strings.Contains(strings.ToLower(finalUrl), "approvals_code") {
		return models.LoginResult{Status: "2FARequired"}, errors.New("tài khoản yêu cầu nhập mã 2FA. Vui lòng thêm key 2FA vào cấu hình")
	}

	if uid == "" {
		if m := regexp.MustCompile(`"c_user"\s*:\s*"?(\d+)"?`).FindStringSubmatch(rawBody); len(m) > 1 {
			uid = m[1]
		}
		if uid == "" {
			if m := regexp.MustCompile(`"session_key"\s*:\s*"([^"]+)"`).FindStringSubmatch(rawBody); len(m) > 1 {
				uid = m[1]
			}
		}
	}

	cookieFull = strings.Join(cookieParts, "; ")

	if uid == "" || cookieFull == "" {
		os.WriteFile("debug_login_response.json", []byte(rawBody), 0644)
		return models.LoginResult{Status: "Failed", RawResponse: truncate(rawBody, 500)},
			errors.New("đăng nhập Web không thành công và API bị từ chối")
	}

	return models.LoginResult{
		UID:        uid,
		CookieFull: cookieFull,
		Status:     "Live",
	}, nil
}

// truncate cắt chuỗi nếu quá dài (để dùng trong debug message)
func truncate(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen]) + "..."
}
