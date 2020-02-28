package service

import "git.anphabe.net/event/anphabe-event-hub/domain/model/scanItem"

type DbConnectionInterface interface {
	InitRepository(name string) (scanItem.RepositoryInterface, error)
}
