//responders with !word triggers
package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

//rolls a variable number and type of dice. same format of parameters as other responders.
func diceRoller(session *discordgo.Session, message *discordgo.MessageCreate, responded chan bool) {
	expr, err := regexp.Compile("^!(ir|r)oll [0-9]{0,3}d[0-9]{1,3}$")
	check(err)
	if expr.MatchString(message.Content) {
		//TODO: get the number and type of dice to roll.
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
			rolledDie := rollDie(numFaces)
			rollTotal += rolledDie
			if len(rollString) > 0 {
				rollString = rollString + " + " + strconv.FormatInt(rolledDie, 10)
			} else {
				rollString = strconv.FormatInt(rolledDie, 10)
			}
		}

		individualRolls, err := regexp.MatchString("i", message.Content) //iroll #d# to show the results of every roll
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

//pulls a 'compliment' from a json list found online from emergencycompliment.com and displays it
func compliment(session *discordgo.Session, message *discordgo.MessageCreate, responded chan bool) {
	expr, err := regexp.Compile("^!compliment$")
	check(err)
	if expr.MatchString(message.Content) {
		var data map[string]interface{}
		err2 := getJSON("https://spreadsheets.google.com/feeds/list/1eEa2ra2yHBXVZ_ctH4J15tFSGEu-VTSunsrvaCAV598/od6/public/values?alt=json", &data)
		check(err2)
		numOptions := len(data["feed"].(map[string]interface{})["entry"].([]interface{})) //get the list of possible 'compliments
		selection := rollDie(int64(numOptions))

		//aaaaand the following line makes me feel sick.
		complimentText := (data["feed"].(map[string]interface{})["entry"].([]interface{})[selection].(map[string]interface{})["title"].(map[string]interface{})["$t"])

		session.ChannelMessageSend(message.ChannelID, (message.Author.Mention() + " " + complimentText.(string)))
	}
}

//removes messages from bots within a range. defaults to 100 messages back to clean
func cleanup(session *discordgo.Session, message *discordgo.MessageCreate, responded chan bool) {
	expr, err := regexp.Compile("^!cleanup[ ]{0,1}[0-9]*$")
	check(err)
	if expr.MatchString(message.Content) {
		stringNum := strings.Split(message.Content, " ")
		numToClean := int64(99)
		if len(stringNum) > 1 {
			numToClean, _ = strconv.ParseInt(stringNum[1], 10, 0)
		}
		check(err)
		numToClean++          //to handle cleaning the invoking command
		if numToClean > 100 { //handle out of range problems.
			numToClean = 100
		} else if numToClean < 0 {
			numToClean = 0
		}

		numCleaned := 0                                                                          //how many messages get removed/cleaned up
		messages, err := session.ChannelMessages(message.ChannelID, int(numToClean), "", "", "") //get the list of messages
		check(err)

		var stringMessages = []string{message.ID} //build a list of message ids to pass to delete
		for _, message := range messages {
			if message.Author.Bot {
				stringMessages = append(stringMessages, message.ID)
				numCleaned++
			}
		}
		err2 := session.ChannelMessagesBulkDelete(message.ChannelID, stringMessages)
		check(err2)

		result := fmt.Sprintf("%d messages deleted.", numCleaned)
		session.ChannelMessageSend(message.ChannelID, result)
		responded <- true
	}

}

//delete a number of messages. same basically as cleanup
func delete(session *discordgo.Session, message *discordgo.MessageCreate, responded chan bool) {
	expr, err := regexp.Compile("^!delete [0-9]+$")
	check(err)
	if expr.MatchString(message.Content) {
		numToDelete, err := strconv.ParseInt(strings.Split(message.Content, " ")[1], 10, 0)
		check(err)
		numToDelete++
		if numToDelete > 100 {
			numToDelete = 100
		} else if numToDelete < 0 {
			numToDelete = 0
		}
		messages, err := session.ChannelMessages(message.ChannelID, int(numToDelete), "", "", "")
		check(err)
		var stringMessages []string
		for _, message := range messages {
			stringMessages = append(stringMessages, message.ID)
		}
		err2 := session.ChannelMessagesBulkDelete(message.ChannelID, stringMessages)
		check(err2)

		result := fmt.Sprintf("%d messages deleted.", numToDelete)
		session.ChannelMessageSend(message.ChannelID, result)
		responded <- true
	}

}

//add current sever to the muted list, which allows only mute commands to be received or sent.
func mute(session *discordgo.Session, message *discordgo.MessageCreate, responded chan bool) {
	expr, err := regexp.Compile("^!mute$")
	check(err)

	//mute
	if expr.MatchString(message.Content) {
		currentServer := getServer(session, message)
		currentRoles := getRoles(session, message)
		if !in(currentServer, config.mutedServers) && in("Bot Admin", currentRoles) { //mute
			config.mutedServers = append(config.mutedServers, currentServer)
			session.ChannelMessageSend(message.ChannelID, "Bot muted.")
			responded <- true
			return
		}
	}
	return
}

func unmute(session *discordgo.Session, message *discordgo.MessageCreate, responded chan bool) {
	expr2, err2 := regexp.Compile("^!unmute$")
	check(err2)

	if expr2.MatchString(message.Content) {
		currentServer := getServer(session, message)
		currentRoles := getRoles(session, message)

		if in(currentServer, config.mutedServers) && in("Bot Admin", currentRoles) { //mute
			for i, serv := range config.mutedServers {
				if serv == currentServer {
					config.mutedServers = append(config.mutedServers[:i], config.mutedServers[i+1:]...)
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
	check(err3)

	if expr3.MatchString(message.Content) { //mutestatus
		currentServer := getServer(session, message)
		mutedStatus := " "
		if in(currentServer, config.mutedServers) {
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
	check(err)
	if expr.MatchString(message.Content) && runIt(8) {
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
