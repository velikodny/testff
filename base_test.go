package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
	Comment  string
	Request  body
	Expected body
}

type body struct {
	JSON       string
	StatusCode int
}

var testCases = []testCase{
	{
		Comment: "a body has wrong check type",
		Request: body{
			JSON: `[
				{
				  "checkType": "DEVICE1",
				  "activityType": "SIGNUP",
				  "checkSessionKey": "string",
				  "activityData": [
					{
					  "kvpKey": "ip.address",
					  "kvpValue": "true",
					  "kvpType": "general.bool"
					}
				  ]
				}
			  ]`},
		Expected: body{
			StatusCode: 500},
	},
	{
		Comment: "without error into a body",
		Request: body{
			JSON: `[
				{
				  "checkType": "DEVICE",
				  "activityType": "SIGNUP",
				  "checkSessionKey": "string",
				  "activityData": [
					{
					  "kvpKey": "ip.address",
					  "kvpValue": "true",
					  "kvpType": "general.bool"
					}
				  ]
				}
			  ]`},
		Expected: body{
			JSON:       "{\"puppy\":true}",
			StatusCode: 200},
	},
	{
		Comment: "a body has empty string",
		Request: body{
			JSON: ""},
		Expected: body{
			StatusCode: 500},
	},
	{
		Comment: "a body has wrong activity type",
		Request: body{
			JSON: `[
				{
				  "checkType": "DEVICE",
				  "activityType": "SIGNUP1",
				  "checkSessionKey": "string",
				  "activityData": [
					{
					  "kvpKey": "ip.address",
					  "kvpValue": "true",
					  "kvpType": "general.bool"
					}
				  ]
				}
			  ]`},
		Expected: body{
			StatusCode: 500},
	},
	{
		Comment: "a body has wrong activity data",
		Request: body{
			JSON: `[
				{
				  "checkType": "DEVICE",
				  "activityType": "SIGNUP",
				  "checkSessionKey": "string1",
				  "activityData": [
					{
					  "kvpKey": "ip.address",
					  "kvpValue": "true",
					  "kvpType": ""
					}
				  ]
				}
			  ]`},
		Expected: body{
			JSON:       "",
			StatusCode: 500},
	},
}

func TestValidCheckType(t *testing.T) {
	validator := Validator{}

	assert.Equal(t, true, validator.validCheckType("BIOMETRIC"))
	assert.Equal(t, false, validator.validCheckType("SIG"))

}

func TestValidActivityType(t *testing.T) {
	validator := Validator{}

	assert.Equal(t, true, validator.validActivityType("LOGIN"))
	assert.Equal(t, false, validator.validActivityType("LOGIN_2"))
}

func TestValidActivityData(t *testing.T) {
	validator := Validator{}
	correctData := []ActivityData{
		ActivityData{
			KvpKey:   "Test 1",
			KvpValue: "1",
			KvpType:  "general.integer",
		},
		ActivityData{
			KvpKey:   "Test 2",
			KvpValue: "true",
			KvpType:  "general.bool",
		},
		ActivityData{
			KvpKey:   "Test 3",
			KvpValue: "2.33",
			KvpType:  "general.float",
		},
		ActivityData{
			KvpKey:   "Test 4",
			KvpValue: "2",
			KvpType:  "general.string",
		},
	}

	assert.Equal(t, true, validator.validActivityData(correctData))

	incorrectDataKey := []ActivityData{
		ActivityData{
			KvpKey:   "Test",
			KvpValue: "true",
			KvpType:  "general.bool",
		},
		ActivityData{
			KvpKey:   "Test",
			KvpValue: "true",
			KvpType:  "general.bool",
		},
	}

	assert.Equal(t, false, validator.validActivityData(incorrectDataKey))

	incorrectKvpType := []ActivityData{
		ActivityData{
			KvpKey:   "Test",
			KvpValue: "food",
			KvpType:  "general.bool",
		},
	}

	assert.Equal(t, false, validator.validActivityData(incorrectKvpType))

	incorrectKvpType = []ActivityData{
		ActivityData{
			KvpKey:   "Test",
			KvpValue: "food",
			KvpType:  "general.integer",
		},
	}

	assert.Equal(t, false, validator.validActivityData(incorrectKvpType))

	incorrectKvpType = []ActivityData{
		ActivityData{
			KvpKey:   "Test",
			KvpValue: "food",
			KvpType:  "general.float",
		},
	}

	assert.Equal(t, false, validator.validActivityData(incorrectKvpType))

	incorrectKvpType = []ActivityData{
		ActivityData{
			KvpKey:   "Test",
			KvpValue: "food",
			KvpType:  "general.non",
		},
	}

	assert.Equal(t, false, validator.validActivityData(incorrectKvpType))
}

func newRecorder(req *http.Request, method string, strPath string, fnHandler func(w http.ResponseWriter, r *http.Request, param httprouter.Params)) *httptest.ResponseRecorder {
	router := httprouter.New()
	router.Handle(method, strPath, fnHandler)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	return rec
}

func TestIsGoodNilBody(t *testing.T) {
	req, err := http.NewRequest("POST", "/isgood", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := newRecorder(req, "POST", "/isgood", isGood)
	if rec.Code != 500 {
		t.Error("Expected response code to be 200")
	}
}
func TestIsGoodErrorCheckType(t *testing.T) {

	for _, item := range testCases {

		body := item.Request.JSON

		req, err := http.NewRequest("POST", "/isgood", strings.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}

		rec := newRecorder(req, "POST", "/isgood", isGood)
		if rec.Code != item.Expected.StatusCode {
			t.Error("Expected response code to be 200")
		}

		if item.Expected.JSON != "" {
			assert.Equal(t, item.Expected.JSON, rec.Body.String())
		}
	}
}
