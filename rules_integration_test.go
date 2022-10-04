//go:build integration

// nolint:errcheck
package libwekan

import (
	"github.com/stretchr/testify/require"
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
	// GIVEN
	tableauSansEtiquette, _ := wekan.GetBoardFromSlug(ctx, "tableau-crp-bfc")
	userDeTest, _ := wekan.InsertUser(ctx, BuildUser(t.Name(), t.Name(), t.Name()))
	wekan.EnsureUserIsActiveBoardMember(ctx, tableauSansEtiquette.ID, userDeTest.ID)

	boardLabelId := BoardLabelID(newId6())
	testBoardLabel := BoardLabel{
		ID:    boardLabelId,
		Name:  BoardLabelName(t.Name()),
		Color: "blue",
	}
	wekan.InsertBoardLabel(ctx, tableauSansEtiquette, testBoardLabel)
	tableauAvecEtiquette, _ := wekan.GetBoardFromSlug(ctx, "tableau-crp-bfc")

	regleAAjouter := tableauAvecEtiquette.BuildRule(userDeTest, BoardLabelName(t.Name()))
	ass.NotEmpty(regleAAjouter)

	// WHEN
	err := wekan.InsertRule(ctx, regleAAjouter)
	ass.Nil(err)

	// THEN
	insertedRule, err := wekan.SelectRuleFromID(ctx, regleAAjouter.ID)
	ass.Nil(err)
	ass.NotNil(insertedRule)
	ass.NotEmpty(insertedRule.Action)
	ass.NotEmpty(insertedRule.Trigger)
	ass.Equal(userDeTest.Username, insertedRule.Action.Username)
	ass.Equal(boardLabelId, insertedRule.Trigger.LabelID)
	ass.Equal(tableauAvecEtiquette.ID, insertedRule.Trigger.BoardID)
}

//func TestRules_SelectRuleFromID_whenRuleExists(t *testing.T) {
//	// WHEN
//	ass := assert.New(t)
//	rule := Rule{ID: "absentId"}
//	wekan.InsertRule(ctx, rule)
//
//	// THEN
//	actualRule, err := wekan.SelectRuleFromID(ctx, rule.ID)
//	ass.Nil(err)
//	ass.Equal(rule, actualRule)
//}

func TestRules_SelectRuleFromID_whenRuleDoesntExist(t *testing.T) {
	// WHEN
	ass := assert.New(t)
	rule := Rule{ID: "absentId"}

	// THEN
	actualRule, err := wekan.SelectRuleFromID(ctx, rule.ID)
	ass.Empty(actualRule)
	ass.Equal(RuleNotFoundError{rule.ID}, err)
}

func TestRules_SelectRulesFromBoardID(t *testing.T) {
	// GIVEN
	ass := assert.New(t)
	tableauCodefiNordInitial, _ := wekan.GetBoardFromSlug(ctx, "tableau-codefi-nord")
	testUser, _ := wekan.InsertUser(ctx, BuildUser(t.Name(), t.Name(), t.Name()))
	wekan.EnsureUserIsActiveBoardMember(ctx, tableauCodefiNordInitial.ID, testUser.ID)

	// creation d'une étiquette
	boardLabelID := BoardLabelID(newId6())
	testBoardLabel := BoardLabel{
		ID:    boardLabelID,
		Name:  BoardLabelName(t.Name()),
		Color: "blue",
	}
	wekan.InsertBoardLabel(ctx, tableauCodefiNordInitial, testBoardLabel)
	tableauCodefiNoardAvecEtiquette, _ := wekan.GetBoardFromSlug(ctx, "tableau-codefi-nord")

	// creation d'une regle pour cette étiquette et ce tableau
	rule := tableauCodefiNoardAvecEtiquette.BuildRule(testUser, BoardLabelName(t.Name()))
	wekan.InsertRule(ctx, rule)

	// WHEN
	rules, err := wekan.SelectRulesFromBoardID(ctx, tableauCodefiNoardAvecEtiquette.ID)
	require.Nil(t, err)

	// THEN
	ass.Len(rules, 1)
	actual := rules[0]
	ass.Equal(tableauCodefiNoardAvecEtiquette.ID, actual.BoardID)
	ass.NotNil(actual.Action)
	ass.Equal(testUser.Username, actual.Action.Username)
	ass.NotNil(actual.Trigger)
	ass.Equal(tableauCodefiNoardAvecEtiquette.ID, actual.Trigger.BoardID)
	ass.Equal(boardLabelID, actual.Trigger.LabelID)
	//selectedRule := sliceSelect(rules, func(r Rule) bool { return r.Action.Username == BuildUser(t.Name(), t.Name(), t.Name()).Username })
	//ass.Len(selectedRule, 1)
}

func TestRules_RemoveRuleWithID(t *testing.T) {}
