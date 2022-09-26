package libwekan

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_insertRule(t *testing.T) {
	ass := assert.New(t)
	board, err := wekan.GetBoardFromSlug(context.Background(), "tableau-crp-bfc")
	ass.Nil(err)
	builtUser := BuildUser("test_rules", "test_rules", "test_rules")
	insertedUser, err := wekan.InsertUser(context.Background(), builtUser)
	ass.Nil(err)
	err = wekan.EnsureUserIsActiveBoardMember(context.Background(), board.ID, insertedUser.ID)
	ass.Nil(err)

	rule := board.BuildRule(insertedUser, BoardLabelName("toto"))
	ass.Empty(rule)
	err = wekan.InsertRule(context.Background(), rule)
	ass.IsType(NewInsertEmptyRuleError(), err)

}
