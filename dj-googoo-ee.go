package main

import (
	//	"bytes"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/nlopes/slack"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Println("usage: dj-googoo-ee slack-bot-token iftttToken slack-channel-id")
		os.Exit(1)
	}
	slackBotToken := os.Args[1]
	iftttToken := os.Args[2]
	slackChannelId := os.Args[3]
	fmt.Println("slack-bot-token: " + slackBotToken)
	fmt.Println("iftttToken: " + iftttToken)
	fmt.Println("slackChannelId: " + slackChannelId)

	api := slack.New(slackBotToken)
	logger := log.New(os.Stdout, "dj-googoo-ee: ", log.Lshortfile|log.LstdFlags)
	slack.SetLogger(logger)
	api.SetDebug(true)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		//fmt.Print("Event Received: ")
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
			fmt.Println("Connected to RTM API.")

		case *slack.ConnectedEvent:
			fmt.Println("Infos:", ev.Info)
			fmt.Println("Connection counter:", ev.ConnectionCount)
			//rtm.SendMessage(rtm.NewOutgoingMessage("Hello #Music, I add YouTube links to my Spotify Playlist: https://open.spotify.com/user/kickassmusac/playlist/66RYEqJXsHCffYPrbaV88e", slackChannelId))

		case *slack.MessageEvent:

			if ev.SubMessage == nil || !strings.Contains(ev.SubMessage.Text, "youtu") || !strings.Contains(ev.SubMessage.Attachments[0].Title, "-") || ev.SubType == "message_replied" {
				continue
			}
			songInfo := strings.Split(ev.SubMessage.Attachments[0].Title, "-")
			fmt.Printf("YouTube Song detected: artist %s title %s\n", songInfo[0], songInfo[1])

			url := fmt.Sprintf("https://maker.ifttt.com/trigger/song_event/with/key/%s", iftttToken)
			values := map[string]string{"value1": songInfo[1], "value2": songInfo[0]}

			payload, _ := json.Marshal(values)

			resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
			if err != nil {
				fmt.Println(err)
			}
			_, readErr := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			if readErr != nil {
				fmt.Println(err)
			}
			fmt.Println("Added song to playlist")
			output := fmt.Sprintf("%s - %s added to the playlist.", songInfo[0], songInfo[1])
			rtm.SendMessage(rtm.NewOutgoingMessage(output, "C1MPYM2AE"))

		default:

			// Ignore other events..
			// fmt.Printf("Unexpected: %v\n", msg.Data)
		}
	}
}
