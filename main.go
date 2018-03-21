package sdbot

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var config *Config //config global object

func init() {
	config = ConfigFromJSON() //build the config file
	config.Ars = ReadFromJSON()
}

func main() {
	discord, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}
	// Register ready as a callback for the ready events.
	discord.AddHandler(ready)
	discord.AddHandler(messageCreate) //and callback for message creation on channels

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

//ready when bot ready, just sdbottains the status of the bot.
func ready(session *discordgo.Session, message *discordgo.MessageCreate) {
	session.UpdateStatus(0, config.Status)
}

//called every time a message received in a channel the bot is in
func messageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {
	currentServer := GetServer(session, message)

	if message.Author.ID == session.State.User.ID || message.Author.Bot { // Ignore all messages created by the bot itself(un-needed given the second check) or another bot
		return
	}

	responded := make(chan bool) //build the response channel

	//begin the bot proper!
	go mute(session, message, responded)
	go unmute(session, message, responded)
	go mutestatus(session, message, responded)
	if In(currentServer, config.MutedServers) { //if the message is from a muted server, and the list wasn't updated by the mute commands, return
		return
	}
	if CheckWebHooks(session, message, "cleverbot") {
		go CleverResponse(session, message, responded)
	} else {
		for _, autoresponse := range config.Ars { //parse through all the json-configured responses
			if autoresponse.ServerSpecific == nil || In(currentServer, autoresponse.ServerSpecific) {
				go autoresponse.CheckResponses(session, message, Textresponse, responded)
				go autoresponse.CheckResponses(session, message, Embedresponse, responded)
				go autoresponse.CheckResponses(session, message, Reactionresponse, responded)
			}
		}
		go notJustTheMen(session, message, responded)
		go diceRoller(session, message, responded)
		go compliment(session, message, responded)
		go delete(session, message, responded)
		go cleanup(session, message, responded)
		go botstatus(session, message, responded)
	}

	<-responded //to synchronize back up with the coroutines
	close(responded)
}
