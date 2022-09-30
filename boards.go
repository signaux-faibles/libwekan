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
	Title                      BoardTitle    `bson:"title"`
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
	Slug                       BoardSlug     `bson:"slug"`
	Watchers                   []interface{} `bson:"watchers"`
	AllowsCardNumber           bool          `bson:"allowsCardNumber"`
	AllowsShowLists            bool          `bson:"allowsShowLists"`
}

type BoardLabelID string
type BoardLabelName string
type BoardLabel struct {
	ID    BoardLabelID   `bson:"_id"`
	Name  BoardLabelName `bson:"name"`
	Color string         `bson:"color"`
}

type BoardMember struct {
	UserID        UserID `bson:"userId"`
	IsAdmin       bool   `bson:"isAdmin"`
	IsActive      bool   `bson:"isActive"`
	IsNoComments  bool   `bson:"isNoComments"`
	IsCommentOnly bool   `bson:"isCommentOnly"`
	IsWorker      bool   `bson:"isWorker"`
}

type BoardID string
type BoardSlug string
type BoardTitle string

// GetLabelByName retourne l'objet BoardLabel correspondant au nom, vide si absent
func (board Board) GetLabelByName(name BoardLabelName) BoardLabel {
	for _, label := range board.Labels {
		if label.Name == name {
			return label
		}
	}
	return BoardLabel{}
}

// GetLabelByID retourne l'objet BoardLabel correspondant à l'ID, vide si absent
func (board Board) GetLabelByID(id BoardLabelID) BoardLabel {
	for _, label := range board.Labels {
		if label.ID == id {
			return label
		}
	}
	return BoardLabel{}
}

