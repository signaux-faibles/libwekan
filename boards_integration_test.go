package libwekan

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_createBoard(t *testing.T) {
	ass := assert.New(t)
	board := newBoard("la board à toto", "la-board-a-toto", "board")
	err := wekan.InsertBoard(context.Background(), board)
	ass.Nil(err)
}

func Test_getBoardFromID(t *testing.T) {
	id := "kSPsxQZGLKR9tknEt"
	title := "Tableau CRP BFC"
	slug := "tableau-crp-bfc"

	ass := assert.New(t)
	board, err := wekan.GetBoardFromID(context.Background(), id)

	ass.Nil(err)
	ass.NotEmpty(board)
	ass.Equal(title, board.Title)
	ass.Equal(slug, board.Slug)
	ass.Equal(id, board.ID)
}

func Test_getBoardFromSlug(t *testing.T) {
	id := "kSPsxQZGLKR9tknEt"
	title := "Tableau CRP BFC"
	slug := "tableau-crp-bfc"

	ass := assert.New(t)
	board, err := wekan.GetBoardFromSlug(context.Background(), slug)

	ass.Nil(err)
	ass.NotEmpty(board)
	ass.Equal(title, board.Title)
	ass.Equal(slug, board.Slug)
	ass.Equal(id, board.ID)
}

func Test_getBoardFromTitle(t *testing.T) {
	id := "kSPsxQZGLKR9tknEt"
	title := "Tableau CRP BFC"
	slug := "tableau-crp-bfc"

	ass := assert.New(t)
	board, err := wekan.GetBoardFromTitle(context.Background(), title)

	ass.Nil(err)
	ass.NotEmpty(board)
	ass.Equal(title, board.Title)
	ass.Equal(slug, board.Slug)
	ass.Equal(id, board.ID)
}