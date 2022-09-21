package libwekan

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_UserIsMember(t *testing.T) {
	ass := assert.New(t)

	userToto := User{
		ID: "toto",
	}
	boardWithUserToto := Board{
		Members: []BoardMember{
			{UserId: "toto"},
		},
	}

	member := boardWithUserToto.UserIsMember(userToto)
	ass.True(member)

	boardWithUserTata := Board{
		Members: []BoardMember{
			{UserId: "tata"},
		},
	}

	member = boardWithUserTata.UserIsMember(userToto)
	ass.False(member)

	boardWithUserTotoAndTata := Board{
		Members: []BoardMember{
			{UserId: "tata"},
			{UserId: "toto"},
		},
	}

	member = boardWithUserTotoAndTata.UserIsMember(userToto)
	ass.True(member)
}
