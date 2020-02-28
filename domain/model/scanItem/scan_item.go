// item_test.go
package scanItem

import (
	"errors"
	"fmt"
)

type ScanItemInterface interface {
	GetKey() string
	GetData() map[string]string
	SetField(field string, data string)
	GetField(field string) (string, error)
	SetFields(fields map[string]string)
	//AddActivity(action string, properties map[string]string)
}

type ScanItem struct {
	Key  string
	Data map[string]string
}

func NewScanItem(key string, data map[string]string) (*ScanItem, error) {
	if "" == key {
		return nil, errors.New("key could not be empty")
	}

	if nil == data {
		data = make(map[string]string)
	}

	return &ScanItem{Key: key, Data: data}, nil
}

func CreateTestScanItem(whichItem string) *ScanItem {

	var item *ScanItem

	switch whichItem {

	case "1a":
		item, _ = NewScanItem("testScanItemKey1", map[string]string{"field1": "data1", "field2": "data2"})

	case "1b":
		item, _ = NewScanItem("testScanItemKey1", map[string]string{"field3": "data3", "field4": "data4"})

	case "2":
		item, _ = NewScanItem("testScanItemKey3", map[string]string{"field5": "data5", "field6": "data6"})
	}

	return item
}

func (i *ScanItem) GetKey() string {
	return i.Key
}

func (i *ScanItem) GetData() map[string]string {
	return i.Data
}

func (i *ScanItem) GetField(field string) (string, error) {
	if val, ok := i.Data[field]; ok {
		return val, nil
	}

	return "", fmt.Errorf("field %s not exist", field)
}

func (i *ScanItem) SetField(field string, data string) {
	i.Data[field] = data
}

func (i *ScanItem) SetFields(fields map[string]string) {
	for k, v := range fields {
		i.Data[k] = v
	}
}

//
//func (i *ScanItem) AddActivity(action string, properties map[string]string) {
//
//}
