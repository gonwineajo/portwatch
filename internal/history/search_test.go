package history

import (
	"testing"
	"time"
)

func TestSearch_ByHost(t *testing.T) {
	entries := []Entry{
		{Host: "192.168.1.1", Event: "opened", Ports: []int{80}, Timestamp: time.Now()},
		{Host: "10.0.0.1", Event: "opened", Ports: []int{443}, Timestamp: time.Now()},
	}
	result := Search(entries, SearchOptions{Host: "192.168"})
	if len(result) != 1 || result[0].Host != "192.168.1.1" {
		t.Fatalf("expected 1 match for host filter, got %v", result)
	}
}

func TestSearch_ByPort(t *testing.T) {
	entries := []Entry{
		{Host: "h1", Event: "opened", Ports: []int{80, 443}, Timestamp: time.Now()},
		{Host: "h2", Event: "closed", Ports: []int{22}, Timestamp: time.Now()},
	}
	result := Search(entries, SearchOptions{Port: 443})
	if len(result) != 1 || result[0].Host != "h1" {
		t.Fatalf("expected 1 match for port filter, got %v", result)
	}
}

func TestSearch_ByEvent(t *testing.T) {
	entries := []Entry{
		{Host: "h1", Event: "opened", Ports: []int{80}, Timestamp: time.Now()},
		{Host: "h2", Event: "closed", Ports: []int{22}, Timestamp: time.Now()},
		{Host: "h3", Event: "opened", Ports: []int{8080}, Timestamp: time.Now()},
	}
	result := Search(entries, SearchOptions{Event: "opened"})
	if len(result) != 2 {
		t.Fatalf("expected 2 opened entries, got %d", len(result))
	}
}

func TestSearch_Combined(t *testing.T) {
	entries := []Entry{
		{Host: "web-01", Event: "opened", Ports: []int{80}, Timestamp: time.Now()},
		{Host: "web-02", Event: "opened", Ports: []int{443}, Timestamp: time.Now()},
		{Host: "web-01", Event: "closed", Ports: []int{80}, Timestamp: time.Now()},
	}
	result := Search(entries, SearchOptions{Host: "web-01", Event: "opened", Port: 80})
	if len(result) != 1 {
		t.Fatalf("expected 1 combined match, got %d", len(result))
	}
}

func TestSearch_Empty(t *testing.T) {
	result := Search(nil, SearchOptions{})
	if result != nil {
		t.Fatalf("expected nil result for empty input")
	}
}
