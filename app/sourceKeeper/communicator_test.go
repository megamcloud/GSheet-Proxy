package sourceKeeper_test

import (
	"fmt"
	"git.anphabe.net/event/anphabe-event-hub/app/sourceKeeper"
	"git.anphabe.net/event/anphabe-event-hub/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	"gopkg.in/h2non/gock.v1"
	"net/http"
)

type fakeDbSource struct {
	Domain      string
	Path        string
	Name        string
	IdField     string
	QueryString string
	Steps       []struct {
		name       string
		givenBody  string
		wantParams map[string]string
		wantReturn []map[string]string
	}
}

var _ = Describe("Communicator\n", func() {
	a := getFakeDbSource1()

	Context(fmt.Sprintf("Given a DBSource will be setup at %s\n", a.Domain), func() {
		defer gock.Off()

		Context(fmt.Sprintf("Given it has %d pages\n", len(a.Steps)), func() {

			for _, tt := range a.Steps {
				gock.New(a.Domain).
					Get(a.Path).
					MatchParams(tt.wantParams).
					Reply(http.StatusOK).
					BodyString(tt.givenBody)
			}

			fakeDbSource := config.DbSource{
				Name:           a.Name,
				IdField:        a.IdField,
				FetchingUrl:    a.Domain + a.Path + a.QueryString,
				FetchingFormat: "json",
				UpdateUrl:      "",
				UpdateMethod:   "",
			}

			stubLogger, _ := fakeLogger()
			keeper := sourceKeeper.NewCommunicator(fakeDbSource, stubLogger)

			Context("Calling to Import()\n", func() {
				_ = keeper.Import(func(repoName string, itemKey string, data []map[string]string) int {

					It("Every items must be sent into callback\n", func() {
						Expect(repoName).To(Equal(fakeDbSource.Name))
						Expect(itemKey).To(Equal(fakeDbSource.IdField))
					})

					return len(data)
				})
			})
		})
	})
})

