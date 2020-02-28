package jsonapiClient

import (
	"github.com/buger/jsonparser"
	"github.com/nahid/gohttp"
	"io/ioutil"
	"net/http"
	"time"
)

type JsonAPIResponse struct {
	Next string
	Data []map[string]string
}

func Get(url string) (*JsonAPIResponse, error) {
	body, err := doGet(url)

	if err != nil {
		return nil, err
	}

	json, err := parseJSON(body)

	if err != nil {
		return nil, err
	}

	return json, nil
}

func doGet(url string) ([]byte, error) {
	//github.com/nahid/gohttp
	var httpClient = &http.Client{
		Timeout: time.Second * 10,
	}

	resp, err := httpClient.Get(url)

	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return body, nil
}

func doGetAsync(url string) ([]byte, error) {
	req := gohttp.NewRequest()
	ch := make(chan *gohttp.AsyncResponse, 1)

	defer close(ch)

	req.AsyncGet(url, ch)

	op := <- ch
	if nil != op.Err {
		return nil, op.Err
	}

	return op.Resp.GetBodyAsByte()
}

func parseJSON(input []byte) (*JsonAPIResponse, error) {

	response := &JsonAPIResponse{
		Next: "",
		Data: []map[string]string{},
	}

	response.Next, _ = jsonparser.GetString(input, "links", "next")

	Data, vType, _, err := jsonparser.Get(input, "data")

	if err != nil {
		return nil, err
	}

	switch vType.String() {
	case "array":
		_, _ = jsonparser.ArrayEach(Data, func(value []byte, DataType jsonparser.ValueType, offset int, err error) {
			if item := parseItem(value); nil != item {
				response.Data = append(response.Data, item)
			}
		})

	case "object":
		if item := parseItem(Data); nil != item {
			response.Data = append(response.Data, item)
		}
	}

	return response, nil
}

func parseItem(input []byte) map[string]string {
	var item map[string]string

	_ = jsonparser.ObjectEach(input, func(key []byte, value []byte, DataType jsonparser.ValueType, offset int) error {
		if nil == item {
			item = make(map[string]string)
		}

		item[string(key)] = string(value)

		return nil
	})

	return item
}
