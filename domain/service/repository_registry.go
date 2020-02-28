package service

import (
	"git.anphabe.net/event/anphabe-event-hub/domain/model/scanItem"
)

type RepositoryRegistryInterface interface {
	GetRepository(name string) (scanItem.RepositoryInterface, error)
	initRepository(name string) (scanItem.RepositoryInterface, error)
	Shutdown()
}

type repoRegistry struct {
	repositories map[string]scanItem.RepositoryInterface
	db           DbConnectionInterface
}

func NewRepositoryRegistry(db DbConnectionInterface) *repoRegistry {
	return &repoRegistry{
		repositories: make(map[string]scanItem.RepositoryInterface),
		db:           db,
	}
}

func (m *repoRegistry) GetRepository(name string) (scanItem.RepositoryInterface, error) {
	if r, ok := m.repositories[name]; ok {
		return r, nil
	}

	return m.initRepository(name)
}

func (m *repoRegistry) Shutdown() {
	for repoName, _ := range m.repositories {
		m.repositories[repoName].CloseDb()
	}
}

func (m *repoRegistry) initRepository(name string) (scanItem.RepositoryInterface, error) {
	var err error

	m.repositories[name], err = m.db.InitRepository(name)

	return m.repositories[name], err
}
