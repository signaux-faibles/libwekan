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
	updatedUser, err := wekan.GetUserFromID(context.Background(), insertedUser.ID)
	templateBoard, err := wekan.GetBoardFromID(context.Background(), updatedUser.Profile.TemplatesBoardId)
	ass.False(updatedUser.LoginDisabled)
	ass.True(templateBoard.UserIsActiveMember(updatedUser))
}
