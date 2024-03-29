package libwekan

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type ActivityID string
type Activity struct {
	ID             ActivityID   `bson:"_id" json:"_id,omitempty"`
	UserID         UserID       `bson:"userId,omitempty" json:"userId,omitempty"`
	Username       Username     `bson:"username,omitempty" json:"username,omitempty"`
	Type           string       `bson:"type,omitempty" json:"type,omitempty"`
	AssigneeID     UserID       `bson:"assigneeId,omitempty" json:"assigneeId,omitempty"`
	MemberID       UserID       `bson:"memberId,omitempty" json:"memberId,omitempty"`
	ActivityType   string       `bson:"activityType,omitempty" json:"activityType,omitempty"`
	ActivityTypeID string       `bson:"activityTypeId,omitempty" json:"activityTypeId,omitempty"`
	BoardID        BoardID      `bson:"boardId,omitempty" json:"boardId,omitempty"`
	BoardLabelID   BoardLabelID `bson:"labelId,omitEmpty" json:"labelId,omitempty"`
	CardTitle      string       `bson:"cardTitle,omitempty" json:"cardTitle,omitempty"`
	ListID         ListID       `bson:"listId,omitempty" json:"listId,omitempty"`
	OldListID      ListID       `bson:"oldListId,omitempty" json:"oldListId,omitempty"`
	ListName       string       `bson:"listName,omitempty" json:"listName,omitempty"`
	CardID         CardID       `bson:"cardId,omitempty" json:"cardId,omitempty"`
	CommentID      CommentID    `bson:"commentId, omitempty"`
	SwimlaneID     SwimlaneID   `bson:"swimlaneId,omitempty" json:"swimlaneId,omitempty"`
	OldSwimlaneID  SwimlaneID   `bson:"oldSwimlaneId,omitempty" json:"oldSwimlaneId,omitempty"`
	SwimlaneName   string       `bson:"swimlaneName,omitempty" json:"swimlaneName,omitempty"`
	CreatedAt      time.Time    `bson:"createdAt" json:"createdAt,omitempty"`
	ModifiedAt     time.Time    `bson:"modifiedAt" json:"modifiedAt,omitempty"`
}

func (activityID ActivityID) Check(ctx context.Context, wekan *Wekan) error {
	_, err := wekan.GetActivityFromID(ctx, activityID)
	return err
}

func (activityID ActivityID) GetDocument(ctx context.Context, wekan *Wekan) (Activity, error) {
	return wekan.GetActivityFromID(ctx, activityID)
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

func newActivityCardJoinAssignee(userID UserID, username Username, memberID UserID, boardID BoardID, listID ListID, cardID CardID, swimlaneID SwimlaneID) Activity {
	return Activity{
		UserID:       userID,
		Username:     username,
		AssigneeID:   memberID,
		BoardID:      boardID,
		CardID:       cardID,
		ListID:       listID,
		SwimlaneID:   swimlaneID,
		ActivityType: "joinAssignee",
	}
}

func newActivityCardUnjoinMember(userID UserID, username Username, memberID UserID, boardID BoardID, listID ListID, cardID CardID, swimlaneID SwimlaneID) Activity {
	return Activity{
		UserID:       userID,
		Username:     username,
		MemberID:     memberID,
		BoardID:      boardID,
		CardID:       cardID,
		ListID:       listID,
		SwimlaneID:   swimlaneID,
		ActivityType: "unjoinMember",
	}
}

func newActivityCardUnjoinAssignee(userID UserID, username Username, assigneeID UserID, boardID BoardID, listID ListID, cardID CardID, swimlaneID SwimlaneID) Activity {
	return Activity{
		UserID:       userID,
		Username:     username,
		AssigneeID:   assigneeID,
		BoardID:      boardID,
		CardID:       cardID,
		ListID:       listID,
		SwimlaneID:   swimlaneID,
		ActivityType: "unjoinAssignee",
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

func (wekan *Wekan) newActivityCreateCardFromCard(ctx context.Context, card Card) (Activity, error) {
	list, err := wekan.GetListFromID(ctx, card.ListID)
	if err != nil {
		return Activity{}, err
	}
	swimlane, err := wekan.GetSwimlaneFromID(ctx, card.SwimlaneID)
	if err != nil {
		return Activity{}, err
	}
	return newActivityCreateCard(card.UserID, list, card, swimlane), nil
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

func (wekan *Wekan) newActivityMoveCardFromMovedCard(ctx context.Context, oldCard Card, userID UserID) (Activity, error) {
	newCard, err := oldCard.ID.GetDocument(ctx, wekan)
	if err != nil {
		return Activity{}, err
	}
	newList, err := newCard.ListID.GetDocument(ctx, wekan)
	if err != nil {
		return Activity{}, err
	}
	swimlane, err := oldCard.SwimlaneID.GetDocument(ctx, wekan)
	if err != nil {
		return Activity{}, err
	}
	return newActivityMoveCard(userID, oldCard, newList, swimlane), nil
}

func newActivityMoveCard(userID UserID, oldCard Card, newList List, swimlane Swimlane) Activity {
	return Activity{
		UserID:        userID,
		ActivityType:  "moveCard",
		BoardID:       oldCard.BoardID,
		CardID:        oldCard.ID,
		CardTitle:     oldCard.Title,
		OldListID:     oldCard.ListID,
		ListID:        newList.ID,
		ListName:      newList.Title,
		SwimlaneID:    swimlane.ID,
		OldSwimlaneID: oldCard.SwimlaneID,
		SwimlaneName:  swimlane.Title,
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

func (wekan *Wekan) SelectActivitiesFromCardID(ctx context.Context, cardID CardID) ([]Activity, error) {
	var activities []Activity
	filter := bson.M{"cardId": cardID}
	sort := options.Find().SetSort(bson.M{"createdAt": 1})
	cur, err := wekan.db.Collection("activities").Find(ctx, filter, sort)
	if err != nil {
		return nil, UnexpectedMongoError{err}
	}
	err = cur.All(ctx, &activities)
	if err != nil {
		return nil, UnexpectedMongoError{err}
	}
	return activities, nil
}

func (wekan *Wekan) GetActivityFromID(ctx context.Context, activityID ActivityID) (Activity, error) {
	var activity Activity
	err := wekan.db.Collection("activities").FindOne(ctx, bson.M{"_id": activityID}).Decode(&activity)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return Activity{}, ActivityNotFoundError{string(activityID)}
		}
		return Activity{}, UnexpectedMongoError{err}
	}
	return activity, nil
}

func (wekan *Wekan) SelectActivitiesFromQuery(ctx context.Context, query bson.M) ([]Activity, error) {
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

func (wekan *Wekan) SelectActivitiesFromBoardID(ctx context.Context, boardID BoardID) ([]Activity, error) {
	return wekan.SelectActivitiesFromQuery(ctx, bson.M{"boardId": boardID})
}
