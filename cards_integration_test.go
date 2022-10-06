//go:build integration

// nolint:errcheck

package libwekan

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCards_InsertCard_withGetCardFromID(t *testing.T) {
	ass := assert.New(t)
	// GIVEN
	card := createTestCard(t, createTestUser(t, "").ID, nil, nil, nil)
	// WHEN
	actualCard, err := wekan.GetCardFromID(ctx, card.ID)
	ass.Nil(err)

	// THEN
	ass.Equal(card, actualCard)
}

func TestCards_GetCardsFromID_whenCardDoesntExists(t *testing.T) {
	_, err := wekan.GetCardFromID(ctx, CardID(t.Name()+"CatchMeIfYouCan"))
	assert.IsType(t, CardNotFoundError{}, err)
}

func TestCards_SelectCardsFromUserID(t *testing.T) {
	// GIVEN
	card := createTestCard(t, createTestUser(t, "").ID, nil, nil, nil)

	// WHEN
	actualCards, err := wekan.SelectCardsFromUserID(ctx, card.UserID)

	// THEN
	assert.Nil(t, err)
	assert.Equal(t, []Card{card}, actualCards)
}

func TestCards_SelectCardsFromBoardID(t *testing.T) {
	// GIVEN
	card := createTestCard(t, createTestUser(t, "").ID, nil, nil, nil)
	wekan.InsertCard(ctx, card)

	// WHEN
	actualCards, err := wekan.SelectCardsFromBoardID(ctx, card.BoardID)

	// THEN
	assert.Nil(t, err)
	assert.Equal(t, []Card{card}, actualCards)
}

func TestCards_SelectCardsFromMemberID(t *testing.T) {
	ass := assert.New(t)
	// GIVEN
	card := createTestCard(t, createTestUser(t, "Owner").ID, nil, nil, nil)
	member := createTestUser(t, "Member")
	wekan.AddMemberToBoard(ctx, card.BoardID, BoardMember{UserID: member.ID, IsActive: true})
	insertedMember, _ := wekan.GetUserFromID(ctx, member.ID)
	wekan.AddMemberToCard(ctx, card.ID, member.ID)

	// WHEN
	insertedCard, _ := wekan.GetCardFromID(ctx, card.ID)
	actualCards, err := wekan.SelectCardsFromMemberID(ctx, insertedMember.ID)

	// THEN
	ass.Nil(err)
	require.Len(t, actualCards, 1)
	ass.Equal(insertedCard, actualCards[0])
}

func TestCards_SelectCardsFromSwimlaneID(t *testing.T) {
	// GIVEN
	card := createTestCard(t, createTestUser(t, "").ID, nil, nil, nil)
	wekan.InsertCard(ctx, card)

	// WHEN
	actualCards, err := wekan.SelectCardsFromSwimlaneID(ctx, card.SwimlaneID)

	// THEN
	assert.Nil(t, err)
	assert.Equal(t, []Card{card}, actualCards)
}

func TestCards_SelectCardsFromListID(t *testing.T) {
	// GIVEN
	card := createTestCard(t, createTestUser(t, "").ID, nil, nil, nil)
	wekan.InsertCard(ctx, card)

	// WHEN
	actualCards, err := wekan.SelectCardsFromListID(ctx, card.ListID)

	// THEN
	assert.Nil(t, err)
	assert.Equal(t, []Card{card}, actualCards)
}

func TestCards_AddCardMembership_WhenBoardIsTheSame(t *testing.T) {
	ass := assert.New(t)
	// GIVEN
	board, swimlanes, lists := createTestBoard(t, "", 1, 1)
	user := createTestUser(t, "User")
	wekan.AddMemberToBoard(ctx, board.ID, BoardMember{UserID: user.ID, IsActive: true})
	member := createTestUser(t, "Member")
	wekan.AddMemberToBoard(ctx, board.ID, BoardMember{UserID: member.ID, IsActive: true})
	card := createTestCard(t, user.ID, &board.ID, &(swimlanes[0].ID), &(lists[0].ID))

	// WHEN
	err := wekan.AddMemberToCard(ctx, card.ID, member.ID)
	ass.Nil(err)

	// THEN
	actualCard, _ := wekan.GetCardFromID(ctx, card.ID)
	ass.Contains(actualCard.Members, member.ID)
}

func TestCards_AddCardMembership_WhenBoardIsNotTheSame(t *testing.T) {
	ass := assert.New(t)
	// GIVEN
	memberBoard, _, _ := createTestBoard(t, "", 1, 1)
	cardBoard, cardSwimlanes, cardLists := createTestBoard(t, "", 1, 1)

	cardOwner := createTestUser(t, "CardOwner")
	wekan.AddMemberToBoard(ctx, cardBoard.ID, BoardMember{UserID: cardOwner.ID, IsActive: true})
	cardMember := createTestUser(t, "CardMember")
	wekan.AddMemberToBoard(ctx, memberBoard.ID, BoardMember{UserID: cardMember.ID, IsActive: true})
	card := createTestCard(t, cardOwner.ID, &cardBoard.ID, &(cardSwimlanes[0].ID), &(cardLists[0].ID))

	// WHEN
	err := wekan.AddMemberToCard(ctx, card.ID, cardMember.ID)
	ass.IsType(ForbiddenOperationError{}, err)

	// THEN
	actualCard, _ := wekan.GetCardFromID(ctx, card.ID)
	ass.NotContains(actualCard.Members, cardMember.ID)
}

