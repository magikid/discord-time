package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

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

	user, err := dg.User("@me")
	log.Printf("id: %s, username: %s", user.ID, user.Username)

	dg.AddHandler(ready)
	dg.AddHandler(incomingMessageHandler)

	dg.Identify.Intents = discordgo.IntentsGuildMessages

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
	guild, err := s.Guild("689512841887481875")
	check(err)
	log.Print(guild.Name)
}

func incomingMessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.Contains(m.Author.Username, "magikid") {
		return
	}

	if m.Content != "ping" {
		return
	}

	channel, err := s.UserChannelCreate(m.Author.ID)
	check(err)
	_, err = s.ChannelMessageSend(channel.ID, "pong")
	check(err)

	log.Printf("Recieved message: %s", m.Content)
}
