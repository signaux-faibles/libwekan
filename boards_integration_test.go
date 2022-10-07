//go:build integration

// nolint:errcheck
package libwekan

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoards_createBoard(t *testing.T) {
	ass := assert.New(t)
	// GIVEN
	board := BuildBoard("la board à toto", "la-board-a-toto", "board")

	// WHEN
	err := wekan.InsertBoard(ctx, board)
	ass.Nil(err)

	// THEN
	actualBoard, err := wekan.GetBoardFromID(ctx, board.ID)
	ass.Equal(board, actualBoard)
}

func TestBoards_getBoardFromID(t *testing.T) {
	ass := assert.New(t)
	// GIVEN
	id := BoardID("kSPsxQZGLKR9tknEt")
	title := BoardTitle("Tableau CRP BFC")
	slug := BoardSlug("tableau-crp-bfc")

	// WHEN
	board, err := wekan.GetBoardFromID(ctx, id)

	// THEN
	ass.Nil(err)
	ass.NotEmpty(board)
	ass.Equal(title, board.Title)
	ass.Equal(slug, board.Slug)
	ass.Equal(id, board.ID)
}

func TestBoards_getBoardFromID_withBadID(t *testing.T) {
	ass := assert.New(t)
	// GIVEN
	id := BoardID("badID")

	// WHEN
	board, err := wekan.GetBoardFromID(ctx, id)

	// THEN
	ass.IsType(BoardNotFoundError{}, err)
	ass.Empty(board)
}

func TestBoards_getBoardFromID_withBadQuery(t *testing.T) {
	ass := assert.New(t)
	// GIVEN

	badWekan := newTestBadWekan("notAdb")

	// WHEN
	board, err := badWekan.GetBoardFromID(ctx, "")

	// THEN
	ass.IsType(UnexpectedMongoError{}, err)
	ass.Empty(board)
}

func TestBoards_getBoardFromSlug(t *testing.T) {
	ass := assert.New(t)
	// GIVEN
	id := BoardID("kSPsxQZGLKR9tknEt")
	title := BoardTitle("Tableau CRP BFC")
	slug := BoardSlug("tableau-crp-bfc")

	// WHEN
	board, err := wekan.GetBoardFromSlug(ctx, slug)

	// THEN
	ass.Nil(err)
	ass.NotEmpty(board)
	ass.Equal(title, board.Title)
	ass.Equal(slug, board.Slug)
	ass.Equal(id, board.ID)
}

func TestBoards_getBoardFromTitle(t *testing.T) {
	ass := assert.New(t)
	// GIVEN
	id := BoardID("kSPsxQZGLKR9tknEt")
	title := BoardTitle("Tableau CRP BFC")
	slug := BoardSlug("tableau-crp-bfc")

	// WHEN
	board, err := wekan.GetBoardFromTitle(ctx, title)

	// THEN
	ass.Nil(err)
	ass.NotEmpty(board)
	ass.Equal(title, board.Title)
	ass.Equal(slug, board.Slug)
	ass.Equal(id, board.ID)
}

func TestBoards_AddMemberToBoard_with_active_member(t *testing.T) {
	ass := assert.New(t)
	// GIVEN
	board, _ := wekan.GetBoardFromSlug(ctx, "tableau-crp-bfc")
	insertedUser := createTestUser(t, "")
	ass.False(board.UserIsMember(insertedUser))
	boardMember := BoardMember{UserID: insertedUser.ID, IsActive: true}

	// WHEN
	err := wekan.AddMemberToBoard(ctx, board.ID, boardMember)
	ass.Nil(err)

	// THEN
	actualBoard, _ := wekan.GetBoardFromSlug(ctx, "tableau-crp-bfc")
	ass.True(actualBoard.UserIsMember(insertedUser))
	ass.True(actualBoard.UserIsActiveMember(insertedUser))
}

