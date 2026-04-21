package history

import "sort"

// NormalizeOptions controls how entries are normalized.
type NormalizeOptions struct {
	// DeduplicateScans removes consecutive scan entries for the same host
	// that carry identical open-port sets (no change events).
	DeduplicateScans bool
	// SortPorts ensures OpenedPorts and ClosedPorts are sorted ascending.
	SortPorts bool
	// DropEmpty removes entries where both OpenedPorts and ClosedPorts are
	// empty and the event is not EventScan.
	DropEmpty bool
}

// Normalize applies a set of cleanup transformations to a slice of entries
// and returns the cleaned result. The original slice is not modified.
func Normalize(entries []Entry, opts NormalizeOptions) []Entry {
	out := make([]Entry, 0, len(entries))

	type scanKey struct {
		host        string
		fingerprint string
	}
	lastScan := make(map[string]string) // host → last fingerprint

	for _, e := range entries {
		e = copyEntry(e)

		if opts.SortPorts {
			sort.Ints(e.OpenedPorts)
			sort.Ints(e.ClosedPorts)
			sort.Ints(e.Ports)
		}

		if opts.DropEmpty &&
			e.Event != EventScan &&
			len(e.OpenedPorts) == 0 &&
			len(e.ClosedPorts) == 0 {
			continue
		}

		if opts.DeduplicateScans && e.Event == EventNoChange {
			fp := portFingerprint(e.Ports)
			if lastScan[e.Host] == fp {
				continue
			}
			lastScan[e.Host] = fp
		}

		out = append(out, e)
	}
	return out
}

// copyEntry performs a shallow copy of an Entry with independent port slices.
func copyEntry(e Entry) Entry {
	cp := e
	if e.Ports != nil {
		cp.Ports = append([]int(nil), e.Ports...)
	}
	if e.OpenedPorts != nil {
		cp.OpenedPorts = append([]int(nil), e.OpenedPorts...)
	}
	if e.ClosedPorts != nil {
		cp.ClosedPorts = append([]int(nil), e.ClosedPorts...)
	}
	return cp
}
