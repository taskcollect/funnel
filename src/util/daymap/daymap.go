package daymap

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/buger/jsonparser"
)

type DaymapAuthMethod struct {
	Username, Password string
	Cookies            []byte
}

/*
	Given a body json bytearray, populate it with the credentials from the
	given DaymapAuthMethod struct. Returns the resulting JSON.

	This should be used to create a JSON object that daymap-getter will use
	to authenticate.
*/
func populatePayloadWithAuth(payload []byte, auth *DaymapAuthMethod) ([]byte, error) {
	payload, err := jsonparser.Set(payload, []byte(strconv.Quote(auth.Username)), "username")
	if err != nil {
		return nil, err
	}

	payload, err = jsonparser.Set(payload, []byte(strconv.Quote(auth.Password)), "password")
	if err != nil {
		return nil, err
	}

	if auth.Cookies != nil {
		payload, err = jsonparser.Set(payload, auth.Cookies, "cookies")
		if err != nil {
			return nil, err
		}
	}

	return payload, nil
}

/*
	Given a payload with all other request params, a DaymapAuthMethod object bearing
	valid credentials, a target url and an expected return type, send a POST request
	to daymap-getter and return results, after validating that the returned type
	matches the expected type provided as a parameter.

	This executes the actual network request, after the payload is populated with
	the necessary fields by populatePayloadWithAuth().
*/
func postDaymapPayload(url string, payload []byte, auth *DaymapAuthMethod, expectedType jsonparser.ValueType) ([]byte, error) {
	// populate payload with auth info
	payload, err := populatePayloadWithAuth(payload, auth)
	if err != nil {
		return nil, errors.New("bad auth structure")
	}

	// send request to target url
	res, err := http.Post(url, "application/json", bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}

	switch res.StatusCode {
	case 401:
		return nil, errors.New("bad auth")
	case 502:
		return nil, errors.New("bad daymap")
	case 500:
		return nil, errors.New("upstream error")
	case 200:
		break
	default:
		return nil, fmt.Errorf("got weird code: %d", res.StatusCode)
	}

	resp, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	_, dataType, _, err := jsonparser.Get(resp, "data")
	if err != nil {
		return nil, err
	}

	// validate expected return type
	if dataType != expectedType {
		println(dataType, "mismatches expected", expectedType)
		return nil, errors.New("unexpected response type")
	}

	// return everything including cookies
	return resp, nil
}

/*
	----- GETTER FUNCTIONS -----

	Effectively, the cycle of any getter function should be:

	1. Make a payload JSON bytearray (containing just {})
		!!! WARNING !!! It's a huge security risk to substitute strings into
		the intitial payload state. jsonparser.Set exists for a reason.
		You should only do it if the string is STATIC eg. not substituted dynamically,
		something like { "booleanValue": true, "stringValue": "static string" }

	2. Set any additional parameters daymap-getter needs via jsonparser.Set.
	   Make sure to overwrite the original bytearray with the return value of
	   jsonparser.Set

	3. Pass all necessary values to this function (postDaymapPayload) and simply
	   return what it returns directly.

	Note that getter functions will get called from HTTP endpoint handlers!
	That means that these shouldn't talk to the database or do anything else
	than literally return the data from the getters.

	Yes, that includes "cookies" too.
*/

/*
	Given a DaymapAuthMethod & start + end times, get all the lessons for the
	user with the provided credentials within the provided timeframe.
*/
func GetLessons(auth *DaymapAuthMethod, startTime, endTime time.Time) ([]byte, error) {
	payload := []byte{'{', '}'}

	// add the timestamps
	payload, err := jsonparser.Set(payload, []byte(fmt.Sprint(startTime.Unix())), "start")
	if err != nil {
		return nil, err
	}

	payload, err = jsonparser.Set(payload, []byte(fmt.Sprint(endTime.Unix())), "end")
	if err != nil {
		return nil, err
	}

	return postDaymapPayload("http://daymap:9000/lessons/", payload, auth, jsonparser.Array)
}

/*
	Given a DaymapAuthMethod & a lesson ID, get the lesson plans for that lesson
	ID using the provided credentials.
*/
func GetLessonPlans(auth *DaymapAuthMethod, lessonID int) ([]byte, error) {
	payload := []byte{'{', '}'}

	// add the lesson id
	payload, err := jsonparser.Set(payload, []byte(fmt.Sprint(lessonID)), "lesson_id")
	if err != nil {
		return nil, err
	}

	return postDaymapPayload("http://daymap:9000/lessons/plans/", payload, auth, jsonparser.Object)
}

/*
	Given a DaymapAuthMethod, get all daymap messages for the user ever.
*/
func GetMessages(auth *DaymapAuthMethod) ([]byte, error) {
	// nothing but the credentials required here
	return postDaymapPayload("http://daymap:9000/messages/", []byte{'{', '}'}, auth, jsonparser.Array)
}
