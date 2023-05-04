package libwekan

import "time"

type CommentID string

type Comment struct {
	ID         CommentID `bson:"_id"`
	BoardID    BoardID   `bson:"boardId"`
	CardID     CardID    `bson:"cardId"`
	CreatedAt  time.Time `bson:"createdAt"`
	ModifiedAt time.Time `bson:"modifiedAt"`
	Text       string    `bson:"text"`
	UserID     UserID    `bson:"userId"`
}
