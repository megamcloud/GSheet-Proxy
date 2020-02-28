package dbSource

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
	"net/http"

	"reflect"
	"testing"
)

func TestNewDbSource(t *testing.T) {
	sut := setupFakeDBSource("https://stub.com/sample")

	assert.Exactly(t, "*dbSource.DbSource", reflect.TypeOf(sut).String())
}

func Test_DbSource_StartFetching__Expect__Always_Return_FetchingProgressPointer(t *testing.T) {
	sut := setupFakeDBSource("https://stub.com/sample")

	assert.Exactly(t, "*dbSource.FetchingProgress", reflect.TypeOf(sut.StartFetching(0, 0)).String())
}

func Test_DbSource_StartFetching__When__CallTo_FetchingProgress_HasNext__Expect__AlwaysReturnTRUE(t *testing.T) {
	sut := setupFakeDBSource("https://stub.com/sample")

	assert.Exactly(t, true, sut.StartFetching(0,0).HasNext())
}

func Test_DbSource_StartFetching__Given__NoProgressStarted__Expect__Always_Return_FetchingProgressPointer(t *testing.T) {
	sut := setupFakeDBSource("https://stub.com/sample")

	assert.Exactly(t, "*dbSource.FetchingProgress", reflect.TypeOf(sut.StartFetching(0,0)).String())
}

func Test_DbSource_StartFetching__Given__ProgressStarted__Expect__Always_Return_FetchingProgressPointer(t *testing.T) {
	sut := setupFakeDBSource("https://stub.com/sample")
	progress := sut.StartFetching(0, 0)

	assert.Exactly(t, progress, sut.StartFetching(0,0))
}

func Test_DbSource_StartFetching_On_Replace_OffsetPlaceHolder(t *testing.T) {
	fakeDomain := "http://stub.com"
	fakePath := "/sample"
	fakeBody := `{"links":{},"data":[]}`

	tests := []struct {
		name           string
		givenPath      string
		givenValue     int
		wantParamName  string
		wantParamValue string
		wantReturn     []map[string]string
		wantErr        bool
	}{
		{
			name:           "__Given__Missing_%offset%_PLACEHOLDER__Expect__UrlWithCorrectOffset_IsRequested",
			givenValue:     -1,
			givenPath:      "/sample?var1=x&var2=y&offset=",
			wantParamName:  "offset",
			wantParamValue: "",
			wantReturn:     []map[string]string{},
			wantErr:        false,
		},

		{
			name:           "__Given__Offset_AtEndingOf_QueryString__Expect__UrlWithCorrectOffset_IsRequested",
			givenValue:     100,
			givenPath:      fakePath + "?var1=x&var2=y&offset=%offset%",
			wantParamName:  "offset",
			wantParamValue: "100",
			wantReturn: []map[string]string{},
			wantErr:        false,
		},

		{
			name:           "__Given__Offset_AtBeginOf_QueryString__Expect__UrlWithCorrectOffset_IsRequested",
			givenValue:     100,
			givenPath:      fakePath + "?offset=%offset%&var1=x&var2=y",
			wantParamName:  "offset",
			wantParamValue: "100",
			wantReturn: []map[string]string{},
			wantErr:        false,
		},

		{
			name:           "__Given__Offset_InMiddleOf_QueryString__Expect__UrlWithCorrectOffset_IsRequested",
			givenValue:     100,
			givenPath:      fakePath + "?var1=x&offset=%offset%&var2=y",
			wantParamName:  "offset",
			wantParamValue: "100",
			wantReturn: []map[string]string{},
			wantErr:        false,
		},

		{
			name:           "__Given__NegativeValueOffset__Expect__UrlWithCorrectOffset_IsRequested",
			givenValue:     -1,
			givenPath:      "/sample?var1=x&var2=y&offset=%offset%",
			wantParamName:  "offset",
			wantParamValue: "-1",
			wantReturn: []map[string]string{},
			wantErr:        false,
		},

		{
			name:           "__Given__ZeroValueOffset__Expect__UrlWithCorrectOffset_IsRequested",
			givenValue:     0,
			givenPath:      fakePath + "?var1=x&var2=y&offset=%offset%",
			wantParamName:  "offset",
			wantParamValue: "0",
			wantReturn: []map[string]string{},
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// CLEAN UP
			defer gock.Off()

			// ARRANGE
			gock.New(fakeDomain).
				Get("/sample").
				MatchParam(tt.wantParamName, tt.wantParamValue).
				Reply(http.StatusOK).
				BodyString(fakeBody)

			sut := setupFakeDBSource(fakeDomain + tt.givenPath)

			// ACT
			got, err := sut.StartFetching(tt.givenValue, 0).FetchNext()

			// ASSERT
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Exactly(t, tt.wantReturn, got)
		})
	}
}

