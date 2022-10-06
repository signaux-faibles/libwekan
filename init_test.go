package libwekan

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

var ctx = context.Background()

func TestWekan_AdminUsername(t *testing.T) {
	// GIVEN
	expectedName := Username("test")
	unitWekan := Wekan{
		adminUsername: expectedName,
	}

	// THEN
	assert.Equal(t, unitWekan.AdminUsername(), expectedName)
}

func TestWekan_AdminID(t *testing.T) {
	// GIVEN
	expectedID := UserID("test")
	unitWekan := Wekan{
		adminUserID: expectedID,
	}

	// THEN
	assert.Equal(t, unitWekan.AdminID(), expectedID)
}

//func TestWekan_IsPrivileged(t *testing.T) {
//	// GIVEN
//	unitWekan := Wekan{
//		privileged: true,
//	}
//
//	// THEN
//	assert.True(t, unitWekan.IsPrivileged())
//}
