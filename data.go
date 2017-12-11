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
	out := []FridgeData{}

	if check := len(in) % 8; check != 0 {
		return nil, errors.New("not right amout of bytes")
	}

	dataPerNode := []string{}

	for i := 0; i < len(in); i += 8 {
		data, err := parseHex(in[i : i+8])
		if err != nil {
			return nil, err
		}
		dataPerNode = append(dataPerNode, data)
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
	in = filterZeros(in)
	fmt.Println(in)
	fmt.Println(len(in))

	per := 2
	if len(in) > 8 { // if multiple SigFox.write happen it gets even weirder
		per = 8
	}

	parts := []string{}
	for i := 0; i < len(in); i += per {
		parts = append(parts, in[i:i+per])
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

func filterZeros(in string) string {
	hadNonZero := false
	lastChar := -1
	for i := 0; i < len(in); i++ {
		if in[i] == '0' && !hadNonZero {
			lastChar = i
		} else {
			hadNonZero = true
		}
	}
	if lastChar >= 0 {
		in = in[lastChar+1:]
	}

	if (len(in) % 2) != 0 {
		in = "0" + in
	}

	return in
}
