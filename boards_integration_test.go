//go:build integration

// nolint:errcheck
package libwekan

import (
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoards_createBoard(t *testing.T) {
	ass := assert.New(t)
	board := newBoard("la board à toto", "la-board-a-toto", "board")
	err := wekan.InsertBoard(ctx, board)
	ass.Nil(err)
}

func TestBoards_getBoardFromID(t *testing.T) {
	id := BoardID("kSPsxQZGLKR9tknEt")
	title := BoardTitle("Tableau CRP BFC")
	slug := BoardSlug("tableau-crp-bfc")

	ass := assert.New(t)
	board, err := wekan.GetBoardFromID(ctx, id)

	ass.Nil(err)
	ass.NotEmpty(board)
	ass.Equal(title, board.Title)
	ass.Equal(slug, board.Slug)
	ass.Equal(id, board.ID)
}

func TestBoards_getBoardFromSlug(t *testing.T) {
	id := BoardID("kSPsxQZGLKR9tknEt")
	title := BoardTitle("Tableau CRP BFC")
	slug := BoardSlug("tableau-crp-bfc")

	ass := assert.New(t)
	board, err := wekan.GetBoardFromSlug(ctx, slug)

	ass.Nil(err)
	ass.NotEmpty(board)
	ass.Equal(title, board.Title)
	ass.Equal(slug, board.Slug)
	ass.Equal(id, board.ID)
}

func TestBoards_getBoardFromTitle(t *testing.T) {
	id := BoardID("kSPsxQZGLKR9tknEt")
	title := BoardTitle("Tableau CRP BFC")
	slug := BoardSlug("tableau-crp-bfc")

	ass := assert.New(t)
	board, err := wekan.GetBoardFromTitle(ctx, string(title))

	ass.Nil(err)
	ass.NotEmpty(board)
	ass.Equal(title, board.Title)
	ass.Equal(slug, board.Slug)
	ass.Equal(id, board.ID)
}

func TestBoards_AddMemberToBoard_with_active_member(t *testing.T) {
	ass := assert.New(t)
	board, err := wekan.GetBoardFromSlug(ctx, "tableau-crp-bfc")

	ass.Nil(err)

	user := BuildUser(t.Name(), t.Name(), t.Name())
	insertedUser, _ := wekan.InsertUser(ctx, user)
	ass.False(board.UserIsMember(insertedUser))

	boardMember := BoardMember{UserID: insertedUser.ID, IsActive: true}
	err = wekan.AddMemberToBoard(ctx, board.ID, boardMember)
	ass.Nil(err)

	actualBoard, _ := wekan.GetBoardFromSlug(ctx, "tableau-crp-bfc")

	ass.True(actualBoard.UserIsMember(insertedUser))
	ass.True(actualBoard.UserIsActiveMember(insertedUser))
}

func TestBoards_AddMemberToBoard_with_inactive_member(t *testing.T) {
	ass := assert.New(t)
	// GIVEN
	board, err := wekan.GetBoardFromSlug(ctx, "tableau-crp-bfc")
	ass.Nil(err)

	user := createTestUser(t, "coucou")
	//insertedUser, _ := wekan.InsertUser(ctx, user)
	//ass.False(board.UserIsMember(insertedUser))

	boardMember := BoardMember{UserID: user.ID, IsActive: false}

	// WHEN
	err = wekan.AddMemberToBoard(ctx, board.ID, boardMember)

	// THEN
	ass.Nil(err)
	actualBoard, _ := wekan.GetBoardFromSlug(ctx, "tableau-crp-bfc")
	ass.True(actualBoard.UserIsMember(user))
	ass.False(actualBoard.UserIsActiveMember(user))
}

func TestBoards_EnableBoardMember(t *testing.T) {
	ass := assert.New(t)
	// GIVEN
	board, err := wekan.GetBoardFromSlug(ctx, "tableau-crp-bfc")

	ass.Nil(err)

	user := BuildUser("test_enable_board_member", "tebm", "Test Enable Board Member")
	insertedUser, _ := wekan.InsertUser(ctx, user)
	ass.False(board.UserIsMember(insertedUser))

	notEnabledUser := BuildUser("test_not_enable_board_member", "tnebm", "Test Not Enable Board Member")
	insertedNotEnabledUser, _ := wekan.InsertUser(ctx, notEnabledUser)
	ass.False(board.UserIsMember(insertedNotEnabledUser))

	wekan.AddMemberToBoard(ctx, board.ID, BoardMember{insertedUser.ID, false, false, false, false, false})
	wekan.AddMemberToBoard(ctx, board.ID, BoardMember{insertedNotEnabledUser.ID, false, false, false, false, false})
	insertedMemberBoard, _ := wekan.GetBoardFromSlug(ctx, "tableau-crp-bfc")
	ass.False(insertedMemberBoard.UserIsActiveMember(insertedUser))
	ass.False(insertedMemberBoard.UserIsActiveMember(insertedNotEnabledUser))

	// WHEN
	err = wekan.EnableBoardMember(ctx, insertedMemberBoard.ID, user.ID)

	// THEN
	ass.Nil(err)
	enabledMemberBoard, _ := wekan.GetBoardFromSlug(ctx, "tableau-crp-bfc")
	ass.True(enabledMemberBoard.UserIsActiveMember(user))
	ass.False(enabledMemberBoard.UserIsActiveMember(notEnabledUser))

	// on vérifie que l'activité a été créée
	expected := newActivityAddBoardMember(wekan.adminUserID, user.ID, board.ID)
	foundActivities, _ := wekan.selectActivitiesFromQuery(ctx, bson.M{"boardId": expected.BoardID, "memberId": expected.MemberID})
	require.NotEmpty(t, foundActivities)
	ass.Len(foundActivities, 1)

	actual := foundActivities[0]

	ass.Condition(activityCompareFunc(expected, actual))
}

func TestBoards_EnableBoardMember_when_user_is_not_on_board(t *testing.T) {
	ass := assert.New(t)
	// GIVEN
	user := BuildUser("an_absent_user", "aau", "An absent user")
	insertedUser, _ := wekan.InsertUser(ctx, user)
	ass.NotNil(insertedUser)
	board, _ := wekan.GetBoardFromSlug(ctx, "tableau-crp-bfc")
	ass.False(board.UserIsMember(user))

	// WHEN
	err := wekan.EnableBoardMember(ctx, board.ID, user.ID)
	ass.Nil(err)

	shouldNotUpdatedBoard, _ := wekan.GetBoardFromSlug(ctx, "tableau-crp-bfc")
	ass.False(shouldNotUpdatedBoard.UserIsMember(user))
	expected := newActivityAddBoardMember(wekan.adminUserID, user.ID, board.ID)
	activities, err := wekan.selectActivitiesFromQuery(ctx, bson.M{"memberId": expected.MemberID, "boardId": expected.BoardID, "activityType": expected.ActivityType})
	ass.Empty(activities, 0)
	ass.Empty(nil)
}

func TestBoards_DisableBoardMember(t *testing.T) {
	ass := assert.New(t)

	// GIVEN
	board, _ := wekan.GetBoardFromSlug(ctx, "tableau-crp-bfc")
	user := BuildUser("test_disable_board_member", "tdbm", "Test Disable Board Member")
	enabledUser := BuildUser("test_not_disable_board_member", "tndbm", "Test Not Disable Board Member")
	insertedUser, _ := wekan.InsertUser(ctx, user)
	insertedEnabledUser, _ := wekan.InsertUser(ctx, enabledUser)
	wekan.AddMemberToBoard(ctx, board.ID, BoardMember{insertedUser.ID, false, true, false, false, false})
	wekan.AddMemberToBoard(ctx, board.ID, BoardMember{insertedEnabledUser.ID, false, true, false, false, false})
	insertedMemberBoard, _ := wekan.GetBoardFromSlug(ctx, "tableau-crp-bfc")

	// WHEN
	err := wekan.DisableBoardMember(ctx, insertedMemberBoard.ID, user.ID)

	// THEN
	ass.Nil(err)
	disabledMemberBoard, err := wekan.GetBoardFromSlug(ctx, "tableau-crp-bfc")
	ass.Nil(err)
	ass.False(disabledMemberBoard.UserIsActiveMember(user))
	ass.True(disabledMemberBoard.UserIsActiveMember(enabledUser))

	// on vérifie que l'activité correspondante a été créée
	expected := newActivityRemoveBoardMember(wekan.adminUserID, user.ID, board.ID)
	foundActivities, _ := wekan.selectActivitiesFromQuery(ctx, bson.M{"boardId": expected.BoardID, "memberId": expected.MemberID, "activityType": "removeBoardMember"})
	require.NotEmpty(t, foundActivities)
	ass.Len(foundActivities, 1)

	actual := foundActivities[0]
	ass.Condition(activityCompareFunc(expected, actual))
}

func TestBoards_EnsureUserIsActiveBoardMember(t *testing.T) {
	// GIVEN
	ass := assert.New(t)
	board, err := wekan.GetBoardFromSlug(ctx, "tableau-crp-bfc")
	user := BuildUser("test_ensure_user_is_active_board_member", "teuiabm", "Test Ensure User Is Active Board Member")
	insertedUser, err := wekan.InsertUser(ctx, user)

	// WHEN
	err = wekan.EnsureUserIsActiveBoardMember(ctx, board.ID, user.ID)

	// THEN
	ass.Nil(err)
	updatedBoard, err := wekan.GetBoardFromSlug(ctx, "tableau-crp-bfc")
	ass.Nil(err)
	ass.True(updatedBoard.UserIsMember(insertedUser))
	ass.True(updatedBoard.UserIsActiveMember(insertedUser))
}

func TestBoards_InsertBoardLabel_whenBoardLabelDontExists(t *testing.T) {
	// GIVEN
	ass := assert.New(t)
	board, _ := wekan.GetBoardFromSlug(ctx, "tableau-crp-bfc")
	boardLabel := BoardLabel{
		ID:    BoardLabelID(newId6()),
		Name:  "test label",
		Color: "orange",
	}

	// WHEN
	err := wekan.InsertBoardLabel(ctx, board, boardLabel)

	// THEN
	ass.Nil(err)
	updatedBoard, err := wekan.GetBoardFromSlug(ctx, "tableau-crp-bfc")
	ass.Nil(err)
	insertedLabel := updatedBoard.GetLabelByID(boardLabel.ID)
	ass.NotEmpty(insertedLabel)
}

func TestBoards_InsertBoardLabel_whenBoardLabelAlreadyExists(t *testing.T) {
	// GIVEN
	ass := assert.New(t)
	testBoardSlug := BoardSlug("tableau-crp-bfc")
	board, _ := wekan.GetBoardFromSlug(ctx, testBoardSlug)
	boardLabel := BoardLabel{
		ID:    BoardLabelID(newId6()),
		Name:  "test label",
		Color: "orange",
	}
	wekan.InsertBoardLabel(ctx, board, boardLabel)
	updatedBoard, _ := wekan.GetBoardFromSlug(ctx, testBoardSlug)

	// WHEN
	err := wekan.InsertBoardLabel(ctx, updatedBoard, boardLabel)

	// THEN
	ass.IsType(BoardLabelAlreadyExistsError{}, err)
}

func TestBoards_InsertBoardLabel_whenBoardIsUnknown(t *testing.T) {
	// WHEN
	ass := assert.New(t)
	boardID := BoardID("fakeID")
	board := Board{ID: boardID}
	boardLabel := BoardLabel{
		ID:    BoardLabelID(newId6()),
		Name:  "test label",
		Color: "orange",
	}

	// THEN
	err := wekan.InsertBoardLabel(ctx, board, boardLabel)
	ass.Equal(UnknownBoardError{board}, err)
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

// nolint:errcheck
func TestBoard_SelectBoardsFromMemberID(t *testing.T) {
	ass := assert.New(t)

	boardBFC, _ := wekan.GetBoardFromSlug(ctx, "tableau-crp-bfc")
	boardNORD, _ := wekan.GetBoardFromSlug(ctx, "tableau-codefi-nord")

	usernameBFC := t.Name() + "bfc"
	usernameNORD := t.Name() + "nord"
	usernameBOTH := t.Name() + "both"

	userBFC := BuildUser(usernameBFC, usernameBFC, usernameBFC)
	userNORD := BuildUser(usernameNORD, usernameNORD, usernameNORD)
	userBOTH := BuildUser(usernameBOTH, usernameBOTH, usernameBOTH)

	insertedUserBFC, _ := wekan.InsertUser(ctx, userBFC)
	insertedUserNORD, _ := wekan.InsertUser(ctx, userNORD)
	insertedUserBOTH, _ := wekan.InsertUser(ctx, userBOTH)

	wekan.AddMemberToBoard(ctx, boardBFC.ID, BoardMember{UserID: insertedUserBFC.ID, IsActive: true})
	wekan.AddMemberToBoard(ctx, boardNORD.ID, BoardMember{UserID: insertedUserNORD.ID, IsActive: true})
	wekan.AddMemberToBoard(ctx, boardBFC.ID, BoardMember{UserID: insertedUserBOTH.ID, IsActive: true})
	wekan.AddMemberToBoard(ctx, boardNORD.ID, BoardMember{UserID: insertedUserBOTH.ID, IsActive: true})

	boardsUserBFC, err := wekan.SelectBoardsFromMemberID(ctx, userBFC.ID)
	ass.Nil(err)
	ass.Len(boardsUserBFC, 2)

	boardsUserNORD, err := wekan.SelectBoardsFromMemberID(ctx, userNORD.ID)
	ass.Len(boardsUserNORD, 2)
	ass.Nil(err)

	boardsUserBOTH, err := wekan.SelectBoardsFromMemberID(ctx, userBOTH.ID)
	ass.Len(boardsUserBOTH, 3)
	ass.Nil(err)

	boardsUserBIDON, err := wekan.SelectBoardsFromMemberID(ctx, "idBidon")
	ass.Len(boardsUserBIDON, 0)
	ass.Nil(err)
}

func createTestUser(t *testing.T, suffix string) User {
	username := t.Name() + suffix
	user := BuildUser(username, username, username)
	user.ID = UserID(username)
	insertedUser, _ := wekan.InsertUser(ctx, user)
	return insertedUser
}

func TestBoards_SelectBoardsFromSlugExpression_withTwoBoards(t *testing.T) {
	ass := assert.New(t)
	boardsTableau, err := wekan.SelectBoardsFromSlugExpression(ctx, "^tableau-.*$")
	ass.Nil(err)
	ass.Len(boardsTableau, 2)
}

func TestBoards_SelectBoardsFromSlugExpression_withNoBoard(t *testing.T) {
	ass := assert.New(t)
	boardsNull, err := wekan.SelectBoardsFromSlugExpression(ctx, "^.*zgorglub$")
	ass.Nil(err)
	ass.Nil(boardsNull)
}

func TestBoards_SelectBoardsFromSlugExpression_withBadRegexp(t *testing.T) {
	ass := assert.New(t)
	boardsNull, err := wekan.SelectBoardsFromSlugExpression(ctx, "[")
	ass.IsType(UnexpectedMongoError{}, err)
	ass.Nil(boardsNull)
}
