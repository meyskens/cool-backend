package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

// SigfoxCallback contains a  SigFox data callback
type SigfoxCallback struct {
	Device string    `json:"device"`
	Data   string    `json:"data"`
	Time   time.Time `json:"time"`
}

// SigfoxUplinkData contains tha uplink callback info
type SigfoxUplinkData struct {
	DownlinkData string `json:"downlinkData"`
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
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

	log.Debugf(ctx, "Got data: %v", info)

	uplink := map[string]SigfoxUplinkData{"test": SigfoxUplinkData{DownlinkData: "ok"}}
	response, _ := json.Marshal(uplink)
	w.Write(response)
}
