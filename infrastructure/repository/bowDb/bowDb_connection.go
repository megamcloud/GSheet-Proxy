package bowDb

import (
	"fmt"
	"git.anphabe.net/event/anphabe-event-hub/domain/model/scanItem"
	"github.com/dgraph-io/badger"
	"github.com/zippoxer/bow"
)

type Connection struct {
	dbFolder string
	conn     *bow.DB
}

func NewBowDbConnection(dbFolder string) *Connection {
	// Open database under directory "test".
	retryOpts := badger.DefaultOptions(dbFolder)
	retryOpts.Truncate = true

	conn, err := bow.Open(dbFolder, bow.SetBadgerOptions(retryOpts))

	if err != nil {
		fmt.Println(err.Error())
		conn.Close()
	}

	return &Connection{
		dbFolder: dbFolder,
		conn:     conn,
	}
}

func (c *Connection) InitRepository(name string) (scanItem.RepositoryInterface, error) {
	return &ScanItemRepository{
		repoName: name,
		conn:     c.conn,
	}, nil
}