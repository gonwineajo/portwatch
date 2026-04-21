package history

import (
	"testing"
)

func chainEntry(host string, event string, opened, closed []int, ts int64) Entry {
	return Entry{
		Host:      host,
		Event:     event,
		Opened:    opened,
		Closed:    closed,
		Timestamp: ts,
	}
}

func TestBuildChains_BasicOpened(t *testing.T) {
	entries := []Entry{
		chainEntry("host-a", EventOpened, []int{80, 443}, nil, 1000),
	}
	chains := BuildChains(entries)
	if len(chains) != 2 {
		t.Fatalf("expected 2 chains, got %d", len(chains))
	}
	if chains[0].Port != 443 || chains[1].Port != 80 {
		t.Errorf("unexpected ports: %d, %d", chains[0].Port, chains[1].Port)
	}
}

func TestBuildChains_OpenThenClose(t *testing.T) {
	entries := []Entry{
		chainEntry("host-a", EventOpened, []int{22}, nil, 1000),
		chainEntry("host-a", EventClosed, nil, []int{22}, 2000),
	}
	chains := BuildChains(entries)
	if len(chains) != 1 {
		t.Fatalf("expected 1 chain, got %d", len(chains))
	}
	c := chains[0]
	if len(c.Steps) != 2 {
		t.Fatalf("expected 2 steps, got %d", len(c.Steps))
	}
	if c.Steps[0].Event != EventOpened || c.Steps[1].Event != EventClosed {
		t.Errorf("unexpected step events: %v", c.Steps)
	}
}

func TestBuildChains_SkipsNoChange(t *testing.T) {
	entries := []Entry{
		chainEntry("host-a", EventNoChange, nil, nil, 1000),
		chainEntry("host-a", EventOpened, []int{8080}, nil, 2000),
	}
	chains := BuildChains(entries)
	if len(chains) != 1 {
		t.Fatalf("expected 1 chain, got %d", len(chains))
	}
}

func TestBuildChains_MultipleHosts(t *testing.T) {
	entries := []Entry{
		chainEntry("host-a", EventOpened, []int{80}, nil, 1000),
		chainEntry("host-b", EventOpened, []int{80}, nil, 1000),
	}
	chains := BuildChains(entries)
	if len(chains) != 2 {
		t.Fatalf("expected 2 chains, got %d", len(chains))
	}
	if chains[0].Host != "host-a" || chains[1].Host != "host-b" {
		t.Errorf("unexpected host order: %s, %s", chains[0].Host, chains[1].Host)
	}
}

func TestChainsByHost_Filters(t *testing.T) {
	entries := []Entry{
		chainEntry("host-a", EventOpened, []int{80}, nil, 1000),
		chainEntry("host-b", EventOpened, []int{443}, nil, 1000),
		chainEntry("host-a", EventOpened, []int{22}, nil, 2000),
	}
	all := BuildChains(entries)
	filtered := ChainsByHost(all, "host-a")
	for _, c := range filtered {
		if c.Host != "host-a" {
			t.Errorf("expected host-a, got %s", c.Host)
		}
	}
	if len(filtered) != 2 {
		t.Errorf("expected 2 chains for host-a, got %d", len(filtered))
	}
}

func TestBuildChains_Empty(t *testing.T) {
	chains := BuildChains(nil)
	if len(chains) != 0 {
		t.Errorf("expected empty chains, got %d", len(chains))
	}
}
