package libwekan

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

// ListID porte bien son nom
type ListID string

type ListWipLimit struct {
	Value   int  `bson:"value"`
	Enabled bool `bson:"enabled"`
	Soft    bool `bson:"soft"`
}

type List struct {
	ID         ListID       `bson:"_id"`
	Title      string       `bson:"title"`
	BoardID    BoardID      `bson:"boardId"`
	Sort       float64      `bson:"sort"`
	Type       string       `bson:"type"`
	Starred    bool         `bson:"starred"`
	Archived   bool         `bson:"archived"`
	SwimlaneID string       `bson:"swimlaneId"`
	Width      string       `bson:"width"`
	CreatedAt  time.Time    `bson:"createdAt"`
	UpdatedAt  time.Time    `bson:"updatedAt"`
	ModifiedAt time.Time    `bson:"modifiedAt"`
	WipLimit   ListWipLimit `bson:"wipLimit"`
}

func BuildList(boardID BoardID, title string, sort float64) List {
	return List{
		ID:      ListID(newId()),
		Title:   title,
		BoardID: boardID,
		Type:    "list",
		Width:   "270px",
		Sort:    sort,
		WipLimit: ListWipLimit{
			Value: 1,
		},
	}
}

func (listID ListID) Check(ctx context.Context, wekan *Wekan) error {
	_, err := wekan.GetListFromID(ctx, listID)
	return err
}

func (listID ListID) GetDocument(ctx context.Context, wekan *Wekan) (List, error) {
	return wekan.GetListFromID(ctx, listID)
}

func (wekan *Wekan) InsertList(ctx context.Context, list List) error {
	if err := wekan.AssertPrivileged(ctx); err != nil {
		return err
	}
	if _, err := wekan.GetBoardFromID(ctx, list.BoardID); err != nil {
		return err
	}

	_, err := wekan.db.Collection("lists").InsertOne(ctx, list)
	if err != nil {
		return UnexpectedMongoError{err}
	}
	return nil
}

func (wekan *Wekan) GetListFromID(ctx context.Context, listID ListID) (List, error) {
	var list List
	err := wekan.db.Collection("lists").FindOne(ctx, bson.M{"_id": listID}).Decode(&list)
	if err != nil {
		return List{}, UnexpectedMongoError{err}
	}
	return list, nil
}

func (wekan *Wekan) SelectListsFromBoardID(ctx context.Context, boardID BoardID) ([]List, error) {
	var lists []List
	cur, err := wekan.db.Collection("lists").Find(ctx, bson.M{"boardId": boardID})
	if err != nil {
		return nil, UnexpectedMongoError{err}
	}
	err = cur.All(ctx, &lists)
	if err != nil {
		return nil, UnexpectedMongoError{err}
	}
	return lists, nil
}
