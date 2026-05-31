package chat

import (
	"errors"
	"log"

	"github.com/latoiste/netspace/db"
)

type Manager struct {
	Hubs map[string]*Hub
	repo *db.Repository
}

func NewManager(repo *db.Repository) *Manager {
	return &Manager{
		Hubs: make(map[string]*Hub, 0),
		repo: repo,
	}
}

func (m *Manager) LocationHub(locationSlug string, locationId int) *Hub {
	hub, ok := m.Hubs[locationSlug]
	if !ok {
		log.Printf("Hub for %v isn't created, creating\n", locationSlug)
		m.Hubs[locationSlug] = NewHub(m.repo, locationId)
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
	go hub.run()
	go hub.broadcastLoop()

	return nil
}

func (m *Manager) ForceLogoutUser(userId string) {
	for _, hub := range m.Hubs {
		client, ok := hub.Clients[userId]
		if ok {
			client.forceLogout()
		}
	}
}
