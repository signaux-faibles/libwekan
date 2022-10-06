package libwekan

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoards_UserIsMember(t *testing.T) {
	ass := assert.New(t)

	userToto := User{
		ID: "toto",
	}
	boardWithUserToto := Board{
		Members: []BoardMember{
			{UserID: "toto"},
		},
	}

	member := boardWithUserToto.UserIsMember(userToto)
	ass.True(member)

	boardWithUserTata := Board{
		Members: []BoardMember{
			{UserID: "tata"},
		},
	}

	member = boardWithUserTata.UserIsMember(userToto)
	ass.False(member)

	boardWithUserTotoAndTata := Board{
		Members: []BoardMember{
			{UserID: "tata"},
			{UserID: "toto"},
		},
	}

	member = boardWithUserTotoAndTata.UserIsMember(userToto)
	ass.True(member)
}

func TestBoards_UserIsActiveMember(t *testing.T) {
	ass := assert.New(t)

	userToto := User{
		ID: "toto",
	}
	boardWithActiveToto := Board{
		Members: []BoardMember{
			{UserID: "toto", IsActive: true},
		},
	}

	isActive := boardWithActiveToto.UserIsActiveMember(userToto)
	ass.True(isActive)

	boardWithInactiveToto := Board{
		Members: []BoardMember{
			{UserID: "toto", IsActive: false},
		},
	}

	isActive = boardWithInactiveToto.UserIsActiveMember(userToto)
	ass.False(isActive)

	boardWithUserTata := Board{
		Members: []BoardMember{
			{UserID: "tata", IsActive: true},
		},
	}

	isActive = boardWithUserTata.UserIsActiveMember(userToto)
	ass.False(isActive)

	boardWithUserTotoAndTata := Board{
		Members: []BoardMember{
			{UserID: "tata", IsActive: true},
			{UserID: "toto", IsActive: false},
		},
	}

	isActive = boardWithUserTotoAndTata.UserIsActiveMember(userToto)
	ass.False(isActive)
}

func TestBoards_GetLabelByName_whenBoardLabelExists(t *testing.T) {
	ass := assert.New(t)
	testId := BoardLabelID(newId6())
	labelName := BoardLabelName("existing label")
	board := Board{
		Labels: []BoardLabel{
			{
				ID:    testId,
				Name:  labelName,
				Color: "orange",
			},
		},
	}
	label := board.GetLabelByName(labelName)
	ass.NotEmpty(label)
}

func TestBoards_GetLabelByName_whenBoardLabelDoesntExists(t *testing.T) {
	ass := assert.New(t)
	board := Board{}
	label := board.GetLabelByName("anotherName")
	ass.Empty(label)
}

func TestBoards_GetLabelByID_whenBoardLabelExists(t *testing.T) {
	ass := assert.New(t)
	testId := BoardLabelID(newId6())
	board := Board{
		Labels: []BoardLabel{
			{
				ID:    testId,
				Name:  "existing label",
				Color: "orange",
			},
		},
	}
	label := board.GetLabelByID(testId)
	ass.NotEmpty(label)
}

func TestBoards_GetLabelByID_whenBoardLabelDoesntExists(t *testing.T) {
	ass := assert.New(t)
	board := Board{}
	label := board.GetLabelByID("anotherID")
	ass.Empty(label)
}

func createTestBoard(t *testing.T, suffix string, swimlanesCount int, listsCount int) (Board, []Swimlane, []List) {
	ctx := context.Background()
	board := BuildBoard(t.Name()+suffix, t.Name()+suffix, "board")
	wekan.InsertBoard(ctx, board)
	var swimlanes []Swimlane
	var lists []List
	for i := 0; i < swimlanesCount; i++ {
		swimlane := BuildSwimlane(board.ID, "swimlane", t.Name()+"swimlane", i)
		swimlanes = append(swimlanes, swimlane)
		wekan.InsertSwimlane(ctx, swimlane)
	}
	for i := 0; i < listsCount; i++ {
		title := fmt.Sprintf("%sList%d", t.Name(), i)
		list := BuildList(board.ID, title, i)
		lists = append(lists, list)
		wekan.InsertList(ctx, list)
	}
	return board, swimlanes, lists
}

// NewBoardLabel retourne un objet BoardLabel
func Test_BuildBoardLabel(t *testing.T) {
	id := t.Name() + "ID"
	name := t.Name() + "Name"
	color := t.Name() + "Color"

	expected := BoardLabel{
		ID:    BoardLabelID(id),
		Name:  BoardLabelName(name),
		Color: color,
	}
	boardLabel := NewBoardLabel(name, expected.Color)

	assert.Equal(t, expected.Name, boardLabel.Name)
	assert.Equal(t, expected.Color, boardLabel.Color)
}
