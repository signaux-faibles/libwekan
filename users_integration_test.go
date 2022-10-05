//go:build integration
// +build integration

// nolint:errcheck
package libwekan

import (
	"go.mongodb.org/mongo-driver/bson"
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

func TestUsers_createUser(t *testing.T) {
	ass := assert.New(t)
	// GIVEN
	username := Username("toto@grand.velo.com")
	initials := "TGV"
	fullname := "Toto Grand-Vélo"
	user := BuildUser(string(username), initials, fullname)

	// WHEN
	insertedUser, err := wekan.InsertUser(ctx, user)
	ass.Nil(err)

	// THEN
	foundUser, err := wekan.GetUserFromUsername(ctx, username)
	ass.Nil(err)
	ass.Equal(username, foundUser.Username)
	ass.Equal(initials, foundUser.Profile.Initials)
	ass.Equal(fullname, foundUser.Profile.Fullname)
	ass.Equal(string(username), foundUser.Emails[0].Address)

	// On vérifie également que l'utilisateur est bien membre de son tableau de templates
	templateBoard, err := wekan.GetBoardFromID(ctx, insertedUser.Profile.TemplatesBoardId)
	ass.Nil(err)
	ass.True(templateBoard.UserIsMember(user))
}

func TestUsers_createDuplicateUser(t *testing.T) {
	ass := assert.New(t)
	// GIVEN
	user := BuildUser("tata@grand.vela.com", "TGV", "Tata Grand-Véla")
	wekan.InsertUser(ctx, user)

	// WHEN
	_, err := wekan.InsertUser(ctx, user)

	// THEN
	ass.IsType(UserAlreadyExistsError{}, err)
}

func TestUsers_DisableUser(t *testing.T) {
	// insertion d'un nouvel utilisateur
	// ajout de l'utilisateur en tant que membre actif sur tableau-crp-bfc
	// ajout de l'utilisateur en tant que membre inactif sur tableau-codefi-nord
	ass := assert.New(t)
	user := BuildUser(t.Name(), t.Name(), t.Name())
	insertedUser, _ := wekan.InsertUser(ctx, user)
	ass.False(insertedUser.LoginDisabled)
	templateBoard, _ := wekan.GetBoardFromID(ctx, insertedUser.Profile.TemplatesBoardId)
	ass.True(templateBoard.UserIsActiveMember(insertedUser))
	bfcBoard, _ := wekan.GetBoardFromSlug(ctx, "tableau-crp-bfc")
	wekan.EnsureUserIsActiveBoardMember(ctx, bfcBoard.ID, insertedUser.ID)
	activatedUserBfcBoard, _ := wekan.GetBoardFromSlug(ctx, "tableau-crp-bfc")
	ass.True(activatedUserBfcBoard.UserIsActiveMember(insertedUser))
	nordBoard, _ := wekan.GetBoardFromSlug(ctx, "tableau-codefi-nord")
	wekan.AddMemberToBoard(ctx, nordBoard.ID, BoardMember{UserID: insertedUser.ID})
	notActivatedUserNordBoard, _ := wekan.GetBoardFromSlug(ctx, "tableau-codefi-nord")
	ass.True(notActivatedUserNordBoard.UserIsMember(insertedUser))
	ass.False(notActivatedUserNordBoard.UserIsActiveMember(insertedUser))

	// désactivation de l'utilisateur
	err := wekan.DisableUser(ctx, insertedUser)
	ass.Nil(err)

	// vérification des hypothèses
	// l'utilisateur est désactivé
	disabledUser, _ := wekan.GetUserFromID(ctx, insertedUser.ID)
	ass.True(disabledUser.LoginDisabled)
	// l'utilisateur est désactivé des boards où il était actif (templateBoard, tableau-crp-bfc)
	disabledUserBfcBoard, _ := wekan.GetBoardFromSlug(ctx, "tableau-crp-bfc")
	ass.False(disabledUserBfcBoard.UserIsActiveMember(disabledUser))
	disabledUserNordBoard, _ := wekan.GetBoardFromSlug(ctx, "tableau-crp-bfc")
	ass.False(disabledUserNordBoard.UserIsActiveMember(disabledUser))
	disabledUserTemplateBoard, _ := wekan.GetBoardFromID(ctx, insertedUser.Profile.TemplatesBoardId)
	ass.False(disabledUserTemplateBoard.UserIsActiveMember(disabledUser))
	// les activities ont été insérées pour les deux boards où il était actif
	activities, _ := wekan.selectActivitiesFromQuery(ctx, bson.M{"memberId": disabledUser.ID, "activityType": "removeBoardMember"})
	ass.Len(activities, 2)
}

func TestUsers_EnableUser(t *testing.T) {
	ass := assert.New(t)
	// GIVEN
	enabledUser := BuildUser(t.Name(), t.Name(), t.Name())
	insertedUser, _ := wekan.InsertUser(ctx, enabledUser)
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
	ass.IsType(UnknownUserError{}, err)
	ass.Len(selectedExistingUsers, 0)
}

func TestUsers_GetUsersFromIDs_WhenUsersExists(t *testing.T) {
	ass := assert.New(t)
	// GIVEN
	var actualUserIDs []UserID
	for _, suffix := range []string{"a", "b", "c", "d", "e"} {
		user := createTestUser(t, suffix)
		actualUserIDs = append(actualUserIDs, user.ID)
	}
	someExistingUserIDs := actualUserIDs[0:3]

	// WHEN
	selectedExistingUsers, err := wekan.GetUsersFromIDs(ctx, someExistingUserIDs)

	// THEN
	ass.Nil(err)
	ass.Len(selectedExistingUsers, 3)
}

func TestUsers_GetUsersFromIDs(t *testing.T) {
	ass := assert.New(t)
	// GIVEN
	var actualUserIDs []UserID
	for _, suffix := range []string{"a", "b", "c", "d", "e"} {
		user := createTestUser(t, suffix)
		actualUserIDs = append(actualUserIDs, user.ID)
	}
	someExistingUserIDs := append(actualUserIDs[0:3], "iAmNotAnID")

	// WHEN
	selectedExistingUsers, err := wekan.GetUsersFromIDs(ctx, someExistingUserIDs)

	// THEN
	ass.IsType(UnknownUserError{}, err)
	ass.Len(selectedExistingUsers, 0)
}
