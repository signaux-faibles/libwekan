//go:build integration

// nolint:errcheck

package libwekan

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func getElement[Element any](elements []Element, fn func(element Element) bool) *Element {
	for _, element := range elements {
		if fn(element) {
			return &element
		}
	}
	return nil
}

// createTestCard crée un objet de type `Card`, l'insère dans la base de test et le retourne
func createTestCard(t *testing.T, userID UserID, boardID *BoardID, swimlaneID *SwimlaneID, listID *ListID) Card {
	ctx := context.Background()
	var board Board
	var swimlanes []Swimlane
	var lists []List

	if boardID == nil || swimlaneID == nil || listID == nil {
		board, swimlanes, lists = createTestBoard(t, "", 1, 1)
		swimlaneID = &swimlanes[0].ID
		listID = &lists[0].ID
	} else {
		board, _ = wekan.GetBoardFromID(ctx, *boardID)
		swimlanes, _ = wekan.GetSwimlanesFromBoardID(ctx, *boardID)
		lists, _ = wekan.SelectListsFromBoardID(ctx, *boardID)
		swimlane := getElement(swimlanes, func(swimlane Swimlane) bool { return swimlane.ID == *swimlaneID })
		list := getElement(lists, func(list List) bool { return list.ID == *listID })
		if swimlane == nil || list == nil {
			return Card{}
		}
	}

	title := t.Name() + "Title"
	description := t.Name() + "Description"
	card := BuildCard(board.ID, *listID, *swimlaneID, title, description, userID)
	wekan.InsertCard(ctx, card)
	return card
}

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

func TestCards_ArchiveCard(t *testing.T) {
	ass := assert.New(t)

	//GIVEN
	board, swimlanes, lists := createTestBoard(t, "", 1, 1)
	card := createTestCard(t, wekan.adminUserID, &board.ID, &(swimlanes[0].ID), &(lists[0].ID))

	//WHEN
	wekan.ArchiveCard(ctx, card.ID)
	actualCard, err := wekan.GetCardFromID(ctx, card.ID)
	ass.Nil(err)

	//THEN
	ass.Greater(actualCard.ModifiedAt, card.ModifiedAt)
	ass.Greater(actualCard.DateLastActivity, card.DateLastActivity)
	ass.True(actualCard.Archived)
}

func TestCards_ArchiveCard_whenCardAlreadyArchived(t *testing.T) {
	ass := assert.New(t)

	//GIVEN
	board, swimlanes, lists := createTestBoard(t, "", 1, 1)
	card := createTestCard(t, wekan.adminUserID, &board.ID, &(swimlanes[0].ID), &(lists[0].ID))
	wekan.ArchiveCard(ctx, card.ID)
	archivedCard, err := wekan.GetCardFromID(ctx, card.ID)
	ass.Nil(err)

	//WHEN
	err = wekan.ArchiveCard(ctx, archivedCard.ID)

	//THEN
	ass.ErrorAs(err, &NothingDoneError{})
}

func TestCards_ArchiveCard_whenCardDoesNotExists(t *testing.T) {
	ass := assert.New(t)

	//GIVEN
	unknownCardID := CardID(t.Name() + "not an actual ID")

	//WHEN
	err := wekan.ArchiveCard(ctx, unknownCardID)

	//THEN
	ass.ErrorAs(err, &CardNotFoundError{})
}

func TestCards_UnarchiveCard(t *testing.T) {
	ass := assert.New(t)

	//GIVEN
	board, swimlanes, lists := createTestBoard(t, "", 1, 1)
	card := createTestCard(t, wekan.adminUserID, &board.ID, &(swimlanes[0].ID), &(lists[0].ID))
	archivedCard, err := wekan.GetCardFromID(ctx, card.ID)
	ass.Nil(err)

	//WHEN
	wekan.UnarchiveCard(ctx, archivedCard.ID)
	actualCard, err := wekan.GetCardFromID(ctx, card.ID)
	ass.Nil(err)

	//THEN
	ass.Greater(actualCard.ModifiedAt, card.ModifiedAt)
	ass.Greater(actualCard.DateLastActivity, card.DateLastActivity)
	ass.False(actualCard.Archived)
}

func TestCards_UnarchiveCard_whenCardNotArchived(t *testing.T) {
	ass := assert.New(t)

	//GIVEN
	board, swimlanes, lists := createTestBoard(t, "", 1, 1)
	card := createTestCard(t, wekan.adminUserID, &board.ID, &(swimlanes[0].ID), &(lists[0].ID))

	//WHEN
	err := wekan.UnarchiveCard(ctx, card.ID)

	//THEN
	ass.ErrorAs(err, &NothingDoneError{})
}

func TestCards_UnarchiveCard_whenCardDoesNotExists(t *testing.T) {
	ass := assert.New(t)

	//GIVEN
	unknownCardID := CardID(t.Name() + "not an actual ID")

	//WHEN
	err := wekan.UnarchiveCard(ctx, unknownCardID)

	//THEN
	ass.ErrorAs(err, &CardNotFoundError{})
}

func TestCards_UpdateCardDescription(t *testing.T) {
	ass := assert.New(t)

	//GIVEN
	board, swimlanes, lists := createTestBoard(t, "", 1, 1)
	card := createTestCard(t, wekan.adminUserID, &board.ID, &(swimlanes[0].ID), &(lists[0].ID))
	description := "description originale"
	wekan.UpdateCardDescription(ctx, card.ID, description)

	// WHEN
	newDescription := "nouvelle description"
	err := wekan.UpdateCardDescription(ctx, card.ID, newDescription)
	ass.Nil(err)

	//THEN
	actualCard, err := wekan.GetCardFromID(ctx, card.ID)
	ass.Nil(err)
	ass.Equal(newDescription, actualCard.Description)
	ass.Greater(actualCard.ModifiedAt, card.ModifiedAt)
	ass.Greater(actualCard.DateLastActivity, card.DateLastActivity)
}

func TestCards_UpdateCardDescription_WhenCardDoesntExists(t *testing.T) {
	ass := assert.New(t)

	//GIVEN
	unknownCardID := CardID(t.Name() + "not an actual ID")

	// WHEN
	err := wekan.UpdateCardDescription(ctx, unknownCardID, "")

	//THEN
	ass.ErrorAs(err, &CardNotFoundError{})
}
