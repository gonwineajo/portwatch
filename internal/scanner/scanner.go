package scanner

import (
	"fmt"
	"net"
	"sort"
	"time"
)

// Result holds the open ports found on a host.
type Result struct {
	Host  string
	Ports []int
	At    time.Time
}

// Scanner performs TCP port scans on a host.
type Scanner struct {
	Timeout    time.Duration
	Concurrent int
}

// New returns a Scanner with sensible defaults.
func New() *Scanner {
	return &Scanner{
		Timeout:    500 * time.Millisecond,
		Concurrent: 100,
	}
}

// Scan checks every port in ports on host and returns open ones.
func (s *Scanner) Scan(host string, ports []int) (*Result, error) {
	type work struct{ port int }
	type hit struct{ port int }

	jobs := make(chan work, len(ports))
	hits := make(chan hit, len(ports))

	for i := 0; i < s.Concurrent; i++ {
		go func() {
			for w := range jobs {
				addr := fmt.Sprintf("%s:%d", host, w.port)
				conn, err := net.DialTimeout("tcp", addr, s.Timeout)
				if err == nil {
					conn.Close()
					hits <- hit{w.port}
				} else {
					hits <- hit{-1}
				}
			}
		}()
	}

	for _, p := range ports {
		jobs <- work{p}
	}
	close(jobs)

	var open []int
	for range ports {
		if h := <-hits; h.port != -1 {
			open = append(open, h.port)
		}
	}
	sort.Ints(open)

	return &Result{Host: host, Ports: open, At: time.Now()}, nil
}
