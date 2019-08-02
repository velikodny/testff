package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

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

func TestIsGood(t *testing.T) {

	body := `[
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
  ]`

	req, err := http.NewRequest("POST", "/isgood", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	router := httprouter.New()
	router.Handle("POST", "/isgood", isGood)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Error("Expected response code to be 200")
	}

	expected := "{\"puppy\":true}"
	if rec.Body.String() != expected {
		t.Error("Response body does not match")
	}

}

func TestIsGoodNilBody(t *testing.T) {
	req, err := http.NewRequest("POST", "/isgood", nil)
	if err != nil {
		t.Fatal(err)
	}

	router := httprouter.New()
	router.Handle("POST", "/isgood", isGood)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != 500 {
		t.Error("Expected response code to be 200")
	}
}
