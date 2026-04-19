package history

import "testing"

func TestPortDeltas_Basic(t *testing.T) {
	entries := []Entry{
		{Opened: []int{80, 443}, Closed: []int{}},
		{Opened: []int{80}, Closed: []int{443}},
	}
	deltas := PortDeltas(entries)

	m := map[int]PortDelta{}
	for _, d := range deltas {
		m[d.Port] = d
	}

	if m[80].TimesOpened != 2 {
		t.Errorf("port 80 opened: want 2, got %d", m[80].TimesOpened)
	}
	if m[443].TimesOpened != 1 || m[443].TimesClosed != 1 {
		t.Errorf("port 443 unexpected: %+v", m[443])
	}
}

func TestPortDeltas_Empty(t *testing.T) {
	if len(PortDeltas(nil)) != 0 {
		t.Error("expected empty result")
	}
}

func TestPortDeltas_Sorted(t *testing.T) {
	entries := []Entry{
		{Opened: []int{8080, 22, 443}},
	}
	deltas := PortDeltas(entries)
	for i := 1; i < len(deltas); i++ {
		if deltas[i].Port < deltas[i-1].Port {
			t.Error("results not sorted")
		}
	}
}