func TestBoards_AddMemberToBoard_withUnknownUser(t *testing.T) {
	ass := assert.New(t)
	// GIVEN
	board, _ := wekan.GetBoardFromSlug(ctx, "tableau-crp-bfc")
	notInsertedUser := BuildUser(t.Name(), t.Name(), t.Name())
	ass.False(board.UserIsMember(notInsertedUser))
	boardMember := BoardMember{UserID: notInsertedUser.ID, IsActive: true}

	// WHEN
	err := wekan.AddMemberToBoard(ctx, board.ID, boardMember)
	ass.IsType(UserNotFoundError{}, err)

	// THEN
	actualBoard, _ := wekan.GetBoardFromSlug(ctx, "tableau-crp-bfc")
	ass.False(actualBoard.UserIsMember(notInsertedUser))
	ass.False(actualBoard.UserIsActiveMember(notInsertedUser))
}

func TestBoards_AddMemberToBoard_with_inactive_member(t *testing.T) {
	ass := assert.New(t)
	// GIVEN
	board, err := wekan.GetBoardFromSlug(ctx, "tableau-crp-bfc")
	ass.Nil(err)
	user := createTestUser(t, "coucou")
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

	insertedUser := createTestUser(t, "enabled")
	ass.False(board.UserIsMember(insertedUser))

	insertedNotEnabledUser := createTestUser(t, "disabled")
	ass.False(board.UserIsMember(insertedNotEnabledUser))

	wekan.AddMemberToBoard(ctx, board.ID, BoardMember{insertedUser.ID, false, false, false, false, false})
	wekan.AddMemberToBoard(ctx, board.ID, BoardMember{insertedNotEnabledUser.ID, false, false, false, false, false})
	insertedMemberBoard, _ := wekan.GetBoardFromSlug(ctx, "tableau-crp-bfc")
	ass.False(insertedMemberBoard.UserIsActiveMember(insertedUser))
	ass.False(insertedMemberBoard.UserIsActiveMember(insertedNotEnabledUser))

	// WHEN
	err = wekan.EnableBoardMember(ctx, insertedMemberBoard.ID, insertedUser.ID)

	// THEN
	ass.Nil(err)
	enabledMemberBoard, _ := wekan.GetBoardFromSlug(ctx, "tableau-crp-bfc")
	ass.True(enabledMemberBoard.UserIsActiveMember(insertedUser))
	ass.False(enabledMemberBoard.UserIsActiveMember(insertedNotEnabledUser))

	// on vérifie que l'activité a été créée
	expected := newActivityAddBoardMember(wekan.adminUserID, insertedUser.ID, board.ID)
	foundActivities, _ := wekan.SelectActivitiesFromQuery(ctx, bson.M{"boardId": expected.BoardID, "memberId": expected.MemberID})
	require.NotEmpty(t, foundActivities)
	ass.Len(foundActivities, 1)

	actual := foundActivities[0]

	ass.Condition(activityCompareFunc(expected, actual))
}

func TestBoards_EnableBoardMember_when_user_is_not_on_board(t *testing.T) {
	ass := assert.New(t)
	// GIVEN
	insertedUser := createTestUser(t, "absent")
	board, _ := wekan.GetBoardFromSlug(ctx, "tableau-crp-bfc")
	ass.False(board.UserIsMember(insertedUser))

	// WHEN
	err := wekan.EnableBoardMember(ctx, board.ID, insertedUser.ID)
	ass.Nil(err)

	// THEN
	shouldNotUpdatedBoard, _ := wekan.GetBoardFromSlug(ctx, "tableau-crp-bfc")
	ass.False(shouldNotUpdatedBoard.UserIsMember(insertedUser))
	expected := newActivityAddBoardMember(wekan.adminUserID, insertedUser.ID, board.ID)
	activities, err := wekan.SelectActivitiesFromQuery(ctx, bson.M{"memberId": expected.MemberID, "boardId": expected.BoardID, "activityType": expected.ActivityType})
	ass.Empty(activities, 0)
	ass.Empty(nil)
}

