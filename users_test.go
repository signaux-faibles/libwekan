package libwekan

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_users_setAdmin(t *testing.T) {
	ass := assert.New(t)
	user := BuildUser(t.Name(), t.Name(), t.Name())
	admin := user.Admin(true)
	noAdmin := user.Admin(false)
	ass.True(admin.IsAdmin)
	ass.False(noAdmin.IsAdmin)
	ass.NotSame(user, admin)
	ass.NotSame(user, noAdmin)
	ass.NotSame(admin, noAdmin)
}
