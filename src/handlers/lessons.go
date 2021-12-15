package handlers

import (
	"io"
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

func parseBodyAndGetDaymapAuth(reader io.ReadCloser, w http.ResponseWriter) ([]byte, *daymap.DaymapAuthMethod, error) {
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return nil, nil, err
	}

	auth, err := userman.GetDaymapAuthMethod(body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return nil, nil, err
	}

	return body, auth, nil
}

// helper method to write the appropriate error header for daymap-getter interfaces
func writeAppropriateError(err error, w http.ResponseWriter) {
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
}

func filterAndWriteData(fullresp []byte, w http.ResponseWriter) error {
	data, _, _, err := jsonparser.Get(fullresp, "data")
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	_, err = w.Write(data)
	return err
}

func (h *BaseHandler) GetLessons(w http.ResponseWriter, r *http.Request) {
	if !EnsureMethod("POST", w, r) {
		return
	}

	body, auth, err := parseBodyAndGetDaymapAuth(r.Body, w)
	if err != nil {
		// error already handled by function
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

	fullresp, err := daymap.GetLessons(
		auth,
		time.Unix(startTimestamp, 0),
		time.Unix(endTimestamp, 0),
	)

	if err != nil {
		writeAppropriateError(err, w)
		return
	}

	err = filterAndWriteData(fullresp, w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// update cookies in bg if everything went well
	go backgroundCookieUpdate(auth, fullresp)
}

func (h *BaseHandler) GetLessonPlans(w http.ResponseWriter, r *http.Request) {
	if !EnsureMethod("POST", w, r) {
		return
	}

	body, auth, err := parseBodyAndGetDaymapAuth(r.Body, w)
	if err != nil {
		// error already handled by function
		return
	}

	lessonID, err := jsonparser.GetInt(body, "lesson_id")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fullresp, err := daymap.GetLessonPlans(
		auth, int(lessonID),
	)

	if err != nil {
		writeAppropriateError(err, w)
		return
	}

	err = filterAndWriteData(fullresp, w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// update cookies in bg if everything went well
	go backgroundCookieUpdate(auth, fullresp)
}
