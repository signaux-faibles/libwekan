//go:build integration

// nolint:errcheck
package libwekan

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUsers_usernameExists(t *testing.T) {
	ass := assert.New(t)
	username := Username("signaux.faibles")

	exists, err := wekan.UsernameExists(ctx, username)
	ass.Nil(err)
	ass.True(exists)
}

//func TestUsers_createUser(t *testing.T) {
//	ass := assert.New(t)
//	// GIVEN
//	username := Username("toto@grand.velo.com")
//	initials := "TGV"
//	fullname := "Toto Grand-Vélo"
//	user := BuildUser(string(username), initials, fullname)
//
//	// WHEN
//	insertedUser, err := wekan.InsertUser(ctx, user)
//	ass.Nil(err)
//
//	// THEN
//	foundUser, err := wekan.GetUserFromUsername(ctx, username)
//	ass.Nil(err)
//	ass.Equal(username, foundUser.Username)
//	ass.Equal(initials, foundUser.Profile.Initials)
//	ass.Equal(fullname, foundUser.Profile.Fullname)
//	ass.Equal(string(username), foundUser.Emails[0].Address)
//
//	// On vérifie également que l'utilisateur est bien membre de son tableau de templates
//	templateBoard, err := wekan.GetBoardFromID(ctx, insertedUser.Profile.TemplatesBoardId)
//	ass.Nil(err)
//	ass.True(templateBoard.UserIsMember(user))
//}

//func TestUsers_createDuplicateUser(t *testing.T) {
//	ass := assert.New(t)
//	// GIVEN
//	user := BuildUser("tata@grand.vela.com", "TGV", "Tata Grand-Véla")
//	wekan.InsertUser(ctx, user)
//
//	// WHEN
//	_, err := wekan.InsertUser(ctx, user)
//
//	// THEN
//	ass.IsType(UserAlreadyExistsError{}, err)
//}

//func TestUsers_DisableUser(t *testing.T) {
//	// insertion d'un nouvel utilisateur
//	// ajout de l'utilisateur en tant que membre actif sur tableau-crp-bfc
//	// ajout de l'utilisateur en tant que membre inactif sur tableau-codefi-nord
//	ass := assert.New(t)
//	user := BuildUser(t.Name(), t.Name(), t.Name())
//	insertedUser, _ := wekan.InsertUser(ctx, user)
//	ass.False(insertedUser.LoginDisabled)
//	templateBoard, _ := wekan.GetBoardFromID(ctx, insertedUser.Profile.TemplatesBoardId)
//	ass.True(templateBoard.UserIsActiveMember(insertedUser))
//	bfcBoard, _ := wekan.GetBoardFromSlug(ctx, "tableau-crp-bfc")
//	wekan.EnsureUserIsActiveBoardMember(ctx, bfcBoard.ID, insertedUser.ID)
//	activatedUserBfcBoard, _ := wekan.GetBoardFromSlug(ctx, "tableau-crp-bfc")
//	ass.True(activatedUserBfcBoard.UserIsActiveMember(insertedUser))
//	nordBoard, _ := wekan.GetBoardFromSlug(ctx, "tableau-codefi-nord")
//	wekan.AddMemberToBoard(ctx, nordBoard.ID, BoardMember{UserID: insertedUser.ID})
//	notActivatedUserNordBoard, _ := wekan.GetBoardFromSlug(ctx, "tableau-codefi-nord")
//	ass.True(notActivatedUserNordBoard.UserIsMember(insertedUser))
//	ass.False(notActivatedUserNordBoard.UserIsActiveMember(insertedUser))
//
//	// désactivation de l'utilisateur
//	err := wekan.DisableUser(ctx, insertedUser)
//	ass.Nil(err)
//
//	// vérification des hypothèses
//	// l'utilisateur est désactivé
//	disabledUser, _ := wekan.GetUserFromID(ctx, insertedUser.ID)
//	ass.True(disabledUser.LoginDisabled)
//	// l'utilisateur est désactivé des boards où il était actif (templateBoard, tableau-crp-bfc)
//	disabledUserBfcBoard, _ := wekan.GetBoardFromSlug(ctx, "tableau-crp-bfc")
//	ass.False(disabledUserBfcBoard.UserIsActiveMember(disabledUser))
//	disabledUserNordBoard, _ := wekan.GetBoardFromSlug(ctx, "tableau-crp-bfc")
//	ass.False(disabledUserNordBoard.UserIsActiveMember(disabledUser))
//	disabledUserTemplateBoard, _ := wekan.GetBoardFromID(ctx, insertedUser.Profile.TemplatesBoardId)
//	ass.False(disabledUserTemplateBoard.UserIsActiveMember(disabledUser))
//	// les activities ont été insérées pour les deux boards où il était actif
//	activities, _ := wekan.selectActivitiesFromQuery(ctx, bson.M{"memberId": disabledUser.ID, "activityType": "removeBoardMember"})
//	ass.Len(activities, 2)
//}

