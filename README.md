# DJ-GOOGOO-EE

A slack bot that listens to the music channel for YouTube links and appends them to a Spotify playlist for office use on the Sonos.
Uses an integration with IFTT maker applet that must be setup before hand.

## Requirements

* GoLang
* Slack bot key
* IFTT maker key
* IFTT applet that waits for a web request then adds songs to a Spotify playlist

## Usage

1. Build it using `go build dj-googoo-ee.go`
2. Run it using `./dj-googoo-ee <SLACK_BOT_KEY> <IFTT_KEY>`
