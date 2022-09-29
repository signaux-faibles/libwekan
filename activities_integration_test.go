package libwekan

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_insertActivity_whenEverythingsFine(t *testing.T) {
	ass := assert.New(t)
	activity := newActivityCreateBoard(UserID(t.Name()+"_userId"), BoardID(t.Name()+"_boardID"))
	insertedActivity, err := wekan.insertActivity(context.Background(), activity)
	ass.Nil(err)

	selectedActivity, err := wekan.selectActivityFromID(context.Background(), insertedActivity.ID)
	ass.Nil(err)
	ass.Equal(insertedActivity, selectedActivity)
}

func Test_insertActivity_withActivityIsSet(t *testing.T) {
	ass := assert.New(t)
	activity := newActivityCreateBoard(UserID(t.Name()+"_userId"), BoardID(t.Name()+"_boardID"))
	insertedActivity, err := wekan.insertActivity(context.Background(), activity)
	ass.Nil(err)

	notInsertedActivity, err := wekan.insertActivity(context.Background(), insertedActivity)
	ass.Empty(notInsertedActivity)
	ass.IsType(AlreadySetActivity{}, err)
}
