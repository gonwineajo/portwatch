package history

import (
	"testing"
	"time"
)

var pivotEntries = []Entry{
	{Host: "host-a", Event: EventOpened, Ports: []int{80, 443}, Timestamp: time.Now()},
	{Host: "host-b", Event: EventOpened, Ports: []int{443, 8080}, Timestamp: time.Now()},
	{Host: "host-a", Event: EventScan, Ports: []int{80}, Timestamp: time.Now()},
	{Host: "host-c", Event: EventClosed, Ports: []int{22}, Timestamp: time.Now()},
	{Host: "host-a", Event: EventNoChange, Ports: []int{80}, Timestamp: time.Now()},
}

func TestPivot_PortsPresent(t *testing.T) {
	m := Pivot(pivotEntries)
	if len(m.Ports) == 0 {
		t.Fatal("expected ports in matrix")
	}
	for i := 1; i < len(m.Ports); i++ {
		if m.Ports[i] <= m.Ports[i-1] {
			t.Errorf("ports not sorted at index %d", i)
		}
	}
}

func TestPivot_SkipsClosedAndNoChange(t *testing.T) {
	m := Pivot(pivotEntries)
	// port 22 came only from a closed event — should not appear
	for _, p := range m.Ports {
		if p == 22 {
			t.Error("port 22 (closed-only) should not appear in matrix")
		}
	}
}

func TestPivot_HostsForPort(t *testing.T) {
	m := Pivot(pivotEntries)
	hosts := m.HostsForPort(443)
	if len(hosts) != 2 {
		t.Fatalf("expected 2 hosts for port 443, got %d", len(hosts))
	}
	if hosts[0] != "host-a" || hosts[1] != "host-b" {
		t.Errorf("unexpected hosts: %v", hosts)
	}
}

func TestPivot_HostsForPort_Missing(t *testing.T) {
	m := Pivot(pivotEntries)
	hosts := m.HostsForPort(9999)
	if hosts != nil {
		t.Errorf("expected nil for unknown port, got %v", hosts)
	}
}

func TestPivot_PortsForHost(t *testing.T) {
	m := Pivot(pivotEntries)
	ports := m.PortsForHost("host-a")
	// host-a has 80 (opened + scan) and 443 (opened)
	if len(ports) != 2 {
		t.Fatalf("expected 2 ports for host-a, got %d: %v", len(ports), ports)
	}
	if ports[0] != 80 || ports[1] != 443 {
		t.Errorf("unexpected ports for host-a: %v", ports)
	}
}

func TestPivot_Empty(t *testing.T) {
	m := Pivot(nil)
	if len(m.Ports) != 0 || len(m.Hosts) != 0 {
		t.Error("expected empty matrix for nil input")
	}
}
