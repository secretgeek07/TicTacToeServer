package main

import (
	"encoding/json"
	"log"
)

type Message struct {
	Action string `json:"action"`
	Row    int    `json:"row,omitempty"`
	Col    int    `json:"col,omitempty"`
}

func HandleIncomingMessage(p *Player, rawMsg []byte) {
	var msg Message
	err := json.Unmarshal(rawMsg, &msg)
	if err != nil {
		log.Println("Error decoding message:", err)
		return
	}

	switch msg.Action {
	case "move":
		if p.match != nil && p.match.active {
			p.match.ProcessMove(p, msg.Row, msg.Col)
		}
	case "join":
		log.Println("Player sent join request.")
	default:
		log.Println("Unknown action:", msg.Action)
	}
}
