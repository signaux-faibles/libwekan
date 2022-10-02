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
	board, _ := wekan.GetBoardFromSlug(ctx, "tableau-crp-bfc")
	builtUser := BuildUser(t.Name(), t.Name(), t.Name())
	insertedUser, _ := wekan.InsertUser(ctx, builtUser)
	wekan.EnsureUserIsActiveBoardMember(ctx, board.ID, insertedUser.ID)
	testBoardLabel := BoardLabel{
		ID:    BoardLabelID(newId6()),
		Name:  BoardLabelName(t.Name()),
		Color: "blue",
	}
	wekan.InsertBoardLabel(ctx, board, testBoardLabel)
	updatedBoard, _ := wekan.GetBoardFromSlug(ctx, "tableau-crp-bfc")

	// THEN
	rule := updatedBoard.BuildRule(insertedUser, BoardLabelName(t.Name()))
	ass.NotEmpty(rule)
	err := wekan.InsertRule(ctx, rule)
	ass.Nil(err)
	insertedRule, _ := wekan.SelectRuleFromID(ctx, rule.ID)
	ass.Equal(rule, insertedRule)
}

func TestRules_SelectRuleFromID_whenRuleExists(t *testing.T) {
	// WHEN
	ass := assert.New(t)
	rule := Rule{ID: "absentId"}
	wekan.InsertRule(ctx, rule)

	// THEN
	actualRule, err := wekan.SelectRuleFromID(ctx, rule.ID)
	ass.Nil(err)
	ass.Equal(rule, actualRule)
}

func TestRules_SelectRuleFromID_whenRuleDoesntExist(t *testing.T) {
	// WHEN
	ass := assert.New(t)
	rule := Rule{ID: "absentId"}

	// THEN
	actualRule, err := wekan.SelectRuleFromID(ctx, rule.ID)
	ass.Empty(actualRule)
	ass.Equal(UnknownRuleError{rule}, err)
}

func TestRules_SelectRulesFromBoardID(t *testing.T) {
	// WHEN
	ass := assert.New(t)
	board, _ := wekan.GetBoardFromSlug(ctx, "tableau-codefi-nord")
	builtUser := BuildUser(t.Name(), t.Name(), t.Name())
	insertedUser, _ := wekan.InsertUser(ctx, builtUser)
	wekan.EnsureUserIsActiveBoardMember(ctx, board.ID, insertedUser.ID)
	testBoardLabel := BoardLabel{
		ID:    BoardLabelID(newId6()),
		Name:  BoardLabelName(t.Name()),
		Color: "blue",
	}
	wekan.InsertBoardLabel(ctx, board, testBoardLabel)
	updatedBoard, _ := wekan.GetBoardFromSlug(ctx, "tableau-codefi-nord")
	rule := updatedBoard.BuildRule(insertedUser, BoardLabelName(t.Name()))
	wekan.InsertRule(ctx, rule)

	// THEN
	rules, err := wekan.SelectRulesFromBoardID(ctx, updatedBoard.ID)
	ass.Nil(err)
	selectedRule := sliceSelect(rules, func(r Rule) bool { return r.Action.Username == builtUser.Username })
	ass.Len(selectedRule, 1)
}

func TestRules_RemoveRuleWithID(t *testing.T) {}