func TestBoards_DisableBoardMember(t *testing.T) {
	ass := assert.New(t)

	// GIVEN
	board, _ := wekan.GetBoardFromSlug(ctx, "tableau-crp-bfc")
	insertedDisabledUser := createTestUser(t, "disabled")
	insertedEnabledUser := createTestUser(t, "enabled")
	wekan.AddMemberToBoard(ctx, board.ID, BoardMember{insertedDisabledUser.ID, false, true, false, false, false})
	wekan.AddMemberToBoard(ctx, board.ID, BoardMember{insertedEnabledUser.ID, false, true, false, false, false})
	insertedMemberBoard, _ := wekan.GetBoardFromSlug(ctx, "tableau-crp-bfc")

	// WHEN
	err := wekan.DisableBoardMember(ctx, insertedMemberBoard.ID, insertedDisabledUser.ID)

	// THEN
	ass.Nil(err)
	disabledMemberBoard, err := wekan.GetBoardFromSlug(ctx, "tableau-crp-bfc")
	ass.Nil(err)
	ass.False(disabledMemberBoard.UserIsActiveMember(insertedDisabledUser))
	ass.True(disabledMemberBoard.UserIsActiveMember(insertedEnabledUser))

	// on vérifie que l'activité correspondante a été créée
	expected := newActivityRemoveBoardMember(wekan.adminUserID, insertedDisabledUser.ID, board.ID)
	foundActivities, _ := wekan.SelectActivitiesFromQuery(ctx, bson.M{"boardId": expected.BoardID, "memberId": expected.MemberID, "activityType": "removeBoardMember"})
	require.NotEmpty(t, foundActivities)
	ass.Len(foundActivities, 1)

	actual := foundActivities[0]
	ass.Condition(activityCompareFunc(expected, actual))
}

func TestBoards_EnsureUserIsActiveBoardMember(t *testing.T) {
	// GIVEN
	ass := assert.New(t)
	board, _ := wekan.GetBoardFromSlug(ctx, "tableau-crp-bfc")
	insertedUser := createTestUser(t, "")

	// WHEN
	modified, err := wekan.EnsureUserIsActiveBoardMember(ctx, board.ID, insertedUser.ID)
	ass.Nil(err)

	// THEN
	updatedBoard, err := wekan.GetBoardFromSlug(ctx, "tableau-crp-bfc")
	ass.Nil(err)
	ass.True(modified)
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
	ass.ErrorAs(err, &BoardNotFoundError{})
}

//func TestBoard_SelectBoardsFromMemberID(t *testing.T) {
//	ass := assert.New(t)
//
//	boardBFC, _ := wekan.GetBoardFromSlug(ctx, "tableau-crp-bfc")
//	boardNORD, _ := wekan.GetBoardFromSlug(ctx, "tableau-codefi-nord")
//
//	usernameBFC := t.Name() + "bfc"
//	usernameNORD := t.Name() + "nord"
//	usernameBOTH := t.Name() + "both"
//
//	userBFC := BuildUser(usernameBFC, usernameBFC, usernameBFC)
//	userNORD := BuildUser(usernameNORD, usernameNORD, usernameNORD)
//	userBOTH := BuildUser(usernameBOTH, usernameBOTH, usernameBOTH)
//
//	insertedUserBFC, _ := wekan.InsertUser(ctx, userBFC)
//	insertedUserNORD, _ := wekan.InsertUser(ctx, userNORD)
//	insertedUserBOTH, _ := wekan.InsertUser(ctx, userBOTH)
//
//	wekan.AddMemberToBoard(ctx, boardBFC.ID, BoardMember{UserID: insertedUserBFC.ID, IsActive: true})
//	wekan.AddMemberToBoard(ctx, boardNORD.ID, BoardMember{UserID: insertedUserNORD.ID, IsActive: true})
//	wekan.AddMemberToBoard(ctx, boardBFC.ID, BoardMember{UserID: insertedUserBOTH.ID, IsActive: true})
//	wekan.AddMemberToBoard(ctx, boardNORD.ID, BoardMember{UserID: insertedUserBOTH.ID, IsActive: true})
//
//	boardsUserBFC, err := wekan.SelectBoardsFromMemberID(ctx, userBFC.ID)
//	ass.Nil(err)
//	ass.Len(boardsUserBFC, 2)
//
//	boardsUserNORD, err := wekan.SelectBoardsFromMemberID(ctx, userNORD.ID)
//	ass.Len(boardsUserNORD, 2)
//	ass.Nil(err)
//
//	boardsUserBOTH, err := wekan.SelectBoardsFromMemberID(ctx, userBOTH.ID)
//	ass.Len(boardsUserBOTH, 3)
//	ass.Nil(err)
//
//	boardsUserBIDON, err := wekan.SelectBoardsFromMemberID(ctx, "idBidon")
//	ass.Len(boardsUserBIDON, 0)
//	ass.Nil(err)
//}

