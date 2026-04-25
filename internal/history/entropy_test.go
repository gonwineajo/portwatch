package history

import (
	"math"
	"testing"
	"time"
)

func entropyEntry(host, event string, ports []int, secsAgo int) Entry {
	return Entry{
		Host:      host,
		Event:     event,
		Ports:     ports,
		Timestamp: time.Now().Add(-time.Duration(secsAgo) * time.Second),
	}
}

func TestAnalyseEntropy_HighEntropyHost(t *testing.T) {
	// hostA opens four different ports equally — maximum spread → high entropy
	entries := []Entry{
		entropyEntry("hostA", EventOpened, []int{80}, 100),
		entropyEntry("hostA", EventOpened, []int{443}, 90),
		entropyEntry("hostA", EventOpened, []int{8080}, 80),
		entropyEntry("hostA", EventOpened, []int{9090}, 70),
	}
	res := AnalyseEntropy(entries, 1)
	if len(res) != 1 {
		t.Fatalf("expected 1 result, got %d", len(res))
	}
	if res[0].Host != "hostA" {
		t.Errorf("unexpected host %s", res[0].Host)
	}
	// uniform distribution over 4 symbols → entropy == 2.0 bits
	if math.Abs(res[0].Entropy-2.0) > 0.001 {
		t.Errorf("expected entropy ~2.0, got %.4f", res[0].Entropy)
	}
}

func TestAnalyseEntropy_LowEntropyHost(t *testing.T) {
	// hostB always opens the same port → entropy == 0
	entries := []Entry{
		entropyEntry("hostB", EventOpened, []int{80}, 50),
		entropyEntry("hostB", EventOpened, []int{80}, 40),
		entropyEntry("hostB", EventOpened, []int{80}, 30),
	}
	res := AnalyseEntropy(entries, 1)
	if len(res) != 1 {
		t.Fatalf("expected 1 result, got %d", len(res))
	}
	if res[0].Entropy != 0.0 {
		t.Errorf("expected entropy 0, got %.4f", res[0].Entropy)
	}
}

func TestAnalyseEntropy_SkipsNoChange(t *testing.T) {
	entries := []Entry{
		entropyEntry("hostC", EventNoChange, []int{80, 443}, 60),
		entropyEntry("hostC", EventOpened, []int{22}, 50),
	}
	res := AnalyseEntropy(entries, 1)
	if len(res) != 1 {
		t.Fatalf("expected 1 result, got %d", len(res))
	}
	if res[0].Events != 1 {
		t.Errorf("expected 1 event (no-change skipped), got %d", res[0].Events)
	}
}

func TestAnalyseEntropy_MinEventsFilter(t *testing.T) {
	entries := []Entry{
		entropyEntry("hostD", EventOpened, []int{80}, 20),
	}
	res := AnalyseEntropy(entries, 5)
	if len(res) != 0 {
		t.Errorf("expected 0 results due to minEvents filter, got %d", len(res))
	}
}

func TestAnalyseEntropy_SortedDescending(t *testing.T) {
	entries := []Entry{
		// hostE: single port → entropy 0
		entropyEntry("hostE", EventOpened, []int{80}, 10),
		entropyEntry("hostE", EventOpened, []int{80}, 9),
		// hostF: two ports equally → entropy 1 bit
		entropyEntry("hostF", EventOpened, []int{80}, 10),
		entropyEntry("hostF", EventOpened, []int{443}, 9),
	}
	res := AnalyseEntropy(entries, 1)
	if len(res) != 2 {
		t.Fatalf("expected 2 results, got %d", len(res))
	}
	if res[0].Host != "hostF" {
		t.Errorf("expected hostF first (higher entropy), got %s", res[0].Host)
	}
}
