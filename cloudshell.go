package main

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

	resp, _ := client.Do(req)
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

func createTerminal(token string, consoleURL string) string {
	fmt.Println("Connecting terminal...")
	url := fmt.Sprintf("%s/terminals?cols=80&rows=30&shell=bash", consoleURL)
	data := `{"tokens": []}`
	result := sendRequest(token, "POST", url, data)
	return result["socketUri"].(string)
}

func Provision(token string) string {
	consoleURL := createConsole(token)
	terminalURI := createTerminal(token, consoleURL)
	return terminalURI
}
