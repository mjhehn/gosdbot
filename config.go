package gosdbot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

//Config is a wrapper for the autoresposnses and mutedServers objects needed by the bot
type Config struct {
	Ars          []*AutoResponse
	MutedServers []string
	Token        string
	Status       string
	OwnerID      string
	CleverUser   string
	CleverID     string
}

//ConfigJSON is a JSON serializeable struct to use to set Config via json parsing
type ConfigJSON struct {
	Token      string
	Status     string
	OwnerID    string
	CleverUser string
	CleverID   string
}

//NewConfig is a server cstor
func NewConfig() *Config {
	sc := new(Config)
	sc.Ars = []*AutoResponse{}
	sc.MutedServers = []string{}
	sc.Token = "bottokengoeshere"
	sc.Status = "default!"
	sc.OwnerID = "0"
	sc.CleverID = "0"
	sc.CleverUser = "0"
	return sc
}

//NewConfigByToken is a server cstor
func NewConfigByToken(token string) *Config {
	sc := new(Config)
	sc.Ars = []*AutoResponse{}
	sc.MutedServers = []string{}
	sc.Token = token
	sc.Status = "default!"
	sc.OwnerID = "0"
	sc.CleverID = "0"
	sc.CleverUser = "0"
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
	c.CleverID = data.CleverID
	c.CleverUser = data.CleverUser
	return c
}

//Mute adds current sever to the muted list, which allows only mute commands to be received or sent.
func (c *Config) Mute(session *discordgo.Session, message *discordgo.MessageCreate, responded chan bool) {
	expr, err := regexp.Compile("^!mute$")
	Check(err)

	//mute
	if expr.MatchString(message.Content) {
		currentServer := GetServer(session, message)
		currentRoles := GetRoles(session, message)
		if !In(currentServer, c.MutedServers) && In("Bot Admin", currentRoles) { //mute
			c.MutedServers = append(c.MutedServers, currentServer)
			session.ChannelMessageSend(message.ChannelID, "Bot muted.")
			responded <- true
			return
		}
	}
	return
}

//Unmute removes the current server from the muted list
func (c *Config) Unmute(session *discordgo.Session, message *discordgo.MessageCreate, responded chan bool) {
	expr2, err2 := regexp.Compile("^!unmute$")
	Check(err2)

	if expr2.MatchString(message.Content) {
		currentServer := GetServer(session, message)
		currentRoles := GetRoles(session, message)

		if In(currentServer, c.MutedServers) && In("Bot Admin", currentRoles) { //mute
			for i, serv := range c.MutedServers {
				if serv == currentServer {
					c.MutedServers = append(c.MutedServers[:i], c.MutedServers[i+1:]...)
				}
			}
			session.ChannelMessageSend(message.ChannelID, "Bot unmuted.")
			responded <- true
			return
		}
	}
	return
}

//Mutestatus as it says, returns the current message channel's server mute state
func (c *Config) Mutestatus(session *discordgo.Session, message *discordgo.MessageCreate, responded chan bool) {
	expr3, err3 := regexp.Compile("^!mutestatus$")
	Check(err3)

	if expr3.MatchString(message.Content) { //mutestatus
		currentServer := GetServer(session, message)
		mutedStatus := " "
		if In(currentServer, c.MutedServers) {
			mutedStatus = "muted"
		} else {
			mutedStatus = "not muted"
		}
		result := fmt.Sprintf("Semi-Decent is currently %v on this server", mutedStatus)
		session.ChannelMessageSend(message.ChannelID, result)
		responded <- true
		return
	}
	return
}

//Botstatus updates the bot's "playing" status notifier
func (c *Config) Botstatus(session *discordgo.Session, message *discordgo.MessageCreate, responded chan bool) {
	expr2, err2 := regexp.Compile("^!status .*$")
	Check(err2)

	if expr2.MatchString(message.Content) {
		if message.Author.ID == c.OwnerID { //
			c.Status = message.Content[8:]
			session.UpdateStatus(0, c.Status)
			session.ChannelMessageDelete(message.ChannelID, message.ID)
			responded <- true
			return
		}
	}
	return
}

//LeaveServer removes the bot from the guildID given to it if the user making the request is an admin
func (c *Config) LeaveServer(session *discordgo.Session, message *discordgo.MessageCreate, responded chan bool) {
	expr, err := regexp.Compile("^!leave [0-9]*$")
	Check(err)

	//mute
	if expr.MatchString(message.Content) {
		currentServer := GetServer(session, message)
		currentRoles := GetRoles(session, message)
		if !In(currentServer, c.MutedServers) && In("Bot Admin", currentRoles) { //mute
			guildID := strings.Split(message.Content, " ")
			err := session.GuildLeave(guildID[1])
			Check(err)
			responded <- true
			return
		}
	}
	return
}
