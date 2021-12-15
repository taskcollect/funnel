package email

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/buger/jsonparser"
)

func GetAllEmails(username, password string, amount uint) ([]byte, error) {
	req, err := http.NewRequest("GET", "http://email:5000/v1/mail", nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("amount", fmt.Sprint(amount))
	req.URL.RawQuery = q.Encode()

	// the email getter doesn't support capitalized headers, so manual setting of the struct is required
	// see: https://github.com/taskcollect/email-getter/issues/1

	req.Header = http.Header{
		"username": []string{username},
		"password": []string{password},
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		b, _ := ioutil.ReadAll(resp.Body)
		log.Println(string(b))
		return nil, errors.New("upstream error: " + resp.Status)
	}

	fullresp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	emailList, eListType, _, err := jsonparser.Get(fullresp, "messages")
	if err != nil {
		return nil, err
	}

	if eListType != jsonparser.Array {
		return nil, errors.New("email-getter returned unexpected response type in messages key")
	}

	return emailList, nil
}
