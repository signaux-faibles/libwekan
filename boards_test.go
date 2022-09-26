package libwekan

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_UserIsMember(t *testing.T) {
	ass := assert.New(t)

	userToto := User{
		ID: "toto",
	}
	boardWithUserToto := Board{
		Members: []BoardMember{
			{UserId: "toto"},
		},
	}

	member := boardWithUserToto.UserIsMember(userToto)
	ass.True(member)

	boardWithUserTata := Board{
		Members: []BoardMember{
			{UserId: "tata"},
		},
	}

	member = boardWithUserTata.UserIsMember(userToto)
	ass.False(member)

	boardWithUserTotoAndTata := Board{
		Members: []BoardMember{
			{UserId: "tata"},
			{UserId: "toto"},
		},
	}

	member = boardWithUserTotoAndTata.UserIsMember(userToto)
	ass.True(member)
}

func Test_UserIsActiveMember(t *testing.T) {
	ass := assert.New(t)

	userToto := User{
		ID: "toto",
	}
	boardWithActiveToto := Board{
		Members: []BoardMember{
			{UserId: "toto", IsActive: true},
		},
	}

	isActive := boardWithActiveToto.UserIsActiveMember(userToto)
	ass.True(isActive)

	boardWithInactiveToto := Board{
		Members: []BoardMember{
			{UserId: "toto", IsActive: false},
		},
	}

	isActive = boardWithInactiveToto.UserIsActiveMember(userToto)
	ass.False(isActive)

	boardWithUserTata := Board{
		Members: []BoardMember{
			{UserId: "tata", IsActive: true},
		},
	}

	isActive = boardWithUserTata.UserIsActiveMember(userToto)
	ass.False(isActive)

	boardWithUserTotoAndTata := Board{
		Members: []BoardMember{
			{UserId: "tata", IsActive: true},
			{UserId: "toto", IsActive: false},
		},
	}

	isActive = boardWithUserTotoAndTata.UserIsActiveMember(userToto)
	ass.False(isActive)
}

func Test_GetLabelByName_whenBoardLabelExists(t *testing.T) {
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

func Test_GetLabelByName_whenBoardLabelDoesntExists(t *testing.T) {
	ass := assert.New(t)
	board := Board{}
	label := board.GetLabelByName("anotherName")
	ass.Empty(label)
}

func Test_GetLabelByID_whenBoardLabelExists(t *testing.T) {
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

func Test_GetLabelByID_whenBoardLabelDoesntExists(t *testing.T) {
	ass := assert.New(t)
	board := Board{}
	label := board.GetLabelByID("anotherID")
	ass.Empty(label)
}
