package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func main() {
	docId := "25845944448416411"
	
	// Test 1: Just empty variables
	data := url.Values{}
	data.Set("doc_id", docId)
	data.Set("variables", "{}")
	
	req, _ := http.NewRequest("POST", "https://www.facebook.com/api/graphql/", strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer res.Body.Close()
	
	body, _ := io.ReadAll(res.Body)
	fmt.Println("Response with empty variables:")
	
	var prettyJSON map[string]interface{}
	json.Unmarshal(body, &prettyJSON)
	pretty, _ := json.MarshalIndent(prettyJSON, "", "  ")
	fmt.Println(string(pretty))
}
