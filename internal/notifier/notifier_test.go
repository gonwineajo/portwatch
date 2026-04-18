package notifier

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/snapshot"
)

func diff(opened, closed []int) snapshot.Diff {
	return snapshot.Diff{Opened: opened, Closed: closed}
}

func TestNotify_Stdout(t *testing.T) {
	n := New(Config{Channel: ChannelStdout})
	var buf bytes.Buffer
	n.out = &buf
	if err := n.Notify("localhost", diff([]int{8080}, nil)); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "opened=") {
		t.Errorf("expected opened in output, got: %s", buf.String())
	}
}

func TestNotify_NoChanges_Silent(t *testing.T) {
	n := New(Config{Channel: ChannelStdout})
	var buf bytes.Buffer
	n.out = &buf
	if err := n.Notify("localhost", diff(nil, nil)); err != nil {
		t.Fatal(err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected no output for empty diff")
	}
}

func TestNotify_Webhook_Success(t *testing.T) {
	var received webhookPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := New(Config{Channel: ChannelWebhook, Target: ts.URL})
	if err := n.Notify("myhost", diff([]int{443}, []int{80})); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(received.Text, "myhost") {
		t.Errorf("expected host in webhook payload, got: %s", received.Text)
	}
}

func TestNotify_Webhook_BadStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n := New(Config{Channel: ChannelWebhook, Target: ts.URL})
	err := n.Notify("host", diff([]int{22}, nil))
	if err == nil {
		t.Error("expected error on 500 response")
	}
}

func TestFormatMessage(t *testing.T) {
	msg := formatMessage("10.0.0.1", diff([]int{22, 80}, []int{8080}))
	if !strings.Contains(msg, "host=10.0.0.1") {
		t.Errorf("missing host: %s", msg)
	}
	if !strings.Contains(msg, "opened=") || !strings.Contains(msg, "closed=") {
		t.Errorf("missing opened/closed: %s", msg)
	}
}
