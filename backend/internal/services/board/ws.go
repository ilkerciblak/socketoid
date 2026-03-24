package board

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"ilkerciblak/socketoid/internal/api/ws"
)

const (
	EventCardCreated = "board.card.created"
	EventCardMoved   = "board.card.moved"
	EventCardDeleted = "board.card.deleted"
	EventCardUpdated = "board.card.updated"
)

type CardCreatedPayload struct {
	CardID string `json:"card_id"`
	Title  string `json:"title"`
	Column string `json:"column"`
}

type CardMovedPayload struct {
	CardID string `json:"card_id"`
	Column string `json:"column"`
}

type CardDeletedPayload struct {
	CardID string `json:"card_id"`
}

type CardUpdatePayload struct {
	CardID string `json:"card_id"`
	Title  string `json:"title"`
	Column string `json:"column"`
}

func createCardHandler(b *Board, h *ws.Hub) ws.HandlerFunc {
	return func(client *ws.Client, payload json.RawMessage) error {
		var card Card
		if err := json.Unmarshal(payload, &card); err != nil {
			return fmt.Errorf(
				"unmarshalling error: %v for\nevent:%s\ndata:%v",
				err.Error(),
				EventCardCreated,
				payload,
			)
		}

		card.CardID = randomID()

		_, err := CreateCard(card, *b)
		if err != nil {
			return err
		}

		h.Broadcast <- payload

		return nil
	}
}

func moveCardHandler(b *Board, h *ws.Hub) ws.HandlerFunc {
	return func(client *ws.Client, payload json.RawMessage) error {
		var card Card
		if err := json.Unmarshal(payload, &card); err != nil {
			return fmt.Errorf(
				"unmarshalling error: %v for\nevent:%s\ndata:%v",
				err.Error(),
				EventCardMoved,
				payload,
			)
		}

		_, err := MoveCard(card.CardID, string(card.Column), *b)
		if err != nil {
			return err
		}
		h.Broadcast <- payload

		return nil
	}
}

func deleteCardHandler(b *Board, h *ws.Hub) ws.HandlerFunc {
	return func(client *ws.Client, payload json.RawMessage) error {
		var req struct {
			CardID string `json:"card_id"`
		}

		if err := json.Unmarshal(payload, &req); err != nil {
			return fmt.Errorf(
				"unmarshalling error: %v for\nevent:%s\ndata:%v",
				err.Error(),
				EventCardDeleted,
				payload,
			)
		}

		err := DeleteCard(req.CardID, *b)
		if err != nil {
			return err
		}
		h.Broadcast <- payload

		return nil
	}
}

func updateCardHandler(b *Board, h *ws.Hub) ws.HandlerFunc {
	return func(client *ws.Client, payload json.RawMessage) error {
		var card Card
		if err := json.Unmarshal(payload, &card); err != nil {
			return fmt.Errorf(
				"unmarshalling error: %v for\nevent:%s\ndata:%v",
				err.Error(),
				EventCardUpdated,
				payload,
			)
		}

		card.CardID = randomID()

		_, err := UpdateCard(card, *b)
		if err != nil {
			return err
		}
		h.Broadcast <- payload

		return nil
	}
}

func randomID() string {
	randomBytes := make([]byte, 32)

	// generated random bytes using crypto/rand
	_, err := rand.Read(randomBytes)
	if err != nil {
		fmt.Printf("[rand.Read]:  %v", err)
	}

	state := base64.URLEncoding.EncodeToString(randomBytes)
	return state
}

func RegisterBoardHandlers(router *ws.Router, hub *ws.Hub) {
	b := InitiateBoard()

	router.Register(EventCardCreated, createCardHandler(b, hub))
	router.Register(EventCardMoved, moveCardHandler(b, hub))
	router.Register(EventCardDeleted, deleteCardHandler(b, hub))
	router.Register(EventCardUpdated, updateCardHandler(b, hub))
}
