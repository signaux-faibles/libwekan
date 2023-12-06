//go:build integration

// nolint:errcheck
package libwekan

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLists_InsertList(t *testing.T) {
	// GIVEN
	ass := assert.New(t)
	board, _, _ := createTestBoard(t, "", 0, 0)
	newList := BuildList(board.ID, t.Name()+"_title", 0)

	// WHEN
	err := wekan.InsertList(ctx, newList)
	ass.NoError(err)

	// THEN
	actualList, err := wekan.GetListFromID(ctx, newList.ID)
	ass.NoError(err)
	ass.Equal(newList.ID, actualList.ID)
	ass.Equal(newList.Title, actualList.Title)
	ass.Equal(newList.Sort, actualList.Sort)
	ass.Equal(newList.Archived, actualList.Archived)
	ass.NotNil(actualList.CreatedAt)
	ass.NotNil(actualList.ModifiedAt)
}

func TestLists_InsertList_WhenBoardDoesntExists(t *testing.T) {
	// GIVEN
	ass := assert.New(t)
	list := BuildList(BoardID(t.Name()+"_notaboardid"), t.Name(), 0)

	// WHEN
	err := wekan.InsertList(ctx, list)

	// THEN
	ass.ErrorAs(err, &BoardNotFoundError{})
}

func TestLists_SelectListsFromBoardID(t *testing.T) {
	// GIVEN
	ass := assert.New(t)
	board, _, _ := createTestBoard(t, "", 0, 4)

	// WHEN
	lists, err := wekan.SelectListsFromBoardID(ctx, board.ID)
	ass.NoError(err)

	// THEN
	ass.Len(lists, 4)
}

func TestLists_SelectListsFromBoardID_whenBoardDoesntExists(t *testing.T) {
	// GIVEN
	ass := assert.New(t)

	// WHEN
	lists, err := wekan.SelectListsFromBoardID(ctx, BoardID(t.Name()+"_notaboardID"))

	// THEN
	ass.ErrorAs(err, &BoardNotFoundError{})
	ass.Len(lists, 0)
}
