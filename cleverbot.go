package gosdbot

import (
	cleverbot "github.com/CleverbotIO/go-cleverbot.io"
	"github.com/bwmarrin/discordgo"
)

//CleverResponse parses requests from channels tagged with the cleverbot webhook and mentioning this bot, and then asks them to cleverbot.
func CleverResponse(session *discordgo.Session, message *discordgo.MessageCreate, responded chan bool, clvrbot *cleverbot.Session) {
	if CheckWebHooks(session, message, "cleverbot") {
		if (len(message.Mentions) == 1 && message.Mentions[0].ID == session.State.User.ID) || CheckWebHooks(session, message, "nomention") {
			//send call to the cleverbot.
			result, err := clvrbot.Ask(message.Content)
			Check(err)
			session.ChannelMessageSend(message.ChannelID, (message.Author.Mention() + " " + result))
			responded <- true
		}
	}
	return
}
