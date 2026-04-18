package notifier

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/user/portwatch/internal/snapshot"
)

// Channel represents a notification delivery method.
type Channel string

const (
	ChannelStdout  Channel = "stdout"
	ChannelWebhook Channel = "webhook"
	ChannelScript  Channel = "script"
)

// Config holds notifier configuration.
type Config struct {
	Channel Channel
	Target  string // webhook URL or script path
}

// Notifier sends notifications about port changes.
type Notifier struct {
	cfg Config
	out io.Writer
}

// New creates a Notifier with the given config.
func New(cfg Config) *Notifier {
	return &Notifier{cfg: cfg, out: os.Stdout}
}

// Notify dispatches a notification for the given diff.
func (n *Notifier) Notify(host string, diff snapshot.Diff) error {
	if len(diff.Opened) == 0 && len(diff.Closed) == 0 {
		return nil
	}
	msg := formatMessage(host, diff)
	switch n.cfg.Channel {
	case ChannelWebhook:
		return sendWebhook(n.cfg.Target, msg)
	case ChannelScript:
		return runScript(n.cfg.Target, host, msg)
	default:
		fmt.Fprintln(n.out, msg)
		return nil
	}
}

func formatMessage(host string, diff snapshot.Diff) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("[portwatch] host=%s", host))
	if len(diff.Opened) > 0 {
		sb.WriteString(fmt.Sprintf(" opened=%v", diff.Opened))
	}
	if len(diff.Closed) > 0 {
		sb.WriteString(fmt.Sprintf(" closed=%v", diff.Closed))
	}
	return sb.String()
}

func runScript(path, host, msg string) error {
	cmd := exec.Command(path, host, msg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
