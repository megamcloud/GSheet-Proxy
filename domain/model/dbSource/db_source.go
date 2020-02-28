package dbSource

import (
	"git.anphabe.net/event/anphabe-event-hub/config"
	"log"
	"net/url"
	"regexp"
	"strconv"
)

const BatchSize int = 25

// DbSourceConfig ...
type DbSource struct {
	Name           string
	FetchingUrl    string
	FetchingFormat string
	UpdateUrl      string
	UpdateMethod   string
}

func NewDBSource(cfg config.DbSource) *DbSource {
	return &DbSource{
		Name:           cfg.Name,
		FetchingUrl:    cfg.FetchingUrl,
		FetchingFormat: cfg.FetchingFormat,
		UpdateUrl:      cfg.UpdateUrl,
		UpdateMethod:   cfg.UpdateMethod,
	}
}

func (i *DbSource) GetFetchingUrl(startOffset int, size int) string {

	if 0 == size {
		size = BatchSize
	}

	var offsetRegex = regexp.MustCompile(`%offset%`)
	var sizeRegex = regexp.MustCompile(`%size%`)

	realUrl := offsetRegex.ReplaceAllString(i.FetchingUrl, strconv.Itoa(startOffset))
	realUrl = sizeRegex.ReplaceAllString(realUrl, strconv.Itoa(size))

	return realUrl
}

func (i *DbSource) GetUpdateUrl(key string, params map[string]string) string {
	var keyRegex = regexp.MustCompile(`%key%`)

	u, err := url.Parse(keyRegex.ReplaceAllString(i.UpdateUrl, key))

	if err != nil {
		log.Fatal(err)
	}

	q := u.Query()
	for key, value := range params {
		q.Set(key, value)
	}

	u.RawQuery = q.Encode()

	return u.String()
}

