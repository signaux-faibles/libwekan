package libwekan

// ListID porte bien son nom
type ListID string

type List struct {
	ID      ListID  `bson:"_id"`
	Title   string  `bson:"title"`
	BoardID BoardID `bson:"boardId"`
}
