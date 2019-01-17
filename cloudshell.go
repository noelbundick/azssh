package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/pkg/browser"
)

const (
	redirectURL string = "https://vscode-redirect.azurewebsites.net/"
	clientID    string = "aebc6443-996d-45c2-90f0-388ff96faa56" // VS Code Azure Account extension
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
	data := `{"properties": { "osType": "linux" } }`
	result := sendRequest(token, "PUT", "https://management.azure.com/providers/Microsoft.Portal/consoles/default?api-version=2017-08-01-preview", data)
	properties := result["properties"].(map[string]interface{})
	return properties["uri"].(string)
}

func createTerminal(token string, consoleURL string) string {
	url := fmt.Sprintf("%s/terminals?cols=80&rows=30&shell=bash", consoleURL)
	data := `{"tokens": []}`
	result := sendRequest(token, "POST", url, data)
	return result["socketUri"].(string)
}

func getToken() string {
	config, err := adal.NewOAuthConfig("https://login.microsoftonline.com/", "common")
	if err != nil {
		log.Fatal(err)
	}

	sender := &http.Client{}
	code, err := adal.InitiateDeviceAuth(sender, *config, clientID, "https://management.core.windows.net/")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("code:", *code.Message)
	browser.OpenURL(*code.VerificationURL)

	token, err := adal.WaitForUserCompletion(sender, code)
	if err != nil {
		log.Fatal(err)
	}
	return token.OAuthToken()
}

func Provision() string {
	token := getToken()
	consoleURL := createConsole(token)
	terminalURI := createTerminal(token, consoleURL)
	return terminalURI
}
