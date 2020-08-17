package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var (
	webhookURL string
	logger     *log.Logger
)

// Discord color values
const (
	ColorRed   = 10038562
	ColorGreen = 3066993
	ColorGrey  = 9807270
)

type alertManAlert struct {
	Annotations struct {
		Message string `json:"message"`
	} `json:"annotations"`
	Labels struct {
		AlertName string `json:"alertname"`
	} `json:"labels"`
}

type alertManOut struct {
	Alerts      []alertManAlert `json:"alerts"`
	Status      string          `json:"status"`
	GroupLabels struct {
		AlertName string `json:"alertname"`
	} `json:"groupLabels"`
}

type discordOut struct {
	Embeds []discordEmbed `json:"embeds"`
}

type discordEmbed struct {
	Title  string              `json:"title"`
	Color  int                 `json:"color"`
	Fields []discordEmbedField `json:"fields"`
}

type discordEmbedField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func sendMessage(amo alertManOut) {
	DO := discordOut{Embeds: make([]discordEmbed, 0)}

	RichEmbed := discordEmbed{
		Title:  amo.GroupLabels.AlertName,
		Fields: make([]discordEmbedField, 0),
	}

	switch amo.Status {
	case "firing":
		RichEmbed.Color = ColorRed
	case "resolved":
		RichEmbed.Color = ColorGreen
	default:
		RichEmbed.Color = ColorGrey
	}

	for _, alert := range amo.Alerts {
		field := discordEmbedField{
			Name:  alert.Labels.AlertName,
			Value: alert.Annotations.Message,
		}
		RichEmbed.Fields = append(RichEmbed.Fields, field)
	}

	DO.Embeds = append(DO.Embeds, RichEmbed)

	DOD, err := json.Marshal(DO)
	if err != nil {
		logger.Fatalf("%+v\n", err)
		return
	}
	logger.Printf("Sending to discord as %s\n", string(DOD))
	response, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(DOD))
	if err != nil {
		logger.Fatalf("%+v\n", response)
		logger.Fatalf("%+v\n", err)
	} else {
		logger.Printf("%+v\n", response)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Fatalf("%+v\n", err)
		return
	}

	amo := alertManOut{}
	err = json.Unmarshal(b, &amo)
	if err != nil {
		logger.Fatalf("%+v\n", err)
		return
	}
	logger.Printf("Received alert %+v\n", amo)
	sendMessage(amo)
}

func main() {
	logger = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	webhookURL = os.Getenv("DISCORD_WEBHOOK")

	if webhookURL == "" {
		logger.Fatalln("Environment variable DISCORD_WEBHOOK not found")
		os.Exit(1)
	}

	logger.Println("Listening on 0.0.0.0:9094")
	logger.Fatalln(http.ListenAndServe(":9094", http.HandlerFunc(handler)))
}
