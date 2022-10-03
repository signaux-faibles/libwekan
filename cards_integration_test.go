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
	card := buildTestCard(t)

	// WHEN
	err := wekan.InsertCard(ctx, card)
	ass.Nil(err)
	actualCard, err := wekan.GetCardFromID(ctx, card.ID)
	ass.Nil(err)

	// THEN
	ass.Equal(card, actualCard)
}

func TestCards_GetCardsFromID_whenCardDoesntExists(t *testing.T) {
	card := buildTestCard(t)
	_, err := wekan.GetCardFromID(ctx, card.ID)
	assert.IsType(t, CardNotFoundError{}, err)
}

func TestCards_SelectCardsFromUserID(t *testing.T) {
	// GIVEN
	card := buildTestCard(t)
	wekan.InsertCard(ctx, card)

	// WHEN
	actualCards, err := wekan.SelectCardsFromUserID(ctx, card.UserID)

	// THEN
	assert.Nil(t, err)
	assert.Equal(t, []Card{card}, actualCards)
}

func TestCards_SelectCardsFromBoardID(t *testing.T) {
	// GIVEN
	card := buildTestCard(t)
	wekan.InsertCard(ctx, card)

	// WHEN
	actualCards, err := wekan.SelectCardsFromBoardID(ctx, card.BoardID)

	// THEN
	assert.Nil(t, err)
	assert.Equal(t, []Card{card}, actualCards)
}

func TestCards_SelectCardsFromMemberID(t *testing.T) {
	// GIVEN
	card := buildTestCard(t)
	memberID := UserID(t.Name() + "MemberID")
	card.Members = []UserID{memberID}
	wekan.InsertCard(ctx, card)

	// WHEN
	actualCards, err := wekan.SelectCardsFromMemberID(ctx, memberID)

	// THEN
	assert.Nil(t, err)
	assert.Equal(t, []Card{card}, actualCards)
}

func TestCards_SelectCardsFromSwimlaneID(t *testing.T) {
	// GIVEN
	card := buildTestCard(t)
	wekan.InsertCard(ctx, card)

	// WHEN
	actualCards, err := wekan.SelectCardsFromSwimlaneID(ctx, card.SwimlaneID)

	// THEN
	assert.Nil(t, err)
	assert.Equal(t, []Card{card}, actualCards)
}

func TestCards_SelectCardsFromListID(t *testing.T) {
	// GIVEN
	card := buildTestCard(t)
	wekan.InsertCard(ctx, card)

	// WHEN
	actualCards, err := wekan.SelectCardsFromListID(ctx, card.ListID)

	// THEN
	assert.Nil(t, err)
	assert.Equal(t, []Card{card}, actualCards)
}
