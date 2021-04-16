package main

import (
	"fmt"
	"sync"
)

// Settings settngs for all buttons/contexts
type Settings struct {
	sync.Mutex `json:"-"`
	pi         map[string]*PropertyInspector `json:"-"`
}

var (
	settings = Settings{
		pi: make(map[string]*PropertyInspector),
	}
)

// Save save setting with sd context
func (s *Settings) Save(ctxStr string, pi *PropertyInspector) {
	s.Lock()
	defer s.Unlock()
	s.pi[ctxStr] = pi
}

// Load setting with specified context
func (s *Settings) Load(ctxStr string) (*PropertyInspector, error) {
	s.Lock()
	defer s.Unlock()
	b, ok := s.pi[ctxStr]
	if !ok {
		return nil, fmt.Errorf("Setting not found for this context")
	}
	return b, nil
}

// PropertyInspector Settings for each button to save persistantly on action instance
type PropertyInspector struct {
	Command   string `json:"command"`
	Connected bool   `json:"connected"`
}
