package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/gorilla/websocket"
	"github.com/pkg/browser"
)

const (
	redirectURL string = "https://vscode-redirect.azurewebsites.net/"
	clientID    string = "aebc6443-996d-45c2-90f0-388ff96faa56" // VS Code Azure Account extension
)

func dial(url string) *websocket.Conn {
	log.Println("connect:", url)
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	return c
}

// read from ws
func pumpOutput(c *websocket.Conn, w io.Writer) {
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
		}
		w.Write(message)
	}
}

// send to ws
func pumpInput(c *websocket.Conn, r io.Reader, done chan interface{}) {
	data := make([]byte, 1)
	for {
		r.Read(data)

		err := c.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			log.Println("send:", err)
		}
	}
}

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

func provision() string {
	token := getToken()
	consoleURL := createConsole(token)
	terminalURI := createTerminal(token, consoleURL)
	return terminalURI
}

func connect(url string) {
	// disable input buffering
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	// do not display entered characters on the screen
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()

	done := make(chan interface{})

	c := dial(url)
	defer c.Close()

	go pumpOutput(c, os.Stdout)
	go pumpInput(c, os.Stdin, done)

	<-done
}

func main() {
	log.SetFlags(0)
	log.SetOutput(ioutil.Discard)

	url := provision()
	connect(url)
}
