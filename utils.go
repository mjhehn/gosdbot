package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"time"
)

func buildResponseList() {
	//TODO: setup list of autoresponses and how to load them up (json? csv? seperate lists for embeds vs text?)
	ars = append(ars, NewAutoResponse("[p|P]ing", []*TextResponse{NewTextResponse(1, "pong")}, nil, nil, nil, false, nil, nil))
	ars = append(ars, NewAutoResponse("[p|P]ling", []*TextResponse{NewTextResponse(1, "pong")}, nil, nil, nil, false, nil, nil))
	ars = append(ars, NewAutoResponse("[p|P]sing", []*TextResponse{NewTextResponse(1, "pong")}, nil, nil, nil, false, nil, nil))

	ars = append(ars, NewAutoResponse("pong", nil, []*EmbedResponse{NewEmbedResponse(1, "https://starscollideani.files.wordpress.com/2014/08/cirno-bsod.png?w=332&h=199")}, nil, nil, false, nil, nil))
	ars = append(ars, NewAutoResponse("handholding", nil, nil, []*ReactionResponse{NewReactionResponse(1, []string{"\U0001F1F1", "\U0001F1EA", "\U0001F1FC", "\U0001F1E9"})}, nil, false, nil, nil))
}

func readResponseList() {
	jsonResponse, err1 := ioutil.ReadFile("responses.json")
	check(err1)

	err := json.Unmarshal(jsonResponse, &ars)
	check(err)

	for i := range ars {
		ars[i].updateRegex()
	}
}

func writeResponseList() {
	jsonResponse, _ := json.Marshal(ars)
	f, _ := os.Create("responses.json")
	defer f.Close()
	f.Write(jsonResponse)
	f.Sync()
}

func check(err error) {
	if err != nil {
		fmt.Println(err.Error)
	}
}

func rollDie(numFaces int64) int64 {
	rnged, err := rand.Int(rand.Reader, big.NewInt(numFaces))
	check(err)
	return rnged.Int64() + 1
}

var myClient = &http.Client{Timeout: 10 * time.Second}

func getJSON(url string, target interface{}) error {
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}
