package main

import (
	"log"
	"sync"
)

type Matchmaker struct {
	queue []*Player
	lock  sync.Mutex
}

var MatchmakerInstance = &Matchmaker{}

func (m *Matchmaker) AddPlayer(p *Player) {
	m.lock.Lock()
	defer m.lock.Unlock()

	log.Printf("Player joined queue. Queue length before adding: %d", len(m.queue))
	if len(m.queue) > 0 {
		opponent := m.queue[0]
		m.queue = m.queue[1:]

		p.opponent = opponent
		opponent.opponent = p

		p.playerMark = "X"
		opponent.playerMark = "O"

		// Create the match instance
		match := NewMatch(p, opponent)
		p.match = match
		opponent.match = match

		log.Println("Matching two players. Starting match...")
		// Send start messages
		startMessageP := []byte(`{"action":"start", "mark":"X"}`)
		startMessageO := []byte(`{"action":"start", "mark":"O"}`)

		log.Println("Going to send start message to both players")
		p.send <- startMessageP
		opponent.send <- startMessageO
		log.Println("Sent start message to both players")

	} else {
		m.queue = append(m.queue, p)
		log.Printf("No opponent found. Player added to queue. New queue length: %d", len(m.queue))

	}
}

func (m *Matchmaker) RemovePlayer(p *Player) {
	m.lock.Lock()
	defer m.lock.Unlock()

	for i, player := range m.queue {
		if player == p {
			m.queue = append(m.queue[:i], m.queue[i+1:]...)
			return
		}
	}

	if p.match != nil && p.match.active {
		p.match.lock.Lock()
		defer p.match.lock.Unlock()

		if p.match.active {
			p.match.active = false

			if p.opponent != nil {
				p.opponent.send <- []byte(`{"action":"opponent_left"}`)
				p.opponent.conn.Close()
			}
		}
	}
}

func (m *Matchmaker) Run() {
	select {}
}
