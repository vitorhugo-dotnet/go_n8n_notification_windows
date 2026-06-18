//go:build windows

package notifier

import (
	"github.com/go-toast/toast"
	"github.com/vitorhugo-java/go_n8n_notification_windows/internal/alert"
)

// Send dispatches a Windows toast notification for the given alert state.
func Send(appID string, s *alert.State) error {
	n := toast.Notification{
		AppID:   appID,
		Title:   s.Title,
		Message: s.Message,
	}

	switch s.Sound {
	case "critical", "alarm":
		n.Audio = toast.LoopingAlarm
	case "call":
		n.Audio = toast.LoopingCall
	case "silent":
		n.Audio = toast.Silent
	default:
		n.Audio = toast.Default
	}

	return n.Push()
}
