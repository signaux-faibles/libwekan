package libwekan

import "time"

type CommentID string

type Comment struct {
	ID         CommentID `bson:"_id" json:"_id,omitempty"`
	BoardID    BoardID   `bson:"boardId" json:"boardId,omitempty"`
	CardID     CardID    `bson:"cardId" json:"cardId,omitempty"`
	CreatedAt  time.Time `bson:"createdAt" json:"createdAt,omitempty"`
	ModifiedAt time.Time `bson:"modifiedAt" json:"modifiedAt,omitempty"`
	Text       string    `bson:"text" json:"text,omitempty"`
	UserID     UserID    `bson:"userId" json:"userId,omitempty"`
}
