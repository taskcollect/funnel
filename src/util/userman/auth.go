package userman

import (
	"strconv"

	"github.com/buger/jsonparser"
)

/*
	Make a JSON with username & password that userman will accept.
*/
func ConstructUsermanAuth(user string, pass string) ([]byte, error) {
	out := []byte("{}")

	out, err := jsonparser.Set(out, []byte(strconv.Quote(user)), "user")
	if err != nil {
		return nil, err
	}

	out, err = jsonparser.Set(out, []byte(strconv.Quote(pass)), "secret")
	if err != nil {
		return nil, err
	}

	return out, nil
}
