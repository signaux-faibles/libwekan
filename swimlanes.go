package libwekan

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type SwimlaneID string

type Swimlane struct {
	ID         SwimlaneID `bson:"_id"`
	Title      string     `bson:"title"`
	BoardID    BoardID    `bson:"boardId"`
	Sort       float64    `bson:"sort"`
	Type       string     `bson:"type"`
	Archived   bool       `bson:"archived"`
	CreatedAt  time.Time  `bson:"createdAt"`
	UpdatedAt  time.Time  `bson:"updateAt"`
	ModifiedAt time.Time  `bson:"modifiedAt"`
}

func BuildSwimlane(boardID BoardID, swimlaneType string, title string, sort float64) Swimlane {
	swimlane := Swimlane{
		ID:         SwimlaneID(newId()),
		Title:      title,
		BoardID:    boardID,
		Sort:       sort,
		Type:       swimlaneType,
		Archived:   false,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		ModifiedAt: time.Now(),
	}
	return swimlane
}

func buildCardTemplateSwimlane(boardId BoardID) Swimlane {
	return BuildSwimlane(boardId, "template-container", "Card Templates", 1)
}
func buildListTemplateSwimlane(boardId BoardID) Swimlane {
	return BuildSwimlane(boardId, "template-container", "List Templates", 2)
}
func buildBoardTemplateSwimlane(boardId BoardID) Swimlane {
	return BuildSwimlane(boardId, "template-container", "Board Templates", 3)
}

func (wekan *Wekan) InsertSwimlane(ctx context.Context, swimlane Swimlane) error {
	if err := wekan.AssertPrivileged(ctx); err != nil {
		return err
	}

	if _, err := wekan.insertActivity(ctx, newActivityCreateSwimlane(wekan.adminUserID, swimlane.BoardID, swimlane.ID)); err != nil {
		return err
	}

	_, err := wekan.db.Collection("swimlanes").InsertOne(ctx, swimlane)
	if err != nil {
		return UnexpectedMongoError{err}
	}
	_, err = wekan.insertActivity(ctx, newActivityCreateSwimlane(wekan.adminUserID, swimlane.BoardID, swimlane.ID))
	if err != nil {
		return UnexpectedMongoError{err}
	}
	return err
}

func (swimlaneID SwimlaneID) Check(ctx context.Context, wekan *Wekan) error {
	_, err := wekan.GetSwimlaneFromID(ctx, swimlaneID)
	return err
}

func (swimlaneID SwimlaneID) GetDocument(ctx context.Context, wekan *Wekan) (Swimlane, error) {
	return wekan.GetSwimlaneFromID(ctx, swimlaneID)
}

func (wekan *Wekan) GetSwimlaneFromID(ctx context.Context, swimlaneID SwimlaneID) (Swimlane, error) {
	var swimlane Swimlane
	if err := wekan.db.Collection("swimlanes").FindOne(ctx, bson.M{"_id": swimlaneID}).Decode(&swimlane); err != nil {
		return Swimlane{}, UnexpectedMongoError{err}
	}
	return swimlane, nil
}

func (wekan *Wekan) GetSwimlanesFromBoardID(ctx context.Context, boardID BoardID) ([]Swimlane, error) {
	var swimlanes []Swimlane
	cur, err := wekan.db.Collection("swimlanes").Find(ctx, bson.M{"boardId": boardID})
	if err != nil {
		return nil, UnexpectedMongoError{err}
	}
	if err := cur.All(ctx, &swimlanes); err != nil {
		return nil, UnexpectedMongoError{err}
	}
	return swimlanes, nil
}
