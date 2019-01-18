package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"

	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/pkg/browser"
)

const (
	redirectURL string = "https://vscode-redirect.azurewebsites.net/"
	clientID    string = "aebc6443-996d-45c2-90f0-388ff96faa56" // VS Code Azure Account extension
)

func getTokenCachePath() string {
	u, err := user.Current()
	if err != nil {
		log.Println(err)
		return ""
	}

	return fmt.Sprintf("%s/.azssh/token.json", u.HomeDir)
}

func saveToken(token adal.Token, tokenPath string) error {
	tokenPathDir := filepath.Dir(tokenPath)
	err := os.MkdirAll(tokenPathDir, os.ModePerm)
	if err != nil {
		log.Println(err)
	}

	return adal.SaveToken(tokenPath, 0600, token)
}

func GetToken() string {
	tokenPath := getTokenCachePath()
	config, err := adal.NewOAuthConfig("https://login.microsoftonline.com/", "common")
	if err != nil {
		log.Fatal(err)
	}

	refreshCallback := func(token adal.Token) error {
		return saveToken(token, tokenPath)
	}

	token, err := adal.LoadToken(tokenPath)
	if err != nil {
		log.Println(err)
	}

	if token != nil {
		if token.IsExpired() {
			spt, err := adal.NewServicePrincipalTokenFromManualToken(*config, clientID, "https://management.core.windows.net/", *token, refreshCallback)
			if err != nil {
				log.Println("token refresh failure:", err)
			}
			err = spt.Refresh()
			if err != nil {
				log.Println("token refresh failure:", err)
			}

			spToken := spt.Token()
			token = &spToken
		}

		return token.OAuthToken()
	}

	sender := &http.Client{}
	code, err := adal.InitiateDeviceAuth(sender, *config, clientID, "https://management.core.windows.net/")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(*code.Message)
	browser.OpenURL(*code.VerificationURL)

	token, err = adal.WaitForUserCompletion(sender, code)
	if err != nil {
		log.Fatal(err)
	}

	spt, err := adal.NewServicePrincipalTokenFromManualToken(*config, clientID, "https://management.core.windows.net/", *token, nil)

	err = saveToken(spt.Token(), tokenPath)
	if err != nil {
		log.Println(err)
	}

	return spt.OAuthToken()
}
