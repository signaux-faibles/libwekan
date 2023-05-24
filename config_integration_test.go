//go:build integration

package libwekan

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfig_MatchingBoardIsPresent(t *testing.T) {
	ass := assert.New(t)

	slugWhichMatch := "tableau-crp-" + t.Name()
	boardID := BoardID(newId())
	err := wekan.InsertBoard(ctx, Board{
		ID:   boardID,
		Slug: BoardSlug(slugWhichMatch),
	})
	ass.NoError(err)

	config, err := wekan.SelectConfig(ctx)
	ass.NoError(err)

	_, ok := config.Boards[boardID]
	ass.True(ok)
}

func TestConfig_NotMatchingBoardIsAbsent(t *testing.T) {
	ass := assert.New(t)

	slugDontMatch := "tableau-pascrp-" + t.Name()
	boardID := BoardID(newId())
	err := wekan.InsertBoard(ctx, Board{
		ID:   boardID,
		Slug: BoardSlug(slugDontMatch),
	})
	ass.NoError(err)

	config, err := wekan.SelectConfig(ctx)
	ass.NoError(err)

	_, ok := config.Boards[boardID]
	ass.False(ok)
}
