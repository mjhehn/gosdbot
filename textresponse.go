package main

import "github.com/bwmarrin/discordgo"

//TextResponse stores the chance a text response will be made and the text itself
type TextResponse struct {
	Chance int
	Text   string
}

//NewTextResponse ...
func NewTextResponse(setChance int, setText string) *TextResponse {
	e := new(TextResponse)
	e.Chance = setChance
	e.Text = setText
	return e
}

func (t *TextResponse) respond(session *discordgo.Session, message *discordgo.MessageCreate) bool {
	session.ChannelMessageSend(message.ChannelID, t.Text)
	return true
}
