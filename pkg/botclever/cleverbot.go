package botclever

import (
	"semi-decent-bot/pkg/botutils"

	cleverbot "github.com/CleverbotIO/go-cleverbot.io"
	"github.com/bwmarrin/discordgo"
)

var clvrbot *cleverbot.Session

func init() {
	var err error
	clvrbot, err = cleverbot.New("IGHvKK0w5ozUyNlp", "Y2PWfAYhH43Aa9Fy1gMjg9OxuYFk3w7B")
	botutils.Check(err)
}

//CleverResponse parses requests from channels tagged with the cleverbot webhook and mentioning this bot, and then asks them to cleverbot.
func CleverResponse(session *discordgo.Session, message *discordgo.MessageCreate, responded chan bool) {
	if botutils.CheckWebHooks(session, message, "cleverbot") {
		if (len(message.Mentions) == 1 && message.Mentions[0].ID == session.State.User.ID) || botutils.CheckWebHooks(session, message, "nomention") {
			//send call to the cleverbot.
			result, err := clvrbot.Ask(message.Content)
			botutils.Check(err)
			session.ChannelMessageSend(message.ChannelID, (message.Author.Mention() + " " + result))
			responded <- true
		}
	}

	return
}
