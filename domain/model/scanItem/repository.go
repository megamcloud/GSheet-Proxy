package scanItem

type RepositoryInterface interface {
	GetRepoName() string
	NewItem(itemKey string, data map[string]string) (*ScanItem, error)
	SetItem(item *ScanItem)
	GetItem(key string) (*ScanItem, bool)
	GetItemDetail(key string) (*ItemDetail, bool)
	Items() []*ScanItem
	Len() int
	// latest item on top
	GetItemActivities(itemKey string) *ItemActivities
	AddItemActivity(itemKey string, activity ItemActivity) *ItemActivities
	CloseDb()
}
