package libwekan

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_BuildList(t *testing.T) {
	list := BuildList(BoardID(t.Name()), t.Name(), 0)
	expected := List{
		ID:      "expectedID",
		Title:   t.Name(),
		BoardID: BoardID(t.Name()),
		Type:    "list",
		Width:   "270px",
		WipLimit: ListWipLimit{
			Value: 1,
		},
	}
	list.ID = expected.ID
	assert.Equal(t, expected, list)
}
