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
	id := BoardID("kSPsxQZGLKR9tknEt")
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
	board, err := wekan.GetBoardFromSlug(context.Background(), BoardSlug(slug))

	ass.Nil(err)
	ass.NotEmpty(board)
	ass.Equal(title, board.Title)
	ass.Equal(slug, board.Slug)
	ass.Equal(id, string(board.ID))
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
	ass.Equal(id, string(board.ID))
}

func Test_AddMemberToBoard(t *testing.T) {
	ass := assert.New(t)
	board, err := wekan.GetBoardFromSlug(context.Background(), "tableau-crp-bfc")

	ass.Nil(err)

	user := BuildUser("test_add_member_to_board", "tamb", "Test Add Member To Board")
	insertedUser, err := wekan.InsertUser(context.Background(), user)
	ass.Nil(err)

	ass.False(board.UserIsMember(insertedUser))

	boardMember := BoardMember{insertedUser.ID, false, true, false, false, false}
	err = wekan.AddMemberToBoard(context.Background(), board.ID, boardMember)
	ass.Nil(err)

	actualBoard, err := wekan.GetBoardFromSlug(context.Background(), "tableau-crp-bfc")
	ass.Nil(err)

	ass.True(actualBoard.UserIsMember(insertedUser))
}

func Test_EnableBoardMember(t *testing.T) {
	ass := assert.New(t)
	board, err := wekan.GetBoardFromSlug(context.Background(), "tableau-crp-bfc")

	ass.Nil(err)

	user := BuildUser("test_enable_board_member", "tebm", "Test Enable Board Member")
	insertedUser, err := wekan.InsertUser(context.Background(), user)
	ass.Nil(err)
	ass.False(board.UserIsMember(insertedUser))
	notEnabledUser := BuildUser("test_not_enable_board_member", "tnebm", "Test Not Enable Board Member")
	insertedNotEnabledUser, err := wekan.InsertUser(context.Background(), notEnabledUser)
	ass.Nil(err)
	ass.False(board.UserIsMember(insertedNotEnabledUser))

	err = wekan.AddMemberToBoard(context.Background(), board.ID, BoardMember{insertedUser.ID, false, false, false, false, false})
	ass.Nil(err)
	err = wekan.AddMemberToBoard(context.Background(), board.ID, BoardMember{insertedNotEnabledUser.ID, false, false, false, false, false})
	ass.Nil(err)
	insertedMemberBoard, err := wekan.GetBoardFromSlug(context.Background(), "tableau-crp-bfc")
	ass.Nil(err)
	ass.False(insertedMemberBoard.UserIsActiveMember(insertedUser))
	ass.False(insertedMemberBoard.UserIsActiveMember(insertedNotEnabledUser))

	err = wekan.EnableBoardMember(context.Background(), insertedMemberBoard.ID, user.ID)
	ass.Nil(err)
	enabledMemberBoard, err := wekan.GetBoardFromSlug(context.Background(), "tableau-crp-bfc")
	ass.Nil(err)
	ass.True(enabledMemberBoard.UserIsActiveMember(user))
	ass.False(enabledMemberBoard.UserIsActiveMember(notEnabledUser))
}

func Test_DisableBoardMember(t *testing.T) {
	ass := assert.New(t)
	board, err := wekan.GetBoardFromSlug(context.Background(), "tableau-crp-bfc")

	ass.Nil(err)

	user := BuildUser("test_disable_board_member", "tdbm", "Test Disable Board Member")
	enabledUser := BuildUser("test_not_disable_board_member", "tndbm", "Test Not Disable Board Member")

	insertedUser, err := wekan.InsertUser(context.Background(), user)
	ass.Nil(err)
	insertedEnabledUser, err := wekan.InsertUser(context.Background(), enabledUser)
	ass.Nil(err)
	ass.False(board.UserIsMember(insertedUser))
	ass.False(board.UserIsMember(insertedEnabledUser))

	err = wekan.AddMemberToBoard(context.Background(), board.ID, BoardMember{insertedUser.ID, false, true, false, false, false})
	ass.Nil(err)
	err = wekan.AddMemberToBoard(context.Background(), board.ID, BoardMember{insertedEnabledUser.ID, false, true, false, false, false})
	ass.Nil(err)

	insertedMemberBoard, err := wekan.GetBoardFromSlug(context.Background(), "tableau-crp-bfc")
	ass.Nil(err)
	ass.True(insertedMemberBoard.UserIsActiveMember(user))
	ass.True(insertedMemberBoard.UserIsActiveMember(enabledUser))

	err = wekan.DisableBoardMember(context.Background(), insertedMemberBoard.ID, user.ID)
	ass.Nil(err)

	disabledMemberBoard, err := wekan.GetBoardFromSlug(context.Background(), "tableau-crp-bfc")
	ass.Nil(err)

	ass.False(disabledMemberBoard.UserIsActiveMember(user))
	ass.True(disabledMemberBoard.UserIsActiveMember(enabledUser))
}

func Test_EnsureUserIsActiveBoardMember(t *testing.T) {
	ass := assert.New(t)
	board, err := wekan.GetBoardFromSlug(context.Background(), "tableau-crp-bfc")
	ass.Nil(err)

	user := BuildUser("test_ensure_user_is_active_board_member", "teuiabm", "Test Ensure User Is Active Board Member")
	insertedUser, err := wekan.InsertUser(context.Background(), user)
	ass.Nil(err)

	ass.False(board.UserIsMember(insertedUser))
	ass.False(board.UserIsActiveMember(insertedUser))

	err = wekan.EnsureUserIsActiveBoardMember(context.Background(), board.ID, user.ID)
	ass.Nil(err)

	updatedBoard, err := wekan.GetBoardFromSlug(context.Background(), "tableau-crp-bfc")
	ass.Nil(err)
	ass.True(updatedBoard.UserIsMember(insertedUser))
	ass.True(updatedBoard.UserIsActiveMember(insertedUser))
}

func Test_InsertBoardLabel_whenBoardLabelDontExists(t *testing.T) {
	ass := assert.New(t)
	board, err := wekan.GetBoardFromSlug(context.Background(), "tableau-crp-bfc")
	ass.Nil(err)
	boardLabel := BoardLabel{
		ID:    BoardLabelID(newId6()),
		Name:  "test label",
		Color: "orange",
	}
	err = wekan.InsertBoardLabel(context.Background(), board, boardLabel)
	ass.Nil(err)
	updatedBoard, err := wekan.GetBoardFromSlug(context.Background(), "tableau-crp-bfc")
	ass.Nil(err)
	insertedLabel := updatedBoard.GetLabelByID(boardLabel.ID)
	ass.NotEmpty(insertedLabel)
}

func Test_InsertBoardLabel_whenBoardLabelAlreadyExists(t *testing.T) {
	ass := assert.New(t)
	board, err := wekan.GetBoardFromSlug(context.Background(), "tableau-codefi-nord")
	ass.Nil(err)
	boardLabel := BoardLabel{
		ID:    BoardLabelID(newId6()),
		Name:  "test label",
		Color: "orange",
	}
	err = wekan.InsertBoardLabel(context.Background(), board, boardLabel)
	ass.Nil(err)
	updatedBoard, err := wekan.GetBoardFromSlug(context.Background(), "tableau-crp-bfc")
	ass.Nil(err)
	err = wekan.InsertBoardLabel(context.Background(), updatedBoard, boardLabel)
	ass.IsType(BoardLabelAlreadyExistsError{}, err)
}
