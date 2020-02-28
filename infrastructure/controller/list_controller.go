package controller

import (
	"git.anphabe.net/event/anphabe-event-hub/app/sourceKeeper"
	"git.anphabe.net/event/anphabe-event-hub/domain/model/scanItem"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ListController struct {
	repo scanItem.RepositoryInterface
}

func ShowHello(c *gin.Context, listening string) {
	var publicUrl string = listening + "/public"

	settings := map[string]string{}

	settings["listen"] = listening
	settings["publicUrl"] = publicUrl
	settings["scanAgentUrl"] = listening + "/vueScanAgent"
	settings["apiUrl"] = publicUrl + "/api"

	c.HTML(http.StatusOK, "hello.tmpl", settings)
}

func ShowItemDetailJSON(c *gin.Context, sourceKeeper *sourceKeeper.Keeper) {
	repoName := c.Param("dbName")
	itemKey := c.Param("itemKey")

	item, found := sourceKeeper.GetItemDetail(repoName, itemKey)

	c.Header("Content-Type", "application/json")
	if found {
		c.JSON(http.StatusOK, item)
	} else {
		c.JSON(http.StatusNotFound, nil)
	}

}

func ShowItemDetailHTML(c *gin.Context, sourceKeeper *sourceKeeper.Keeper) {
	repoName := c.Param("dbName")
	itemKey := c.Param("itemKey")

	item, found := sourceKeeper.GetItemDetail(repoName, itemKey)

	if found {
		c.HTML(http.StatusOK, "found.tmpl", gin.H{
			"item": item,
		})
	} else {
		c.HTML(http.StatusNotFound, "not_found.tmpl", nil)
	}
}



func ShowRepositoryJSON(c *gin.Context, sourceKeeper *sourceKeeper.Keeper) {
	repoName := c.Param("dbName")
	responses := []map[string]string{}

	for _, item := range sourceKeeper.GetItems(repoName) {
		responseItem := item.Data
		responseItem["Key"] = item.Key

		responses = append(responses, responseItem)
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, responses)
}

func ShowRepositoryHTML(c *gin.Context) {
	repoName := c.Param("dbName")

	c.HTML(http.StatusOK, "repository_items.tmpl", gin.H{
		"repoName": repoName,
	})
}

//
// url: /admin/qr-check/:dbName?Key=%qrData%&activity=asdfadsf
// params:
//	Key: unique id
//  ActivityName: xyz
//  extra_fields: key1=val1,key2=val2
func ScanCheckJSON(c *gin.Context, sourceKeeper *sourceKeeper.Keeper) {
	params := map[string]string{}

	for _, param := range c.Params {
		params[param.Key] = param.Value
	}

	values := c.Request.URL.Query()
	for key, _ := range values {
		params[key] = values.Get(key)
	}

	repoName := params["dbName"]
	activityName := params["activityName"]
	itemKey, found := params["itemKey"]

	if !found {
		itemKey = params["key"]
	}

	delete(params, "dbName")
	delete(params, "key")
	delete(params, "itemKey")
	delete(params, "activityName")

	c.Header("Content-Type", "application/json")
	if item, found := sourceKeeper.ScanItem(repoName, itemKey, activityName, params); found {
		c.JSON(http.StatusOK, item)
	} else {
		c.JSON(http.StatusNotFound, nil)
	}
}

//
// url: /admin/qr-check/:dbName/:itemKey?activity=asdfadsf
// url: /admin/qr-check/:dbName?Key=%qrData%&activity=asdfadsf
// params:
//	itemKey: unique id (barcode)
// 	Key: unique id (barcode)
//  ActivityName: xyz
//  extra_fields: key1=val1,key2=val2
func ScanCheckHTML(c *gin.Context, sourceKeeper *sourceKeeper.Keeper) {
	params := map[string]string{}

	for _, param := range c.Params {
		params[param.Key] = param.Value
	}

	values := c.Request.URL.Query()
	for key, _ := range values {
		params[key] = values.Get(key)
	}

	repoName := params["dbName"]
	activityName := params["activityName"]
	itemKey, foundKey := params["itemKey"]

	if !foundKey {
		itemKey = params["key"]
	}

	delete(params, "dbName")
	delete(params, "key")
	delete(params, "itemKey")
	delete(params, "activityName")

	item, found := sourceKeeper.ScanItem(repoName, itemKey, activityName, params)

	if found {
		c.HTML(http.StatusOK, "found.tmpl", extractMap(item))
	} else {
		c.HTML(http.StatusNotFound, "not_found.tmpl", gin.H{ "key" : itemKey})
	}
}

func StartImport(c *gin.Context, sourceKeeper *sourceKeeper.Keeper) {
	repoName := c.Param("dbName")

	c.Header("Content-Type", "application/json")
	if found := sourceKeeper.StartImport(repoName); found {
		c.JSON(http.StatusOK, "ok" )
	} else {
		c.JSON(http.StatusNotFound, nil)
	}
}

func extractMap(item *scanItem.ItemDetail) map[string]interface{} {
	result := make(map[string]interface{})

	result["item"] = item

	for key, value := range item.Data {
		result[key] = value
	}

	return result
}