func getFakeDbSource1() fakeDbSource {
	a := fakeDbSource{
		Domain:      "http://stub1.com",
		Path:        "/sample",
		Name:        "dbTest1",
		IdField:     "field_qrcode",
		QueryString: "?var1=x&var2=y&offset=%offset%&limit=%size%",
		Steps: []struct {
			name       string
			givenBody  string
			wantParams map[string]string
			wantReturn []map[string]string
		}{
			{
				name: "Fetch_Step1__",
				givenBody: `
				{
				   "links":{
					  "self":"http://stub1.com/sample?var1=x&var2=y&offset=0&limit=2",
					  "next":"http://stub1.com/sample?var1=x&var2=y&offset=2&limit=2"
				   },
				   "data":[
					  {
						 "row_":12,
						 "field_qrcode": "101",
						 "field_companyName":"Abbott Laboratories S.A.",
						 "field_Sales":"cThanh",
					  },
					  {
						 "row_":17,
						 "field_qrcode": "102",
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
						"field_qrcode":      "101",
						"field_companyName": "Abbott Laboratories S.A.",
						"field_Sales":       "cThanh",
					},
					{
						"row_":              "17",
						"field_qrcode":      "102",
						"field_companyName": "Central Group",
						"field_Sales":       "cThanh",
					},
				},
			},

			{
				name: "Fetch_Step2__",
				givenBody: `
				{
				   "links":{
					  "self":"http://stub1.com/sample?var1=x&var2=y&offset=2&limit=2",
					  "next":"http://stub1.com/sample?var1=x&var2=y&offset=4&limit=2"
				   },
				   "data":[
					  {
						 "row_":12,
						 "field_qrcode": "103",
						 "field_companyName":"Abbott Laboratories S.A.",
						 "field_Sales":"cThanh",
					  },
					  {
						 "row_":17,
						 "field_qrcode": "104",
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
						"field_qrcode":      "103",
						"field_companyName": "Abbott Laboratories S.A.",
						"field_Sales":       "cThanh",
					},
					{
						"row_":              "17",
						"field_qrcode":      "104",
						"field_companyName": "Central Group",
						"field_Sales":       "cThanh",
					},
				},
			},

			{
				name: "Fetch_Step3__",
				givenBody: `
				{
				   "links":{
					  "self":"http://stub1.com/sample?var1=x&var2=y&offset=4&limit=2",
					  "next":""
				   },
				   "data":[
					  {
						 "row_":12,
						 "field_qrcode": "105",
						 "field_companyName":"Abbott Laboratories S.A.",
						 "field_Sales":"cThanh",
					  },
					  {
						 "row_":17,
	                     "field_qrcode": "106",
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
						"field_qrcode":      "105",
						"field_companyName": "Abbott Laboratories S.A.",
						"field_Sales":       "cThanh",
					},
					{
						"row_":              "17",
						"field_qrcode":      "106",
						"field_companyName": "Central Group",
						"field_Sales":       "cThanh",
					},
				},
			},
		},
	}

	return a
}

func getFakeDbSource2() fakeDbSource {
	a := fakeDbSource{
		Domain:      "http://stub2.com",
		Path:        "/sample",
		Name:        "dbTest2",
		IdField:     "field_qrcode",
		QueryString: "?var1=x&var2=y&offset=%offset%&limit=%size%",
		Steps: []struct {
			name       string
			givenBody  string
			wantParams map[string]string
			wantReturn []map[string]string
		}{
			{
				name: "Fetch_Step1__",
				givenBody: `
				{
				   "links":{
					  "self":"http://stub2.com/sample?var1=x&var2=y&offset=0&limit=2",
					  "next":"http://stub2.com/sample?var1=x&var2=y&offset=2&limit=2"
				   },
				   "data":[
					  {
						 "row_":12,
						 "field_qrcode": "201",
						 "field_companyName":"Abbott Laboratories S.A.",
						 "field_Sales":"cThanh",
					  },
					  {
						 "row_":17,
						 "field_qrcode": "202",
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
						"field_qrcode":      "201",
						"field_companyName": "Abbott Laboratories S.A.",
						"field_Sales":       "cThanh",
					},
					{
						"row_":              "17",
						"field_qrcode":      "202",
						"field_companyName": "Central Group",
						"field_Sales":       "cThanh",
					},
				},
			},

			{
				name: "Fetch_Step2__",
				givenBody: `
				{
				   "links":{
					  "self":"http://stub2.com/sample?var1=x&var2=y&offset=2&limit=2",
					  "next":"http://stub2.com/sample?var1=x&var2=y&offset=4&limit=2"
				   },
				   "data":[
					  {
						 "row_":12,
						 "field_qrcode": "203",
						 "field_companyName":"Abbott Laboratories S.A.",
						 "field_Sales":"cThanh",
					  },
					  {
						 "row_":17,
						 "field_qrcode": "204",
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
						"field_qrcode":      "203",
						"field_companyName": "Abbott Laboratories S.A.",
						"field_Sales":       "cThanh",
					},
					{
						"row_":              "17",
						"field_qrcode":      "204",
						"field_companyName": "Central Group",
						"field_Sales":       "cThanh",
					},
				},
			},

			{
				name: "Fetch_Step3__",
				givenBody: `
				{
				   "links":{
					  "self":"http://stub2.com/sample?var1=x&var2=y&offset=4&limit=2",
					  "next":""
				   },
				   "data":[
					  {
						 "row_":12,
						 "field_qrcode": "205",
						 "field_companyName":"Abbott Laboratories S.A.",
						 "field_Sales":"cThanh",
					  },
					  {
						 "row_":17,
	                     "field_qrcode": "206",
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
						"field_qrcode":      "205",
						"field_companyName": "Abbott Laboratories S.A.",
						"field_Sales":       "cThanh",
					},
					{
						"row_":              "17",
						"field_qrcode":      "206",
						"field_companyName": "Central Group",
						"field_Sales":       "cThanh",
					},
				},
			},
		},
	}

	return a
}

func fakeLogger() (*zap.Logger, *observer.ObservedLogs) {
	zapCore, observedLogs := observer.New(zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	}))

	return zap.New(zapCore), observedLogs
}
