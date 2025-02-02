package tools

import (
	"fmt"
	"github.com/gen2brain/beeep"
	"net/http"
	"net/url"
	"strings"
)

func SendDesktopNotification(title, message string, verbose bool) error {
	if verbose {
		fmt.Println("Sending desktop notification")
		fmt.Println("Title: ", title)
		fmt.Println("Message: ", message)
	}
	return beeep.Notify(title, message, "")
}

func SendTelegramNotification(token, chatID, message string, telegramCustomAPI string, verbose bool) error {
	baseURL := "https://api.telegram.org"
	if telegramCustomAPI != "" {
		customAPI := telegramCustomAPI
		if !strings.HasPrefix(customAPI, "http://") && 
		   !strings.HasPrefix(customAPI, "https://") {
			customAPI = "https://" + customAPI
		}
		baseURL = customAPI
	}

	if verbose {
		fmt.Println("Sending Telegram notification")
		maskedToken := fmt.Sprintf("********%s", token[3:])
		fmt.Println("Token: ", maskedToken)
		fmt.Println("Chat ID: ", chatID)
		fmt.Println("Message: ", message)
		fmt.Println("Using API URL: ", baseURL)
	}

	apiURL := fmt.Sprintf("%s/bot%s/sendMessage", baseURL, token)
	params := url.Values{}
	params.Add("chat_id", chatID)
	params.Add("text", message)

	if verbose {
		fmt.Println("Sending request to: ", apiURL)
	}
	resp, err := http.PostForm(apiURL, params)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if verbose {
		fmt.Println("Telegram response: ", resp.Status)
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("Telegram response: %s", resp.Status)
	}
	return nil
}