//go:build integration
// +build integration

// nolint:errcheck
package libwekan

import (
	"context"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
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
	title := BoardTitle("Tableau CRP BFC")
	slug := BoardSlug("tableau-crp-bfc")

	ass := assert.New(t)
	board, err := wekan.GetBoardFromID(context.Background(), id)

	ass.Nil(err)
	ass.NotEmpty(board)
	ass.Equal(title, board.Title)
	ass.Equal(slug, board.Slug)
	ass.Equal(id, board.ID)
}

func Test_getBoardFromSlug(t *testing.T) {
	id := BoardID("kSPsxQZGLKR9tknEt")
	title := BoardTitle("Tableau CRP BFC")
	slug := BoardSlug("tableau-crp-bfc")

	ass := assert.New(t)
	board, err := wekan.GetBoardFromSlug(context.Background(), slug)

	ass.Nil(err)
	ass.NotEmpty(board)
	ass.Equal(title, board.Title)
	ass.Equal(slug, board.Slug)
	ass.Equal(id, board.ID)
}

func Test_getBoardFromTitle(t *testing.T) {
	id := BoardID("kSPsxQZGLKR9tknEt")
	title := BoardTitle("Tableau CRP BFC")
	slug := BoardSlug("tableau-crp-bfc")

	ass := assert.New(t)
	board, err := wekan.GetBoardFromTitle(context.Background(), string(title))

	ass.Nil(err)
	ass.NotEmpty(board)
	ass.Equal(title, board.Title)
	ass.Equal(slug, board.Slug)
	ass.Equal(id, board.ID)
}

func Test_AddMemberToBoard(t *testing.T) {
	ass := assert.New(t)
	board, err := wekan.GetBoardFromSlug(context.Background(), "tableau-crp-bfc")

	ass.Nil(err)

	user := BuildUser("test_add_member_to_board", "tamb", "Test Add Member To Board")
	insertedUser, _ := wekan.InsertUser(context.Background(), user)
	ass.False(board.UserIsMember(insertedUser))

	boardMember := BoardMember{insertedUser.ID, false, true, false, false, false}
	err = wekan.AddMemberToBoard(context.Background(), board.ID, boardMember)
	ass.Nil(err)

	actualBoard, _ := wekan.GetBoardFromSlug(context.Background(), "tableau-crp-bfc")

	ass.True(actualBoard.UserIsMember(insertedUser))

	activities, _ := wekan.selectActivitiesFromQuery(context.Background(), bson.M{"memberId": insertedUser.ID, "activityType": "addBoardMember"})
	ass.Len(activities, 2)
}

func Test_EnableBoardMember(t *testing.T) {
	ass := assert.New(t)
	// GIVEN
	board, err := wekan.GetBoardFromSlug(context.Background(), "tableau-crp-bfc")

	ass.Nil(err)

	user := BuildUser("test_enable_board_member", "tebm", "Test Enable Board Member")
	insertedUser, _ := wekan.InsertUser(context.Background(), user)
	ass.False(board.UserIsMember(insertedUser))

	notEnabledUser := BuildUser("test_not_enable_board_member", "tnebm", "Test Not Enable Board Member")
	insertedNotEnabledUser, _ := wekan.InsertUser(context.Background(), notEnabledUser)
	ass.False(board.UserIsMember(insertedNotEnabledUser))

	addMemberToBoard(board.ID, BoardMember{insertedUser.ID, false, false, false, false, false}, t)
	addMemberToBoard(board.ID, BoardMember{insertedNotEnabledUser.ID, false, false, false, false, false}, t)
	insertedMemberBoard, _ := wekan.GetBoardFromSlug(context.Background(), "tableau-crp-bfc")
	ass.False(insertedMemberBoard.UserIsActiveMember(insertedUser))
	ass.False(insertedMemberBoard.UserIsActiveMember(insertedNotEnabledUser))

	// WHEN
	err = wekan.EnableBoardMember(context.Background(), insertedMemberBoard.ID, user.ID)
	ass.Nil(err)

	// THEN
	enabledMemberBoard, _ := wekan.GetBoardFromSlug(context.Background(), "tableau-crp-bfc")

	ass.True(enabledMemberBoard.UserIsActiveMember(user))
	ass.False(enabledMemberBoard.UserIsActiveMember(notEnabledUser))

	// on vérifie que l'activité a été créée
	expected := newActivityAddBoardMember(wekan.adminUserID, user.ID, board.ID)
	foundActivities, _ := wekan.selectActivitiesFromQuery(context.Background(), bson.M{"boardId": expected.BoardID, "memberId": expected.MemberID})
	require.NotEmpty(t, foundActivities)
	ass.Len(foundActivities, 1)

	actual := foundActivities[0]

	ass.Condition(activityCompareFunc(expected, actual))
}

func Test_EnableBoardMember_when_user_is_not_on_board(t *testing.T) {
	ass := assert.New(t)
	// GIVEN
	user := BuildUser("an_absent_user", "aau", "An absent user")
	insertedUser, _ := wekan.InsertUser(context.Background(), user)
	ass.NotNil(insertedUser)
	board, _ := wekan.GetBoardFromSlug(context.Background(), "tableau-crp-bfc")
	ass.False(board.UserIsMember(user))

	// WHEN
	err := wekan.EnableBoardMember(context.Background(), board.ID, user.ID)
	ass.Nil(err)

	shouldNotUpdatedBoard, _ := wekan.GetBoardFromSlug(context.Background(), "tableau-crp-bfc")
	ass.False(shouldNotUpdatedBoard.UserIsMember(user))
	expected := newActivityAddBoardMember(wekan.adminUserID, user.ID, board.ID)
	activities, err := wekan.selectActivitiesFromQuery(context.Background(), bson.M{"memberId": expected.MemberID, "boardId": expected.BoardID, "activityType": expected.ActivityType})
	ass.Empty(activities, 0)
	ass.Empty(nil)
}

