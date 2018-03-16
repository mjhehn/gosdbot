package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/bwmarrin/discordgo"
)

//prints the results of an error if it exists
func check(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

//rollDie 'rolls a die', but more important is a bounded number generator from 1-numFaces
func rollDie(numFaces int64) int64 {
	rnged, err := rand.Int(rand.Reader, big.NewInt(numFaces))
	check(err)
	return rnged.Int64() + 1
}

var myClient = &http.Client{Timeout: 10 * time.Second} //to help handle getting json from the web
func getJSON(url string, target interface{}) error {
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

//takes in a chance (1:chance) and rolls to see if the number generator for that range hits 0
func runIt(chance int) bool {
	rnged, _ := rand.Int(rand.Reader, big.NewInt(int64(chance)))
	if rnged.Int64() == int64(0) {
		return true
	}
	return false
}

//retreive the guild/server by name from the current message and session
func getServer(session *discordgo.Session, msg *discordgo.MessageCreate) string {
	guild := getGuild(session, msg)
	if guild != nil {
		return guild.Name
	}
	return ""
}

//retreive the actual server/guild object
func getGuild(session *discordgo.Session, msg *discordgo.MessageCreate) *discordgo.Guild {
	channel, err := session.State.Channel(msg.ChannelID)
	if err != nil {
		channel, err = session.Channel(msg.ChannelID)
		if err != nil {
			check(err)
			return nil
		}
	}

	// Attempt to get the guild from the state,
	// If there is an error, fall back to the restapi.
	guild, err := session.State.Guild(channel.GuildID)
	if err != nil {
		guild, err = session.Guild(channel.GuildID)
		if err != nil {
			check(nil)
			return nil
		}
	}
	return guild
}

//retreive the roles of a user based on the message the sent+the current session
func getRoles(session *discordgo.Session, msg *discordgo.MessageCreate) []string {
	currentGuild := getGuild(session, msg)

	user, err := session.GuildMember(currentGuild.ID, msg.Author.ID)
	check(err)
	userRoles := user.Roles
	var roleList []string
	for _, role := range userRoles {
		roleObject, err := session.State.Role(currentGuild.ID, role)
		check(err)
		roleList = append(roleList, roleObject.Name)
	}
	return roleList
}

//checks if s is in the list
func in(s string, list []string) bool {
	for _, item := range list {
		if s == item {
			return true
		}
	}
	return false
}

//retreives server emoji by name if it exists.
func getEmoji(session *discordgo.Session, msg *discordgo.MessageCreate, name string) *discordgo.Emoji {
	guild := getGuild(session, msg)
	guildEmojis := guild.Emojis
	if guildEmojis != nil {
		for _, guildEmoji := range guildEmojis {
			if guildEmoji.Name == name {
				return guildEmoji
			}
		}
	}

	return nil
}
