//go:build windows

package notifier

import (
	"os"
	"testing"
	"time"

	"seall/internal/database/models"
)

// Explicit opt-in: this test displays a real Windows notification.
func TestAutoScannerWindowsNotificationLive(t *testing.T) {
	if os.Getenv("SEALL_TEST_WINDOWS_NOTIFICATION") != "1" {
		t.Skip("set SEALL_TEST_WINDOWS_NOTIFICATION=1 to display a real notification")
	}
	n := NewNotifier()
	n.SetSettings(t.TempDir(), &models.NotificationSettings{}, nil)
	result := make(chan error, 1)
	n.push = func(title, message, _ string) error {
		if title != "Seall: Auto Scanner" {
			t.Errorf("unexpected scan notification title: %q", title)
		}
		err := defaultPush(title, message, "")
		result <- err
		return err
	}
	n.Notify(AutoScanner, "Notification branding test. No library scan was run.")
	select {
	case err := <-result:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(30 * time.Second):
		t.Fatal("Windows notification did not finish")
	}
}
