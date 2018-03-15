package main

import (
	"crypto/rand"
	"math/big"
	"regexp"

	"github.com/bwmarrin/discordgo"
)

//AutoResponse stores all elements needed for an Autoresponse
type AutoResponse struct {
	trigger        string
	regex          *regexp.Regexp
	responses      []*TextResponse
	embeds         []*EmbedResponse
	reactions      []*ReactionResponse
	mentions       []string
	cleanup        bool
	userSpecific   []string
	serverSpecific []string
}

//NewAutoResponse build an AutoResponse object, takes arguments for everything and then creates a compiled regex from the trigger
func NewAutoResponse(trigger string, responses []*TextResponse, embeds []*EmbedResponse, reactions []*ReactionResponse, mentions []string, cleanup bool, userSpecific []string, serverSpecific []string) *AutoResponse {
	a := new(AutoResponse)
	a.trigger = trigger
	a.regex, _ = regexp.Compile(trigger)
	a.responses = responses
	a.embeds = embeds
	a.reactions = reactions
	a.mentions = mentions
	a.cleanup = cleanup
	a.userSpecific = userSpecific
	a.serverSpecific = serverSpecific
	return a
}

//check for a textResponse to make
func (a *AutoResponse) checkTextResponses(session *discordgo.Session, message *discordgo.MessageCreate, responded chan bool) {
	for _, response := range a.responses {
		if response != nil && a.regex.MatchString(message.Content) && runIt(response.chance) {
			response.respond(session, message)
			responded <- true
			return
		}
	}
}

func (a *AutoResponse) checkEmbedResponses(session *discordgo.Session, message *discordgo.MessageCreate, responded chan bool) {
	for _, response := range a.embeds {
		if response != nil && a.regex.MatchString(message.Content) && runIt(response.chance) {
			responded <- response.respond(session, message)
			return
		}
	}
}

func (a *AutoResponse) checkReactionResponses(session *discordgo.Session, message *discordgo.MessageCreate, responded chan bool) {
	//check for an reaction response to make
	for _, response := range a.reactions {
		if response != nil && a.regex.MatchString(message.Content) && runIt(response.chance) {
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
