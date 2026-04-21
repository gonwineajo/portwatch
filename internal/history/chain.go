package history

// Chain represents a sequence of state transitions for a single port on a host.
type Chain struct {
	Host  string
	Port  int
	Steps []ChainStep
}

// ChainStep is a single transition in a port's lifecycle.
type ChainStep struct {
	Timestamp int64
	Event     string
}

// BuildChains reconstructs per-host-port event chains from a slice of entries.
// Each chain tracks the ordered sequence of events (opened/closed) for a
// specific (host, port) pair, making it easy to reason about port lifecycles.
func BuildChains(entries []Entry) []Chain {
	type key struct {
		host string
		port int
	}

	chainMap := map[key]*Chain{}

	for _, e := range entries {
		if e.Event == EventNoChange {
			continue
		}
		ports := e.Opened
		event := EventOpened
		if e.Event == EventClosed {
			ports = e.Closed
			event = EventClosed
		}
		for _, p := range ports {
			k := key{host: e.Host, port: p}
			if _, ok := chainMap[k]; !ok {
				chainMap[k] = &Chain{Host: e.Host, Port: p}
			}
			chainMap[k].Steps = append(chainMap[k].Steps, ChainStep{
				Timestamp: e.Timestamp,
				Event:     event,
			})
		}
	}

	result := make([]Chain, 0, len(chainMap))
	for _, c := range chainMap {
		result = append(result, *c)
	}
	sortChains(result)
	return result
}

// ChainsByHost returns only the chains matching the given host.
func ChainsByHost(chains []Chain, host string) []Chain {
	var out []Chain
	for _, c := range chains {
		if c.Host == host {
			out = append(out, c)
		}
	}
	return out
}

func sortChains(chains []Chain) {
	for i := 1; i < len(chains); i++ {
		for j := i; j > 0; j-- {
			a, b := chains[j-1], chains[j]
			if a.Host > b.Host || (a.Host == b.Host && a.Port > b.Port) {
				chains[j-1], chains[j] = chains[j], chains[j-1]
			} else {
				break
			}
		}
	}
}
