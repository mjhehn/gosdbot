//contains most of the functions and methods associated directly with the Autoresponse struct and its most relevant interfaces.
package main

import (
	"regexp"

	"github.com/bwmarrin/discordgo"
)

const (
	textresponse     = 0
	embedresponse    = 1
	reactionresponse = 2
)

//Used to help generalize the code between the embed, text, and reaction responses.
type responder interface {
	respond(session *discordgo.Session, message *discordgo.MessageCreate, mentions []string) bool
	chance() bool
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

func (a *AutoResponse) checkResponses(session *discordgo.Session, message *discordgo.MessageCreate, checkType int, responded chan bool) {
	var selectedResponse responder
	var responses []responder
	switch checkType {
	case 0:
		for _, r := range a.Responses {
			responses = append(responses, r)
		}
	case 1:
		for _, r := range a.Embeds {
			responses = append(responses, r)
		}
	case 2:
		for _, r := range a.Reactions {
			responses = append(responses, r)
		}
	default:
		for _, r := range a.Responses {
			responses = append(responses, r)
		}
		break
	}

	//if this is user specific, and the invoker isn't in that list
	if a.UserSpecific != nil && !in(message.Author.Username, a.UserSpecific) {
		return
	}

	if responses != nil && a.Regex.MatchString(message.Content) {
		for i := range responses {
			if responses[i].chance() {
				selectedResponse = responses[i]
			}
		}
		if selectedResponse != nil || responses == nil {
			if a.Cleanup {
				session.ChannelMessageDelete(message.ChannelID, message.ID)
			}
			a.addReactions(session, message)
			responded <- selectedResponse.respond(session, message, a.Mentions)
		}
		return
	}
}

func (a *AutoResponse) addReactions(session *discordgo.Session, message *discordgo.MessageCreate) {
	if a.Reactions != nil {
		for _, response := range a.Reactions {
			if response.chance() {
				response.respond(session, message, nil)
			}
		}
	}
}

func (a *AutoResponse) updateRegex() {
	a.Regex, _ = regexp.Compile(a.Trigger)
}
