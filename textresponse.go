package main

import "github.com/bwmarrin/discordgo"

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

func (t *TextResponse) respond(session *discordgo.Session, message *discordgo.MessageCreate) bool {
	session.ChannelMessageSend(message.ChannelID, t.text)
	return true
}
