//go:build integration

// nolint:errcheck

package libwekan

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCards_InsertCard_withGetCardFromID(t *testing.T) {
	ass := assert.New(t)
	// GIVEN
	boardID := BoardID(t.Name())
	listID := ListID(t.Name())
	swimlaneID := SwimlaneID(t.Name())
	title := t.Name()
	description := t.Name()
	userID := UserID(t.Name())
	card := BuildCard(boardID, listID, swimlaneID, title, description, userID)

	// WHEN
	err := wekan.InsertCard(ctx, card)
	ass.Nil(err)
	actualCard, err := wekan.GetCardFromID(ctx, card.ID)
	ass.Nil(err)

	// THEN
	ass.Equal(card, actualCard)
}

func TestCards_GetCardsFromID_whenCardDoesntExists(t *testing.T) {
	cardID := CardID(t.Name())
	_, err := wekan.GetCardFromID(ctx, cardID)
	assert.IsType(t, CardNotFoundError{}, err)
}

func TestCards_SelectCardsFromUserID(t *testing.T) {
	_, err := wekan.SelectCardsFromUserID(ctx, "")
	assert.IsType(t, NotImplemented{}, err)
}

func TestCards_SelectCardsFromBoardID(t *testing.T) {
	_, err := wekan.SelectCardsFromBoardID(ctx, "")
	assert.IsType(t, NotImplemented{}, err)
}

func TestCards_SelectCardsFromMemberID(t *testing.T) {
	_, err := wekan.SelectCardsFromMemberID(ctx, "")
	assert.IsType(t, NotImplemented{}, err)
}
func TestCards_SelectCardsFromSwimlaneID(t *testing.T) {
	_, err := wekan.SelectCardsFromSwimlaneID(ctx, "")
	assert.IsType(t, NotImplemented{}, err)
}

func TestCards_SelectCardsFromListID(t *testing.T) {
	_, err := wekan.SelectCardsFromListID(ctx, "")
	assert.IsType(t, NotImplemented{}, err)
}
