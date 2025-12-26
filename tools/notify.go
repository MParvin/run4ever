package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"
	"net/url"
	"strings"

	"github.com/gen2brain/beeep"
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
		return fmt.Errorf("telegram response: %s", resp.Status)
	}
	return nil
}

// SendSlackNotification sends a notification to Slack
func SendSlackNotification(webhookURL, message string, verbose bool) error {
	if verbose {
		fmt.Println("Sending Slack notification")
		fmt.Println("Message: ", message)
	}

	payload := map[string]string{
		"text": message,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal Slack payload: %w", err)
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send Slack notification: %w", err)
	}
	defer resp.Body.Close()

	if verbose {
		fmt.Println("Slack response: ", resp.Status)
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("slack response: %s", resp.Status)
	}
	return nil
}

// SendEmailNotification sends an email notification
func SendEmailNotification(to, from, password, smtpHost string, smtpPort int, subject, message string, verbose bool) error {
	if verbose {
		fmt.Println("Sending email notification")
		fmt.Println("To: ", to)
		fmt.Println("From: ", from)
		fmt.Println("SMTP: ", smtpHost, ":", smtpPort)
	}

	// Setup authentication
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Compose email
	emailBody := fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s\r\n", to, subject, message)

	// Send email
	addr := fmt.Sprintf("%s:%d", smtpHost, smtpPort)
	err := smtp.SendMail(addr, auth, from, []string{to}, []byte(emailBody))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	if verbose {
		fmt.Println("Email sent successfully")
	}
	return nil
}
