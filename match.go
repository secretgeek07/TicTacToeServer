package main

import (
	"encoding/json"
	"log"
	"sync"
)

type Match struct {
	player1     *Player
	player2     *Player
	active      bool
	currentTurn *Player
	lock        sync.Mutex
}

func NewMatch(p1, p2 *Player) *Match {
	match := &Match{
		player1: p1,
		player2: p2,
		active:  true,
	}
	if p1.playerMark == "X" {
		match.currentTurn = p1
	} else {
		match.currentTurn = p2
	}

	// match.currentTurn.StartMoveTimer()
	return match
}

func (m *Match) EndGameDueToTimeout() {
	m.lock.Lock()
	defer m.lock.Unlock()

	if !m.active {
		return
	}
	m.active = false

	log.Println("Ending game due to timeout")

	if m.player1 != nil {
		m.player1.send <- []byte(`{"action":"timeout"}`)
		m.player1.conn.Close()
	}
	if m.player2 != nil {
		m.player2.send <- []byte(`{"action":"timeout"}`)
		m.player2.conn.Close()
	}
}

func (m *Match) ProcessMove(player *Player, row, col int) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if !m.active {
		return
	}

	if player != m.currentTurn {
		player.send <- []byte(`{"action":"not_your_turn"}`)
		return
	}
	if player.opponent != nil {
		moveMsg := Message{
			Action: "move",
			Row:    row,
			Col:    col,
		}
		bytes, _ := json.Marshal(moveMsg)
		player.opponent.send <- bytes

		player.StopMoveTimer()
		// player.opponent.StartMoveTimer()
		m.currentTurn = player.opponent
	}
}
func (m *Match) PlayerLeft(leavingPlayer *Player) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if !m.active {
		return
	}

	m.active = false

	if leavingPlayer.opponent != nil {
		leavingPlayer.opponent.send <- []byte(`{"action":"opponent_left"}`)
		leavingPlayer.opponent.conn.Close()
	}
}
