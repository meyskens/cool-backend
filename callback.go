package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/bigquery"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

var sigfoxToken = os.Getenv("SIGFOX_API_TOKEN")

// SigfoxCallback contains a  SigFox data callback
type SigfoxCallback struct {
	Device  string  `json:"device"`
	Data    string  `json:"data"`
	Time    int     `json:"time"`
	SNR     float64 `json:"snr"`
	ACK     bool    `json:"ack"`
	Station string  `json:"station"`
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

	writeMessageToDatabase(ctx, info)

	callback := make([]byte, 8)
	callback[0] = byte(getSleep())
	uplink := map[string]SigfoxUplinkData{info.Device: SigfoxUplinkData{DownlinkData: hex.EncodeToString(callback)}}
	response, _ := json.Marshal(uplink)

	log.Debugf(ctx, "Send callback %s", response)
	w.Header().Add("Content-Type", "application/json")
	w.Write(response)
}

func writeMessageToDatabase(ctx context.Context, message SigfoxCallback) {
	projectID := appengine.AppID(ctx)

	// Create the BigQuery service.
	bq, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		log.Debugf(ctx, "could not create service: %v", err)
		return
	}

	uploader := bq.Dataset("cooling").Table("messages").Uploader()
	items := []*SigfoxCallback{
		&message,
	}
	if err := uploader.Put(ctx, items); err != nil {
		log.Debugf(ctx, "Uploader error %v", err)
	}
}

func getSleep() int8 {
	loc, _ := time.LoadLocation("Europe/Brussels")
	now := time.Now().In(loc)

	if now.Hour() >= 11 && now.Hour() <= 13 { // during lunch hours
		return 5
	}

	return 20
}
