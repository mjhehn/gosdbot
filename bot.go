package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var ars = []*AutoResponse{}

type responder interface {
	respond(session *discordgo.Session, message *discordgo.MessageCreate) bool
}

func main() {
	discord, err := discordgo.New("Bot " + "NDE1NzMxMDQ0MDgwNDE4ODE2.DW6Lpw.LYY0ZKfCKvC3YMJbCe1u0XpIF7M")
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}
	// Register ready as a callback for the ready events.
	discord.AddHandler(ready)

	discord.AddHandler(messageCreate)

	err = discord.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	readResponseList()
	//buildResponseList()
	//writeResponseList()

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	discord.Close()
}

func ready(session *discordgo.Session, message *discordgo.MessageCreate) {
	session.UpdateStatus(0, "Being a bot")
}

//called every time a message received in a channel the bot is in
func messageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself(un-needed given the second check) or another bot
	if message.Author.ID == session.State.User.ID || message.Author.Bot {
		return
	}

	responded := make(chan bool)

	for _, autoresponse := range ars {
		go autoresponse.checkTextResponses(session, message, responded)
		go autoresponse.checkEmbedResponses(session, message, responded)
		go autoresponse.checkReactionResponses(session, message, responded)
	}
	go diceRoller(session, message, responded)
	<-responded //to synchronize back up with the coroutines
}

func diceRoller(session *discordgo.Session, message *discordgo.MessageCreate, responded chan bool) {
	expr, err := regexp.Compile("^!(ir|r)oll [0-9]{0,3}d[0-9]{1,3}$")
	check(err)
	if expr.MatchString(message.Content) {
		//TODO: get the number and type of dice to roll.
		roll := strings.Split(message.Content, " ")
		dice := strings.Split(roll[1], "d") //base case for no specified number of dice to roll

		var numDice int64 = 1
		if dice[0] != "" {
			numDice, _ = strconv.ParseInt(dice[0], 10, 64)
		}
		numFaces, _ := strconv.ParseInt(dice[1], 10, 64)

		var rollString string
		var rollTotal int64
		for i := int64(0); i < numDice; i++ {
			rolledDie := rollDie(numFaces)
			rollTotal += rolledDie
			if len(rollString) > 0 {
				rollString = rollString + " + " + strconv.FormatInt(rolledDie, 10)
			} else {
				rollString = strconv.FormatInt(rolledDie, 10)
			}
		}

		individualRolls, err := regexp.MatchString("i", message.Content)
		check(err)
		var result string
		if individualRolls {
			result = fmt.Sprintf("Rolled %dd%d for %d. (%s)", numDice, numFaces, rollTotal, rollString)
		} else {
			result = fmt.Sprintf("Rolled %dd%d for %d.", numDice, numFaces, rollTotal)
		}
		session.ChannelMessageSend(message.ChannelID, result)
		responded <- true
		return
	}
}

func buildResponseList() {
	//TODO: setup list of autoresponses and how to load them up (json? csv? seperate lists for embeds vs text?)
	ars = append(ars, NewAutoResponse("[p|P]ing", []*TextResponse{NewTextResponse(1, "pong")}, nil, nil, nil, false, nil, nil))
	ars = append(ars, NewAutoResponse("[p|P]ling", []*TextResponse{NewTextResponse(1, "pong")}, nil, nil, nil, false, nil, nil))
	ars = append(ars, NewAutoResponse("[p|P]sing", []*TextResponse{NewTextResponse(1, "pong")}, nil, nil, nil, false, nil, nil))

	ars = append(ars, NewAutoResponse("pong", nil, []*EmbedResponse{NewEmbedResponse(1, "https://starscollideani.files.wordpress.com/2014/08/cirno-bsod.png?w=332&h=199")}, nil, nil, false, nil, nil))
	ars = append(ars, NewAutoResponse("handholding", nil, nil, []*ReactionResponse{NewReactionResponse(1, []string{"\U0001F1F1", "\U0001F1EA", "\U0001F1FC", "\U0001F1E9"})}, nil, false, nil, nil))
}

func readResponseList() {
	jsonResponse, err1 := ioutil.ReadFile("responses.json")
	check(err1)

	err := json.Unmarshal(jsonResponse, &ars)
	check(err)

	for i := range ars {
		ars[i].updateRegex()
	}
}

func writeResponseList() {
	jsonResponse, _ := json.Marshal(ars)
	f, _ := os.Create("responses.json")
	defer f.Close()
	f.Write(jsonResponse)
	f.Sync()
}

func check(err error) {
	if err != nil {
		fmt.Println(err.Error)
	}
}

func rollDie(numFaces int64) int64 {
	rnged, err := rand.Int(rand.Reader, big.NewInt(numFaces))
	check(err)
	return rnged.Int64() + 1
}
