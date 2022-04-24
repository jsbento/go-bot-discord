package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	Token string
)

const DogAPIURL = "https://dog.ceo/api/"

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("Error creating Discord session,", err)
		return
	}

	dg.AddHandler(messageCreate)
	dg.Identify.Intents = discordgo.IntentGuildMessages

	// Open websocket connection to Discord to begin listening
	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening connection,", err)
		return
	}

	// Wait until CTRL-C or other termination signal
	fmt.Println("Bot running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Close Discord session
	dg.Close()
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// Function called whenever a message is created in any channel the bot has access to
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages from the bot
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "!dog" {
		// Call KuteGo API for dr-who Gopher
		res, err := http.Get(DogAPIURL + "breeds/image/random")
		if err != nil {
			fmt.Println(err)
		}
		defer res.Body.Close()

		if res.StatusCode == 200 {
			body, err := ioutil.ReadAll(res.Body)
			var data Response
			errJS := json.Unmarshal(body, &data)
			if errJS != nil {
				fmt.Println("Error parsing JSON response!")
			}
			_, err = s.ChannelMessageSend(m.ChannelID, data.Message)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Error: Can't get random dog picture!")
		}
	}
}
