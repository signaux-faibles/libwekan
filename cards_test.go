package libwekan

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

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

func getElement[Element any](elements []Element, fn func(element Element) bool) *Element {
	for _, element := range elements {
		if fn(element) {
			return &element
		}
	}
	return nil
}

func TestCard_AddMember(t *testing.T) {
	card := BuildCard(
		BoardID(t.Name()+"boardID"),
		ListID(t.Name()+"boardID"),
		SwimlaneID(t.Name()+"boardID"),
		(t.Name() + "title"),
		(t.Name() + "description"),
		UserID(t.Name()+"userID"),
	)
	memberID := UserID(t.Name() + "memberID")
	card.AddMember(memberID)
	assert.Len(t, card.Members, 1)
}

func TestCard_AddMember_Duplicate(t *testing.T) {
	card := BuildCard(
		BoardID(t.Name()+"boardID"),
		ListID(t.Name()+"boardID"),
		SwimlaneID(t.Name()+"boardID"),
		(t.Name() + "title"),
		(t.Name() + "description"),
		UserID(t.Name()+"userID"),
	)
	memberID := UserID(t.Name() + "memberID")
	card.AddMember(memberID)
	card.AddMember(memberID)
	assert.Len(t, card.Members, 1)
}
