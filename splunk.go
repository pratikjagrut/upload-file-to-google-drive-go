package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"
)

var (
	ACCESS_TOKEN  string = os.Getenv("ACCESS_TOKEN")
	SPLUNK_INDEX  string = os.Getenv("SPLUNK_INDEX")
	SPLUNK_URL    string = fmt.Sprintf("%s/services/collector", os.Getenv("SPLUNK_URL"))
	SPLUNK_TOKEN  string = os.Getenv("SPLUNK_TOKEN")
	SPLUNK_SOURCE string = os.Getenv("SPLUNK_SOURCE")
)

func read_content(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return data, err
}

func parse_traffic_for_count(traffic []byte) int {
	re := regexp.MustCompile(`HTTP`)
	return len(re.FindAllString(string(traffic), -1))
}

func UploadDataToSplunk(fileUploadStatus string) error {
	splunk_data := make(map[string]string)
	splunk_data["Exfiltration Channel"] = "google drive"
	splunk_data["Data Content"] = "credit card"
	splunk_data["Test Date"] = time.Now().String()
	splunk_data["Test Status"] = fileUploadStatus
	data, err := os.ReadFile(UPLOAD_FILE)
	if err != nil {
		return err
	}
	splunk_data["Count"] = strconv.Itoa(parse_traffic_for_count(data))

	splunk_data_json, err := json.Marshal(splunk_data)
	if err != nil {
		return err
	}
	// fmt.Println(splunk_data_json)
	splunk_auth_header := fmt.Sprintf("Splunk %s", SPLUNK_TOKEN)
	http_content := make(map[string]string)
	http_content["index"] = SPLUNK_INDEX
	http_content["event"] = string(splunk_data_json)
	http_content["source"] = SPLUNK_SOURCE
	http_json, err := json.Marshal(http_content)
	if err != nil {
		return err
	}
	request, err := http.NewRequest("POST", SPLUNK_URL, bytes.NewBuffer(http_json))
	if err != nil {
		return err
	}
	request.Header.Set("Authorization", splunk_auth_header)
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	resp.Body.Close()
	fmt.Println("Splunk ingestion complete")
	return nil
}
