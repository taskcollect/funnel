package userman

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"main/util/daymap"
	"net/http"
	"strconv"

	"github.com/buger/jsonparser"
)

/*
	Fetch stored cookies from userman given a username and password.
	Returns an empty json ( {} ) if no cookies.
*/
func FetchDaymapCookies(user string, password string) ([]byte, error) {
	// request body
	// note already populated with credentials
	rb := []byte(`{"creds":true}`)

	rb, err := jsonparser.Set(rb, []byte(strconv.Quote(user)), "user")
	if err != nil {
		return nil, err
	}

	rb, err = jsonparser.Set(rb, []byte(strconv.Quote(password)), "secret")
	if err != nil {
		return nil, err
	}

	// request body reader
	rbr := bytes.NewReader(rb)

	// construct the request object
	req, err := http.NewRequest(http.MethodGet, "http://userman:2000/v1/get", rbr)
	if err != nil {
		return nil, err
	}

	// execute the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// get the cookies
	cookies, cType, _, err := jsonparser.Get(raw, "creds.daymap")
	if err != nil && cType != jsonparser.NotExist {
		return nil, err
	}

	if cType == jsonparser.NotExist {
		return nil, nil
	}

	if cType != jsonparser.Object {
		return nil, errors.New("got non object cookies from userman")
	}

	return cookies, nil
}

/*
	Given a JSON request, extract the username, password
	from request, contact userman for cookies and return a
	DaymapAuthMethod struct

	Relies on FetchDaymapCookies to contact userman for cookies.
*/
func GetDaymapAuthMethod(body []byte) (*daymap.DaymapAuthMethod, error) {
	user, err := jsonparser.GetString(body, "username")
	if err != nil {
		return nil, err
	}

	pass, pType, _, err := jsonparser.Get(body, "password")
	if err != nil {
		return nil, err
	}
	if pType != jsonparser.NotExist && pType != jsonparser.String {
		return nil, errors.New("password was not a string")
	}

	// get the cookies from userman
	cookies, err := FetchDaymapCookies(user, string(pass))
	if err != nil {
		// this is not fatal, let's keep the show going if something goes wrong there
		log.Printf("warning, error in cookie fetch for %s: %s", user, err.Error())
		cookies = nil
	}

	return &daymap.DaymapAuthMethod{
		Username: user,
		Password: string(pass),
		Cookies:  cookies,
	}, nil
}

/*
	Given a username, password, and cookie JSON, send a request to userman
	to update the user's stored daymap cookies in the database.

	Effectively updates the "daymap" property in the creds object of the user.
*/
func UpdateDaymapCookies(user string, pass string, cookies []byte) error {
	payload, err := ConstructUsermanAuth(user, pass)
	if err != nil {
		return err
	}

	payload, err = jsonparser.Set(payload, cookies, "creds.daymap")
	if err != nil {
		return err
	}

	pr := bytes.NewReader(payload)

	resp, err := http.Post("http://userman:2000/v1/alter", "application/json", pr)

	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("userman returned %s", resp.Status)
	}

	return nil
}

/*
	Given a username, password & daymap-getter response json,
	use UpdateDaymapCookies to send the data to userman and update
	the user's cookies.
*/
func ParseAndUpdateCookies(user string, pass string, fullresp []byte) error {
	cookies, cType, _, err := jsonparser.Get(fullresp, "cookies")
	if err != nil && cType != jsonparser.NotExist {
		return err
	}

	return UpdateDaymapCookies(user, pass, cookies)
}
