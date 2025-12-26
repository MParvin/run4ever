package tools

import (
	"testing"
)

func TestShouldNotify(t *testing.T) {
	tests := []struct {
		name       string
		notifyOn   string
		exitStatus int
		expected   bool
	}{
		{"always true", "always", 0, true},
		{"always true on failure", "always", 1, true},
		{"success on success", "success", 0, true},
		{"success on failure", "success", 1, false},
		{"failure on success", "failure", 0, false},
		{"failure on failure", "failure", 1, true},
		{"empty notifyOn", "", 0, false},
		{"invalid notifyOn", "invalid", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := shouldNotify(tt.notifyOn, tt.exitStatus)
			if result != tt.expected {
				t.Errorf("shouldNotify(%s, %d) = %v, want %v", tt.notifyOn, tt.exitStatus, result, tt.expected)
			}
		})
	}
}

func TestStatusToString(t *testing.T) {
	tests := []struct {
		exitStatus int
		expected   string
	}{
		{0, "Success"},
		{1, "Failure"},
		{2, "Failure"},
		{255, "Failure"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := statusToString(tt.exitStatus)
			if result != tt.expected {
				t.Errorf("statusToString(%d) = %s, want %s", tt.exitStatus, result, tt.expected)
			}
		})
	}
}

// Test SendDesktopNotification - this will only work if desktop notifications are available
func TestSendDesktopNotification(t *testing.T) {
	// This test will be skipped if desktop notifications are not available
	// or if running in a headless environment
	t.Skip("Skipping desktop notification test - requires GUI environment")
}

// Test SendTelegramNotification with invalid token (skip network tests)
func TestSendTelegramNotificationInvalid(t *testing.T) {
	t.Skip("Skipping network-dependent test")
	// This would make a real network request which can timeout
	// err := SendTelegramNotification("invalid_token", "123456", "test message", "", true)
	// if err == nil {
	// 	t.Error("Expected error with invalid token")
	// }
}

// Test SendTelegramNotification with invalid chat ID (skip network tests)
func TestSendTelegramNotificationInvalidChatID(t *testing.T) {
	t.Skip("Skipping network-dependent test")
}

// Test SendTelegramNotification with custom API (skip network tests)
func TestSendTelegramNotificationCustomAPI(t *testing.T) {
	t.Skip("Skipping network-dependent test")
}

// Test doNotify function
func TestDoNotify(t *testing.T) {
	// Test desktop notification path (will be skipped in headless environments)
	t.Run("desktop", func(t *testing.T) {
		// This will attempt to send a desktop notification
		// In a headless environment, it might fail gracefully
		// doNotify doesn't return an error, so we just test it doesn't panic
		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("doNotify panicked: %v", r)
				}
			}()
			doNotify("always", "desktop", false, "Test Title", "Test Message")
		}()
	})

	// Test telegram notification path with invalid credentials
	t.Run("telegram_invalid", func(t *testing.T) {
		// doNotify calls SendTelegramNotification which may return an error
		// but doNotify itself doesn't return it, so we test it doesn't panic
		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("doNotify panicked: %v", r)
				}
			}()
			doNotify("always", "telegram", false, "Test Title", "Test Message")
		}()
	})

	// Test slack notification path
	t.Run("slack", func(t *testing.T) {
		// doNotify calls SendSlackNotification which may return an error
		// but doNotify itself doesn't return it, so we test it doesn't panic
		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("doNotify panicked: %v", r)
				}
			}()
			// Set slack webhook URL to empty to test error handling
			slackWebhookURL = ""
			doNotify("always", "slack", false, "Test Title", "Test Message")
		}()
	})

	// Test email notification path
	t.Run("email", func(t *testing.T) {
		// doNotify calls SendEmailNotification which may return an error
		// but doNotify itself doesn't return it, so we test it doesn't panic
		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("doNotify panicked: %v", r)
				}
			}()
			// Set email parameters to empty to test error handling
			emailTo = ""
			emailFrom = ""
			emailPassword = ""
			emailSMTPHost = ""
			doNotify("always", "email", false, "Test Title", "Test Message")
		}()
	})

	// Test invalid method
	t.Run("invalid_method", func(t *testing.T) {
		// Should not panic, though no notification will be sent
		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("doNotify panicked: %v", r)
				}
			}()
			doNotify("always", "invalid", false, "Test Title", "Test Message")
		}()
	})
}

// TestSendSlackNotification tests Slack notification with invalid webhook
func TestSendSlackNotification(t *testing.T) {
	t.Run("invalid_webhook_url", func(t *testing.T) {
		// Test with empty webhook URL
		err := SendSlackNotification("", "test message", false)
		if err == nil {
			t.Error("Expected error with empty webhook URL")
		}

		// Test with invalid webhook URL
		err = SendSlackNotification("not-a-valid-url", "test message", false)
		if err == nil {
			t.Error("Expected error with invalid webhook URL")
		}
	})

	t.Run("malformed_webhook_url", func(t *testing.T) {
		// Test with malformed URL
		err := SendSlackNotification("http://[invalid", "test message", false)
		if err == nil {
			t.Error("Expected error with malformed webhook URL")
		}
	})
}

// TestSendEmailNotification tests Email notification with invalid parameters
func TestSendEmailNotification(t *testing.T) {
	tests := []struct {
		name      string
		to        string
		from      string
		password  string
		smtpHost  string
		smtpPort  int
		expectErr bool
	}{
		{
			name:      "empty_to",
			to:        "",
			from:      "sender@example.com",
			password:  "password",
			smtpHost:  "smtp.example.com",
			smtpPort:  587,
			expectErr: true,
		},
		{
			name:      "empty_from",
			to:        "recipient@example.com",
			from:      "",
			password:  "password",
			smtpHost:  "smtp.example.com",
			smtpPort:  587,
			expectErr: true,
		},
		{
			name:      "empty_smtp_host",
			to:        "recipient@example.com",
			from:      "sender@example.com",
			password:  "password",
			smtpHost:  "",
			smtpPort:  587,
			expectErr: true,
		},
		{
			name:      "invalid_smtp_host",
			to:        "recipient@example.com",
			from:      "sender@example.com",
			password:  "password",
			smtpHost:  "invalid-host-that-does-not-exist-12345",
			smtpPort:  587,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SendEmailNotification(tt.to, tt.from, tt.password, tt.smtpHost, tt.smtpPort, "Test Subject", "Test Message", false)
			if tt.expectErr && err == nil {
				t.Error("Expected error but got nil")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}
