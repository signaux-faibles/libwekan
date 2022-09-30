//go:build integration

// nolint:errcheck
package libwekan

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRules_InsertRule_whenLabelDoesntExist(t *testing.T) {
	ass := assert.New(t)
	board, err := wekan.GetBoardFromSlug(ctx, "tableau-crp-bfc")
	ass.Nil(err)
	builtUser := BuildUser(t.Name(), t.Name(), t.Name())
	insertedUser, err := wekan.InsertUser(ctx, builtUser)
	ass.Nil(err)
	err = wekan.EnsureUserIsActiveBoardMember(ctx, board.ID, insertedUser.ID)
	ass.Nil(err)
	updatedBoard, err := wekan.GetBoardFromSlug(ctx, "tableau-crp-bfc")
	ass.Nil(err)
	rule := updatedBoard.BuildRule(insertedUser, "toto")
	ass.Empty(rule)
	err = wekan.InsertRule(ctx, rule)
	ass.IsType(InsertEmptyRuleError{}, err)
}

func TestRulesInsertRule_whenUserIsNotMember(t *testing.T) {
	ass := assert.New(t)
	board, err := wekan.GetBoardFromSlug(ctx, "tableau-crp-bfc")
	ass.Nil(err)
	builtUser := BuildUser(t.Name(), t.Name(), t.Name())
	insertedUser, err := wekan.InsertUser(ctx, builtUser)
	ass.Nil(err)
	ass.False(board.UserIsMember(insertedUser))
	rule := board.BuildRule(insertedUser, "toto")
	ass.Empty(rule)
	err = wekan.InsertRule(ctx, rule)
	ass.IsType(InsertEmptyRuleError{}, err)
}

func TestRules_InsertRule_whenEverythingsFine(t *testing.T) {
	ass := assert.New(t)
	board, err := wekan.GetBoardFromSlug(ctx, "tableau-crp-bfc")
	ass.Nil(err)
	builtUser := BuildUser(t.Name(), t.Name(), t.Name())
	insertedUser, err := wekan.InsertUser(ctx, builtUser)
	ass.Nil(err)
	err = wekan.EnsureUserIsActiveBoardMember(ctx, board.ID, insertedUser.ID)
	ass.Nil(err)
	testBoardLabel := BoardLabel{
		ID:    BoardLabelID(newId6()),
		Name:  BoardLabelName(t.Name()),
		Color: "blue",
	}
	err = wekan.InsertBoardLabel(ctx, board, testBoardLabel)
	ass.Nil(err)
	updatedBoard, err := wekan.GetBoardFromSlug(ctx, "tableau-crp-bfc")
	ass.Nil(err)
	rule := updatedBoard.BuildRule(insertedUser, BoardLabelName(t.Name()))
	ass.NotEmpty(rule)
	err = wekan.InsertRule(ctx, rule)
	ass.Nil(err)
}
