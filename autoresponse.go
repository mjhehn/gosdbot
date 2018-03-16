package main

import (
	"regexp"

	"github.com/bwmarrin/discordgo"
)

//TODO: implement cleanup feature
type responder interface {
	respond(session *discordgo.Session, message *discordgo.MessageCreate, mentions []string) bool
}

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
	var selectedResponse responder
	if a.UserSpecific != nil && !in(message.Author.Username, a.UserSpecific) {
		return
	}

	if a.Responses != nil && a.Regex.MatchString(message.Content) {
		for _, response := range a.Responses {
			if response.chance() {
				selectedResponse = response
			}
		}
		if selectedResponse != nil {
			if a.Cleanup {
				session.ChannelMessageDelete(message.ChannelID, message.ID)
			}
			responded <- selectedResponse.respond(session, message, a.Mentions)
		}
		return
	}
}

func (a *AutoResponse) checkEmbedResponses(session *discordgo.Session, message *discordgo.MessageCreate, responded chan bool) {
	if a.UserSpecific != nil && !in(message.Author.Username, a.UserSpecific) {
		return
	}

	var selectedResponse responder
	if a.Embeds != nil && a.Regex.MatchString(message.Content) {
		for _, response := range a.Embeds {
			if response.chance() {
				selectedResponse = response
			}
		}
		if selectedResponse != nil {
			if a.Cleanup {
				session.ChannelMessageDelete(message.ChannelID, message.ID)
			}
			responded <- selectedResponse.respond(session, message, a.Mentions)
		}
		return
	}
}

func (a *AutoResponse) checkReactionResponses(session *discordgo.Session, message *discordgo.MessageCreate, responded chan bool) {
	if a.UserSpecific != nil && !in(message.Author.Username, a.UserSpecific) {
		return
	}

	//check for an reaction response to make
	var selectedResponse responder
	if a.Reactions != nil && a.Regex.MatchString(message.Content) {
		for _, response := range a.Reactions {
			if response.chance() {
				selectedResponse = response
			}
		}
		if selectedResponse != nil {
			responded <- selectedResponse.respond(session, message, a.Mentions)
		}
		return
	}
}

func (a *AutoResponse) updateRegex() {
	a.Regex, _ = regexp.Compile(a.Trigger)
}
