package azssh

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func sendRequest(token string, method string, url string, payload string) map[string]interface{} {
	client := &http.Client{}

	req, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("url:", url)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	log.Println("body:", string(body))

	var result map[string]interface{}
	json.Unmarshal(body, &result)
	return result
}

func createConsole(token string) string {
	fmt.Println("Requesting a Cloud Shell.")
	data := `{"properties": { "osType": "linux" } }`
	result := sendRequest(token, "PUT", "https://management.azure.com/providers/Microsoft.Portal/consoles/default?api-version=2017-08-01-preview", data)
	properties := result["properties"].(map[string]interface{})
	return properties["uri"].(string)
}

func createTerminal(token string, consoleURL string, shellType string, initialSize TerminalSize) (string, string) {
	fmt.Println("Connecting terminal...")
	url := fmt.Sprintf("%s/terminals?cols=%d&rows=%d&shell=%s", consoleURL, initialSize.Cols, initialSize.Rows, shellType)
	data := `{"tokens": []}`
	result := sendRequest(token, "POST", url, data)
	return result["id"].(string), result["socketUri"].(string)
}

func resizeTerminal(token string, consoleURL string, terminalID string, resize <-chan TerminalSize) {
	for {
		newSize := <-resize
		url := fmt.Sprintf("%s/terminals/%s/size?cols=%d&rows=%d", consoleURL, terminalID, newSize.Cols, newSize.Rows)
		sendRequest(token, "POST", url, "")
	}
}

// ProvisionCloudShell sets up a Cloud Shell and a websocket to connect into it
func ProvisionCloudShell(token string, shellType string, initialSize TerminalSize, resize <-chan TerminalSize) string {
	consoleURL := createConsole(token)
	terminalID, websocketURI := createTerminal(token, consoleURL, shellType, initialSize)
	go resizeTerminal(token, consoleURL, terminalID, resize)
	return websocketURI
}
