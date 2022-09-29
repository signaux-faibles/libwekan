package libwekan

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type ActivityID string
type Activity struct {
	ID             ActivityID `bson:"_id"`
	UserID         UserID     `bson:"userId,omitempty"`
	Username       Username   `bson:"username"`
	Type           string     `bson:"type"`
	MemberID       UserID     `bson:"memberId,omitempty"`
	ActivityType   string     `bson:"activityType"`
	ActivityTypeID string     `bson:"activityTypeID"`
	BoardID        BoardID    `bson:"boardId,omitempty"`
	ListID         ListID     `bson:"listId,omitempty"`
	CardID         CardID     `bson:"cardId,omitempty"`
	CommentID      CommentID  `bson:"commentId, omitempty"`
	SwimlaneID     SwimlaneID `bson:"swimlaneID,omitempty"`
	CreatedAt      time.Time  `bson:"createdAt"`
	ModifiedAt     time.Time  `bson:"modifiedAt"`
}

func (activity Activity) withID(t time.Time) (Activity, error) {
	if activity.ID != "" {
		return activity, AlreadySetActivity{}
	}
	activity.ID = ActivityID(newId())
	activity.CreatedAt = t.In(time.UTC).Truncate(time.Millisecond)
	activity.ModifiedAt = t.In(time.UTC).Truncate(time.Millisecond)
	return activity, nil
}

func newActivityCreateBoard(userID UserID, boardID BoardID) Activity {
	return Activity{
		BoardID:        boardID,
		ActivityTypeID: string(boardID),
		UserID:         userID,
		ActivityType:   "createBoard",
		Type:           "board",
	}
}

func newActivityCreateSwimlane(userID UserID, boardID BoardID, swimlaneID SwimlaneID) Activity {
	return Activity{
		SwimlaneID:   swimlaneID,
		UserID:       userID,
		BoardID:      boardID,
		ActivityType: "createSwimlane",
		Type:         "list",
	}
}

func newActivityAddBoardMember(userID UserID, memberID UserID, boardID BoardID) Activity {
	return Activity{
		UserID:       userID,
		MemberID:     memberID,
		BoardID:      boardID,
		ActivityType: "addBoardMember",
		Type:         "member",
	}
}

func newActivityJoinMember(userID UserID, username Username, memberID UserID, boardID BoardID, listID ListID, cardID CardID, swimlaneID SwimlaneID) Activity {
	return Activity{
		UserID:       userID,
		Username:     username,
		MemberID:     memberID,
		BoardID:      boardID,
		CardID:       cardID,
		ListID:       listID,
		SwimlaneID:   swimlaneID,
		ActivityType: "joinMember",
	}
}

func newAddComment(userID UserID, boardID BoardID, cardID CardID, commentID CommentID, listID ListID, swimlaneID SwimlaneID) Activity {
	return Activity{
		UserID:       userID,
		BoardID:      boardID,
		CardID:       cardID,
		CommentID:    commentID,
		ListID:       listID,
		SwimlaneID:   swimlaneID,
		ActivityType: "addComment",
	}
}

func (wekan *Wekan) insertActivity(ctx context.Context, activity Activity) (Activity, error) {
	insertable, err := activity.withID(time.Now())
	if err != nil {
		return Activity{}, err
	}
	_, err = wekan.db.Collection("activities").InsertOne(ctx, insertable)
	if err != nil {
		return Activity{}, UnexpectedMongoError{err}
	}
	return insertable, nil
}

func (wekan *Wekan) selectActivityFromID(ctx context.Context, activityID ActivityID) (Activity, error) {
	var activity Activity
	err := wekan.db.Collection("activities").FindOne(ctx, bson.M{"_id": activityID}).Decode(&activity)
	if err != nil {
		return Activity{}, UnexpectedMongoError{err}
	}
	return activity, nil
}
