package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var ars []*AutoResponse

//TODO: consider splitting these class definitions off into their own files
//AutoResponse stores all elements needed for an Autoresponse
type AutoResponse struct {
	trigger        string
	responses      []*TextResponse
	embeds         []*EmbedResponse
	reactions      []*ReactionResponse
	mentions       []string
	cleanup        bool
	userSpecific   []string
	serverSpecific []string
}

//NewAutoResponse ...
func NewAutoResponse(trigger string, responses []*TextResponse, embeds []*EmbedResponse, reactions []*ReactionResponse, mentions []string, cleanup bool, userSpecific []string, serverSpecific []string) *AutoResponse {
	a := new(AutoResponse)
	a.trigger = trigger
	a.responses = responses
	a.embeds = embeds
	a.reactions = reactions
	a.mentions = mentions
	a.cleanup = cleanup
	a.userSpecific = userSpecific
	a.serverSpecific = serverSpecific
	return a
}

//TextResponse stores the chance a text response will be made and the text itself
type TextResponse struct {
	chance int
	text   string
}

//NewTextResponse ...
func NewTextResponse(setChance int, setText string) *TextResponse {
	e := new(TextResponse)
	e.chance = setChance
	e.text = setText
	return e
}

//ReactionResponse stores the chance a text response will be made and the text itself
type ReactionResponse struct {
	chance int
	emojis []string
}

//NewReactionResponse ...
func NewReactionResponse(setChance int, setText []string) *ReactionResponse {
	e := new(ReactionResponse)
	e.chance = setChance
	e.emojis = setText
	return e
}

func main() {
	discord, err := discordgo.New("Bot " + "NDE1NzMxMDQ0MDgwNDE4ODE2.DW6Lpw.LYY0ZKfCKvC3YMJbCe1u0XpIF7M")
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}
	// Register ready as a callback for the ready events.
	discord.AddHandler(ready)

	discord.AddHandler(messageCreate)

	err = discord.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	discord.Close()
}

func ready(session *discordgo.Session, message *discordgo.MessageCreate) {
	session.UpdateStatus(0, "Being a bot")
	ars = append(ars, NewAutoResponse("ping", []*TextResponse{NewTextResponse(1, "pong")}, nil, nil, nil, false, nil, nil))
	ars = append(ars, NewAutoResponse("pong", nil, []*EmbedResponse{NewEmbedResponse(1, "https://starscollideani.files.wordpress.com/2014/08/cirno-bsod.png?w=332&h=199")}, nil, nil, false, nil, nil))
	ars = append(ars, NewAutoResponse("handholding", nil, nil, []*ReactionResponse{NewReactionResponse(1, []string{"\U0001F1F1", "\U0001F1EA", "\U0001F1FC", "\U0001F1E9"})}, nil, false, nil, nil))
}

//called every time a message received in a channel the bot is in
func messageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {
	//TODO: setup list of autoresponses and how to load them up (json? csv? seperate lists for embeds vs text?)
	//TODO: add regex for checks
	//TODO: add concurrency

	// Ignore all messages created by the bot itself(un-needed given the second check) or another bot
	if message.Author.ID == session.State.User.ID || message.Author.Bot {
		return
	}

	for _, autoresponse := range ars {
		//check for a textResponse to make
		//TODO: add RNG chance to do the particular trigger.
		for _, response := range autoresponse.responses {
			if response != nil && message.Content == autoresponse.trigger {
				session.ChannelMessageSend(message.ChannelID, response.text)
				return
			}
		}

		//check for an embedded response to make
		//TODO: add RNG chance to do the particular trigger.
		for _, embed := range autoresponse.embeds {
			if embed != nil && message.Content == autoresponse.trigger {
				session.ChannelMessageSendEmbed(message.ChannelID, embed.getEmbed())
				return
			}
		}

		//TODO: add the RNG trigger
		for _, reaction := range autoresponse.reactions {
			if reaction != nil && message.Content == autoresponse.trigger {
				for _, emoji := range reaction.emojis {
					session.MessageReactionAdd(message.ChannelID, message.ID, emoji)
				}
				return
			}
		}
	}

	//if message.Content == "ping" {
	//textresponse := NewTextResponse(3, "pong")
	//session.ChannelMessageSend(message.ChannelID, textresponse.text)

	//embed := NewEmbedResponse(3, "https://starscollideani.files.wordpress.com/2014/08/cirno-bsod.png?w=332&h=199")
	//session.ChannelMessageSendEmbed(message.ChannelID, embed.getEmbed())

	//session.MessageReactionAdd(message.ChannelID, message.ID, "\U0001F1F1")
	//}

}
