package libwekan

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_usernameExists(t *testing.T) {
	ass := assert.New(t)
	username := Username("signaux.faibles")

	exists, err := wekan.UsernameExists(context.Background(), username)
	ass.Nil(err)
	ass.True(exists)
}

func Test_createUser(t *testing.T) {
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

func Test_createDuplicateUser(t *testing.T) {
	ass := assert.New(t)

	user := BuildUser("tata@grand.vela.com", "TGV", "Tata Grand-Véla")
	_, err := wekan.InsertUser(context.Background(), user)

	ass.Nil(err)

	user = BuildUser("tata@grand.vela.com", "TGV", "Tata Grand-Véla")
	_, err = wekan.InsertUser(context.Background(), user)

	_, ok := err.(UserAlreadyExistsError)
	ass.True(ok)
}
