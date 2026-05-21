package chat

import (
	"errors"
	"log"
)

type Manager struct {
	Hubs map[string]*Hub
}

func NewManager() *Manager {
	return &Manager{
		Hubs: make(map[string]*Hub, 0),
	}
}

func (m *Manager) LocationHub(locationSlug string) *Hub {
	hub, ok := m.Hubs[locationSlug]
	if !ok {
		log.Printf("Hub for %v isn't created, creating\n", locationSlug)
		m.Hubs[locationSlug] = NewHub()
		hub = m.Hubs[locationSlug]

		m.StartHub(locationSlug)
	}

	return hub
}

func (m *Manager) StartHub(locationSlug string) error {
	hub, ok := m.Hubs[locationSlug]
	if !ok {
		return errors.New("Invalid hub location")
	}
	go hub.Run()
	return nil
}

// func (m *Manager) StopHub(locationSlug string) error {
// 	hub, ok := m.Hubs[locationSlug]
// 	if !ok {
// 		return errors.New("Invalid hub location")
// 	}

// }
