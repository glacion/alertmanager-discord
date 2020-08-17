package discord

import "bigo.digital/alertmanager-discord/pkg/alertmanager"

// Colors
const (
	colorRed   = 10038562
	colorGreen = 3066993
	colorGrey  = 9807270
)

// Payload Struct for the message to send to discord
type Payload struct {
	Embeds []embed `json:"embeds"`
}

type embed struct {
	Title  string       `json:"title"`
	Color  int          `json:"color"`
	Fields []embedField `json:"fields"`
}

type embedField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// MakeMessage Create discord message from a Alertmanager payload
func MakeMessage(payload alertmanager.Payload) Payload {
	message := Payload{Embeds: []embed{}}

	embedMessage := embed{
		Title:  payload.GroupLabels.AlertName,
		Fields: []embedField{},
	}

	switch payload.Status {
	case "firing":
		embedMessage.Color = colorRed
	case "resolved":
		embedMessage.Color = colorGreen
	default:
		embedMessage.Color = colorGrey
	}

	for _, alert := range payload.Alerts {
		field := embedField{
			Name:  alert.Labels.AlertName,
			Value: alert.Annotations.Message,
		}
		embedMessage.Fields = append(embedMessage.Fields, field)
	}

	message.Embeds = append(message.Embeds, embedMessage)
	return message
}
