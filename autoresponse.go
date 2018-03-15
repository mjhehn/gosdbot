package main

import (
	"crypto/rand"
	"math/big"
	"regexp"

	"github.com/bwmarrin/discordgo"
)

//AutoResponse stores all elements needed for an Autoresponse
type AutoResponse struct {
	Trigger        string
	Regex          *regexp.Regexp
	Responses      []*TextResponse
	Embeds         []*EmbedResponse
	Reactions      []*ReactionResponse
	Mentions       []string
	Cleanup        bool
	UserSpecific   []string
	ServerSpecific []string
}

//NewAutoResponse build an AutoResponse object, takes arguments for everything and then creates a compiled regex from the Trigger
func NewAutoResponse(Trigger string, Responses []*TextResponse, Embeds []*EmbedResponse, Reactions []*ReactionResponse, Mentions []string, Cleanup bool, UserSpecific []string, ServerSpecific []string) *AutoResponse {
	a := new(AutoResponse)
	a.Trigger = Trigger
	a.Regex, _ = regexp.Compile(Trigger)
	a.Responses = Responses
	a.Embeds = Embeds
	a.Reactions = Reactions
	a.Mentions = Mentions
	a.Cleanup = Cleanup
	a.UserSpecific = UserSpecific
	a.ServerSpecific = ServerSpecific
	return a
}

//check for a textResponse to make
func (a *AutoResponse) checkTextResponses(session *discordgo.Session, message *discordgo.MessageCreate, responded chan bool) {
	for _, response := range a.Responses {
		if response != nil && a.Regex.MatchString(message.Content) && runIt(response.Chance) {
			response.respond(session, message)
			responded <- true
			return
		}
	}
}

func (a *AutoResponse) checkEmbedResponses(session *discordgo.Session, message *discordgo.MessageCreate, responded chan bool) {
	for _, response := range a.Embeds {
		if response != nil && a.Regex.MatchString(message.Content) && runIt(response.Chance) {
			responded <- response.respond(session, message)
			return
		}
	}
}

func (a *AutoResponse) checkReactionResponses(session *discordgo.Session, message *discordgo.MessageCreate, responded chan bool) {
	//check for an reaction response to make
	for _, response := range a.Reactions {
		if response != nil && a.Regex.MatchString(message.Content) && runIt(response.Chance) {
			response.respond(session, message)
			responded <- true
			return
		}
	}
}

func runIt(chance int) bool {
	rnged, _ := rand.Int(rand.Reader, big.NewInt(int64(chance)))
	if rnged.Int64() == int64(0) {
		return true
	}
	return false
}

func (a *AutoResponse) updateRegex() {
	a.Regex, _ = regexp.Compile(a.Trigger)
}
