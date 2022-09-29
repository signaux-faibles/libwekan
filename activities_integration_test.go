package libwekan

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
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

func Test_insertActivity_withActivityIsAlreadySet(t *testing.T) {
	ass := assert.New(t)
	activity := newActivityCreateBoard(UserID(t.Name()+"_userId"), BoardID(t.Name()+"_boardID"))
	insertedActivity, err := wekan.insertActivity(context.Background(), activity)
	ass.Nil(err)

	notInsertedActivity, err := wekan.insertActivity(context.Background(), insertedActivity)
	ass.Empty(notInsertedActivity)
	ass.IsType(AlreadySetActivityError{}, err)
}

func Test_selectActivitiesFromBoardID(t *testing.T) {
	ass := assert.New(t)

	board, err := wekan.GetBoardFromSlug(context.Background(), "tableau-codefi-nord")
	ass.Nil(err)
	ass.NotEmpty(board)

	activityToInsert := newActivityAddBoardMember("userID_de_test", "memberId_de_test", board.ID)
	insertedActivity, _ := wekan.insertActivity(context.Background(), activityToInsert)
	activities, err := wekan.selectActivitiesFromQuery(context.Background(), bson.M{"boardId": board.ID})
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
