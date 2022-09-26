package libwekan

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_InsertRule_whenLabelDoesntExist(t *testing.T) {
	ass := assert.New(t)
	board, err := wekan.GetBoardFromSlug(context.Background(), "tableau-crp-bfc")
	ass.Nil(err)
	builtUser := BuildUser(t.Name(), t.Name(), t.Name())
	insertedUser, err := wekan.InsertUser(context.Background(), builtUser)
	ass.Nil(err)
	err = wekan.EnsureUserIsActiveBoardMember(context.Background(), board.ID, insertedUser.ID)
	ass.Nil(err)
	updatedBoard, err := wekan.GetBoardFromSlug(context.Background(), "tableau-crp-bfc")
	ass.Nil(err)
	rule := updatedBoard.BuildRule(insertedUser, "toto")
	ass.Empty(rule)
	err = wekan.InsertRule(context.Background(), rule)
	ass.IsType(InsertEmptyRuleError{}, err)
}

func TestWekan_InsertRule_whenUserIsNotMember(t *testing.T) {
	ass := assert.New(t)
	board, err := wekan.GetBoardFromSlug(context.Background(), "tableau-crp-bfc")
	ass.Nil(err)
	builtUser := BuildUser(t.Name(), t.Name(), t.Name())
	insertedUser, err := wekan.InsertUser(context.Background(), builtUser)
	ass.Nil(err)
	ass.False(board.UserIsMember(insertedUser))
	rule := board.BuildRule(insertedUser, "toto")
	ass.Empty(rule)
	err = wekan.InsertRule(context.Background(), rule)
	ass.IsType(InsertEmptyRuleError{}, err)
}

func TestWekan_InsertRule_whenEverythingsFine(t *testing.T) {
	ass := assert.New(t)
	board, err := wekan.GetBoardFromSlug(context.Background(), "tableau-crp-bfc")
	ass.Nil(err)
	builtUser := BuildUser(t.Name(), t.Name(), t.Name())
	insertedUser, err := wekan.InsertUser(context.Background(), builtUser)
	ass.Nil(err)
	err = wekan.EnsureUserIsActiveBoardMember(context.Background(), board.ID, insertedUser.ID)
	ass.Nil(err)
	testBoardLabel := BoardLabel{
		ID:    BoardLabelID(newId6()),
		Name:  BoardLabelName(t.Name()),
		Color: "blue",
	}
	err = wekan.InsertBoardLabel(context.Background(), board, testBoardLabel)
	ass.Nil(err)
	updatedBoard, err := wekan.GetBoardFromSlug(context.Background(), "tableau-crp-bfc")
	rule := updatedBoard.BuildRule(insertedUser, BoardLabelName(t.Name()))
	ass.NotEmpty(rule)

	err = wekan.InsertRule(context.Background(), rule)
	ass.Nil(err)
}
