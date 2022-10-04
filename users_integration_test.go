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
	username := Username("toto@grand.velo.com")
	initials := "TGV"
	fullname := "Toto Grand-Vélo"

	user := BuildUser(string(username), initials, fullname)
	insertedUser, err := wekan.InsertUser(ctx, user)
	ass.Nil(err)

	foundUser, err := wekan.GetUserFromUsername(ctx, username)
	ass.Nil(err)
	ass.Equal(username, foundUser.Username)
	ass.Equal(initials, foundUser.Profile.Initials)
	ass.Equal(fullname, foundUser.Profile.Fullname)
	ass.Equal(string(username), foundUser.Emails[0].Address)

	templateBoard, err := wekan.GetBoardFromID(ctx, insertedUser.Profile.TemplatesBoardId)
	ass.Nil(err)
	ass.True(templateBoard.UserIsMember(user))
}

func TestUsers_createDuplicateUser(t *testing.T) {
	ass := assert.New(t)

	user := BuildUser("tata@grand.vela.com", "TGV", "Tata Grand-Véla")
	_, err := wekan.InsertUser(ctx, user)

	ass.Nil(err)

	user = BuildUser("tata@grand.vela.com", "TGV", "Tata Grand-Véla")
	_, err = wekan.InsertUser(ctx, user)

	_, ok := err.(UserAlreadyExistsError)
	ass.True(ok)
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
	enabledUser := BuildUser(t.Name(), t.Name(), t.Name())
	insertedUser, _ := wekan.InsertUser(ctx, enabledUser)
	wekan.DisableUser(ctx, insertedUser)
	err := wekan.EnableUser(ctx, insertedUser)
	ass.Nil(err)
	updatedUser, _ := wekan.GetUserFromID(ctx, insertedUser.ID)
	templateBoard, _ := wekan.GetBoardFromID(ctx, updatedUser.Profile.TemplatesBoardId)
	ass.False(updatedUser.LoginDisabled)
	ass.True(templateBoard.UserIsActiveMember(updatedUser))
}

func TestUsers_GetUsersFromUsernames(t *testing.T) {
	ass := assert.New(t)
	for _, i := range []string{"a", "b", "c", "d", "e"} {
		username := t.Name() + i
		user := BuildUser(username, username, username)
		wekan.InsertUser(ctx, user)
	}

	existingUsernames := []Username{Username(t.Name() + "a"), Username(t.Name() + "b"), Username(t.Name() + "c")}
	selectedExistingUsers, err := wekan.GetUsersFromUsernames(ctx, existingUsernames)
	ass.Nil(err)
	ass.Len(selectedExistingUsers, 3)

	notExistingUsernames := []Username{Username(t.Name() + "m"), Username(t.Name() + "n"), Username(t.Name() + "l")}
	selectedNotExistingUsers, err := wekan.GetUsersFromUsernames(ctx, notExistingUsernames)
	ass.IsType(UnknownUserError{}, err)
	ass.Equal("l'utilisateur n'est pas connu (usernames in (TestUsers_GetUsersFromUsernamesl, TestUsers_GetUsersFromUsernamesm, TestUsers_GetUsersFromUsernamesn))", err.Error())
	ass.Len(selectedNotExistingUsers, 0)

	someExistingUsernames := []Username{Username(t.Name() + "a"), Username(t.Name() + "b"), Username(t.Name() + "l")}
	selectedSomeExistingUsers, err := wekan.GetUsersFromUsernames(ctx, someExistingUsernames)
	ass.IsType(UnknownUserError{}, err)
	ass.Equal("l'utilisateur n'est pas connu (usernames in (TestUsers_GetUsersFromUsernamesl))", err.Error())
	ass.Len(selectedSomeExistingUsers, 0)
}

func TestUsers_GetUsersFromIDs(t *testing.T) {
	ass := assert.New(t)

	var actualUserIDS []UserID
	for _, i := range []string{"a", "b", "c", "d", "e"} {
		username := t.Name() + i
		user := BuildUser(username, username, username)
		wekan.InsertUser(ctx, user)
		actualUserIDS = append(actualUserIDS, user.ID)
	}

	existingUserIDs := actualUserIDS[0:3]
	selectedExistingUsers, err := wekan.GetUsersFromIDs(ctx, existingUserIDs)
	ass.Nil(err)
	ass.Len(selectedExistingUsers, 3)

	notExistingUserIDs := []UserID{"notAnID", "notAnotherID", "pinchMeIfItsAnID"}
	selectedNotExistingUsers, err := wekan.GetUsersFromIDs(ctx, notExistingUserIDs)
	ass.IsType(UnknownUserError{}, err)
	ass.Equal("l'utilisateur n'est pas connu (ids in (notAnID, notAnotherID, pinchMeIfItsAnID))", err.Error())
	ass.Len(selectedNotExistingUsers, 0)

	someExistingUserIDs := append(actualUserIDS[0:2], "notAnID")
	selectedSomeExistingUsers, err := wekan.GetUsersFromIDs(ctx, someExistingUserIDs)
	ass.IsType(UnknownUserError{}, err)
	ass.Equal("l'utilisateur n'est pas connu (ids in (notAnID))", err.Error())
	ass.Len(selectedSomeExistingUsers, 0)
}
