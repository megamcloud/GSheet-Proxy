package scanItem

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateTestScanItem_Expect_Return1a(t *testing.T) {
	var inputKey string = "testScanItemKey1"
	var inputData = map[string]string{"field1": "data1", "field2": "data2"}
	var expected, _ = NewScanItem(inputKey, inputData)

	assert.Equal(t, expected, CreateTestScanItem("1a"))
}

func TestCreateTestScanItem_Expect_Return1b(t *testing.T) {
	var inputKey string = "testScanItemKey1"
	var inputData = map[string]string{"field3": "data3", "field4": "data4"}
	var expected, _ = NewScanItem(inputKey, inputData)

	assert.Equal(t, expected, CreateTestScanItem("1b"))
}


func TestCreateTestScanItem_Expect_Return2(t *testing.T) {
	var inputKey string = "testScanItemKey3"
	var inputData = map[string]string{"field5": "data5", "field6": "data6"}
	var expected, _ = NewScanItem(inputKey, inputData)

	assert.Equal(t, expected, CreateTestScanItem("2"))
}


func TestNewScanItem_When_InputKeyIsEmpty_Expect_ReturnError(t *testing.T) {
	var inputKey string = ""
	var inputData = map[string]string{"field1": "data1", "field2": "data2"}
	var _, err = NewScanItem(inputKey, inputData)

	if assert.Error(t, err) {
		assert.Equal(t, errors.New("key could not be empty"), err)
	}
}

func TestNewScanItem_When_InputDataIsNil_Expect_NoError(t *testing.T) {
	var inputKey string = "myKey"
	var inputData map[string]string = nil

	var _, err = NewScanItem(inputKey, inputData)

	assert.Exactly(t, nil, err)
}

func TestNewScanItem_When_InputDataIsNil_Expect_ReturnData_WillEmpty(t *testing.T) {
	var inputKey string = "myKey"
	var inputData map[string]string = nil
	var sut, _ = NewScanItem(inputKey, inputData)

	assert.Exactly(t, 0, len(sut.GetData()))
}

func TestNewScanItem_When_AllInputsAreOk_Expect_Return_ScanItemInterfaceObject(t *testing.T) {
	var inputKey string = "myKey"
	var inputData = map[string]string{"field1": "data1", "field2": "data2"}
	var sut, _ = NewScanItem(inputKey, inputData)

	assert.Implements(t, (*ScanItemInterface)(nil), sut)
}

func TestScanItem_GetKey_Given_ScanItem_Expect_KeyAlwaysReturn(t *testing.T) {
	var inputKey string = "myKey"
	var inputData = map[string]string{"field1": "data1", "field2": "data2"}
	var sut, _ = NewScanItem(inputKey, inputData)

	assert.Exactly(t, inputKey, sut.GetKey())
}

func TestScanItem_GetData_Given_ScanItemCreatedWithNilData_Expect_Return_EmptyData(t *testing.T) {
	var inputKey string = "myKey"
	var inputData map[string]string = nil
	var sut, _ = NewScanItem(inputKey, inputData)

	assert.Exactly(t, 0, len(sut.GetData()))
}

func TestScanItem_GetData_Given_ScanItemCreatedNormally_Expect_Return_CorrectData(t *testing.T) {
	var inputKey string = "myKey"
	var inputData = map[string]string{"field1": "data1", "field2": "data2"}
	var sut, _ = NewScanItem(inputKey, inputData)

	assert.Exactly(t, inputData, sut.GetData())
}

func TestScanItem_GetField_Given_FieldExist_Expect_Return_CorrectFieldData(t *testing.T) {
	var inputKey string = "myKey"
	var inputData = map[string]string{"field1": "data1", "field2": "data2"}
	var sut, _ = NewScanItem(inputKey, inputData)

	var actualField, _ = sut.GetField("field1")

	assert.Exactly(t, "data1", actualField)
}

func TestScanItem_GetField_Given_FieldNotExist(t *testing.T) {
	var inputKey string = "myKey"
	var inputData = map[string]string{"field1": "data1", "field2": "data2"}
	var sut, _ = NewScanItem(inputKey, inputData)

	var _, err = sut.GetField("field0")

	if assert.Error(t, err) {
		assert.Equal(t, fmt.Errorf("field %s not exist", "field0"), err)
	}
}

func TestScanItem_SetField_Given_FieldIsNotSet_Expect_FieldDataIsSet(t *testing.T) {
	var itemKey string = "myKey"
	var fieldKey = "field1"
	var fieldData = "data1"

	var sut, _ = NewScanItem(itemKey, nil)
	sut.SetField(fieldKey, fieldData)

	actualData, _ := sut.GetField(fieldKey)

	assert.Exactly(t, fieldData, actualData)
}

func TestScanItem_SetField_Given_FieldIsAlreadySet_Expect_FieldDataIsOverwrite(t *testing.T) {
	var itemKey string = "myKey"
	var itemData = map[string]string{"field1": "already set"}
	var fieldKey = "field1"
	var fieldData = "data1"

	var sut, _ = NewScanItem(itemKey, itemData)
	actualData, _ := sut.GetField(fieldKey)
	assert.Exactly(t, "already set", actualData)

	sut.SetField(fieldKey, fieldData)
	actualData, _ = sut.GetField(fieldKey)

	assert.Exactly(t, fieldData, actualData)
}


func TestScanItem_SetFields_Given_FieldIsNotSet_Expect_FieldDataIsSet(t *testing.T) {
	var itemKey string = "myKey"
	var fields = map[string]string{"field1": "data1", "field2": "data2"}

	var sut, _ = NewScanItem(itemKey, nil)
	sut.SetFields(fields)

	actualData1, _ :=sut. GetField("field1")
	actualData2, _ := sut.GetField("field2")

	assert.Exactly(t, "data1", actualData1)
	assert.Exactly(t, "data2", actualData2)
}

func TestScanItem_SetFields_Given_FieldIsAlreadySet_Expect_FieldDataIsOverwrite(t *testing.T) {
	var itemKey string = "myKey"
	var itemData = map[string]string{"field1": "already set", "field2": "already set"}
	var fields = map[string]string{"field1": "data1", "field2": "data2"}

	var sut, _ = NewScanItem(itemKey, itemData)
	actualData, _ := sut.GetField("field1")
	assert.Exactly(t, "already set", actualData)
	actualData, _ = sut.GetField("field2")
	assert.Exactly(t, "already set", actualData)


	sut.SetFields(fields)

	actualData1, _ := sut.GetField("field1")
	actualData2, _ := sut.GetField("field2")

	assert.Exactly(t, "data1", actualData1)
	assert.Exactly(t, "data2", actualData2)
}