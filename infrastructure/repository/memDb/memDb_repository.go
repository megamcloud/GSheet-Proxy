// go_cache_test.go
package memDb

import (
	"bufio"
	"fmt"
	"git.anphabe.net/event/anphabe-event-hub/domain/model/scanItem"
	"github.com/patrickmn/go-cache"
	"encoding/gob"
	"os"
)

type ScanItemRepository struct {
	dbName   string
	dbFolder string
	itemStorage     *cache.Cache
	activityStorage *cache.Cache
}

func (s *ScanItemRepository) NewItem(itemKey string, data map[string]string) (*scanItem.ScanItem, error) {
	item, err := scanItem.NewScanItem(itemKey, data)

	if nil != err {
		return nil, err
	}

	s.SetItem(item)

	return item, nil
}

func (s *ScanItemRepository) SetItem(item *scanItem.ScanItem) {
	s.itemStorage.Set(item.GetKey(), item, 0)
}

func (s *ScanItemRepository) GetRepoName() string {
	return s.dbName
}

func (s *ScanItemRepository) GetItem(itemKey string) (*scanItem.ScanItem, bool) {
	if item, found := s.itemStorage.Get(itemKey); found {
		return item.(*scanItem.ScanItem), found
	}

	return nil, false
}

func (s *ScanItemRepository) GetItemDetail(itemKey string) (*scanItem.ItemDetail, bool) {
	if item, found := s.GetItem(itemKey); found {
		activities := s.GetItemActivities(itemKey)
		return &scanItem.ItemDetail{
			ScanItem:   *item,
			Activities: activities.Activities,
		}, true
	}

	return nil, false
}

func (s *ScanItemRepository) AddItemActivity(itemKey string, activity scanItem.ItemActivity) *scanItem.ItemActivities {

	item := s.GetItemActivities(itemKey)

	if nil != item {
		item.Activities = append([]scanItem.ItemActivity{activity}, item.Activities...)
		s.activityStorage.Set(itemKey, item.Activities, 0)
	}

	return item
}

func (s *ScanItemRepository) GetItemActivities(itemKey string) *scanItem.ItemActivities {
	if _, found := s.GetItem(itemKey); found {
		itemActivities := &scanItem.ItemActivities{Key:itemKey, Activities: nil}

		if data , found := s.activityStorage.Get(itemKey); found {
			itemActivities.Activities = data.([]scanItem.ItemActivity)
		}

		return itemActivities
	}

	return nil
}

func (s *ScanItemRepository) Items() []*scanItem.ScanItem {
	var result []*scanItem.ScanItem

	for _, item := range s.itemStorage.Items() {
		result = append(result, item.Object.(*scanItem.ScanItem))
	}

	return result
}

func (s *ScanItemRepository) Len() int {
	return s.itemStorage.ItemCount()
}

func (s *ScanItemRepository) CloseDb() {
	s.SaveToFile()
	s.itemStorage.Flush()
	s.activityStorage.Flush()
}

func (s *ScanItemRepository) SaveToFile() {
	fName := s.dbFolder + "/" + s.dbName

	s.saveStruct(s.itemStorage,  fName + "_item.mem")
	s.saveStruct(s.activityStorage,  fName + "_activity.mem")
}

func (s *ScanItemRepository) saveStruct(mem *cache.Cache, fileName string) {
	memDump := memDumpStruct{ mem.Items()}

	if len(memDump.Items) > 0 {
		if fp, err := os.Create(fileName); err == nil {
			defer fp.Close()

			gob.Register(memDumpStruct{})
			gob.Register(&scanItem.ScanItem{})
			gob.Register(scanItem.ItemActivities{})
			gob.Register([]scanItem.ItemActivity{})
			gob.Register(scanItem.ItemActivity{})

			enc := gob.NewEncoder(bufio.NewWriter(fp))
			if err := enc.Encode(memDump) ; err != nil {
				fmt.Println(fileName)
			}
		}
	}

}