func TestUsers_EnableUser(t *testing.T) {
	ass := assert.New(t)
	// GIVEN
	enabledUser := BuildUser(t.Name(), t.Name(), t.Name())
	wekan.InsertUser(ctx, enabledUser)
	insertedUser, _ := enabledUser.ID.GetDocument(ctx, &wekan)
	wekan.DisableUser(ctx, insertedUser)

	// WHEN
	err := wekan.EnableUser(ctx, insertedUser)
	ass.Nil(err)

	// THEN
	updatedUser, _ := wekan.GetUserFromID(ctx, insertedUser.ID)
	templateBoard, _ := wekan.GetBoardFromID(ctx, updatedUser.Profile.TemplatesBoardId)
	ass.False(updatedUser.LoginDisabled)
	ass.True(templateBoard.UserIsActiveMember(updatedUser))
}

func TestUsers_GetUsersFromUsernames_WhenAllExists(t *testing.T) {
	ass := assert.New(t)
	// GIVEN
	for _, suffix := range []string{"a", "b", "c", "d", "e"} {
		createTestUser(t, suffix)
	}
	existingUsernames := []Username{Username(t.Name() + "a"), Username(t.Name() + "b"), Username(t.Name() + "c")}

	// WHEN
	selectedExistingUsers, err := wekan.GetUsersFromUsernames(ctx, existingUsernames)
	ass.Nil(err)

	// THEN
	ass.Len(selectedExistingUsers, 3)
}

func TestUsers_GetUsersFromUsernames_WhenSomeDoesntExist(t *testing.T) {
	ass := assert.New(t)
	// GIVEN
	for _, suffix := range []string{"a", "b", "c", "d", "e"} {
		createTestUser(t, suffix)
	}
	someNotExistingUsernames := []Username{Username(t.Name() + "a"), Username(t.Name() + "b"), Username(t.Name() + "l")}

	// WHEN
	selectedExistingUsers, err := wekan.GetUsersFromUsernames(ctx, someNotExistingUsernames)

	// THEN
	ass.IsType(UserNotFoundError{}, err)
	ass.Len(selectedExistingUsers, 0)
}

func TestUsers_GetUsersFromIDs(t *testing.T) {

	ass := assert.New(t)
	// GIVEN
	var actualUserIDs []UserID
	for _, suffix := range []string{"a", "b", "c", "d", "e"} {
		user := createTestUser(t, suffix)
		actualUserIDs = append(actualUserIDs, user.ID)
	}
	someExistingUserIDs := actualUserIDs[0:3]

	type args struct {
		IDs []UserID
	}
	type want struct {
		errorType interface{}
		counter   int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{"un tableau d'IDs d'utilisateurs existant", args{someExistingUserIDs}, want{nil, len(someExistingUserIDs)}},
		{"un tableau d'IDs d'utilisateurs existant + 1 ID d'utilisateur non existant", args{append(someExistingUserIDs, "iAmNotAnID")}, want{UserNotFoundError{}, 0}},
		{"un tableau d'IDs vide", args{[]UserID{}}, want{nil, 0}},
		{"un tableau d'IDs nil", args{nil}, want{nil, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// WHEN
			selectedExistingUsers, err := wekan.GetUsersFromIDs(ctx, tt.args.IDs)

			// THEN
			if tt.want.errorType != nil {
				ass.IsType(UserNotFoundError{}, err)
			} else {
				ass.NoError(err)
			}
			ass.Len(selectedExistingUsers, tt.want.counter)
		})
	}

}

