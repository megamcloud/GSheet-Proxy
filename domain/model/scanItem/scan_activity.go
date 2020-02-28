package scanItem

import "time"

type ItemActivities struct {
	Key        string
	Activities []ItemActivity
}

type ItemActivity struct {
	Action  string
	Data    map[string]string
	Created time.Time
}

type ItemDetail struct {
	ScanItem
	Activities []ItemActivity
}

func NewActivity(action string, properties map[string]string ) ItemActivity {
	return ItemActivity{
		Action:  action,
		Data:    properties,
		Created: time.Now(),
	}
}