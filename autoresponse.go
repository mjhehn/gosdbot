//contains most of the functions and methods associated directly with the Autoresponse struct and its most relevant interfaces.
package main

import (
	"regexp"

	"github.com/bwmarrin/discordgo"
)

//Used to help generalize the code between the embed, text, and reaction responses.
type responder interface {
	respond(session *discordgo.Session, message *discordgo.MessageCreate, mentions []string) bool
}

//AutoResponse stores all elements needed for an Autoresponse
type AutoResponse struct {
	Trigger        string              //regex needed to invoke the response as a string
	Regex          *regexp.Regexp      //the regex compiled. handled by NewAutoResponse
	Responses      []*TextResponse     //list of textresposne objects possible to reply with
	Embeds         []*EmbedResponse    //list of embeds possible to reply with
	Reactions      []*ReactionResponse //list of reactions possible to reply with
	Mentions       []string            //list of mentions. currently only works with self-mentions
	Cleanup        bool                //whether the bot should delete the message that invoked the responses
	UserSpecific   []string            //list of users teh autoresponse should run for
	ServerSpecific []string            //list of servers the autoresponse can run on
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

//check for a textResponse to push to the guild/channel
func (a *AutoResponse) checkTextResponses(session *discordgo.Session, message *discordgo.MessageCreate, responded chan bool) {
	var selectedResponse responder
	//if this is user specific, and the invoker isn't in that list
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
			if a.Reactions != nil {
				for _, response := range a.Reactions {
					if response.chance() {
						response.respond(session, message, nil)
					}
				}
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
