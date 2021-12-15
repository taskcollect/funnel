package handlers

import (
	"log"
	"main/util/daymap"
	"main/util/email"
	"net/http"
	"sync"

	"github.com/buger/jsonparser"
)

func (h *BaseHandler) GetAllMessages(w http.ResponseWriter, r *http.Request) {
	if !EnsureMethod("POST", w, r) {
		return
	}

	// this auth object is meant for daymap only but hey, it has a username & password!
	_, daymapAuth, err := parseBodyAndGetDaymapAuth(r.Body, w)
	if err != nil {
		// error already handled by function
		return
	}

	out := []byte{'{', '}'}

	type getterRetVal struct {
		key string
		val []byte
		err error
	}

	// channel where each individual getter returns its data content when it wants to
	data_ch := make(chan getterRetVal)
	// syncgroup that ensures all getters are done before everything finishes
	var wg sync.WaitGroup

	// daymap messages worker
	wg.Add(1)

	go func() {
		defer wg.Done() // tell the waitgroup we're done after this finishes
		daymapFullresp, err := daymap.GetMessages(daymapAuth)

		if err != nil {
			data_ch <- getterRetVal{"", nil, err}
			return
		}

		// just the daymapMessages
		daymapMessages, _, _, err := jsonparser.Get(daymapFullresp, "data")
		if err != nil {
			data_ch <- getterRetVal{"", nil, err}
			return
		}

		// update cookies in bg if everything went well
		go backgroundCookieUpdate(daymapAuth, daymapFullresp)

		data_ch <- getterRetVal{"daymap", daymapMessages, nil}
	}()

	// email messages worker
	wg.Add(1)
	go func() {
		defer wg.Done()

		// get last 100 emails
		emails, err := email.GetAllEmails(daymapAuth.Username, daymapAuth.Password, 100)
		if err != nil {
			data_ch <- getterRetVal{"", nil, err}
			return
		}

		data_ch <- getterRetVal{"emails", emails, nil}
	}()

	// channel to signal when everything is done, may exit with an error
	done_ch := make(chan error)

	/*
		Read from the datachannel and safely insert it into the final response.
		When the channel has been closed, close the done channel, marking
		the end of execution.
	*/
	go func() {
		defer close(done_ch) // close the done channel on exit

		// wait for something on the data channel
		for rv := range data_ch {
			// we got something from one of the worker functions

			if rv.err != nil {
				// a getter came to an error

				// FIXME: send the client back something sensible here
				w.WriteHeader(http.StatusInternalServerError)
				log.Println(rv.err.Error())

				done_ch <- rv.err
				return
			}

			// try modifying the output json with the returned getter data
			out, err = jsonparser.Set(out, rv.val, rv.key)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				done_ch <- err
				return
			}
		}
	}()

	/*
		This goroutine's sole purpose in life is to close the datachannel
		whenever the waitgroup is done, that's all folks.

		When this happens, the above goroutine will stop waiting on the data
		channel and close the done channel, meaning the end of execution.
	*/
	go func() {
		defer close(data_ch) // close the datachannel on exit
		wg.Wait()            // wait out all the goroutines
	}()

	err = <-done_ch // wait until the done signal
	if err != nil {
		// there's an error in the channel, so don't write a response
		// this shouldn't be handled here
		return
	}
	w.Write(out)
}
