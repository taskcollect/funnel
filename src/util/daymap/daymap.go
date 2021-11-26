package daymap

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
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
func PopulateBodyWithAuth(body []byte, auth *DaymapAuthMethod) ([]byte, error) {
	body, err := jsonparser.Set(body, []byte(strconv.Quote(auth.Username)), "username")
	if err != nil {
		return nil, err
	}

	body, err = jsonparser.Set(body, []byte(strconv.Quote(auth.Password)), "password")
	if err != nil {
		return nil, err
	}

	if auth.Cookies != nil {
		body, err = jsonparser.Set(body, auth.Cookies, "cookies")
		if err != nil {
			return nil, err
		}
	}

	return body, nil
}

/*
	Given a DaymapAuthMethod & start + end times, get all the lessons for the
	user with the provided credentials within the provided timeframe.

	Relies on PopulateBodyWithAuth, and talks to daymap-getter.
*/
func GetLessons(auth *DaymapAuthMethod, startTime, endTime time.Time) ([]byte, error) {
	body := []byte{'{', '}'}

	// add the auth info
	body, err := PopulateBodyWithAuth(body, auth)
	if err != nil {
		return nil, errors.New("bad auth structure")
	}

	// add the timestamps
	body, err = jsonparser.Set(body, []byte(fmt.Sprint(startTime.Unix())), "start")
	if err != nil {
		return nil, err
	}

	body, err = jsonparser.Set(body, []byte(fmt.Sprint(endTime.Unix())), "end")
	if err != nil {
		return nil, err
	}

	log.Println(string(body))

	// construct the request object
	res, err := http.Post("http://daymap:9000/lessons/", "application/json", bytes.NewReader(body))
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

	data, dType, _, err := jsonparser.Get(resp, "data")
	if err != nil {
		return nil, err
	}

	if dType != jsonparser.Array {
		println(dType)
		return nil, errors.New("data in response was not an array (?)")
	}

	return data, nil
}
