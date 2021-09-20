package main

import "sync"

type SafeIterator struct {
	mu sync.Mutex
	index int
	values []string
}

func (i *SafeIterator) Next() string {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.index += 1
	return i.values[i.index]
}


var AvailablePorts = SafeIterator{
	values: []string{
		"50000",
		"50001",
		"50002",
		"50003",
		"50004",
		"50005",
		"50006",
		"50007",
		"50008",
		"50009",
		"50010",
		"50011",
		"50012",
		"50013",
		"50014",
		"50015",
		"50016",
		"50017",
		"50018",
		"50019",
		"50020",
		"50021",
		"50022",
		"50023",
		"50024",
		"50025",
		"50026",
	},
}

