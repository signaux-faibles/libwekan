package libwekan

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// Board représente un objet de la collection `boards`
type Board struct {
	ID                         string        `bson:"_id"`
	Title                      string        `bson:"title"`
	Permission                 string        `bson:"permission"`
	Sort                       float64       `bson:"sort"`
	Archived                   bool          `bson:"archived"`
	CreatedAt                  time.Time     `bson:"createAt"`
	ModifiedAt                 time.Time     `bson:"modifiedAt"`
	Stars                      int           `bson:"stars"`
	Labels                     []BoardLabel  `bson:"labels"`
	Members                    []BoardMember `bson:"members"`
	Color                      string        `bson:"color"`
	SubtasksDefaultBoardId     *string       `bson:"subtasksDefaultBoardId"`
	SubtasksDefaultListId      *string       `bson:"subtasksDefaultListId"`
	DateSettingsDefaultBoardId *string       `bson:"dateSettingsDefaultBoardId"`
	DateSettingsDefaultListId  *string       `bson:"dateSettingsDefaultListId"`
	AllowsSubtasks             bool          `bson:"allowsSubtasks"`
	AllowsAttachments          bool          `json:"allowsAttachments"`
	AllowsChecklists           bool          `json:"allowsChecklists"`
	AllowsComments             bool          `json:"allowsComments"`
	AllowsDescriptionTitle     bool          `json:"allowsDescriptionTitle"`
	AllowsDescriptionText      bool          `json:"allowsDescriptionText"`
	AllowsActivities           bool          `json:"allowsActivities"`
	AllowsLabels               bool          `json:"allowsLabels"`
	AllowsAssignee             bool          `json:"allowsAssignee"`
	AllowsMembers              bool          `json:"allowsMembers"`
	AllowsRequestedBy          bool          `json:"allowsRequestedBy"`
	AllowsAssignedBy           bool          `json:"allowsAssignedBy"`
	AllowsReceivedDate         bool          `json:"allowsReceivedDate"`
	AllowsStartDate            bool          `json:"allowsStartDate"`
	AllowsEndDate              bool          `json:"allowsEndDate"`
	AllowsDueDate              bool          `json:"allowsDueDate"`
	PresentParentTask          string        `bson:"presentParentTask"`
	IsOvertime                 bool          `bson:"isOvertime"`
	Type                       string        `bson:"type"`
	Slug                       string        `bson:"slug"`
	Watchers                   []interface{} `bson:"watchers"`
	AllowsCardNumber           bool          `bson:"allowsCardNumber"`
	AllowsShowLists            bool          `bson:"allowsShowLists"`
}

type BoardLabel struct {
	ID    string `bson:"_id"`
	Name  string `bson:"name"`
	Color string `bson:"color"`
}

type BoardMember struct {
	UserId        string `bson:"userId"`
	IsAdmin       bool   `bson:"isAdmin"`
	IsActive      bool   `bson:"isActive"`
	IsNoComments  bool   `bson:"isNoComments"`
	IsCommentOnly bool   `bson:"isCommentOnly"`
	IsWorker      bool   `bson:"isWorker"`
}

// ListAllBoards GetBoardFromSlug GetBoardFromID retourne l'objet board à partir du champ .slug
func (w Wekan) ListAllBoards(ctx context.Context) ([]Board, error) {
	var boards []Board
	cursor, err := w.db.Collection("boards").Find(ctx, bson.M{"type": "board"})
	if err != nil {
		return nil, err
	}
	for cursor.Next(ctx) {
		var board Board
		if err := cursor.Decode(&board); err != nil {
			return nil, err
		}
		boards = append(boards, board)
	}
	return boards, nil
}

// GetBoardFromSlug GetBoardFromID retourne l'objet board à partir du champ .slug
func (w Wekan) GetBoardFromSlug(ctx context.Context, slug string) (Board, error) {
	var board Board
	err := w.db.Collection("boards").FindOne(ctx, bson.M{"slug": slug}).Decode(&board)
	return board, err
}

// GetBoardFromTitle GetBoardFromID retourne l'objet board à partir du champ .title
func (w Wekan) GetBoardFromTitle(ctx context.Context, title string) (Board, error) {
	var board Board
	err := w.db.Collection("boards").FindOne(ctx, bson.M{"title": title}).Decode(&board)
	return board, err
}

// GetBoardFromID retourne l'objet board à partir du champ ._id
func (w Wekan) GetBoardFromID(ctx context.Context, id string) (Board, error) {
	var board Board
	err := w.db.Collection("boards").FindOne(ctx, bson.M{"_id": id}).Decode(&board)
	return board, err
}

// UserIsMember teste si l'utilisateur fait partie de l'array .members
func (b Board) UserIsMember(user User) bool {
	for _, boardMember := range b.Members {
		if boardMember.UserId == user.ID {
			return true
		}
	}
	return false
}

// EnsureMember ajoute l'utilisateur à la board
func (wekan Wekan) EnsureUserIsBoardMember(ctx context.Context, board Board, user User) error {
	if !board.UserIsMember(user) {
		_, err := wekan.db.Collection("boards").UpdateOne(ctx, bson.M{"_id": board.ID},
			bson.M{
				"$push": bson.M{
					"members": BoardMember{
						user.ID, false, true, false, false, false,
					},
				},
			})
		return err
	}
	return nil
}

func newBoard(title string, slug string, boardType string) Board {
	board := Board{
		ID:         newId(),
		Title:      title,
		Permission: "private",
		Type:       boardType,
		Slug:       slug,
		Archived:   false,
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
		Stars:      0,
		Labels: []BoardLabel{
			{"green", "n4eJyZ", ""},
			{"yellow", "x57Yyo", ""},
			{"orange", "Axx4ce", ""},
			{"red", "9dSf3v", ""},
			{"purple", "4GgshQ", ""},
			{"blue", "uZwNq7", ""},
		},
		Members:                []BoardMember{},
		Color:                  "belize",
		AllowsSubtasks:         true,
		AllowsAttachments:      true,
		AllowsChecklists:       true,
		AllowsComments:         true,
		AllowsDescriptionTitle: true,
		AllowsDescriptionText:  true,
		AllowsActivities:       true,
		AllowsLabels:           true,
		AllowsAssignee:         true,
		AllowsMembers:          true,
		AllowsRequestedBy:      true,
		AllowsAssignedBy:       true,
		AllowsReceivedDate:     true,
		AllowsStartDate:        true,
		AllowsEndDate:          true,
		AllowsDueDate:          true,
		PresentParentTask:      "no-parent",
		IsOvertime:             false,
		Sort:                   0,
		AllowsCardNumber:       false,
		AllowsShowLists:        true,
	}

	return board

}

func (wekan Wekan) InsertBoard(ctx context.Context, board Board) error {
	_, err := wekan.db.Collection("boards").InsertOne(ctx, board)
	return err
}
