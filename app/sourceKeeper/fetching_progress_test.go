package sourceKeeper

import (
	"errors"
	"git.anphabe.net/event/anphabe-event-hub/domain/model/dbSource"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
	"net/http"
	"testing"
)

func Test_FetchingProgress(t *testing.T) {
	fakeDomain := "http://stub.com"
	fakePath := "/sample"
	fakeUrl := fakeDomain + fakePath + "?var1=x&var2=y&offset=%offset%&limit=%size%"

	tests := []struct {
		name       string
		givenBody  string
		wantParams map[string]string
		wantReturn []map[string]string
		hasNext    bool
		wantErr    bool
	}{
		{
			name: "Fetch_Step1__",
			givenBody: `
				{  
				   "links":{  
					  "self":"http://stub.com/sample?var1=x&var2=y&offset=0&limit=2", 
					  "next":"http://stub.com/sample?var1=x&var2=y&offset=2&limit=2"
				   },
				   "data":[        
					  {  
						 "row_":12,
						 "field_companyName":"Abbott Laboratories S.A.",
						 "field_Sales":"cThanh",
					  },
					  {  
						 "row_":17,
						 "field_companyName":"Central Group",
						 "field_Sales":"cThanh",
					  },
				   ]
				}
				`,
			wantParams: map[string]string{
				"offset": "0",
				"limit":  "2",
			},
			wantReturn: []map[string]string{
				{
					"row_":              "12",
					"field_companyName": "Abbott Laboratories S.A.",
					"field_Sales":       "cThanh",
				},
				{
					"row_":              "17",
					"field_companyName": "Central Group",
					"field_Sales":       "cThanh",
				},
			},
			hasNext: true,
			wantErr: false,
		},

		{
			name: "Fetch_Step2__",
			givenBody: `
				{  
				   "links":{  
					  "self":"http://stub.com/sample?var1=x&var2=y&offset=2&limit=2", 
					  "next":"http://stub.com/sample?var1=x&var2=y&offset=4&limit=2"
				   },
				   "data":[        
					  {  
						 "row_":12,
						 "field_companyName":"Abbott Laboratories S.A.",
						 "field_Sales":"cThanh",
					  },
					  {  
						 "row_":17,
						 "field_companyName":"Central Group",
						 "field_Sales":"cThanh",
					  },
				   ]
				}
				`,
			wantParams: map[string]string{
				"offset": "2",
				"limit":  "2",
			},
			wantReturn: []map[string]string{
				{
					"row_":              "12",
					"field_companyName": "Abbott Laboratories S.A.",
					"field_Sales":       "cThanh",
				},
				{
					"row_":              "17",
					"field_companyName": "Central Group",
					"field_Sales":       "cThanh",
				},
			},
			hasNext: true,
			wantErr: false,
		},

		{
			name: "Fetch_Step3__",
			givenBody: `
				{  
				   "links":{  
					  "self":"http://stub.com/sample?var1=x&var2=y&offset=4&limit=2", 
					  "next":""
				   },
				   "data":[        
					  {  
						 "row_":12,
						 "field_companyName":"Abbott Laboratories S.A.",
						 "field_Sales":"cThanh",
					  },
					  {  
						 "row_":17,
						 "field_companyName":"Central Group",
						 "field_Sales":"cThanh",
					  },
				   ]
				}
				`,
			wantParams: map[string]string{
				"offset": "4",
				"limit":  "2",
			},
			wantReturn: []map[string]string{
				{
					"row_":              "12",
					"field_companyName": "Abbott Laboratories S.A.",
					"field_Sales":       "cThanh",
				},
				{
					"row_":              "17",
					"field_companyName": "Central Group",
					"field_Sales":       "cThanh",
				},
			},
			hasNext: false,
			wantErr: false,
		},
	}

	// CLEAN UP
	defer gock.Off()

	// ARRANGE

	for _, tt := range tests {
		gock.New(fakeDomain).
			Get(fakePath).
			MatchParams(tt.wantParams).
			Reply(http.StatusOK).
			BodyString(tt.givenBody)
	}

	sut := newFetchingProgress(setupFakeDBSource(fakeUrl).GetFetchingUrl(0, 2))

	for _, tt := range tests {

		// ACT
		got, err := sut.FetchNext()

		// ASSERT
		if (err != nil) != tt.wantErr {
			t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
			return
		}

		assert.Exactly(t, tt.wantReturn, got)
		assert.Exactly(t, tt.hasNext, sut.HasNext())
	}
}

func TestFetchingProgress_FetchNext_Given__Server_WillError__When__Fetch__Expect__ReturnError(t *testing.T) {
	defer gock.Off()

	fakeDomain := "http://stub.com"
	fakePath := "/sample"
	fakeError := "there is some error"

	gock.New(fakeDomain).
		Get(fakePath).
		ReplyError(errors.New(fakeError))

	sut := newFetchingProgress(setupFakeDBSource(fakeDomain+fakePath).GetFetchingUrl(0, 2))
	_, err := sut.FetchNext()

	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), fakeError)
	}
}

func setupFakeDBSource(fetchingUrl string) *dbSource.DbSource {
	return &dbSource.DbSource{
		Name:           "testSource",
		FetchingUrl:    fetchingUrl,
		FetchingFormat: "json",
		UpdateUrl:      "",
		UpdateMethod:   "",
	}
}
