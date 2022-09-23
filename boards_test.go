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

func Test_UserIsActiveMember(t *testing.T) {
	ass := assert.New(t)

	userToto := User{
		ID: "toto",
	}
	boardWithActiveToto := Board{
		Members: []BoardMember{
			{UserId: "toto", IsActive: true},
		},
	}

	isActive := boardWithActiveToto.UserIsActiveMember(userToto)
	ass.True(isActive)

	boardWithInactiveToto := Board{
		Members: []BoardMember{
			{UserId: "toto", IsActive: false},
		},
	}

	isActive = boardWithInactiveToto.UserIsActiveMember(userToto)
	ass.False(isActive)

	boardWithUserTata := Board{
		Members: []BoardMember{
			{UserId: "tata", IsActive: true},
		},
	}

	isActive = boardWithUserTata.UserIsActiveMember(userToto)
	ass.False(isActive)

	boardWithUserTotoAndTata := Board{
		Members: []BoardMember{
			{UserId: "tata", IsActive: true},
			{UserId: "toto", IsActive: false},
		},
	}

	isActive = boardWithUserTotoAndTata.UserIsActiveMember(userToto)
	ass.False(isActive)
}
