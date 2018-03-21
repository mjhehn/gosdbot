//Package botresponse contains most of the structs, functions and methods associated directly with the Autoresponse struct and its most relevant interfaces.
package botresponse

import (
	"encoding/json"
	"io/ioutil"
	"regexp"
	"sdbot/pkg/botutils"

	"github.com/bwmarrin/discordgo"
)

//help improve readability of the 'general' response checker method.
const (
	Textresponse     = 0
	Embedresponse    = 1
	Reactionresponse = 2
)

//Used to help generalize the code between the embed, text, and reaction responses.
type respondable interface {
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

//CheckResponses acts as a generic checker and applier of responses possible from the autoresponse it is called upon.
//takes:
//pointer to the current session
//pointer to the message created by messagecreate event
//an integer telling the check whether it's checking through text responses, embedded resposnes, or (emoji) reaction responses (0, 1, 2 respectively)
//a bool channel to note when a goroutine has made an answer.
func (a *AutoResponse) CheckResponses(session *discordgo.Session, message *discordgo.MessageCreate, checkType int, responded chan bool) {
	var selectedResponse respondable
	var responses []respondable

	defer func() { //handle failure due to closed channel
		if r := recover(); r != nil {
			//fmt.Println("Recovered in f", r)
		}
	}()

	switch checkType { //switch depending on what type of responses we're checking. integers with constants.
	case Textresponse:
		for _, r := range a.Responses {
			responses = append(responses, r)
		}
	case Embedresponse:
		for _, r := range a.Embeds {
			responses = append(responses, r)
		}
	case Reactionresponse:
		for _, r := range a.Reactions {
			responses = append(responses, r)
		}
		return
	}

	//if this is user specific, and the message author isn't in that list
	if a.UserSpecific != nil && !botutils.In(message.Author.Username, a.UserSpecific) {
		return
	}

	if responses != nil && a.Regex.MatchString(message.Content) { //check for the trigger, and if there even are responses to make
		for i := range responses {
			if responses[i].chance() { //if rolled successfully to respond with this responses
				selectedResponse = responses[i]
			}
		}
		if selectedResponse != nil { //check that we got a response
			if a.Cleanup { //check if this response is supposed to delete teh command message
				session.ChannelMessageDelete(message.ChannelID, message.ID)
			}
			a.addReactions(session, message)                                    //add reactions, if any.
			responded <- selectedResponse.respond(session, message, a.Mentions) //send the message, and push true to the channel
		}
		return
	}
	return
}

//addReactions adds relevant reaction items to a message.
func (a *AutoResponse) addReactions(session *discordgo.Session, message *discordgo.MessageCreate) {
	if a.Reactions != nil {
		for _, response := range a.Reactions {
			if response.chance() {
				response.respond(session, message, nil) //apply reactions to message
			}
		}
	}
}

//updateRegex: update the regex field to a compiled expression
func (a *AutoResponse) updateRegex() {
	a.Regex, _ = regexp.Compile(a.Trigger)
}

//ReadFromJSON builds a list of autoresponses based on a json file
func ReadFromJSON() []*AutoResponse {
	var ars []*AutoResponse
	jsonResponse, err1 := ioutil.ReadFile("responses.json")
	botutils.Check(err1)

	err := json.Unmarshal(jsonResponse, &ars)
	botutils.Check(err)

	for i := range ars {
		ars[i].updateRegex()
	}
	return ars
}
