package libwekan

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func getElement[Element any](elements []Element, fn func(element Element) bool) *Element {
	for _, element := range elements {
		if fn(element) {
			return &element
		}
	}
	return nil
}

func TestCard_AddMember(t *testing.T) {
	card := BuildCard(
		BoardID(t.Name()+"boardID"),
		ListID(t.Name()+"boardID"),
		SwimlaneID(t.Name()+"boardID"),
		(t.Name() + "title"),
		(t.Name() + "description"),
		UserID(t.Name()+"userID"),
	)
	memberID := UserID(t.Name() + "memberID")
	card.AddMember(memberID)
	assert.Len(t, card.Members, 1)
}

func TestCard_AddMember_Duplicate(t *testing.T) {
	card := BuildCard(
		BoardID(t.Name()+"boardID"),
		ListID(t.Name()+"boardID"),
		SwimlaneID(t.Name()+"boardID"),
		(t.Name() + "title"),
		(t.Name() + "description"),
		UserID(t.Name()+"userID"),
	)
	memberID := UserID(t.Name() + "memberID")
	card.AddMember(memberID)
	card.AddMember(memberID)
	assert.Len(t, card.Members, 1)
}
