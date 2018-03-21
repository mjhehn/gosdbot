package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	gosdbot "github.com/mjhehn/semi-decent-bot"

	"github.com/bwmarrin/discordgo"
)

var config *gosdbot.Config //config global object

func init() {
	config = gosdbot.ConfigFromJSON() //build the config file
	config.Ars = gosdbot.ReadFromJSON()
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

//ready when bot ready, just gosdbottains the status of the bot.
func ready(session *discordgo.Session, message *discordgo.MessageCreate) {
	session.UpdateStatus(0, config.Status)
}

//called every time a message received in a channel the bot is in
func messageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {
	currentServer := gosdbot.GetServer(session, message)

	if message.Author.ID == session.State.User.ID || message.Author.Bot { // Ignore all messages created by the bot itself(un-needed given the second check) or another bot
		return
	}

	responded := make(chan bool) //build the response channel

	//begin the bot proper!
	go config.Botstatus(session, message, responded)
	go config.Mute(session, message, responded)
	go config.Unmute(session, message, responded)
	go config.Mutestatus(session, message, responded)
	if gosdbot.In(currentServer, config.MutedServers) { //if the message is from a muted server, and the list wasn't updated by the mute commands, return
		return
	}
	if gosdbot.CheckWebHooks(session, message, "cleverbot") {
		go gosdbot.CleverResponse(session, message, responded)
	} else {
		for _, autoresponse := range config.Ars { //parse through all the json-configured responses
			if autoresponse.ServerSpecific == nil || gosdbot.In(currentServer, autoresponse.ServerSpecific) {
				go autoresponse.CheckResponses(session, message, gosdbot.Textresponse, responded)
				go autoresponse.CheckResponses(session, message, gosdbot.Embedresponse, responded)
				go autoresponse.CheckResponses(session, message, gosdbot.Reactionresponse, responded)
			}
		}
		go gosdbot.NotJustTheMen(session, message, responded)
		go gosdbot.DiceRoller(session, message, responded)
		go gosdbot.Compliment(session, message, responded)
		go gosdbot.Delete(session, message, responded)
		go gosdbot.Cleanup(session, message, responded)
	}

	<-responded //to synchronize back up with the coroutines
	close(responded)
}
