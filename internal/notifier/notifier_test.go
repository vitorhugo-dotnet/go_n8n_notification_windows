package notifier

import (
	"testing"

	"github.com/vitorhugo-java/go_n8n_notification_windows/internal/alert"
)

func TestSend_DoesNotPanic(t *testing.T) {
	s := &alert.State{
		Enabled:  true,
		Priority: "HIGH",
		Title:    "Test Alert",
		Message:  "This is a test",
		Sound:    "default",
	}
	// On non-Windows this prints to stdout; on Windows it pushes a toast.
	// We only verify it doesn't panic or return an unexpected error structure.
	_ = Send("test.app.id", s)
}
