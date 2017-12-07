package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// FridgeData is the data ata send for one fridge on one timepoints
type FridgeData struct {
	FridgeID     string    `bigquery:"fridgeid"`
	Time         time.Time `bigquery:"time"`
	Temperature  float64   `bigquery:"temperature"`
	Humidity     float64   `bigquery:"humidity"`
	DoorOpenings int64     `bigquery:"dooropenings"`
}

func parseInput(in string, timeSent time.Time) ([]FridgeData, error) {
	in = makeDataSane(in)
	data, err := parseHex(in)
	if err != nil {
		return nil, err
	}
	out := []FridgeData{}

	if check := len(data) % 8; check != 0 {
		return nil, errors.New("not right amout of bytes")
	}

	dataPerNode := []string{}

	for i := 0; i < len(data); i += 8 {
		dataPerNode = append(dataPerNode, data[i:i+8])
	}

	for _, d := range dataPerNode {
		// 2 021 01 03
		// 0 123 45 56
		temp, _ := strconv.ParseFloat(string(d[1])+string(d[2])+"."+string(d[3]), 64)
		humidity, _ := strconv.ParseFloat(string(d[4])+string(d[5]), 64)
		doorsOpen, _ := strconv.ParseInt(string(d[6])+string(d[7]), 10, 64)
		out = append(out, FridgeData{
			FridgeID:     string(d[0]),
			Time:         timeSent,
			Temperature:  temp,
			Humidity:     humidity,
			DoorOpenings: doorsOpen,
		})
	}

	return out, nil
}

func parseHex(in string) (string, error) {
	i, err := strconv.ParseInt(in, 16, 64)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return strconv.FormatInt(i, 10), nil
}

func makeDataSane(in string) string {
	parts := []string{}
	for i := 0; i < len(in); i += 2 {
		parts = append(parts, in[i:i+2])
	}
	reverse(parts)

	if parts[0] == "00" {
		parts = parts[1:]
	}

	return strings.Join(parts, "")
}

func reverse(ss []string) {
	last := len(ss) - 1
	for i := 0; i < len(ss)/2; i++ {
		ss[i], ss[last-i] = ss[last-i], ss[i]
	}
}