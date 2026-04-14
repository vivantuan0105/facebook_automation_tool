package providers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

func (f *FacebookProvider) UploadVideoFB(cookie, filePath string) (string, error) {
	// The standard react_composer attachment upload endpoint natively handles both images and videos (short-medium).
	// We use it as Phase 1 (Single Phase) instead of the broken /video-upload URL.
	
	reqGet, err := http.NewRequest("GET", "https://www.facebook.com/", nil)
	if err != nil {
		return "", fmt.Errorf("lỗi tạo request GET: %v", err)
	}
	reqGet.Header.Set("Cookie", cookie)
	reqGet.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
	reqGet.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml")

	respGet, err := f.client.Do(reqGet)
	if err != nil {
		return "", fmt.Errorf("lỗi kết nối Facebook homepage: %v", err)
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

	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("không thể mở file video: %v", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("fb_dtsg", fbDtsg)

	part, err := writer.CreateFormFile("farr", filepath.Base(filePath))
	if err != nil {
		return "", fmt.Errorf("lỗi tạo form file: %v", err)
	}

	io.Copy(part, file)
	writer.Close()

	reqUpload, err := http.NewRequest("POST", uploadURL, body)
	if err != nil {
		return "", fmt.Errorf("lỗi tạo request upload: %v", err)
	}
	reqUpload.Header.Set("Content-Type", writer.FormDataContentType())
	reqUpload.Header.Set("Cookie", cookie)
	reqUpload.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")

	respUpload, err := f.client.Do(reqUpload)
	if err != nil {
		return "", fmt.Errorf("lỗi kết nối upload server: %v", err)
	}
	defer respUpload.Body.Close()

	upBytes, _ := io.ReadAll(respUpload.Body)
	upText := string(upBytes)

	if respUpload.StatusCode != 200 {
		return "", fmt.Errorf("Upload server trả về %d: %s", respUpload.StatusCode, upText)
	}

	fbidRegex := regexp.MustCompile(`"fbid":"(\d+)"`)
	m := fbidRegex.FindStringSubmatch(upText)
	if len(m) > 1 {
		return m[1], nil // Target FBID to be used as video_id
	}

	idRegex := regexp.MustCompile(`"video_id":"(\d+)"`)
	if m := idRegex.FindStringSubmatch(upText); len(m) > 1 {
		return m[1], nil
	}

	return "", fmt.Errorf("không tìm thấy fbid/video_id trong response JSON: %s", upText)
}

func createVideoPostPayload(actorId, message, videoID string) map[string]interface{} {
	attachments := []interface{}{}
	if videoID != "" {
		attachments = append(attachments, map[string]interface{}{
			"video": map[string]interface{}{
				"id": videoID,
				"notify_when_processed": true,
				"was_created_via_unified_video_flow": map[string]bool{
					"was_created_via_unified_video_flow": true,
				},
			},
		})
	}

	idempotencyToken := generatePseudoUUID() + "_FEED"
	
	return map[string]interface{}{
		"input": map[string]interface{}{
			"actor_id":                 actorId,
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

func (f *FacebookProvider) PostVideoToFacebook(cookie string, message string, videoID string, docId string) (string, error) {
	if docId == "" {
		return "", errors.New("CHƯA CẤU HÌNH DOC_ID ĐĂNG BÀI! Vui lòng vào Cài đặt (Settings) và điền dãy số GraphQL Post Doc_ID.")
	}

	req, err := http.NewRequest("GET", "https://www.facebook.com/", nil)
	if err != nil { return "", err }
	req.Header.Set("Cookie", cookie)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml")

	resp, err := f.client.Do(req)
	if err != nil { return "", err }
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	htmlText := string(bodyBytes)

	fbDtsg := ""
	dtsgRegex1 := regexp.MustCompile(`\["DTSGInitialData",\[\],{"token":"([^"]+)"}`)
	if m := dtsgRegex1.FindStringSubmatch(htmlText); len(m) > 1 {
		fbDtsg = m[1]
	}

	actorId := "0"
	actorRegex := regexp.MustCompile(`"USER_ID":"(\d+)"`)
	if m := actorRegex.FindStringSubmatch(htmlText); len(m) > 1 {
		actorId = m[1]
	}

	variables := createVideoPostPayload(actorId, message, videoID)
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
	if err != nil { return "", err }
	defer res.Body.Close()

	fbBodyBytes, _ := io.ReadAll(res.Body)
	fbBody := string(fbBodyBytes)

	if res.StatusCode != 200 {
		return "", fmt.Errorf("GraphQL HTTP %d: %s", res.StatusCode, fbBody)
	}

	return "Đăng video thành công! (FB đang xử lý ngầm tại Background)", nil
}
