package history

import (
	"testing"
	"time"
)

var clusterBase = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func clusterEntry(host string, ports []int, offset time.Duration) Entry {
	return Entry{
		Timestamp: clusterBase.Add(offset),
		Host:      host,
		Event:     "scan",
		Ports:     ports,
	}
}

func TestClusterByPorts_GroupsIdenticalHosts(t *testing.T) {
	entries := []Entry{
		clusterEntry("host-a", []int{80, 443}, 0),
		clusterEntry("host-b", []int{80, 443}, 0),
		clusterEntry("host-c", []int{22}, 0),
	}
	results := ClusterByPorts(entries)
	if len(results) != 2 {
		t.Fatalf("expected 2 clusters, got %d", len(results))
	}
	// largest cluster first
	if len(results[0].Hosts) != 2 {
		t.Errorf("expected 2 hosts in first cluster, got %d", len(results[0].Hosts))
	}
}

func TestClusterByPorts_UsesLatestScan(t *testing.T) {
	entries := []Entry{
		clusterEntry("host-a", []int{80}, 0),
		clusterEntry("host-a", []int{80, 443}, time.Hour), // newer
		clusterEntry("host-b", []int{80, 443}, 0),
	}
	results := ClusterByPorts(entries)
	if len(results) != 1 {
		t.Fatalf("expected 1 cluster, got %d", len(results))
	}
	if len(results[0].Hosts) != 2 {
		t.Errorf("expected host-a and host-b in same cluster")
	}
}

func TestClusterByPorts_IgnoresNonScan(t *testing.T) {
	entries := []Entry{
		clusterEntry("host-a", []int{80}, 0),
		{Timestamp: clusterBase, Host: "host-b", Event: "opened", Ports: []int{80}},
	}
	results := ClusterByPorts(entries)
	if len(results) != 1 {
		t.Fatalf("expected 1 cluster (only scan events), got %d", len(results))
	}
	if results[0].Hosts[0] != "host-a" {
		t.Errorf("unexpected host: %s", results[0].Hosts[0])
	}
}

func TestClusterByPorts_Empty(t *testing.T) {
	results := ClusterByPorts(nil)
	if len(results) != 0 {
		t.Errorf("expected empty result, got %d", len(results))
	}
}

func TestClusterByPorts_PortsAreSorted(t *testing.T) {
	entries := []Entry{
		clusterEntry("host-a", []int{443, 22, 80}, 0),
	}
	results := ClusterByPorts(entries)
	if len(results) != 1 {
		t.Fatalf("expected 1 cluster")
	}
	ports := results[0].Ports
	for i := 1; i < len(ports); i++ {
		if ports[i] < ports[i-1] {
			t.Errorf("ports not sorted: %v", ports)
		}
	}
}
