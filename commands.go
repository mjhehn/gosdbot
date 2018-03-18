//responders with !word triggers
package main

import (
	"crypto/rand"
	"fmt"
	"godiscordbot/pkg/botutils"
	"math/big"
	"regexp"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

//rolls a variable number and type of dice. same format of parameters as other responders.
func diceRoller(session *discordgo.Session, message *discordgo.MessageCreate, responded chan bool) {
	expr, err := regexp.Compile("^!(ir|r)oll [0-9]{0,3}d[0-9]{1,3}$")
	botutils.Check(err)
	if expr.MatchString(message.Content) {
		roll := strings.Split(message.Content, " ")
		dice := strings.Split(roll[1], "d") //base case for no specified number of dice to roll

		var numDice int64 = 1
		if dice[0] != "" { //basically, if a number of dice is specified
			numDice, _ = strconv.ParseInt(dice[0], 10, 64)
		}
		numFaces, _ := strconv.ParseInt(dice[1], 10, 64)

		var rollString string
		var rollTotal int64
		for i := int64(0); i < numDice; i++ {
			rolledDie := RollDie(numFaces)
			rollTotal += rolledDie
			if len(rollString) > 0 {
				rollString = rollString + " + " + strconv.FormatInt(rolledDie, 10)
			} else {
				rollString = strconv.FormatInt(rolledDie, 10)
			}
		}

		individualRolls, err := regexp.MatchString("i", message.Content) //iroll #d# to show the results of every roll
		botutils.Check(err)
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

//pulls a 'compliment' from a json list found online from emergencycompliment.com and displays it
func compliment(session *discordgo.Session, message *discordgo.MessageCreate, responded chan bool) {
	expr, err := regexp.Compile("^!compliment$")
	botutils.Check(err)
	if expr.MatchString(message.Content) {
		var data map[string]interface{}
		err2 := botutils.GetJSON("https://spreadsheets.google.com/feeds/list/1eEa2ra2yHBXVZ_ctH4J15tFSGEu-VTSunsrvaCAV598/od6/public/values?alt=json", &data)
		botutils.Check(err2)
		numOptions := len(data["feed"].(map[string]interface{})["entry"].([]interface{})) //get the list of possible 'compliments
		selection := RollDie(int64(numOptions))

		//aaaaand the following line makes me feel sick.
		complimentText := (data["feed"].(map[string]interface{})["entry"].([]interface{})[selection].(map[string]interface{})["title"].(map[string]interface{})["$t"])

		session.ChannelMessageSend(message.ChannelID, (message.Author.Mention() + " " + complimentText.(string)))
	}
}

//removes messages from bots within a range. defaults to 100 messages back to clean
func cleanup(session *discordgo.Session, message *discordgo.MessageCreate, responded chan bool) {
	expr, err := regexp.Compile("^!cleanup[ ]{0,1}[0-9]*$")
	botutils.Check(err)
	if expr.MatchString(message.Content) {
		stringNum := strings.Split(message.Content, " ")
		numToClean := int64(99)
		if len(stringNum) > 1 {
			numToClean, _ = strconv.ParseInt(stringNum[1], 10, 0)
		}
		botutils.Check(err)
		numToClean++          //to handle cleaning the invoking command
		if numToClean > 100 { //handle out of range problems.
			numToClean = 100
		} else if numToClean < 0 {
			numToClean = 0
		}

		numCleaned := 0                                                                          //how many messages get removed/cleaned up
		messages, err := session.ChannelMessages(message.ChannelID, int(numToClean), "", "", "") //get the list of messages
		botutils.Check(err)

		var stringMessages = []string{message.ID} //build a list of message ids to pass to delete
		for _, message := range messages {
			if message.Author.Bot {
				stringMessages = append(stringMessages, message.ID)
				numCleaned++
			}
		}
		err2 := session.ChannelMessagesBulkDelete(message.ChannelID, stringMessages)
		botutils.Check(err2)

		result := fmt.Sprintf("%d messages deleted.", numCleaned)
		session.ChannelMessageSend(message.ChannelID, result)
		responded <- true
	}

}

//delete a number of messages. same basically as cleanup
func delete(session *discordgo.Session, message *discordgo.MessageCreate, responded chan bool) {
	expr, err := regexp.Compile("^!delete [0-9]+$")
	botutils.Check(err)
	if expr.MatchString(message.Content) {
		numToDelete, err := strconv.ParseInt(strings.Split(message.Content, " ")[1], 10, 0)
		botutils.Check(err)
		numToDelete++
		if numToDelete > 100 {
			numToDelete = 100
		} else if numToDelete < 0 {
			numToDelete = 0
		}
		messages, err := session.ChannelMessages(message.ChannelID, int(numToDelete), "", "", "")
		botutils.Check(err)
		var stringMessages []string
		for _, message := range messages {
			stringMessages = append(stringMessages, message.ID)
		}
		err2 := session.ChannelMessagesBulkDelete(message.ChannelID, stringMessages)
		botutils.Check(err2)

		result := fmt.Sprintf("%d messages deleted.", numToDelete)
		session.ChannelMessageSend(message.ChannelID, result)
		responded <- true
	}

}

//add current sever to the muted list, which allows only mute commands to be received or sent.
func mute(session *discordgo.Session, message *discordgo.MessageCreate, responded chan bool) {
	expr, err := regexp.Compile("^!mute$")
	botutils.Check(err)

	//mute
	if expr.MatchString(message.Content) {
		currentServer := botutils.GetServer(session, message)
		currentRoles := botutils.GetRoles(session, message)
		if !botutils.In(currentServer, config.MutedServers) && botutils.In("Bot Admin", currentRoles) { //mute
			config.MutedServers = append(config.MutedServers, currentServer)
			session.ChannelMessageSend(message.ChannelID, "Bot muted.")
			responded <- true
			return
		}
	}
	return
}

func unmute(session *discordgo.Session, message *discordgo.MessageCreate, responded chan bool) {
	expr2, err2 := regexp.Compile("^!unmute$")
	botutils.Check(err2)

	if expr2.MatchString(message.Content) {
		currentServer := botutils.GetServer(session, message)
		currentRoles := botutils.GetRoles(session, message)

		if botutils.In(currentServer, config.MutedServers) && botutils.In("Bot Admin", currentRoles) { //mute
			for i, serv := range config.MutedServers {
				if serv == currentServer {
					config.MutedServers = append(config.MutedServers[:i], config.MutedServers[i+1:]...)
				}
			}
			session.ChannelMessageSend(message.ChannelID, "Bot unmuted.")
			responded <- true
			return
		}
	}
	return
}

func mutestatus(session *discordgo.Session, message *discordgo.MessageCreate, responded chan bool) {
	expr3, err3 := regexp.Compile("^!mutestatus$")
	botutils.Check(err3)

	if expr3.MatchString(message.Content) { //mutestatus
		currentServer := botutils.GetServer(session, message)
		mutedStatus := " "
		if botutils.In(currentServer, config.MutedServers) {
			mutedStatus = "muted"
		} else {
			mutedStatus = "not muted"
		}
		result := fmt.Sprintf("Semi-Decent is currently %v on this server", mutedStatus)
		session.ChannelMessageSend(message.ChannelID, result)
		responded <- true
		return
	}
	return
}

//prequel meme
func notJustTheMen(session *discordgo.Session, message *discordgo.MessageCreate, responded chan bool) {
	expr, err := regexp.Compile("[mM]en")
	botutils.Check(err)
	if expr.MatchString(message.Content) && botutils.RunIt(8) {
		var menWord string
		messageLow := strings.Split(strings.ToLower(message.Content), " ")
		for _, word := range messageLow {
			if len(word) > 20 {
				return
			} else if expr.MatchString(word) {
				menWord = word
			}
		}
		result := fmt.Sprintf("and not just the %v, but the %v, and the %v too.", menWord, strings.Replace(menWord, "men", "women", 1), strings.Replace(menWord, "men", "children", 1))
		session.ChannelMessageSend(message.ChannelID, result)
		responded <- true
		return
	}
}

//RollDie 'rolls a die', but more important is a bounded number generator from 1-numFaces
func RollDie(numFaces int64) int64 {
	rnged, err := rand.Int(rand.Reader, big.NewInt(numFaces))
	botutils.Check(err)
	return rnged.Int64() + 1
}
