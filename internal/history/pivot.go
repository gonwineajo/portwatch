package history

import "sort"

// PortHostMatrix holds a cross-tabulation of ports vs hosts.
// Each entry maps port -> set of hosts that had that port open.
type PortHostMatrix struct {
	// Ports lists all observed ports in sorted order.
	Ports []int
	// Hosts lists all observed hosts in sorted order.
	Hosts []string
	// Cells maps port -> host -> open count.
	Cells map[int]map[string]int
}

// Pivot builds a port-by-host matrix from the provided entries.
// Only "opened" and "scan" events contribute to the open counts.
// Closed events are ignored so the matrix reflects positive presence.
func Pivot(entries []Entry) PortHostMatrix {
	cells := map[int]map[string]int{}
	hostSet := map[string]struct{}{}

	for _, e := range entries {
		if e.Event == EventNoChange || e.Event == EventClosed {
			continue
		}
		for _, p := range e.Ports {
			if cells[p] == nil {
				cells[p] = map[string]int{}
			}
			cells[p][e.Host]++
			hostSet[e.Host] = struct{}{}
		}
	}

	ports := make([]int, 0, len(cells))
	for p := range cells {
		ports = append(ports, p)
	}
	sort.Ints(ports)

	hosts := make([]string, 0, len(hostSet))
	for h := range hostSet {
		hosts = append(hosts, h)
	}
	sort.Strings(hosts)

	return PortHostMatrix{
		Ports: ports,
		Hosts: hosts,
		Cells: cells,
	}
}

// HostsForPort returns all hosts that had the given port open, sorted.
func (m PortHostMatrix) HostsForPort(port int) []string {
	hostMap, ok := m.Cells[port]
	if !ok {
		return nil
	}
	hosts := make([]string, 0, len(hostMap))
	for h := range hostMap {
		hosts = append(hosts, h)
	}
	sort.Strings(hosts)
	return hosts
}

// PortsForHost returns all ports observed open on the given host, sorted.
func (m PortHostMatrix) PortsForHost(host string) []int {
	var ports []int
	for _, p := range m.Ports {
		if m.Cells[p][host] > 0 {
			ports = append(ports, p)
		}
	}
	return ports
}
