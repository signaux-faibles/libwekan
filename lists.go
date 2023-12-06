package libwekan

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// ListID porte bien son nom
type ListID string

type ListWipLimit struct {
	Value   int  `bson:"value" json:"value,omitempty"`
	Enabled bool `bson:"enabled" json:"enabled,omitempty"`
	Soft    bool `bson:"soft" json:"soft,omitempty"`
}

type List struct {
	ID         ListID       `bson:"_id" json:"_id,omitempty"`
	Title      string       `bson:"title" json:"title,omitempty"`
	BoardID    BoardID      `bson:"boardId" json:"boardId,omitempty"`
	Sort       float64      `bson:"sort" json:"sort,omitempty"`
	Type       string       `bson:"type" json:"type,omitempty"`
	Starred    bool         `bson:"starred" json:"starred,omitempty"`
	Archived   bool         `bson:"archived" json:"archived,omitempty"`
	SwimlaneID string       `bson:"swimlaneId" json:"swimlaneId,omitempty"`
	Width      string       `bson:"width" json:"width,omitempty"`
	CreatedAt  time.Time    `bson:"createdAt" json:"createdAt,omitempty"`
	UpdatedAt  time.Time    `bson:"updatedAt" json:"updatedAt,omitempty"`
	ModifiedAt time.Time    `bson:"modifiedAt" json:"modifiedAt,omitempty"`
	WipLimit   ListWipLimit `bson:"wipLimit" json:"wipLimit,omitempty"`
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
	err := boardID.Check(ctx, wekan)
	if err != nil {
		return nil, err
	}
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
