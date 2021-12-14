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

/*
   In the background, just update the cookies returned by daymap getter if we can
   we want this to be in the bg since the user would have to wait for this to
   run in order to get the data they need.

   This function should be a goroutine.
*/
func backgroundCookieUpdate(auth *daymap.DaymapAuthMethod, data []byte) {
	err := userman.ParseAndUpdateCookies(auth.Username, auth.Password, data)
	if err != nil {
		log.Printf("failed to update cookies for %s: %s", auth.Username, err.Error())
		return
	}
}

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

	// update cookies in bg if everything went well
	go backgroundCookieUpdate(auth, data)
}

func (h *BaseHandler) GetLessonPlans(w http.ResponseWriter, r *http.Request) {
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

	lessonID, err := jsonparser.GetInt(body, "lesson_id")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data, err := daymap.GetLessonPlans(
		auth, int(lessonID),
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

	// update cookies in bg if everything went well
	go backgroundCookieUpdate(auth, data)
}
