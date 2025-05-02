package main

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Player struct {
	conn       *websocket.Conn
	send       chan []byte
	opponent   *Player
	match      *Match
	moveTimer  *time.Timer
	playerMark string // "X" or "O"
}

func NewPlayer(conn *websocket.Conn) *Player {
	return &Player{
		conn: conn,
		send: make(chan []byte),
	}
}

func (p *Player) Listen() {
	go p.readMessages()
	go p.writeMessages()
}

func (p *Player) readMessages() {
	defer func() {
		if p.match != nil {
			p.match.PlayerLeft(p)
		}

		MatchmakerInstance.RemovePlayer(p)
		p.conn.Close()
	}()

	for {
		_, message, err := p.conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}
		log.Printf("Received: %s", message)
		HandleIncomingMessage(p, message)
	}
}

func (p *Player) writeMessages() {
	for {
		message, ok := <-p.send
		if !ok {
			return
		}
		err := p.conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println("Write error:", err)
			return
		}
	}
}

func (p *Player) StartMoveTimer() {
	if p.moveTimer != nil {
		p.moveTimer.Stop()
	}
	p.moveTimer = time.AfterFunc(30*time.Second, func() {
		log.Println("Timeout! Player took too long.")
		if p.match != nil {
			p.match.EndGameDueToTimeout()
		}
	})
}

func (p *Player) StopMoveTimer() {
	if p.moveTimer != nil {
		p.moveTimer.Stop()
	}
}