func Test_DbSource_StartFetching_On_Replace_SizePlaceHolder(t *testing.T) {
	fakeDomain := "http://stub.com"
	fakePath := "/sample"
	fakeBody := `{"links":{},"data":[]}`

	tests := []struct {
		name           string
		givenPath      string
		givenValue     int
		wantParamName  string
		wantParamValue string
		wantReturn     []map[string]string
		wantErr        bool
	}{
		{
			name:           "__Given__Missing_%size%_PLACEHOLDER__Expect__UrlWithCorrectSize_IsRequested",
			givenValue:     -1,
			givenPath:      "/sample?var1=x&var2=y&limit=",
			wantParamName:  "limit",
			wantParamValue: "",
			wantReturn:     []map[string]string{},
			wantErr:        false,
		},

		{
			name:           "__Given__Offset_AtEndingOf_QueryString__Expect__UrlWithCorrectSize_IsRequested",
			givenValue:     100,
			givenPath:      fakePath + "?var1=x&var2=y&limit=%size%",
			wantParamName:  "limit",
			wantParamValue: "100",
			wantReturn:     []map[string]string{},
			wantErr:        false,
		},

		{
			name:           "__Given__Offset_AtBeginOf_QueryString__Expect__UrlWithCorrectSize_IsRequested",
			givenValue:     100,
			givenPath:      fakePath + "?limit=%size%&var1=x&var2=y",
			wantParamName:  "limit",
			wantParamValue: "100",
			wantReturn:     []map[string]string{},
			wantErr:        false,
		},

		{
			name:           "__Given__Offset_InMiddleOf_QueryString__Expect__UrlWithCorrectSize_IsRequested",
			givenValue:     100,
			givenPath:      fakePath + "?var1=x&limit=%size%&var2=y",
			wantParamName:  "limit",
			wantParamValue: "100",
			wantReturn:     []map[string]string{},
			wantErr:        false,
		},

		{
			name:           "__Given__NegativeValueOffset__Expect__UrlWithCorrectSize_IsRequested",
			givenValue:     -1,
			givenPath:      "/sample?var1=x&var2=y&limit=%size%",
			wantParamName:  "limit",
			wantParamValue: "-1",
			wantReturn:     []map[string]string{},
			wantErr:        false,
		},

		{
			name:           "__Given__ZeroValueSize__Expect__UrlWith_DefaultBatchSize_IsRequested",
			givenValue:     0,
			givenPath:      fakePath + "?var1=x&var2=y&limit=%size%",
			wantParamName:  "limit",
			wantParamValue: "25",
			wantReturn:     []map[string]string{},
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// CLEAN UP
			defer gock.Off()

			// ARRANGE
			gock.New(fakeDomain).
				Get("/sample").
				MatchParam(tt.wantParamName, tt.wantParamValue).
				Reply(http.StatusOK).
				BodyString(fakeBody)

			sut := setupFakeDBSource(fakeDomain + tt.givenPath)

			// ACT
			got, err := sut.StartFetching(0, tt.givenValue).FetchNext()

			// ASSERT
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Exactly(t, tt.wantReturn, got)
		})
	}
}


func TestDbSource_StartFetching_RealSource(t *testing.T) {
	dbSource := &DbSource{
		Name:           "testSource",
		FetchingUrl:    "https://script.google.com/macros/s/AKfycbyKxlzZMiVlF01ZGPAXYsY0ARV-L8V04QCgONo5kIbTAwkfOC4C/exec?path=/sample&order=field_qrcode&offset=%offset%&limit=%size%",
		FetchingFormat: "json",
		UpdateUrl:      "",
		UpdateMethod:   "",
	}

	progress := dbSource.StartFetching(0, 100)

}

func setupFakeDBSource(fetchingUrl string) *DbSource {
	return &DbSource{
		Name:           "testSource",
		FetchingUrl:    fetchingUrl,
		FetchingFormat: "json",
		UpdateUrl:      "",
		UpdateMethod:   "",
	}
}
