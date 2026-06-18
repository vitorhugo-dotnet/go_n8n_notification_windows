//go:build !windows

package notifier

import (
	"fmt"

	"github.com/vitorhugo-java/go_n8n_notification_windows/internal/alert"
)

// Send prints the alert to stdout on non-Windows platforms (development stub).
func Send(appID string, s *alert.State) error {
	fmt.Printf("[NOTIFICATION] appID=%s title=%q message=%q priority=%s sound=%s\n",
		appID, s.Title, s.Message, s.Priority, s.Sound)
	return nil
}
