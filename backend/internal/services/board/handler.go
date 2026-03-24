package board

import "fmt"

func CreateCard(card Card, board Board) (*Card, error) {
	board.mu.Lock()
	defer board.mu.Unlock()

	if _, exists := board.CardList[card.CardID]; exists {
		return nil, fmt.Errorf("card already exists with id")
	}

	board.CardList[card.CardID] = card

	return &card, nil

}

func UpdateCard(card Card, board Board) (*Card, error) {
	board.mu.Lock()
	defer board.mu.Unlock()

	if _, exists := board.CardList[card.CardID]; !exists {
		return nil, fmt.Errorf("card does not exists with id")
	}

	board.CardList[card.CardID] = card

	return &card, nil
}

func DeleteCard(cardID string, board Board) error {
	board.mu.Lock()
	defer board.mu.Unlock()

	if _, exists := board.CardList[cardID]; !exists {
		return fmt.Errorf("card does not exists with id")
	}

	delete(board.CardList, cardID)

	return nil
}

func MoveCard(cardID, column string, board Board) (*Card, error) {
	board.mu.Lock()
	defer board.mu.Unlock()
	card, exists := board.CardList[cardID]
	if !exists {
		return nil, fmt.Errorf("card does not exists with id")
	}

	card.Column = ColumnName(column)

	board.CardList[cardID] = card

	return &card, nil

}
