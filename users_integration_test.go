package libwekan

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUsers_usernameExists(t *testing.T) {
	ass := assert.New(t)
	username := Username("signaux.faibles")

	exists, err := wekan.UsernameExists(context.Background(), username)
	ass.Nil(err)
	ass.True(exists)
}

func TestUsers_createUser(t *testing.T) {
	ass := assert.New(t)
	username := Username("toto@grand.velo.com")
	initials := "TGV"
	fullname := "Toto Grand-Vélo"

	user := BuildUser(string(username), initials, fullname)
	insertedUser, err := wekan.InsertUser(context.Background(), user)
	ass.Nil(err)

	foundUser, err := wekan.GetUserFromUsername(context.Background(), username)
	ass.Nil(err)
	ass.Equal(username, foundUser.Username)
	ass.Equal(initials, foundUser.Profile.Initials)
	ass.Equal(fullname, foundUser.Profile.Fullname)
	ass.Equal(string(username), foundUser.Emails[0].Address)

	templateBoard, err := wekan.GetBoardFromID(context.Background(), insertedUser.Profile.TemplatesBoardId)
	ass.Nil(err)
	ass.True(templateBoard.UserIsMember(user))
}

func TestUsers_createDuplicateUser(t *testing.T) {
	ass := assert.New(t)

	user := BuildUser("tata@grand.vela.com", "TGV", "Tata Grand-Véla")
	_, err := wekan.InsertUser(context.Background(), user)

	ass.Nil(err)

	user = BuildUser("tata@grand.vela.com", "TGV", "Tata Grand-Véla")
	_, err = wekan.InsertUser(context.Background(), user)

	_, ok := err.(UserAlreadyExistsError)
	ass.True(ok)
}

func TestUsers_DisableUser(t *testing.T) {
	ass := assert.New(t)
	disabledUser := BuildUser(t.Name(), t.Name(), t.Name())
	insertedUser, _ := wekan.InsertUser(context.Background(), disabledUser)
	ass.False(insertedUser.LoginDisabled)
	templateBoard, _ := wekan.GetBoardFromID(context.Background(), insertedUser.Profile.TemplatesBoardId)
	ass.True(templateBoard.UserIsActiveMember(insertedUser))
	board, _ := wekan.GetBoardFromSlug(context.Background(), "tableau-crp-bfc")
	wekan.EnsureUserIsActiveBoardMember(context.Background(), board.ID, insertedUser.ID)
	ass.True(templateBoard.UserIsActiveMember(insertedUser))
	wekan.DisableUser(context.Background(), insertedUser)
	updatedUser, _ := wekan.GetUserFromID(context.Background(), insertedUser.ID)
	ass.True(updatedUser.LoginDisabled)
	updatedBoard, _ := wekan.GetBoardFromSlug(context.Background(), "tableau-crp-bfc")
	ass.False(updatedBoard.UserIsActiveMember(updatedUser))
	updatedTemplateBoard, _ := wekan.GetBoardFromID(context.Background(), insertedUser.Profile.TemplatesBoardId)
	ass.False(updatedTemplateBoard.UserIsActiveMember(updatedUser))
}

func TestUsers_EnableUser(t *testing.T) {
	ass := assert.New(t)
	enabledUser := BuildUser(t.Name(), t.Name(), t.Name())
	insertedUser, _ := wekan.InsertUser(context.Background(), enabledUser)
	wekan.DisableUser(context.Background(), insertedUser)
	err := wekan.EnableUser(context.Background(), insertedUser)
	ass.Nil(err)
	updatedUser, _ := wekan.GetUserFromID(context.Background(), insertedUser.ID)
	templateBoard, _ := wekan.GetBoardFromID(context.Background(), updatedUser.Profile.TemplatesBoardId)
	ass.False(updatedUser.LoginDisabled)
	ass.True(templateBoard.UserIsActiveMember(updatedUser))
}

func TestUsers_GetUsersFromUsernames(t *testing.T) {
	ass := assert.New(t)
	for _, i := range []string{"a", "b", "c", "d", "e"} {
		username := t.Name() + i
		user := BuildUser(username, username, username)
		wekan.InsertUser(context.Background(), user)
	}

	existingUsernames := []Username{Username(t.Name() + "a"), Username(t.Name() + "b"), Username(t.Name() + "c")}
	selectedExistingUsers, err := wekan.GetUsersFromUsernames(context.Background(), existingUsernames)
	ass.Nil(err)
	ass.Len(selectedExistingUsers, 3)

	notExistingUsernames := []Username{Username(t.Name() + "m"), Username(t.Name() + "n"), Username(t.Name() + "l")}
	selectedNotExistingUsers, err := wekan.GetUsersFromUsernames(context.Background(), notExistingUsernames)
	ass.IsType(UnknownUserError{}, err)
	ass.Equal("l'utilisateur n'est pas connu (usernames in (TestUsers_GetUsersFromUsernamesl, TestUsers_GetUsersFromUsernamesm, TestUsers_GetUsersFromUsernamesn))", err.Error())
	ass.Len(selectedNotExistingUsers, 0)

	someExistingUsernames := []Username{Username(t.Name() + "a"), Username(t.Name() + "b"), Username(t.Name() + "l")}
	selectedSomeExistingUsers, err := wekan.GetUsersFromUsernames(context.Background(), someExistingUsernames)
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
		wekan.InsertUser(context.Background(), user)
		actualUserIDS = append(actualUserIDS, user.ID)
	}

	existingUserIDs := actualUserIDS[0:3]
	selectedExistingUsers, err := wekan.GetUsersFromIDs(context.Background(), existingUserIDs)
	ass.Nil(err)
	ass.Len(selectedExistingUsers, 3)

	notExistingUserIDs := []UserID{"notAnID", "notAnotherID", "pinchMeIfItsAnID"}
	selectedNotExistingUsers, err := wekan.GetUsersFromIDs(context.Background(), notExistingUserIDs)
	ass.IsType(UnknownUserError{}, err)
	ass.Equal("l'utilisateur n'est pas connu (ids in (notAnID, notAnotherID, pinchMeIfItsAnID))", err.Error())
	ass.Len(selectedNotExistingUsers, 0)

	someExistingUserIDs := append(actualUserIDS[0:2], "notAnID")
	selectedSomeExistingUsers, err := wekan.GetUsersFromIDs(context.Background(), someExistingUserIDs)
	ass.IsType(UnknownUserError{}, err)
	ass.Equal("l'utilisateur n'est pas connu (ids in (notAnID))", err.Error())
	ass.Len(selectedSomeExistingUsers, 0)
}
