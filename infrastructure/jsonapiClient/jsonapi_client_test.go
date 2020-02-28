package jsonapiClient

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
	"net/http"
	"testing"
)

func Test_Get__On_ParsingResponse(t *testing.T) {

	tests := []struct {
		name    string
		given   string
		want    *JsonAPIResponse
		wantErr bool
	}{
		{
			name: "__Given__ResponseHas_MultiplesItems_ArrayData__When__ParseData__Expect__ReturnArrayData",
			given: `
				{  
				   "links":{  
					  "self":"i dont care now", 
					  "next":"i dont care now"
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
					  {  
						 "row_":20,
						 "field_companyName":"INSEE Vietnam",
						 "field_Sales":"cThanh",
					  },
				   ]
				}
				`,
			want: &JsonAPIResponse{
				Next: "i dont care now",
				Data: []map[string]string{
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
					{
						"row_":              "20",
						"field_companyName": "INSEE Vietnam",
						"field_Sales":       "cThanh",
					},
				},
			},
			wantErr: false,
		},

		{
			name: "__Given__ResponseHas_SingleItem_ArrayData__When__ParseData__Expect__ReturnCorrectStruct",
			given: `
				{
				  "links":{
					 "self":"i dont care now",
					  "next":"i dont care now"
				  },
				  "data":[
					 {
						"row_":12,
						"field_companyName":"Abbott Laboratories S.A.",
						"field_Sales":"cThanh",
					 },
				  ]
				}
				`,
			want: &JsonAPIResponse{
				Next: "i dont care now",
				Data: []map[string]string{
					{
						"row_":              "12",
						"field_companyName": "Abbott Laboratories S.A.",
						"field_Sales":       "cThanh",
					},
				},
			},
			wantErr: false,
		},

		{
			name: "__Given__ResponseHas_SingleItem_ObjectData__When__ParseData__Expect__ReturnCorrectStruct",
			given: `
				{
					"links":{
					   "self":"i dont care now",
					   "next":"i dont care now"
					},
					"data": {
						"row_":12,
						"field_companyName":"Abbott Laboratories S.A.",
						"field_Sales":"cThanh",
					},
				}
				`,
			want: &JsonAPIResponse{
				Next: "i dont care now",
				Data: []map[string]string{
					{
						"row_":              "12",
						"field_companyName": "Abbott Laboratories S.A.",
						"field_Sales":       "cThanh",
					},
				},
			},
			wantErr: false,
		},

		{
			name: "__Given__ResponseHas_Empty_ArrayData__When__ParseData__Expect__ReturnCorrectStruct",
			given: `
				{
				  "links":{
					 "self":"i dont care now",
					  "next":"i dont care now"
				  },
				  "data": [],
				}
				`,
			want: &JsonAPIResponse{
				Next: "i dont care now",
				Data: []map[string]string{},
			},
			wantErr: false,
		},

		{
			name: "__Given__ResponseHas_Empty_ObjectData__When__ParseData__Expect__ReturnCorrectStruct",
			given: `
				{
				  "links":{
					 "self":"i dont care now",
					  "next":"i dont care now"
				  },
				  "data": {},
				}
				`,
			want: &JsonAPIResponse{
				Next: "i dont care now",
				Data: []map[string]string{},
			},
			wantErr: false,
		},

		{
			name: "__Given__ResponseHas_Empty_StringData__When__ParseData__Expect__ReturnCorrectStruct",
			given: `
				{
				  "links":{
					 "self":"i dont care now",
					  "next":"i dont care now"
				  },
				  "data": "",
				}
				`,
			want: &JsonAPIResponse{
				Next: "i dont care now",
				Data: []map[string]string{},
			},
			wantErr: false,
		},

		{
			name: "__Given__ResponseHas_MissingData__When__ParseData__Expect__ReturnError",
			given: `
				{
				  "links":{
					 "self":"i dont care now",
					  "next":"i dont care now"
				  },
				}
				`,
			want:    nil,
			wantErr: true,
		},

		/////////////////////////

		{
			name: "__Given__ResponseHas_MissingLinks__When__ParsingLinks__Expect__ReturnEmptyNext",
			given: `
				{
				  "data": {},
				}
				`,
			want: &JsonAPIResponse{
				Next: "",
				Data: []map[string]string{},
			},
			wantErr: false,
		},

		{
			name: "__Given__ResponseHas_EmptyLinksObject__When__ParsingLinks__Expect__ReturnEmptyNext",
			given: `
				{
				  "links":{
				  },
				
				  "data": {},
				}
				`,
			want: &JsonAPIResponse{
				Next: "",
				Data: []map[string]string{},
			},
			wantErr: false,
		},

		{
			name: "__Given__ResponseHas_EmptyLinksArray__When__ParsingLinks__Expect__ReturnEmptyNext",
			given: `
				{
				  "links":[],
				  "data": {},
				}
				`,
			want: &JsonAPIResponse{
				Next: "",
				Data: []map[string]string{},
			},
			wantErr: false,
		},

		{
			name: "__Given__ResponseHas_EmptyLinksString__When__ParsingLinks__Expect__ReturnEmptyNext",
			given: `
				{
				  "links":"",
				  "data": {},
				}
				`,
			want: &JsonAPIResponse{
				Next: "",
				Data: []map[string]string{},
			},
			wantErr: false,
		},

		{
			name: "__Given__ResponseHas_LinksHasEmptyNext__When__ParsingLinks__Expect__ReturnEmptyNext",
			given: `
				{
				  "links":{
					  "next": ""
				   },
				  "data": {},
				}
				`,
			want: &JsonAPIResponse{
				Next: "",
				Data: []map[string]string{},
			},
			wantErr: false,
		},

		{
			name: "__Given__ResponseHas_LinksHasCorrectNext__When__ParsingLinks__Expect__ReturnCorrectNext",
			given: `
				{
				  "links":{
					  "next": "https://stub.com/next"
				   },
				  "data": {},
				}
				`,
			want: &JsonAPIResponse{
				Next: "https://stub.com/next",
				Data: []map[string]string{},
			},
			wantErr: false,
		},

		//{
		//	name: "",
		//	given: `
		//
		//		`,
		//	want: &JsonAPIResponse{
		//
		//	},
		//	wantErr: false,
		//},
	}

	fakeDomain := "http://stub.com"
	fakePath := "/sample"

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ARRANGE
			gock.New(fakeDomain).
				Get(fakePath).
				Reply(http.StatusOK).
				BodyString(tt.given)

			// ACT
			got, err := Get(fakeDomain + fakePath)

			// ASSERT
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Exactly(t, tt.want, got)

			// CLEAN UP
			gock.Off()
		})
	}
}