func TestBoards_SelectDomainRegexp_withTwoBoards(t *testing.T) {
	ass := assert.New(t)
	wekanTest := wekan
	wekanTest.slugDomainRegexp = "^tableau.*"
	boardsTableau, err := wekanTest.SelectDomainBoards(ctx)
	ass.Nil(err)
	ass.Len(boardsTableau, 2)
}

func TestBoards_SelectDomainRegexp_withNoBoard(t *testing.T) {
	ass := assert.New(t)
	wekanTest := wekan
	wekanTest.slugDomainRegexp = "^zgorglub.*"
	boardsNull, err := wekanTest.SelectDomainBoards(ctx)
	ass.Nil(err)
	ass.Nil(boardsNull)
}

func TestBoards_SelectDomainRegexp_withBadRegexp(t *testing.T) {
	ass := assert.New(t)
	wekanTest := wekan
	wekanTest.slugDomainRegexp = "[" // regexp invalide
	boardsNull, err := wekanTest.SelectDomainBoards(ctx)
	ass.IsType(UnexpectedMongoError{}, err)
	ass.Nil(boardsNull)
}

func TestBoards_EnsureUserIsInactiveBoardMember_WhenUserIsActive(t *testing.T) {
	ass := assert.New(t)
	// GIVEN
	user := createTestUser(t, "")
	board, _, _ := createTestBoard(t, "", 1, 1)
	wekan.AddMemberToBoard(ctx, board.ID, BoardMember{UserID: user.ID, IsActive: true})

	// WHEN
	modified, err := wekan.EnsureUserIsInactiveBoardMember(ctx, board.ID, user.ID)
	ass.Nil(err)

	// THEN
	actualBoard, _ := wekan.GetBoardFromID(ctx, board.ID)
	ass.True(modified)
	ass.False(actualBoard.UserIsActiveMember(user))
}

func TestBoards_EnsureUserIsInactiveBoardMember_WhenUserDoesntExist(t *testing.T) {
	ass := assert.New(t)
	// GIVEN
	board, _, _ := createTestBoard(t, "", 1, 1)
	userID := UserID(t.Name() + "notAUserID")
	// WHEN
	modified, err := wekan.EnsureUserIsInactiveBoardMember(ctx, board.ID, userID)

	// THEN
	ass.False(modified)
	ass.IsType(UserNotFoundError{}, err)
}

func TestBoards_EnsureUserIsInactiveBoardMember_WhenBoardDoesntExist(t *testing.T) {
	ass := assert.New(t)
	// GIVEN
	user := createTestUser(t, "")
	badBoardID := BoardID(t.Name() + "notABoardID")
	wekan.AddMemberToBoard(ctx, badBoardID, BoardMember{UserID: user.ID, IsActive: true})

	// WHEN
	modified, err := wekan.EnsureUserIsInactiveBoardMember(ctx, badBoardID, user.ID)

	// THEN
	ass.False(modified)
	ass.IsType(BoardNotFoundError{}, err)
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

func TestBoards_DisableBoardMember_cant_disable_admin(t *testing.T) {
	ass := assert.New(t)
	admin, err := wekan.GetUserFromUsername(ctx, "signaux.faibles")
	err = wekan.DisableBoardMember(ctx, "fakeBoardId", admin.ID)
	ass.IsType(ForbiddenOperationError{}, err)
	ass.ErrorIs(err, ProtectedUserError{admin.ID})
}

func createTestBoard(t *testing.T, suffix string, swimlanesCount int, listsCount int) (Board, []Swimlane, []List) {
	ctx := context.Background()
	board := BuildBoard(t.Name()+suffix, t.Name()+suffix, "board")
	wekan.InsertBoard(ctx, board)
	var swimlanes []Swimlane
	var lists []List
	for i := 0; i < swimlanesCount; i++ {
		swimlane := BuildSwimlane(board.ID, "swimlane", t.Name()+"swimlane", i)
		swimlanes = append(swimlanes, swimlane)
		wekan.InsertSwimlane(ctx, swimlane)
	}
	for i := 0; i < listsCount; i++ {
		title := fmt.Sprintf("%sList%d", t.Name(), i)
		list := BuildList(board.ID, title, i)
		lists = append(lists, list)
		wekan.InsertList(ctx, list)
	}
	return board, swimlanes, lists
}
