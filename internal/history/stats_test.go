package history

import (
	"testing"
	"time"
)

var statsBase = time.Now()

func statsEntries() []Entry {
	return []Entry{
		{Host: "host1", Ports: []int{80, 443}, Event: "opened", Timestamp: statsBase},
		{Host: "host1", Ports: []int{80}, Event: "opened", Timestamp: statsBase.Add(time.Minute)},
		{Host: "host2", Ports: []int{443, 8080}, Event: "opened", Timestamp: statsBase},
		{Host: "host1", Ports: []int{80}, Event: "closed", Timestamp: statsBase.Add(2 * time.Minute)},
	}
}

func TestTopPorts_AllHosts(t *testing.T) {
	entries := statsEntries()
	top := TopPorts(entries, "", 3)
	if len(top) == 0 {
		t.Fatal("expected results")
	}
	if top[0].Port != 80 {
		t.Errorf("expected port 80 first, got %d", top[0].Port)
	}
	if top[0].Count != 2 {
		t.Errorf("expected count 2, got %d", top[0].Count)
	}
}

func TestTopPorts_FilterHost(t *testing.T) {
	top := TopPorts(statsEntries(), "host2", 10)
	if len(top) != 2 {
		t.Fatalf("expected 2 ports for host2, got %d", len(top))
	}
}

func TestTopPorts_SkipClosed(t *testing.T) {
	entries := []Entry{
		{Host: "h", Ports: []int{22}, Event: "closed", Timestamp: statsBase},
	}
	top := TopPorts(entries, "", 10)
	if len(top) != 0 {
		t.Errorf("closed events should be excluded")
	}
}

func TestTopPorts_Limit(t *testing.T) {
	top := TopPorts(statsEntries(), "", 1)
	if len(top) != 1 {
		t.Errorf("expected 1 result, got %d", len(top))
	}
}

func TestHostActivity(t *testing.T) {
	activity := HostActivity(statsEntries())
	if activity["host1"] != 3 {
		t.Errorf("expected 3 events for host1, got %d", activity["host1"])
	}
	if activity["host2"] != 1 {
		t.Errorf("expected 1 event for host2, got %d", activity["host2"])
	}
}
