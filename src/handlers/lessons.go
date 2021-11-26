package handlers

import (
	"io/ioutil"
	"log"
	"main/util/daymap"
	"main/util/userman"
	"net/http"
	"time"

	"github.com/buger/jsonparser"
)

func (h *BaseHandler) GetLessons(w http.ResponseWriter, r *http.Request) {
	if !EnsureMethod("POST", w, r) {
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	auth, err := userman.GetDaymapAuthMethod(body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	startTimestamp, err := jsonparser.GetInt(body, "start")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	endTimestamp, err := jsonparser.GetInt(body, "end")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data, err := daymap.GetLessons(
		auth,
		time.Unix(startTimestamp, 0),
		time.Unix(endTimestamp, 0),
	)
	if err != nil {

		switch err.Error() {
		case "bad auth":
			w.WriteHeader(http.StatusUnauthorized)
			return
		case "bad daymap":
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return

	}
	w.Write(data)

	// in the background, just update the cookies returned by daymap getter if we can
	// we want this to be in the bg since the user would have to wait for this to
	// run in order to get the data they need
	go func() {
		err = userman.ParseAndUpdateCookies(auth.Username, auth.Password, data)
		if err != nil {
			log.Printf("failed to update cookies: %s", err.Error())
			return
		}
	}()
}