// ListAllBoards retourne toutes les boards
func (wekan *Wekan) ListAllBoards(ctx context.Context) ([]Board, error) {
	var boards []Board
	cursor, err := wekan.db.Collection("boards").Find(ctx, bson.M{"type": "board"})
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

// GetBoardFromSlug retourne l'objet board à partir du champ .slug
func (wekan *Wekan) GetBoardFromSlug(ctx context.Context, slug BoardSlug) (Board, error) {
	var board Board
	err := wekan.db.Collection("boards").FindOne(ctx, bson.M{"slug": slug}).Decode(&board)
	return board, err
}

// GetBoardFromTitle GetBoardFromID retourne l'objet board à partir du champ title
func (wekan *Wekan) GetBoardFromTitle(ctx context.Context, title string) (Board, error) {
	var board Board
	err := wekan.db.Collection("boards").FindOne(ctx, bson.M{"title": title}).Decode(&board)
	return board, err
}

// GetBoardFromID retourne l'objet board à partir du champ ._id
func (wekan *Wekan) GetBoardFromID(ctx context.Context, id BoardID) (Board, error) {
	var board Board
	err := wekan.db.Collection("boards").FindOne(ctx, bson.M{"_id": id}).Decode(&board)
	return board, err
}

// getMember teste si l'utilisateur fait partie de l'array members
func (board Board) getMember(userID UserID) BoardMember {
	for _, boardMember := range board.Members {
		if boardMember.UserID == userID {
			return boardMember
		}
	}
	return BoardMember{}
}

// UserIsMember teste si l'utilisateur est membre de la board, activé ou non
func (board Board) UserIsMember(user User) bool {
	return board.getMember(user.ID) != BoardMember{}
}

// UserIsActiveMember teste si l'utilisateur est activé sur la board, s'il est absent il est alors considéré comme inactif
func (board Board) UserIsActiveMember(user User) bool {
	return board.getMember(user.ID).IsActive
}

// AddMemberToBoard ajoute un objet BoardMember sur la board
func (wekan *Wekan) AddMemberToBoard(ctx context.Context, boardID BoardID, boardMember BoardMember) error {
	if err := wekan.CheckAdminUserIsAdmin(ctx); err != nil {
		return err
	}

	if boardMember.IsActive {
		activity := newActivityAddBoardMember(wekan.adminUserID, boardMember.UserID, boardID)
		_, err := wekan.insertActivity(ctx, activity)
		if err != nil {
			return err
		}
	}

	_, err := wekan.db.Collection("boards").UpdateOne(ctx, bson.M{"_id": boardID},
		bson.M{
			"$push": bson.M{
				"members": boardMember,
			},
		})
	if err != nil {
		return UnexpectedMongoError{err}
	}
	return nil
}

// EnableBoardMember active l'utilisateur dans la propriété `member` d'une board
func (wekan *Wekan) EnableBoardMember(ctx context.Context, boardID BoardID, userID UserID) error {
	if err := wekan.CheckAdminUserIsAdmin(ctx); err != nil {
		return err
	}

	updateResults, err := wekan.db.Collection("boards").UpdateOne(ctx,
		bson.M{"_id": boardID},
		bson.M{
			"$set": bson.M{"members.$[member].isActive": true},
		},
		&options.UpdateOptions{
			ArrayFilters: &options.ArrayFilters{
				Filters: bson.A{bson.M{"member.userId": userID}}},
		},
	)
	if err != nil {
		return UnexpectedMongoError{err}
	}
	if updateResults.ModifiedCount == 1 {
		activity := newActivityAddBoardMember(wekan.adminUserID, userID, boardID)
		_, err = wekan.insertActivity(context.Background(), activity)
		return err
	}
	return nil
}

// DisableBoardMember desactive l'utilisateur dans la propriété `member` d'une board
func (wekan *Wekan) DisableBoardMember(ctx context.Context, boardID BoardID, userID UserID) error {
	if wekan.adminUserID == userID {
		return ForbiddenOperationError{"On ne "}
	}
	if err := wekan.CheckAdminUserIsAdmin(ctx); err != nil {
		return err
	}

	updateStats, err := wekan.db.Collection("boards").UpdateOne(ctx, bson.M{"_id": boardID},
		bson.M{
			"$set": bson.M{"members.$[member].isActive": false},
		},
		&options.UpdateOptions{
			ArrayFilters: &options.ArrayFilters{
				Filters: bson.A{bson.M{"member.userId": userID}}},
		},
	)
	if err != nil {
		return UnexpectedMongoError{err}
	}
	if updateStats.ModifiedCount == 1 {
		activity := newActivityRemoveBoardMember(wekan.adminUserID, userID, boardID)
		_, err = wekan.insertActivity(context.Background(), activity)
		return err
	}
	return nil
}

// EnsureUserIsActiveBoardMember fait en sorte de rendre l'utilisateur participant et actif à une board
func (wekan *Wekan) EnsureUserIsActiveBoardMember(ctx context.Context, boardID BoardID, userID UserID) error {
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
func (wekan *Wekan) EnsureUserIsInactiveBoardMember(ctx context.Context, boardID BoardID, userID UserID) error {
	if err := wekan.CheckAdminUserIsAdmin(ctx); err != nil {
		return err
	}

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

func (wekan *Wekan) EnsureUserIsBoardAdmin(ctx context.Context, boardID BoardID, userID UserID) error {
	if err := wekan.CheckAdminUserIsAdmin(ctx); err != nil {
		return err
	}

	err := wekan.EnsureUserIsActiveBoardMember(ctx, boardID, userID)
	if err != nil {
		return err
	}
	_, err = wekan.db.Collection("boards").UpdateOne(ctx, bson.M{"_id": boardID},
		bson.M{
			"$set": bson.M{"members.$[member].isAdmin": true},
		},
		&options.UpdateOptions{
			ArrayFilters: &options.ArrayFilters{
				Filters: bson.A{bson.M{"member.userId": userID}}},
		},
	)
	if err != nil {
		return UnexpectedMongoError{err}
	}
	return nil
}

func newBoard(title string, slug string, boardType string) Board {
	board := Board{
		ID:         BoardID(newId()),
		Title:      BoardTitle(title),
		Permission: "private",
		Type:       boardType,
		Slug:       BoardSlug(slug),
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

func (wekan *Wekan) InsertBoard(ctx context.Context, board Board) error {
	if err := wekan.CheckAdminUserIsAdmin(ctx); err != nil {
		return err
	}

	_, err := wekan.insertActivity(ctx, newActivityCreateBoard(wekan.adminUserID, board.ID))
	if err != nil {
		return UnexpectedMongoError{err}
	}
	_, err = wekan.db.Collection("boards").InsertOne(ctx, board)
	if err != nil {
		return UnexpectedMongoError{err}
	}
	return nil
}

func (wekan *Wekan) InsertBoardLabel(ctx context.Context, board Board, boardLabel BoardLabel) error {
	if err := wekan.CheckAdminUserIsAdmin(ctx); err != nil {
		return err
	}

	if board.GetLabelByName(boardLabel.Name) != (BoardLabel{}) {
		return BoardLabelAlreadyExistsError{boardLabel, board}
	}
	_, err := wekan.db.Collection("boards").UpdateOne(ctx, bson.M{"_id": board.ID},
		bson.M{
			"$push": bson.M{
				"labels": boardLabel,
			},
		})
	return err
}

func (wekan *Wekan) SelectBoardsFromMemberID(ctx context.Context, memberID UserID) ([]Board, error) {
	var boards []Board
	query := bson.M{
		"members.userId": memberID,
	}
	cur, err := wekan.db.Collection("boards").Find(ctx, query)
	if err != nil {
		return nil, UnexpectedMongoError{err}
	}
	err = cur.All(ctx, &boards)
	if err != nil {
		return nil, UnexpectedMongoError{err}
	}
	return boards, nil
}
