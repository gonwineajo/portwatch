package history

import (
	"testing"
)

func makeSig(host string, ports ...int) HostSignature {
	return HostSignature{Host: host, Ports: ports}
}

func TestDiffSignatures_NewPort(t *testing.T) {
	baseline := []HostSignature{makeSig("h1", 80, 443)}
	current := []HostSignature{makeSig("h1", 80, 443, 8080)}

	diffs := DiffSignatures(baseline, current)
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(diffs))
	}
	if len(diffs[0].Added) != 1 || diffs[0].Added[0] != 8080 {
		t.Errorf("expected added=[8080], got %v", diffs[0].Added)
	}
	if len(diffs[0].Removed) != 0 {
		t.Errorf("expected no removed, got %v", diffs[0].Removed)
	}
}

func TestDiffSignatures_RemovedPort(t *testing.T) {
	baseline := []HostSignature{makeSig("h1", 80, 443, 22)}
	current := []HostSignature{makeSig("h1", 80, 443)}

	diffs := DiffSignatures(baseline, current)
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(diffs))
	}
	if len(diffs[0].Removed) != 1 || diffs[0].Removed[0] != 22 {
		t.Errorf("expected removed=[22], got %v", diffs[0].Removed)
	}
}

func TestDiffSignatures_NoChange(t *testing.T) {
	baseline := []HostSignature{makeSig("h1", 80, 443)}
	current := []HostSignature{makeSig("h1", 80, 443)}

	diffs := DiffSignatures(baseline, current)
	if len(diffs) != 0 {
		t.Errorf("expected no diffs, got %d", len(diffs))
	}
}

func TestDiffSignatures_BrandNewHost(t *testing.T) {
	baseline := []HostSignature{}
	current := []HostSignature{makeSig("h2", 22, 80)}

	diffs := DiffSignatures(baseline, current)
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(diffs))
	}
	if diffs[0].Host != "h2" {
		t.Errorf("expected host h2, got %s", diffs[0].Host)
	}
	if len(diffs[0].Added) != 2 {
		t.Errorf("expected 2 added ports, got %v", diffs[0].Added)
	}
	if len(diffs[0].Removed) != 0 {
		t.Errorf("expected no removed, got %v", diffs[0].Removed)
	}
}

func TestDiffSignatures_MultipleHosts_Sorted(t *testing.T) {
	baseline := []HostSignature{
		makeSig("zhost", 80),
		makeSig("ahost", 443),
	}
	current := []HostSignature{
		makeSig("zhost", 80, 8080),
		makeSig("ahost", 443, 22),
	}

	diffs := DiffSignatures(baseline, current)
	if len(diffs) != 2 {
		t.Fatalf("expected 2 diffs, got %d", len(diffs))
	}
	if diffs[0].Host != "ahost" || diffs[1].Host != "zhost" {
		t.Errorf("expected sorted by host, got %s, %s", diffs[0].Host, diffs[1].Host)
	}
}
