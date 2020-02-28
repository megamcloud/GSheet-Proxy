package memDb

import (
	"bufio"
	"encoding/gob"
	"git.anphabe.net/event/anphabe-event-hub/domain/model/scanItem"
	"github.com/patrickmn/go-cache"
	"os"
)

type Connection struct {
	dbFolder string
	connections map[string]*cache.Cache
}

type memDumpStruct struct {
	 Items map[string]cache.Item
}

func NewMemDbConnection(dbFolder string) *Connection {
	return &Connection{
		dbFolder: dbFolder,
		connections: make(map[string]*cache.Cache),
	}
}

func (db *Connection) InitRepository(name string) (scanItem.RepositoryInterface, error) {
	fName := db.dbFolder + "/" + name

	db.connections[name] = db.fromFile( fName + "_item.mem")
	db.connections[name+"_activity"] = db.fromFile( fName + "_activity.mem")

	repo := &ScanItemRepository{
		dbName:  name,
		dbFolder: db.dbFolder,
		itemStorage:     db.connections[name],
		activityStorage: db.connections[name+"_activity"],
	}

	return repo, nil
}

func (db *Connection) fromFile(fileName string) *cache.Cache {

	if fp, err := os.Open(fileName); err == nil {

		defer func() { _ = fp.Close()}()

		var memDump memDumpStruct

		gob.Register(memDumpStruct{})
		gob.Register(&scanItem.ScanItem{})
		gob.Register(scanItem.ItemActivities{})
		gob.Register([]scanItem.ItemActivity{})
		gob.Register(scanItem.ItemActivity{})

		dec := gob.NewDecoder(bufio.NewReader(fp))

		if err := dec.Decode(&memDump) ; err == nil {
			if len(memDump.Items) > 0 {
				return cache.NewFrom(0,0, memDump.Items)
			}
		}
	}

	return cache.New(0,0, )
}
