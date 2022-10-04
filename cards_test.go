package libwekan

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func buildTestCard(t *testing.T) Card {
	boardID := BoardID(t.Name()) + "BoardID"
	listID := ListID(t.Name()) + "ListID"
	swimlaneID := SwimlaneID(t.Name()) + "SwimlaneID"
	title := t.Name() + "Title"
	description := t.Name() + "Description"
	userID := UserID(t.Name()) + "UserID"
	return BuildCard(boardID, listID, swimlaneID, title, description, userID)
}

func TestCard_AddMember(t *testing.T) {
	card := buildTestCard(t)
	memberID := UserID(t.Name() + "memberID")
	card.AddMember(memberID)
	assert.Len(t, card.Members, 1)
}

func TestCard_AddMember_Duplicate(t *testing.T) {
	card := buildTestCard(t)
	memberID := UserID(t.Name() + "memberID")
	card.AddMember(memberID)
	card.AddMember(memberID)
	assert.Len(t, card.Members, 1)
}