// helpers internes aux tests
func createTestUser(t *testing.T, suffix string) User {
	username := t.Name() + suffix
	user := BuildUser(username, username, username)
	user.ID = UserID(username)
	wekan.InsertUser(ctx, user)
	insertedUser, _ := user.ID.GetDocument(ctx, &wekan)
	return insertedUser
}

func TestUsers_EnsureUserIsBoardAdmin(t *testing.T) {
	ass := assert.New(t)
	// GIVEN
	board, _, _ := createTestBoard(t, "", 1, 1)
	admin := createTestUser(t, "User")
	wekan.AddMemberToBoard(ctx, board.ID, BoardMember{UserID: admin.ID, IsActive: true, IsAdmin: false})

	// WHEN
	modified, err := wekan.EnsureUserIsBoardAdmin(ctx, board.ID, admin.ID)
	ass.Nil(err)

	// THEN
	ass.True(modified)
	actualBoard, _ := board.ID.GetDocument(ctx, &wekan)
	ass.True(actualBoard.GetMember(admin.ID).IsAdmin)
}

func TestUsers_EnsureUserIsBoardAdmin_WhenUserIsAlreadyAdmin(t *testing.T) {
	ass := assert.New(t)
	// GIVEN
	board, _, _ := createTestBoard(t, "", 1, 1)
	admin := createTestUser(t, "User")
	wekan.AddMemberToBoard(ctx, board.ID, BoardMember{UserID: admin.ID, IsActive: true, IsAdmin: true})

	// WHEN
	modified, err := wekan.EnsureUserIsBoardAdmin(ctx, board.ID, admin.ID)
	ass.Nil(err)

	// THEN
	ass.False(modified)
	actualBoard, _ := board.ID.GetDocument(ctx, &wekan)
	ass.True(actualBoard.GetMember(admin.ID).IsAdmin)
}

func TestUsers_EnsureMemberOutOfCard(t *testing.T) {
	ass := assert.New(t)
	// GIVEN
	board, swimlanes, lists := createTestBoard(t, "", 1, 1)
	user := createTestUser(t, "User")
	member := createTestUser(t, "Member")
	card := createTestCard(t, user.ID, &board.ID, &(swimlanes[0].ID), &(lists[0].ID))
	wekan.AddMemberToBoard(ctx, board.ID, BoardMember{UserID: member.ID, IsActive: true})
	wekan.AddMemberToCard(ctx, card, member, member)

	// WHEN
	modified, err := wekan.EnsureMemberOutOfCard(ctx, card, member, member)
	ass.Nil(err)
	activities, err := wekan.SelectActivitiesFromCardID(ctx, card.ID)
	ass.Nil(err)

	// THEN
	ass.Len(activities, 3)
	ass.True(modified)
	ass.Nil(err)
	actualCard, err := card.ID.GetDocument(ctx, &wekan)
	ass.NotContains(actualCard.Members, member.ID)
}

func TestUsers_EnsureAssigneeOutOfCard(t *testing.T) {
	ass := assert.New(t)
	// GIVEN
	board, swimlanes, lists := createTestBoard(t, "", 1, 1)
	user := createTestUser(t, "User")
	assignee := createTestUser(t, "assignee")
	card := createTestCard(t, user.ID, &board.ID, &(swimlanes[0].ID), &(lists[0].ID))
	err := wekan.AddMemberToBoard(ctx, board.ID, BoardMember{UserID: assignee.ID, IsActive: true})
	ass.NoError(err)

	err = wekan.AddAssigneeToCard(ctx, card, assignee, assignee)
	ass.NoError(err)

	// WHEN
	modified, err := wekan.EnsureAssigneeOutOfCard(ctx, card, assignee, assignee)
	ass.Nil(err)
	activities, err := wekan.SelectActivitiesFromCardID(ctx, card.ID)
	ass.Nil(err)

	// THEN
	ass.Len(activities, 3)
	ass.True(modified)
	ass.Nil(err)
	actualCard, err := card.ID.GetDocument(ctx, &wekan)
	ass.NotContains(actualCard.Members, assignee.ID)
}

