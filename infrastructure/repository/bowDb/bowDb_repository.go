package bowDb

import (
	"git.anphabe.net/event/anphabe-event-hub/domain/model/scanItem"
	"github.com/zippoxer/bow"
	"log"
)

type ScanItemRepository struct {
	repoName string
	conn     *bow.DB
}

type bowItem struct {
	Key  string `bow:"key"`
	Data map[string]string
}

type bowActivity struct {
	Key  string `bow:"key"`
	Data []scanItem.ItemActivity
}

func (r *ScanItemRepository) NewItem(itemKey string, data map[string]string) (*scanItem.ScanItem, error) {
	item, err := scanItem.NewScanItem(itemKey, data)

	if nil != err {
		return nil, err
	}

	r.SetItem(item)

	return item, nil
}

func (r *ScanItemRepository) SetItem(item *scanItem.ScanItem) {
	_ = r.getBucket().Put(bowItem{
		Key:  item.GetKey(),
		Data: item.GetData(),
	})
}

func (r *ScanItemRepository) GetItem(key string) (*scanItem.ScanItem, bool) {
	var item bowItem

	if err := r.getBucket().Get(key, &item); nil == err {
		return &scanItem.ScanItem{
			Key:  item.Key,
			Data: item.Data,
		}, true
	}

	return nil, false
}

func (r *ScanItemRepository) GetItemDetail(itemKey string) (*scanItem.ItemDetail, bool) {
	if item, found := r.GetItem(itemKey); found {
		activities := r.GetItemActivities(itemKey)
		return &scanItem.ItemDetail{
			ScanItem:   *item,
			Activities: activities.Activities,
		}, true
	}

	return nil, false
}

func (r *ScanItemRepository) getBucket() *bow.Bucket {
	return r.conn.Bucket(r.repoName)
}

func (r *ScanItemRepository) getActivityBucket() *bow.Bucket {
	return r.conn.Bucket(r.repoName + "_activity")
}

func (r *ScanItemRepository) Items() []*scanItem.ScanItem {
	var result []*scanItem.ScanItem

	iter := r.getBucket().Iter()
	defer iter.Close()

	var item bowItem
	for iter.Next(&item) {
		result = append(result, &scanItem.ScanItem{
			Key:  item.Key,
			Data: item.Data,
		})

		item = bowItem{}
	}

	if iter.Err() != nil {
		log.Fatal(iter.Err())
	}

	return result
}


//////////////////////

func (r *ScanItemRepository) AddItemActivity(itemKey string, activity scanItem.ItemActivity) *scanItem.ItemActivities {

	item := r.GetItemActivities(itemKey)

	if nil != item {
		item.Activities = append([]scanItem.ItemActivity{activity}, item.Activities...)

		_ = r.getActivityBucket().Put(bowActivity{
			Key:  item.Key,
			Data: item.Activities,
		})
	}

	return item
}

func (r *ScanItemRepository) GetItemActivities(itemKey string) *scanItem.ItemActivities {

	if item, found := r.GetItem(itemKey); found {

		var dbItem bowActivity

		itemActivities := &scanItem.ItemActivities{Key:item.Key, Activities: nil}

		if err := r.getActivityBucket().Get(itemKey, &dbItem); nil == err {
			itemActivities.Activities = dbItem.Data
		}

		return itemActivities
	}

	return nil
}

func (r *ScanItemRepository) GetRepoName() string {
	return r.repoName
}

func (r *ScanItemRepository) Len() int {
	return 0
}

func (r *ScanItemRepository) CloseDb() {
	if err := r.conn.Close(); err != nil {
		log.Fatal(err.Error())
	}
}
