package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

//ServerConfig is a wrapper for the autoresposnses and mutedServers object
type ServerConfig struct {
	ars          []*AutoResponse
	mutedServers []string
}

//NewServerConfig is a server cstor
func NewServerConfig() *ServerConfig {
	sc := new(ServerConfig)
	sc.ars = []*AutoResponse{}
	sc.mutedServers = []string{}
	sc.readResponseList()
	//buildResponseList()
	//writeResponseList()

	return sc
}

func buildResponseList() {
	config.ars = append(config.ars, NewAutoResponse("[p|P]ing", []*TextResponse{NewTextResponse(1, "pong")}, nil, nil, nil, false, nil, nil))
	config.ars = append(config.ars, NewAutoResponse("[p|P]ling", []*TextResponse{NewTextResponse(1, "pong")}, nil, nil, nil, false, nil, nil))
	config.ars = append(config.ars, NewAutoResponse("[p|P]sing", []*TextResponse{NewTextResponse(1, "pong")}, nil, nil, nil, false, nil, nil))

	config.ars = append(config.ars, NewAutoResponse("pong", nil, []*EmbedResponse{NewEmbedResponse(1, "https://stconfig.arscollideani.files.wordpress.com/2014/08/cirno-bsod.png?w=332&h=199")}, nil, nil, false, nil, nil))
	config.ars = append(config.ars, NewAutoResponse("handholding", nil, nil, []*ReactionResponse{NewReactionResponse(1, []string{"\U0001F1F1", "\U0001F1EA", "\U0001F1FC", "\U0001F1E9"})}, nil, false, nil, nil))
}

//build list of autoresponses based on a json file
func (c *ServerConfig) readResponseList() {
	jsonResponse, err1 := ioutil.ReadFile("responses.json")
	check(err1)

	err := json.Unmarshal(jsonResponse, &c.ars)
	check(err)

	for i := range c.ars {
		c.ars[i].updateRegex()
	}
}

func writeResponseList() {
	jsonResponse, _ := json.Marshal(config.ars)
	f, _ := os.Create("responses.json")
	defer f.Close()
	f.Write(jsonResponse)
	f.Sync()
}
