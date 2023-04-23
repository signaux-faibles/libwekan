package libwekan

import (
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

func TestBoard_HasLabelName_NotExistingLabel(t *testing.T) {
	board := BuildBoard(t.Name(), t.Name(), "board")
	assert.False(t, board.HasLabelName("notExistingLabel"))
}

func TestBoard_HasLabelName_ExistingLabel(t *testing.T) {
	board := BuildBoard(t.Name(), t.Name(), "board")
	name := BoardLabelName("testLabel")
	boardLabel := NewBoardLabel(string(name), "red")
	board.Labels = append(board.Labels, boardLabel)

	assert.True(t, board.HasLabelName(name))
}
