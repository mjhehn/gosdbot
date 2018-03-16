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

func (t *TextResponse) respond(session *discordgo.Session, message *discordgo.MessageCreate, mentions []string) bool {
	allMentions := ""
	if mentions != nil {
		for i := range mentions {
			if mentions[i] == "self" {
				allMentions += " " + message.Author.Mention()
			} else {
				//otherwise mentions don't work! stupid usernames.
			}
		}
	}
	session.ChannelMessageSend(message.ChannelID, (allMentions + " " + t.Text))
	return true
}

func (t *TextResponse) chance() bool {
	return runIt(t.Chance)
}
