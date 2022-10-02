//go:build integration

// nolint:errcheck

package libwekan

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCards_InsertCard(t *testing.T) {
	err := wekan.InsertCard(ctx, Card{})
	assert.IsType(t, NotImplemented{}, err)
}

func TestCards_SelectCardsFromUserID(t *testing.T) {
	_, err := wekan.SelectCardsFromUserID(ctx, "")
	assert.IsType(t, NotImplemented{}, err)
}

func TestCards_SelectCardsFromBoardID(t *testing.T) {
	_, err := wekan.SelectCardsFromBoardID(ctx, "")
	assert.IsType(t, NotImplemented{}, err)
}

func TestCards_SelectCardsFromMemberID(t *testing.T) {
	_, err := wekan.SelectCardsFromMemberID(ctx, "")
	assert.IsType(t, NotImplemented{}, err)
}
func TestCards_SelectCardsFromSwimlaneID(t *testing.T) {
	_, err := wekan.SelectCardsFromSwimlaneID(ctx, "")
	assert.IsType(t, NotImplemented{}, err)
}

func TestCards_SelectCardsFromListID(t *testing.T) {
	_, err := wekan.SelectCardsFromListID(ctx, "")
	assert.IsType(t, NotImplemented{}, err)
}
