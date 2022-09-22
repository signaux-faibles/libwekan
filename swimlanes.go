package libwekan

import (
	"context"
	"time"
)

type Swimlane struct {
	ID         string    `bson:"_id"`
	Title      string    `bson:"title"`
	BoardID    BoardID   `bson:"boardId"`
	Sort       int       `bson:"sort"`
	Type       string    `bson:"type"`
	Archived   bool      `bson:"archived"`
	CreatedAt  time.Time `bson:"createdAt"`
	UpdatedAt  time.Time `bson:"updateAt"`
	ModifiedAt time.Time `bson:"modifiedAt"`
}

func newTemplateSwimlaneContainer(boardId BoardID, title string, sort int) Swimlane {
	newSwimlane := Swimlane{
		ID:         newId(),
		Title:      "Card Templates",
		BoardID:    boardId,
		Sort:       sort,
		Type:       "template-container",
		Archived:   false,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		ModifiedAt: time.Now(),
	}
	return newSwimlane
}

func newCardTemplateSwimlane(boardId BoardID) Swimlane {
	return newTemplateSwimlaneContainer(boardId, "Card Templates", 1)
}
func newListTemplateSwimlane(boardId BoardID) Swimlane {
	return newTemplateSwimlaneContainer(boardId, "List Templates", 2)
}
func newBoardTemplateSwimlane(boardId BoardID) Swimlane {
	return newTemplateSwimlaneContainer(boardId, "Board Templates", 3)
}

func (wekan Wekan) InsertSwimlane(ctx context.Context, swimlane Swimlane) error {
	_, err := wekan.db.Collection("swimlanes").InsertOne(ctx, swimlane)
	return err
}
