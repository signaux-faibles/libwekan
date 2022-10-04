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
	Sort       int          `bson:"sort"`
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

func BuildList(title string, boardID BoardID, sort int) List {
	return List{
		ID:      ListID(newId()),
		Title:   title,
		BoardID: boardID,
		Type:    "list",
		Width:   "270px",
		WipLimit: ListWipLimit{
			Value: 1,
		},
	}
}
func (wekan *Wekan) InsertList(ctx context.Context, list List) error {
	_, err := wekan.db.Collection("list").InsertOne(ctx, list)
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
	cur, err := wekan.db.Collection("lists").Find(ctx, bson.M{})
	if err != nil {
		return nil, UnexpectedMongoError{err}
	}
	err = cur.All(ctx, &lists)
	if err != nil {
		return nil, UnexpectedMongoError{err}
	}
	return lists, nil
}
