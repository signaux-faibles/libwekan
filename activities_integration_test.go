//go:build integration

// nolint:errcheck
package libwekan

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestActivities_insertActivity_whenEverythingsFine(t *testing.T) {
	ass := assert.New(t)
	activity := newActivityCreateBoard(UserID(t.Name()+"_userId"), BoardID(t.Name()+"_boardID"))
	insertedActivity, err := wekan.insertActivity(ctx, activity)
	ass.Nil(err)

	selectedActivity, err := wekan.selectActivityFromID(ctx, insertedActivity.ID)
	ass.Nil(err)
	ass.Equal(insertedActivity, selectedActivity)
}

func TestActivies_insertActivity_withActivityIsAlreadySet(t *testing.T) {
	ass := assert.New(t)
	activity := newActivityCreateBoard(UserID(t.Name()+"_userId"), BoardID(t.Name()+"_boardID"))
	insertedActivity, err := wekan.insertActivity(ctx, activity)
	ass.Nil(err)

	notInsertedActivity, err := wekan.insertActivity(ctx, insertedActivity)
	ass.Empty(notInsertedActivity)
	ass.IsType(AlreadySetActivityError{}, err)
}

func TestActivities_selectActivitiesFromBoardID(t *testing.T) {
	ass := assert.New(t)
	board, _ := wekan.GetBoardFromSlug(ctx, "tableau-codefi-nord")

	activityToInsert := newActivityAddBoardMember("userID_de_test", "memberId_de_test", board.ID)
	insertedActivity, _ := wekan.insertActivity(ctx, activityToInsert)
	activities, err := wekan.selectActivitiesFromQuery(ctx, bson.M{"boardId": board.ID})
	ass.Nil(err)
	ass.NotNil(activities)
	ass.Contains(activities, insertedActivity)
}

func (wekan *Wekan) selectActivityFromID(ctx context.Context, activityID ActivityID) (Activity, error) {
	var activity Activity
	err := wekan.db.Collection("activities").FindOne(ctx, bson.M{"_id": activityID}).Decode(&activity)
	if err != nil {
		return Activity{}, UnexpectedMongoError{err}
	}
	return activity, nil
}

func (wekan *Wekan) selectActivitiesFromQuery(ctx context.Context, query bson.M) ([]Activity, error) {
	var activities []Activity
	cur, err := wekan.db.Collection("activities").Find(ctx, query)
	if err != nil {
		return nil, UnexpectedMongoError{err}
	}
	if err := cur.All(ctx, &activities); err != nil {
		return nil, UnexpectedMongoError{err}
	}
	return activities, nil
}
