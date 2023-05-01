package server

import (
	"fmt"
	"log"
	"strconv"
)

type room struct {
	id        string
	client1   *client
	client2   *client
	gameField []playerType
	turn      playerType
}

func (r room) getGameDtoState(player *client) wsGameStateMessage {
	var you *client
	var enemy *client
	if player == r.client1 {
		you = r.client1
		enemy = r.client2
	} else {
		you = r.client2
		enemy = r.client1
	}
	enemyRole := playerNone
	if enemy != nil {
		enemyRole = enemy.role
	}
	return wsGameStateMessage{
		r.id,
		you.role,
		enemyRole,
		r.gameField,
		r.turn,
	}
}

func (r room) sendCurStateMessage() error {
	fmt.Println("room now", r)

	if r.client1 != nil {
		msg1 := r.getGameDtoState(r.client1)
		if err := r.client1.connection.WriteJSON(msg1); err != nil {
			return err
		}
	}

	if r.client2 != nil {
		msg2 := r.getGameDtoState(r.client2)
		if err := r.client2.connection.WriteJSON(msg2); err != nil {
			return err
		}
		return nil
	}

	return nil
}

func (r room) handlerSetRole(player *client, arg string) {
	player.role = playerType(arg[0])
	fmt.Println("room now", arg)
	if err := r.sendCurStateMessage(); err != nil {
		log.Println("write error:", err)
	}
}

func (r room) handlerMakeTurn(arg string) {
	fieldId, err := strconv.Atoi(arg)
	if err != nil {
		log.Println("write error:", err)
	}

	r.gameField[fieldId] = r.turn

	if r.turn == playerX {
		r.turn = playerO
	} else {
		r.turn = playerX
	}

	if err = r.sendCurStateMessage(); err != nil {
		log.Println("write error:", err)
	}
}
