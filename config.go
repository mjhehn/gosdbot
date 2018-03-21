package sdbot

import (
	"encoding/json"
	"io/ioutil"
)

//Config is a wrapper for the autoresposnses and mutedServers objects needed by the bot
type Config struct {
	Ars          []*AutoResponse
	MutedServers []string
	Token        string
	Status       string
	OwnerID      string
}

//ConfigJSON is a JSON serializeable struct to use to set Config via json parsing
type ConfigJSON struct {
	Token   string
	Status  string
	OwnerID string
}

//NewConfig is a server cstor
func NewConfig() *Config {
	sc := new(Config)
	sc.Ars = []*AutoResponse{}
	sc.MutedServers = []string{}
	sc.Token = "bottokengoeshere"
	sc.Status = "default!"
	sc.OwnerID = "0"
	return sc
}

//NewConfigByToken is a server cstor
func NewConfigByToken(token string) *Config {
	sc := new(Config)
	sc.Ars = []*AutoResponse{}
	sc.MutedServers = []string{}
	sc.Token = token
	sc.Status = "default!"
	sc.OwnerID = "0" //fill owner ID here
	return sc
}

//ConfigFromJSON builds a list of autoresponses based on a json file
func ConfigFromJSON() *Config {
	var data *ConfigJSON
	c := NewConfig()
	jsonConfig, err1 := ioutil.ReadFile("config.json")
	Check(err1)

	err := json.Unmarshal(jsonConfig, &data)
	Check(err)

	c.MutedServers = []string{}
	c.OwnerID = data.OwnerID
	c.Token = data.Token
	c.Status = data.Status
	return c
}