func Test_Get__Given__Server_WillRedirect__When__Fetch__Expect__FollowAndReturnContent(t *testing.T) {
	defer gock.Off()

	fakeDomain := "http://stub.com"
	fakePath := "/sample"
	redirectPath := "/redirect"

	jsonString := `
				{  
				   "links":{  
					  "self":"i dont care now", 
					  "next":"i dont care now"
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
				`

	want := &JsonAPIResponse{
		Next: "i dont care now",
		Data: []map[string]string{
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
	}

	gock.New(fakeDomain).
		Get(redirectPath).
		ReplyFunc(func(response *gock.Response) {
			response.
				Status(http.StatusFound).
				SetHeader("Location", fakePath)
		})

	gock.New(fakeDomain).
		Get(fakePath).
		Reply(http.StatusOK).
		BodyString(jsonString)


	got, _ := Get(fakeDomain + redirectPath)

	assert.Exactly(t, want, got)
}


func Test_Get__Given__Server_WillRedirect__When__Fetch__Expect__ReturnError(t *testing.T) {

	defer gock.Off()

	fakeDomain := "http://stub.com"
	fakePath := "/sample"
	fakeError := "there is some error"

	gock.New(fakeDomain).
		Get(fakePath).
		ReplyError(errors.New(fakeError))

	_, err := Get(fakeDomain + fakePath)

	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), fakeError)
	}
}