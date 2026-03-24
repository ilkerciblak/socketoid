// Package board
package board

import (
	"crypto/rand"
	"encoding/base64"
	"sync"
)

type Board struct {
	mu       *sync.RWMutex
	CardList map[string]Card `json:"card_list"`
}

func InitiateBoard() *Board {
	return &Board{
		mu:       &sync.RWMutex{},
		CardList: map[string]Card{},
	}
}



type Card struct {
	CardID string     `json:"card_id"`
	Column ColumnName `json:"column"`
	Title  string     `json:"title"`
}

func NewCard(column, title string) *Card {
	randomBytes := make([]byte, 32)

	// generated random bytes using crypto/rand
	_, _ = rand.Read(randomBytes)

	id := base64.URLEncoding.EncodeToString(randomBytes)

	return &Card{
		CardID: id,
		Column: ColumnName(column),
		Title:  title,
	}

}

type ColumnName string

const (
	Todo       ColumnName = "to_do"
	InProgress ColumnName = "in_progress"
	Done       ColumnName = "done"
)