func Test_DisableBoardMember(t *testing.T) {
	ass := assert.New(t)

	// GIVEN
	board, err := wekan.GetBoardFromSlug(context.Background(), "tableau-crp-bfc")

	ass.Nil(err)

	user := BuildUser("test_disable_board_member", "tdbm", "Test Disable Board Member")
	enabledUser := BuildUser("test_not_disable_board_member", "tndbm", "Test Not Disable Board Member")

	insertedUser, _ := wekan.InsertUser(context.Background(), user)
	insertedEnabledUser, _ := wekan.InsertUser(context.Background(), enabledUser)
	ass.False(board.UserIsMember(insertedUser))
	ass.False(board.UserIsMember(insertedEnabledUser))

	addMemberToBoard(board.ID, BoardMember{insertedUser.ID, false, true, false, false, false}, t)
	addMemberToBoard(board.ID, BoardMember{insertedEnabledUser.ID, false, true, false, false, false}, t)

	insertedMemberBoard, _ := wekan.GetBoardFromSlug(context.Background(), "tableau-crp-bfc")
	ass.True(insertedMemberBoard.UserIsActiveMember(user))
	ass.True(insertedMemberBoard.UserIsActiveMember(enabledUser))

	// WHEN
	err = wekan.DisableBoardMember(context.Background(), insertedMemberBoard.ID, user.ID)
	ass.Nil(err)

	// THEN
	disabledMemberBoard, err := wekan.GetBoardFromSlug(context.Background(), "tableau-crp-bfc")
	ass.Nil(err)

	ass.False(disabledMemberBoard.UserIsActiveMember(user))
	ass.True(disabledMemberBoard.UserIsActiveMember(enabledUser))

	// on vérifie que l'activité correspondante a été créée
	expected := newActivityRemoveBoardMember(wekan.adminUserID, user.ID, board.ID)
	foundActivities, _ := wekan.selectActivitiesFromQuery(context.Background(), bson.M{"boardId": expected.BoardID, "memberId": expected.MemberID, "activityType": "removeBoardMember"})
	require.NotEmpty(t, foundActivities)
	ass.Len(foundActivities, 1)

	actual := foundActivities[0]
	ass.Condition(activityCompareFunc(expected, actual))
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

func activityCompareFunc(a1 Activity, a2 Activity) func() bool {
	return func() bool {
		return a1.ActivityType == a2.ActivityType &&
			a1.ActivityTypeID == a2.ActivityTypeID &&
			a1.BoardID == a2.BoardID &&
			a1.BoardLabelID == a2.BoardLabelID &&
			a1.CardID == a2.CardID &&
			a1.CommentID == a2.CommentID &&
			a1.ListID == a2.ListID &&
			a1.MemberID == a2.MemberID &&
			a1.SwimlaneID == a2.SwimlaneID &&
			a1.Type == a2.Type &&
			a1.Username == a2.Username &&
			a1.UserID == a2.UserID
	}
}

func TestBoard_SelectBoardsFromMemberID(t *testing.T) {
	ass := assert.New(t)

	boardBFC, _ := wekan.GetBoardFromSlug(context.Background(), "tableau-crp-bfc")
	boardNORD, _ := wekan.GetBoardFromSlug(context.Background(), "tableau-codefi-nord")

	usernameBFC := t.Name() + "bfc"
	usernameNORD := t.Name() + "nord"
	usernameBOTH := t.Name() + "both"

	userBFC := BuildUser(usernameBFC, usernameBFC, usernameBFC)
	userNORD := BuildUser(usernameNORD, usernameNORD, usernameNORD)
	userBOTH := BuildUser(usernameBOTH, usernameBOTH, usernameBOTH)

	insertedUserBFC, _ := wekan.InsertUser(context.Background(), userBFC)
	insertedUserNORD, _ := wekan.InsertUser(context.Background(), userNORD)
	insertedUserBOTH, _ := wekan.InsertUser(context.Background(), userBOTH)

	addMemberToBoard(boardBFC.ID, BoardMember{UserID: insertedUserBFC.ID, IsActive: true}, t)
	addMemberToBoard(boardNORD.ID, BoardMember{UserID: insertedUserNORD.ID, IsActive: true}, t)
	addMemberToBoard(boardBFC.ID, BoardMember{UserID: insertedUserBOTH.ID, IsActive: true}, t)
	addMemberToBoard(boardNORD.ID, BoardMember{UserID: insertedUserBOTH.ID, IsActive: true}, t)

	boardsUserBFC, err := wekan.SelectBoardsFromMemberID(context.Background(), userBFC.ID)
	ass.Nil(err)
	ass.Len(boardsUserBFC, 2)

	boardsUserNORD, err := wekan.SelectBoardsFromMemberID(context.Background(), userNORD.ID)
	ass.Len(boardsUserNORD, 2)
	ass.Nil(err)

	boardsUserBOTH, err := wekan.SelectBoardsFromMemberID(context.Background(), userBOTH.ID)
	ass.Len(boardsUserBOTH, 3)
	ass.Nil(err)

	boardsUserBIDON, err := wekan.SelectBoardsFromMemberID(context.Background(), "idBidon")
	ass.Len(boardsUserBIDON, 0)
	ass.Nil(err)
}

func addMemberToBoard(boardID BoardID, boardMember BoardMember, t *testing.T) {
	req := require.New(t)
	err := wekan.AddMemberToBoard(context.Background(), boardID, boardMember)
	req.Nilf(err, "il faut corriger la méthode AddMemberToBoard qui est testée ailleurs")
}
