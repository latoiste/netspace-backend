package chat

import (
	"log"
	"sync"

	"github.com/latoiste/netspace/db"
)

type Manager struct {
	// mu guards Hubs. The map is read/written from multiple goroutines: each WS
	// upgrade (LocationHub) may create an entry, while admin requests
	// (ForceLogoutUser, TotalOnline) iterate it. Without this lock those races
	// could corrupt the map or miss a freshly created hub.
	mu   sync.RWMutex
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
	m.mu.Lock()
	defer m.mu.Unlock()

	hub, ok := m.Hubs[locationSlug]
	if !ok {
		log.Printf("Hub for %v isn't created, creating\n", locationSlug)
		hub = NewHub(m.repo, locationId)
		m.Hubs[locationSlug] = hub

		// run() is the single goroutine that owns hub state (Clients map, groups,
		// blocks) and also performs broadcast fan-out, so there are no
		// cross-goroutine races on that state.
		go hub.run()
	}

	return hub
}

func (m *Manager) ForceLogoutUser(userId string) {
	// Don't touch hub.Clients here — that map is owned by each hub's run()
	// goroutine. Hand the request off via the channel so the lookup and the
	// socket teardown happen there, race-free. Snapshot the hub list under the
	// lock first so we don't iterate the map while another goroutine mutates it.
	for _, hub := range m.snapshotHubs() {
		hub.forceLogout <- userId
	}
}

// TotalOnline sums the live-socket count across every hub. Used by admin
// analytics as the source of truth for "User Aktif Sekarang" — it reflects
// actual WebSocket connections, not the DB isactive flag (which drifts when a
// process dies before its disconnect handlers run).
func (m *Manager) TotalOnline() int {
	total := 0
	for _, hub := range m.snapshotHubs() {
		total += hub.OnlineCount()
	}
	return total
}

// snapshotHubs returns a copy of the current hub list under the read lock, so
// callers can iterate without holding the lock or racing map mutations.
func (m *Manager) snapshotHubs() []*Hub {
	m.mu.RLock()
	defer m.mu.RUnlock()

	hubs := make([]*Hub, 0, len(m.Hubs))
	for _, hub := range m.Hubs {
		hubs = append(hubs, hub)
	}
	return hubs
}
