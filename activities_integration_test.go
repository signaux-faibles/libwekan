//go:build integration

// nolint:errcheck
package libwekan

import (
	"go.mongodb.org/mongo-driver/bson"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestActivities_insertActivity_whenEverythingsFine(t *testing.T) {
	ass := assert.New(t)
	// GIVEN
	activity := newActivityCreateBoard(UserID(t.Name()+"_userId"), BoardID(t.Name()+"_boardID"))

	// WHEN
	insertedActivity, err := wekan.insertActivity(ctx, activity)
	ass.Nil(err)

	// THEN
	selectedActivity, err := wekan.SelectActivityFromID(ctx, insertedActivity.ID)
	ass.Nil(err)
	ass.Equal(insertedActivity, selectedActivity)
}

func TestActivities_InsertActivity_WithActivityIsAlreadySet(t *testing.T) {
	ass := assert.New(t)
	// GIVEN
	activity := newActivityCreateBoard(UserID(t.Name()+"_userId"), BoardID(t.Name()+"_boardID"))

	// WHEN
	insertedActivity, err := wekan.insertActivity(ctx, activity)
	ass.Nil(err)
	notInsertedActivity, err := wekan.insertActivity(ctx, insertedActivity)

	// THEN
	ass.IsType(AlreadySetActivityError{}, err)
	ass.Empty(notInsertedActivity)
}

func TestActivities_selectActivitiesFromBoardID(t *testing.T) {
	ass := assert.New(t)
	// GIVEN
	board, _, _ := createTestBoard(t, "", 1, 1)
	userID := UserID(t.Name() + "userID")
	memberID := UserID(t.Name() + "memberID")
	activityToInsert := newActivityAddBoardMember(userID, memberID, board.ID)
	insertedActivity, _ := wekan.insertActivity(ctx, activityToInsert)

	// WHEN
	activities, err := wekan.SelectActivitiesFromBoardID(ctx, board.ID)

	// THEN
	ass.Nil(err)
	ass.NotNil(activities)
	ass.Contains(activities, insertedActivity)
}

func TestActivities_selectActivitiesFromQuery(t *testing.T) {
	ass := assert.New(t)
	// GIVEN
	board, _, _ := createTestBoard(t, "", 1, 1)
	userID := UserID(t.Name() + "userID")
	memberID := UserID(t.Name() + "memberID")
	activityToInsert := newActivityAddBoardMember(userID, memberID, board.ID)
	insertedActivity, _ := wekan.insertActivity(ctx, activityToInsert)

	// WHEN
	activities, err := wekan.SelectActivitiesFromQuery(ctx, bson.M{"boardId": board.ID})

	// THEN
	ass.Nil(err)
	ass.NotNil(activities)
	ass.Contains(activities, insertedActivity)
}

func TestActivities_GetActivityFromID(t *testing.T) {
	ass := assert.New(t)
	// GIVEN
	activity := newActivityCreateBoard("", "")
	insertedActivity, _ := wekan.insertActivity(ctx, activity)

	// WHEN
	actualActivityFromMethod, err := wekan.GetActivityFromID(ctx, insertedActivity.ID)
	ass.Nil(err)
	actualActivityFromID, err := insertedActivity.ID.GetDocument(ctx, &wekan)
	ass.Nil(err)

	// THEN
	ass.Equal(insertedActivity, actualActivityFromID)
	ass.Equal(insertedActivity, actualActivityFromMethod)
}

func TestActivities_GetActivityFromID_WhenActivityDoesntExist(t *testing.T) {
	ass := assert.New(t)
	// GIVEN
	activityID := ActivityID(t.Name() + "absent")

	// WHEN
	actualActivityFromMethod, errFromMethod := wekan.GetActivityFromID(ctx, activityID)
	actualActivityFromID, errFromID := activityID.GetDocument(ctx, &wekan)

	// THEN
	ass.Empty(actualActivityFromMethod)
	ass.Empty(actualActivityFromID)
	ass.IsType(UnknownActivityError{}, errFromMethod)
	ass.IsType(UnknownActivityError{}, errFromID)
}

func TestActivities_CheckActivityFromID(t *testing.T) {
	ass := assert.New(t)
	// GIVEN
	activity := newActivityCreateBoard("", "")

	// WHEN
	insertedActivity, _ := wekan.insertActivity(ctx, activity)

	// THEN
	ass.Nil(insertedActivity.ID.Check(ctx, &wekan))

}

func TestActivities_CheckActivityFromID_WhenActivityDoesntExist(t *testing.T) {
	ass := assert.New(t)
	// GIVEN
	activityID := ActivityID(t.Name() + "absent")

	// THEN
	ass.IsType(UnknownActivityError{}, activityID.Check(ctx, &wekan))
}
