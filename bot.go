package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var config *ServerConfig //config global object

func main() {
	config = NewServerConfig()                                                                            //build the config file
	discord, err := discordgo.New("Bot " + "NDE1NzMxMDQ0MDgwNDE4ODE2.DW6Lpw.LYY0ZKfCKvC3YMJbCe1u0XpIF7M") //dev bot
	//discord, err := discordgo.New("Bot " + "NDA4ODQ1OTY4MDEyNzM4NTYw.DVV_9w.wRs92tvW30aAmW8JOMgqB2GzFQY") //main bot

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

//ready when bot ready, just maintains the status of the bot.
func ready(session *discordgo.Session, message *discordgo.MessageCreate) {
	session.UpdateStatus(0, "with 100% Python!")
}

//called every time a message received in a channel the bot is in
func messageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {
	currentServer := getServer(session, message)

	if message.Author.ID == session.State.User.ID || message.Author.Bot { // Ignore all messages created by the bot itself(un-needed given the second check) or another bot
		return
	}

	responded := make(chan bool) //build the response channel

	//begin the bot proper!
	go mute(session, message, responded)
	go unmute(session, message, responded)
	go mutestatus(session, message, responded)
	if in(currentServer, config.mutedServers) { //if the message is from a muted server, and the list wasn't updated by the mute commands, return
		return
	}

	for _, autoresponse := range config.ars { //parse through all the json-configured responses
		if autoresponse.ServerSpecific == nil || in(currentServer, autoresponse.ServerSpecific) {
			go autoresponse.checkResponses(session, message, textresponse, responded)
			go autoresponse.checkResponses(session, message, embedresponse, responded)
			go autoresponse.checkResponses(session, message, reactionresponse, responded)
		}
	}
	go notJustTheMen(session, message, responded)
	go diceRoller(session, message, responded)
	go compliment(session, message, responded)
	go delete(session, message, responded)
	go cleanup(session, message, responded)
	<-responded //to synchronize back up with the coroutines
	close(responded)
}