func TestUsers_EnsureAssigneeOutOfCard_WhenAssigneesPropertyExists(t *testing.T) {
	ass := assert.New(t)
	// GIVEN
	board, swimlanes, lists := createTestBoard(t, "", 1, 1)
	user := createTestUser(t, "User")
	assignee := createTestUser(t, "assignee")
	card := createTestCard(t, user.ID, &board.ID, &(swimlanes[0].ID), &(lists[0].ID))
	err := wekan.AddMemberToBoard(ctx, board.ID, BoardMember{UserID: assignee.ID, IsActive: true})
	ass.NoError(err)
	err = wekan.AddAssigneeToCard(ctx, card, assignee, assignee)
	ass.NoError(err)
	_, err = wekan.EnsureAssigneeOutOfCard(ctx, card, assignee, assignee)
	ass.NoError(err)

	// WHEN
	err = wekan.AddAssigneeToCard(ctx, card, assignee, assignee)
	ass.NoError(err)
	activities, err := wekan.SelectActivitiesFromCardID(ctx, card.ID)
	ass.Nil(err)

	// THEN
	ass.Len(activities, 4)
	ass.Nil(err)
	actualCard, err := card.ID.GetDocument(ctx, &wekan)
	ass.NotContains(actualCard.Members, assignee.ID)
}

func TestUsers_EnsureMemberOutOfCard_WhenUserIsNotBoardMember(t *testing.T) {
	ass := assert.New(t)
	// GIVEN
	board, swimlanes, lists := createTestBoard(t, "", 1, 1)
	user := createTestUser(t, "User")
	member := createTestUser(t, "Member")
	card := createTestCard(t, user.ID, &board.ID, &(swimlanes[0].ID), &(lists[0].ID))
	wekan.AddMemberToBoard(ctx, board.ID, BoardMember{UserID: member.ID, IsActive: true})

	// WHEN
	modified, err := wekan.EnsureMemberOutOfCard(ctx, card, member, member)
	activities, errActivities := wekan.SelectActivitiesFromCardID(ctx, card.ID)
	ass.Nil(errActivities)

	// THEN
	ass.Len(activities, 1)
	ass.Nil(err)
	ass.False(modified)
	actualCard, err := card.ID.GetDocument(ctx, &wekan)
	ass.NotContains(actualCard.Members, member.ID)
}

func TestUsers_EnsureMemberInCard_WhenUserIsActiveBoardMember(t *testing.T) {
	ass := assert.New(t)
	// GIVEN
	board, swimlanes, lists := createTestBoard(t, "", 1, 1)
	user := createTestUser(t, "User")
	member := createTestUser(t, "Member")
	card := createTestCard(t, user.ID, &board.ID, &(swimlanes[0].ID), &(lists[0].ID))
	wekan.AddMemberToBoard(ctx, board.ID, BoardMember{UserID: member.ID, IsActive: true})

	// WHEN
	modified, err := wekan.EnsureMemberInCard(ctx, card, member, member)
	ass.Nil(err)
	activities, err := wekan.SelectActivitiesFromCardID(ctx, card.ID)
	ass.Nil(err)

	// THEN
	ass.Len(activities, 2)
	ass.True(modified)
	actualCard, _ := card.ID.GetDocument(ctx, &wekan)
	ass.Contains(actualCard.Members, member.ID)
}

func TestUsers_EnsureMemberInCard_WhenUserIsInactiveBoardMember(t *testing.T) {
	ass := assert.New(t)
	// GIVEN
	board, swimlanes, lists := createTestBoard(t, "", 1, 1)
	user := createTestUser(t, "User")
	member := createTestUser(t, "Member")
	card := createTestCard(t, user.ID, &board.ID, &(swimlanes[0].ID), &(lists[0].ID))
	wekan.AddMemberToBoard(ctx, board.ID, BoardMember{UserID: member.ID, IsActive: false})

	// WHEN
	modified, err := wekan.EnsureMemberInCard(ctx, card, member, member)
	activities, errActivities := wekan.SelectActivitiesFromCardID(ctx, card.ID)
	ass.Nil(errActivities)

	// THEN
	ass.Len(activities, 1)
	ass.False(modified)
	ass.ErrorIs(err, UserIsNotMemberError{member.ID})
	actualCard, _ := card.ID.GetDocument(ctx, &wekan)
	ass.NotContains(actualCard.Members, member.ID)
}
