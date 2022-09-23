package libwekan

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Board représente un objet de la collection `boards`
type Board struct {
	ID                         BoardID       `bson:"_id"`
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
	UserId        UserID `bson:"userId"`
	IsAdmin       bool   `bson:"isAdmin"`
	IsActive      bool   `bson:"isActive"`
	IsNoComments  bool   `bson:"isNoComments"`
	IsCommentOnly bool   `bson:"isCommentOnly"`
	IsWorker      bool   `bson:"isWorker"`
}

type BoardID string

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
func (w Wekan) GetBoardFromID(ctx context.Context, id BoardID) (Board, error) {
	var board Board
	err := w.db.Collection("boards").FindOne(ctx, bson.M{"_id": id}).Decode(&board)
	return board, err
}

// getMember teste si l'utilisateur fait partie de l'array .members
func (b Board) getMember(userID UserID) BoardMember {
	for _, boardMember := range b.Members {
		if boardMember.UserId == userID {
			return boardMember
		}
	}
	return BoardMember{}
}

// UserIsMember teste si l'utilisateur est membre de la board, activé ou non
func (b Board) UserIsMember(user User) bool {
	return b.getMember(user.ID) != BoardMember{}
}

// UserIsActiveMember teste si l'utilisateur est activé sur la board, s'il est absent il est alors considéré comme inactif
func (b Board) UserIsActiveMember(user User) bool {
	return b.getMember(user.ID).IsActive
}

// AddUserToBoard ajoute un objet BoardMember sur la board
func (wekan Wekan) AddMemberToBoard(ctx context.Context, boardID BoardID, boardMember BoardMember) error {
	_, err := wekan.db.Collection("boards").UpdateOne(ctx, bson.M{"_id": boardID},
		bson.M{
			"$push": bson.M{
				"members": boardMember,
			},
		})
	return err
}

// EnableUserInBoard active l'utilisateur dans la propriété `member` d'une board
func (wekan Wekan) EnableBoardMember(ctx context.Context, boardID BoardID, userID UserID) error {
	_, err := wekan.db.Collection("boards").UpdateOne(ctx, bson.M{"_id": boardID},
		bson.M{
			"$set": bson.M{"members.$[member].isActive": true},
		},
		&options.UpdateOptions{
			ArrayFilters: &options.ArrayFilters{
				Filters: bson.A{bson.M{"member.userId": userID}}},
		},
	)
	return err
}

// DisableBoardMember desactive l'utilisateur dans la propriété `member` d'une board
func (wekan Wekan) DisableBoardMember(ctx context.Context, boardID BoardID, userID UserID) error {
	_, err := wekan.db.Collection("boards").UpdateOne(ctx, bson.M{"_id": boardID},
		bson.M{
			"$set": bson.M{"members.$[member].isActive": false},
		},
		&options.UpdateOptions{
			ArrayFilters: &options.ArrayFilters{
				Filters: bson.A{bson.M{"member.userId": userID}}},
		},
	)
	return err
}

// EnsureActiveUserInBoard fait en sorte de rendre l'utilisateur participant et actif à une board
func (wekan Wekan) EnsureUserIsActiveBoardMember(ctx context.Context, boardID BoardID, userID UserID) error {
	board, err := wekan.GetBoardFromID(ctx, boardID)
	if err != nil {
		return err
	}
	user, err := wekan.GetUserFromID(ctx, userID)
	if err != nil {
		return err
	}
	if board.UserIsActiveMember(user) {
		return nil // l'utilisateur est déjà membre actif pas d'action requise
	}
	if board.UserIsMember(user) {
		return wekan.EnableBoardMember(ctx, board.ID, user.ID)
	}
	return wekan.AddMemberToBoard(ctx, board.ID, BoardMember{user.ID, false, true, false, false, false})
}

// EnsureUserIsInactiveBoardMember fait en sorte de désactiver un utilisateur sur une board lorsqu'il est participant
func (wekan Wekan) EnsureUserIsInactiveBoardMember(ctx context.Context, boardID BoardID, userID UserID) error {
	board, err := wekan.GetBoardFromID(ctx, boardID)
	if err != nil {
		return err
	}
	user, err := wekan.GetUserFromID(ctx, userID)
	if err != nil {
		return err
	}
	if board.UserIsActiveMember(user) {
		return wekan.DisableBoardMember(ctx, board.ID, user.ID)
	}
	return nil
}

// RemoveUserFromBoard ajoute l'utilisateur à la board
func (wekan Wekan) EnsureUsersOnBoard(ctx context.Context, boardID BoardID, userID []UserID) (Board, error) {
	// board, err := wekan.GetBoardFromID(ctx, boardID)
	// if err != nil {
	// 	return Board{}, err
	// }

	// user, err := wekan.GetUserFromID(ctx, userID)
	// if err != nil {
	// 	return Board{}, err
	// }

	// _, err = wekan.db.Collection("boards").UpdateOne(ctx, bson.M{"_id": board.ID},
	// 	bson.M{
	// 		"$pull": bson.M{
	// 			"members": bson.M{
	// 				"userId": user.ID,
	// 			},
	// 		},
	// 	})
	// if err != nil {
	// 	return Board{}, err
	// }
	// return wekan.GetBoardFromID(ctx, board.ID)
	return Board{}, nil
}

func newBoard(title string, slug string, boardType string) Board {
	board := Board{
		ID:         BoardID(newId()),
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
