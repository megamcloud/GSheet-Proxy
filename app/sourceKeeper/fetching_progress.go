package sourceKeeper

import "git.anphabe.net/event/anphabe-event-hub/infrastructure/jsonapiClient"

type FetchingProgress struct {
	nextUrl       string
	data          []map[string]string
}

func newFetchingProgress(url string) *FetchingProgress {
	return &FetchingProgress{
		nextUrl:       url,
		data:          []map[string]string{},
	}
}

func (p *FetchingProgress) HasNext() bool {
	return "" != p.nextUrl
}

func (p *FetchingProgress) FetchNext() ([]map[string]string, error) {

	response, err := jsonapiClient.Get(p.nextUrl)

	if err != nil {
		return nil, err
	}

	p.nextUrl = response.Next
	p.data = response.Data

	return p.data, nil
}
