package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

func check(e error) {
	if e != nil {
		log.Fatalf("error: %v", e)
	}
}

func main() {
	bot_token := os.Getenv("BOT_TOKEN")
	dg, err := discordgo.New(fmt.Sprintf("Bot %s", bot_token))
	check(err)

	dg.AddHandler(ready)
	dg.AddHandler(incomingMessageHandler)

	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsDirectMessages

	err = dg.Open()
	check(err)

	fmt.Println("discord_time is now running!")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	log.Print("Bot ready")
}

func incomingMessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	go convertTime(s, m)

	if m.Content != "ping" {
		return
	}

	dmChannel, err := s.UserChannelCreate(m.Author.ID)
	check(err)
	_, err = s.ChannelMessageSend(dmChannel.ID, "pong")
	check(err)
}

func convertTime(s *discordgo.Session, m *discordgo.MessageCreate) {
	timeMatcher := regexp.MustCompile(`(\d{1,2}):?(\d{2})? ?([a,A,p,P][m,M])? (?P<timezone>[A-Za-z]{3,})`)

	if timeMatcher.MatchString(m.Content) {
		matches := timeMatcher.FindStringSubmatch(m.Content)

		hours := matches[1]
		minutes := matches[2]
		if len(minutes) < 2 {
			minutes = "00"
		}
		ampm := matches[3]
		timezone := matches[4]

		longform := fmt.Sprintf("%v:%v %v %v", hours, minutes, strings.ToUpper(ampm), strings.ToUpper(timezone))

		currentTime, err := time.Parse("3:04 PM MST", longform)
		if err != nil {
			log.Printf("Error parsing time: %v", err.Error())
			return
		}

		_, err = s.ChannelMessageSendReply(m.ChannelID, currentTime.UTC().Format("3:04 PM UTC"), m.Reference())
		check(err)
	}
}