func TestCards_RemoveMemberFromCard_WhenUserIsMember(t *testing.T) {
	ass := assert.New(t)
	// GIVEN
	board, swimlanes, lists := createTestBoard(t, "", 1, 1)
	user := createTestUser(t, "User")
	wekan.AddMemberToBoard(ctx, board.ID, BoardMember{UserID: user.ID, IsActive: true})
	member := createTestUser(t, "Member")
	wekan.AddMemberToBoard(ctx, board.ID, BoardMember{UserID: member.ID, IsActive: true})
	card := createTestCard(t, user.ID, &board.ID, &(swimlanes[0].ID), &(lists[0].ID))
	wekan.AddMemberToCard(ctx, card.ID, member.ID)

	// WHEN
	err := wekan.RemoveMemberFromCard(ctx, card.ID, member.ID)
	ass.Nil(err)

	// Then
	actualCard, _ := wekan.GetCardFromID(ctx, card.ID)
	ass.NotContains(actualCard.Members, card.ID)
}

func TestCards_RemoveMemberFromCard_WhenUserIsNotMember(t *testing.T) {
	ass := assert.New(t)
	// GIVEN
	board, swimlanes, lists := createTestBoard(t, "", 1, 1)
	user := createTestUser(t, "User")
	wekan.AddMemberToBoard(ctx, board.ID, BoardMember{UserID: user.ID, IsActive: true})
	member := createTestUser(t, "Member")
	wekan.AddMemberToBoard(ctx, board.ID, BoardMember{UserID: member.ID, IsActive: true})
	card := createTestCard(t, user.ID, &board.ID, &(swimlanes[0].ID), &(lists[0].ID))

	// WHEN
	err := wekan.RemoveMemberFromCard(ctx, card.ID, member.ID)

	// Then
	actualCard, _ := wekan.GetCardFromID(ctx, card.ID)
	ass.IsType(NothingDoneError{}, err)
	ass.NotContains(actualCard.Members, card.ID)
}

func TestCards_AddLabelToCard_whenLabelIsOnBard(t *testing.T) {
	ass := assert.New(t)
	// GIVEN
	board, swimlanes, lists := createTestBoard(t, "", 1, 1)
	card := createTestCard(t, wekan.adminUserID, &board.ID, &(swimlanes[0].ID), &(lists[0].ID))
	boardLabel := NewBoardLabel(t.Name()+"_BoardLabel", "red")
	wekan.InsertBoardLabel(ctx, board, boardLabel)

	// WHEN
	err := wekan.AddLabelToCard(ctx, card.ID, boardLabel.ID)
	ass.NoError(err)

	// THEN
	actualCard, _ := wekan.GetCardFromID(ctx, card.ID)
	ass.Contains(actualCard.LabelIDs, boardLabel.ID)
}

func TestCards_AddLabelToCard_whenLabelIsNotOnBard(t *testing.T) {
	ass := assert.New(t)
	// GIVEN
	boardWithoutLabel, swimlanes, lists := createTestBoard(t, "_withoutLabel", 1, 1)
	card := createTestCard(t, wekan.adminUserID, &boardWithoutLabel.ID, &(swimlanes[0].ID), &(lists[0].ID))

	boardLabel := NewBoardLabel(t.Name()+"_BoardLabel", "red")
	boardWithLabel, _, _ := createTestBoard(t, "_withoutLabel", 1, 1)
	wekan.InsertBoardLabel(ctx, boardWithLabel, boardLabel)

	// WHEN
	err := wekan.AddLabelToCard(ctx, card.ID, boardLabel.ID)

	// THEN
	ass.IsType(err, BoardLabelNotFoundError{})
	actualCard, _ := wekan.GetCardFromID(ctx, card.ID)
	ass.NotContains(actualCard.LabelIDs, boardLabel.ID)
}

func TestCards_AddLabelToCard_whenLabelAlreadyOnBard(t *testing.T) {
	ass := assert.New(t)
	// GIVEN
	board, swimlanes, lists := createTestBoard(t, "", 1, 1)
	card := createTestCard(t, wekan.adminUserID, &board.ID, &(swimlanes[0].ID), &(lists[0].ID))
	boardLabel := NewBoardLabel(t.Name()+"_BoardLabel", "red")
	wekan.InsertBoardLabel(ctx, board, boardLabel)
	wekan.AddLabelToCard(ctx, card.ID, boardLabel.ID)

	// WHEN
	err := wekan.AddLabelToCard(ctx, card.ID, boardLabel.ID)

	// THEN
	ass.ErrorAs(err, &NothingDoneError{})
	actualCard, _ := wekan.GetCardFromID(ctx, card.ID)
	ass.Contains(actualCard.LabelIDs, boardLabel.ID)
	ass.Len(actualCard.LabelIDs, 1)
}
