package libwekan

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_usernameExists(t *testing.T) {
	ass := assert.New(t)
	username := "signaux.faibles"

	exists, err := wekan.UsernameExists(context.Background(), username)
	ass.Nil(err)
	ass.True(exists)
}

func Test_createUser(t *testing.T) {
	ass := assert.New(t)
	username := "toto@grand.velo.com"
	initials := "TGV"
	fullname := "Toto Grand-Vélo"

	user := BuildUser(username, initials, fullname)
	insertedUser, err := wekan.InsertUser(context.Background(), user)
	ass.Nil(err)

	foundUser, err := wekan.GetUser(context.Background(), username)
	ass.Nil(err)
	ass.Equal(foundUser.Username, username)
	ass.Equal(foundUser.Profile.Initials, initials)
	ass.Equal(foundUser.Profile.Fullname, fullname)
	ass.Equal(foundUser.Emails[0].Address, username)

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
