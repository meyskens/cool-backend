package main

import (
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

var sigfoxToken = os.Getenv("SIGFOX_API_TOKEN")

// SigfoxCallback contains a  SigFox data callback
type SigfoxCallback struct {
	Device string `json:"device"`
	Data   string `json:"data"`
	Time   int    `json:"time"`
}

// SigfoxUplinkData contains tha uplink callback info
type SigfoxUplinkData struct {
	DownlinkData string `json:"downlinkData"`
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Authorization") != "token "+sigfoxToken {
		unauthorizedError(w, r)
		return
	}
	ctx := appengine.NewContext(r)
	if r.Method != "POST" {
		http.NotFound(w, r)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		interalServerError(w, r, err)
		return
	}

	info := SigfoxCallback{}
	err = json.Unmarshal(body, &info)
	if err != nil {
		interalServerError(w, r, err)
		return
	}

	log.Debugf(ctx, "Got info %v", info)

	go writeMessageToDatabase(r, info)

	callback := make([]byte, 8)
	callback[0] = 15 // default timeout
	uplink := map[string]SigfoxUplinkData{info.Device: SigfoxUplinkData{DownlinkData: hex.EncodeToString(callback)}}
	response, _ := json.Marshal(uplink)

	log.Debugf(ctx, "Send callback %s", response)
	w.Header().Add("Content-Type", "application/json")
	w.Write(response)
}

func writeMessageToDatabase(r *http.Request, message SigfoxCallback) {
	ctx := appengine.NewContext(r)
	key := datastore.NewIncompleteKey(ctx, "Message", nil)
	datastore.Put(ctx, key, message)
}
