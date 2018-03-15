package main

import "github.com/bwmarrin/discordgo"

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

func (r *ReactionResponse) respond(session *discordgo.Session, message *discordgo.MessageCreate) bool {
	for _, emoji := range r.emojis {
		session.MessageReactionAdd(message.ChannelID, message.ID, emoji)
	}
	return true
}
