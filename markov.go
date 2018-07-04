package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"

	"markov/chain"
	"markov/user"

	"github.com/bwmarrin/discordgo"
)

const (
	m   = "463160554027221004"
	len = 2
)

var (
	punctuation = regexp.MustCompile(`[""''\?,!]`)
	c           = chain.NewChain(len)
)

func main() {

	token, err := ioutil.ReadFile("token.txt")
	if err != nil {
		fmt.Println("Error reading token")
		return
	}

	oops := user.User{
		Username: "Cristianop1",
	}

	for i := 0; i < 15; i++ {
		oops.GenerateComments(i)
	}

	for _, comment := range oops.Comments {
		new := punctuation.ReplaceAllString(comment.Content, "")
		c.AddComment(new)
	}

	dg, _ := discordgo.New("Bot " + string(token))

	dg.AddHandler(onReady)

	dg.Open()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()

}

func onReady(s *discordgo.Session, e *discordgo.Ready) {
	fmt.Println("Bot is ready")
	for {
		embed := &discordgo.MessageEmbed{
			Author: &discordgo.MessageEmbedAuthor{
				IconURL: "https://cdn.discordapp.com/app-icons/463133518869037056/5b80155a57f4223252f41c36875564c8.png",
				Name:    "OOPS! Studio",
			},
			Description: c.Generate(rand.Intn(30)),
			Color:       0xffb656,
		}
		s.ChannelMessageSendEmbed("463897602207907840", embed)
		time.Sleep(time.Second * 10)
	}
}
