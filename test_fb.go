package main
import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"regexp"
)
func main() {
	cookie := "datr=6M1JaPhd1g6T7V1vZuvpmKC0;xs=1%3A0FwF4gkw-dE3cw%3A2%3A1776045645%3A-1%3A-1;sb=Q07cafyMWDeguux1Y8IExiSZ;fr=1xDujY5gF4655mvfN.AWfs8C-gp2bpncuY2lxPJ3mDzYxyduoF6mR7BaUM96ko8B9pJKw.BobKSj..AAA.0.0.Bp3E5E.AWeSbM-PkJXpFO10lVymdjbJwdY;c_user=61571360758847;pas=61571360758847%3AhAIRYYu4Yh;locale=en_US"
	reqGet, _ := http.NewRequest("GET", "https://www.facebook.com/", nil)
	reqGet.Header.Set("Cookie", cookie)
	reqGet.Header.Set("User-Agent", "Mozilla/5.0")
	respGet, _ := http.DefaultClient.Do(reqGet)
	bodyBytes, _ := io.ReadAll(respGet.Body)
	respGet.Body.Close()
	htmlText := string(bodyBytes)
	fbDtsg := ""
	if m := regexp.MustCompile(`\["DTSGInitialData",\[\],{"token":"([^"]+)"}`).FindStringSubmatch(htmlText); len(m) > 1 {
		fbDtsg = m[1]
	}
	fmt.Println("DTSG: ", fbDtsg)
	
	uid := "61571360758847"

	// test composer photo/upload endpoint with video extension
	url := fmt.Sprintf("https://upload.facebook.com/ajax/react_composer/attachments/photo/upload?__a=1&__req=1&__user=%s", uid)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("fb_dtsg", fbDtsg)
	
	// mock file Upload
	part, _ := writer.CreateFormFile("farr", "test.mp4")
	part.Write([]byte("fake video data fake video data fake video data"))
	writer.Close()
	
	req, _ := http.NewRequest("POST", url, body)
	req.Header.Set("Cookie", cookie)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("User-Agent", "Mozilla/5.0")
	resp, _ := http.DefaultClient.Do(req)
	r, _ := io.ReadAll(resp.Body)
	fmt.Println(resp.StatusCode)
	fmt.Println(string(r))
}
