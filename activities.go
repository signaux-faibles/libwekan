package libwekan

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type ActivityID string
type Activity struct {
	ID             ActivityID   `bson:"_id"`
	UserID         UserID       `bson:"userId,omitempty"`
	Username       Username     `bson:"username,omitempty"`
	Type           string       `bson:"type,omitempty"`
	MemberID       UserID       `bson:"memberId,omitempty"`
	ActivityType   string       `bson:"activityType,omitempty"`
	ActivityTypeID string       `bson:"activityTypeID,omitempty"`
	BoardID        BoardID      `bson:"boardId,omitempty"`
	BoardLabelID   BoardLabelID `bson:"labelId,omitEmpty"`
	CardTitle      string       `bson:"cardTitle,omitempty"`
	ListID         ListID       `bson:"listId,omitempty"`
	ListName       string       `bson:"listName,omitempty"`
	CardID         CardID       `bson:"cardId,omitempty"`
	CommentID      CommentID    `bson:"commentId, omitempty"`
	SwimlaneID     SwimlaneID   `bson:"swimlaneID,omitempty"`
	SwimlaneName   string       `bson:"swimlaneName,omitempty"`
	CreatedAt      time.Time    `bson:"createdAt"`
	ModifiedAt     time.Time    `bson:"modifiedAt"`
}

func (activityID ActivityID) Check(ctx context.Context, wekan *Wekan) error {
	_, err := wekan.GetActivityFromID(ctx, activityID)
	return err
}

func (activity Activity) withIDandDates(t time.Time) (Activity, error) {
	if activity.ID != "" {
		return activity, AlreadySetActivityError{}
	}
	activity.ID = ActivityID(newId())
	activity.CreatedAt = toMongoTime(t)
	activity.ModifiedAt = toMongoTime(t)
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

func newActivityRemoveBoardMember(userID UserID, memberID UserID, boardID BoardID) Activity {
	return Activity{
		UserID:       userID,
		MemberID:     memberID,
		BoardID:      boardID,
		ActivityType: "removeBoardMember",
		Type:         "member",
	}
}

func newActivityCardJoinMember(userID UserID, username Username, memberID UserID, boardID BoardID, listID ListID, cardID CardID, swimlaneID SwimlaneID) Activity {
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

func newActivityAddComment(userID UserID, boardID BoardID, cardID CardID, commentID CommentID, listID ListID, swimlaneID SwimlaneID) Activity {
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

func newActivityAddedLabel(userID UserID, boardLabelID BoardLabelID, boardID BoardID, swimlaneID SwimlaneID) Activity {
	return Activity{
		UserID:       userID,
		BoardLabelID: boardLabelID,
		ActivityType: "addedLabel",
		BoardID:      boardID,
		SwimlaneID:   swimlaneID,
	}
}

func newActivityCreateCard(userID UserID, list List, card Card, swimlane Swimlane) Activity {
	return Activity{
		UserID:       userID,
		ActivityType: "createCard",
		BoardID:      card.BoardID,
		ListName:     list.Title,
		ListID:       list.ID,
		CardID:       card.ID,
		CardTitle:    card.Title,
		SwimlaneName: swimlane.Title,
		SwimlaneID:   swimlane.ID,
	}
}

func (wekan *Wekan) insertActivity(ctx context.Context, activity Activity) (Activity, error) {
	insertable, err := activity.withIDandDates(time.Now())
	if err != nil {
		return Activity{}, err
	}
	_, err = wekan.db.Collection("activities").InsertOne(ctx, insertable)
	if err != nil {
		return Activity{}, UnexpectedMongoError{err}
	}
	return insertable, nil
}

func (wekan *Wekan) GetActivityFromID(ctx context.Context, activityID ActivityID) (Activity, error) {
	var activity Activity
	err := wekan.db.Collection("activities").FindOne(ctx, bson.M{"_id": activityID}).Decode(&activity)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return Activity{}, UnknownActivityError{string(activityID)}
		}
		return Activity{}, UnexpectedMongoError{err}
	}
	return activity, nil
}
