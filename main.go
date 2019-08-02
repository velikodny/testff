package main

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Data a base data of a handler
type Data struct {
	CheckType       string         `json:"checkType"`
	ActivityType    string         `json:"activityType"`
	CheckSessionKey string         `json:"checkSessionKey"`
	ActivityData    []ActivityData `json:"activityData"`
}

// ActivityData to describe activity data
type ActivityData struct {
	KvpKey   string `json:"kvpKey"`
	KvpValue string `json:"kvpValue"`
	KvpType  string `json:"kvpType"`
}

// Error a custom error
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Puppy struct {
	Puppy bool `json:"puppy"`
}

func main() {
	router := httprouter.New()
	router.POST("/isgood", isGood)
	http.ListenAndServe(":8888", router)
}

func isGood(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// var req []Data
	req := new([]Data)
	validator := new(Validator)
	if r.Body == nil {
		returnError(w, &Error{Message: "no pass some body content"})
		return
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		returnError(w, &Error{Message: "fault to decode passed data"})
		return
	}

	for _, item := range *req {
		switch {
		case !validator.validCheckType(item.CheckType):
			fallthrough
		case !validator.validActivityType(item.ActivityType):
			break
		case !validator.validActivityData(item.ActivityData):
			break
		}
	}

	if validator.err != nil {
		returnError(w, validator.err)
		return
	}

	puppy := Puppy{
		Puppy: true,
	}
	resp, _ := json.Marshal(puppy)

	w.WriteHeader(200)
	w.Write(resp)
}

func returnError(w http.ResponseWriter, err *Error) {
	jsonError := []byte{}
	if err != nil {
		jsonError, _ = json.Marshal(err)
	}
	w.WriteHeader(500)
	w.Write(jsonError)
}
