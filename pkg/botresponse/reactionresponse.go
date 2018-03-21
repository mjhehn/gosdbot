package botresponse

import (
	"semi-decent-bot/pkg/botutils"

	"github.com/bwmarrin/discordgo"
)

//ReactionResponse stores the chance a text response will be made and the text itself
type ReactionResponse struct {
	Chance int
	Emojis []string
}

//NewReactionResponse ...
func NewReactionResponse(setChance int, setText []string) *ReactionResponse {
	e := new(ReactionResponse)
	e.Chance = setChance
	e.Emojis = setText
	return e
}

func (r *ReactionResponse) respond(session *discordgo.Session, message *discordgo.MessageCreate, mentions []string) bool {
	for _, emoji := range r.Emojis {
		emojiObject := botutils.GetEmoji(session, message, emoji)
		if emojiObject != nil {
			session.MessageReactionAdd(message.ChannelID, message.ID, emojiObject.APIName())
		} else {
			session.MessageReactionAdd(message.ChannelID, message.ID, emoji)
		}
	}
	return true
}

func (r *ReactionResponse) chance() bool {
	return botutils.RunIt(r.Chance)
}
