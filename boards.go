package libwekan

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Board représente un objet de la collection `boards`
type Board struct {
	ID                         BoardID       `bson:"_id" json:"_id,omitempty"`
	Title                      BoardTitle    `bson:"title" json:"title,omitempty"`
	Permission                 string        `bson:"permission" json:"permission,omitempty"`
	Sort                       float64       `bson:"sort" json:"sort,omitempty"`
	Archived                   bool          `bson:"archived" json:"archived,omitempty"`
	CreatedAt                  time.Time     `bson:"createdAt" json:"createdAt,omitempty"`
	ModifiedAt                 time.Time     `bson:"modifiedAt" json:"modifiedAt,omitempty"`
	Stars                      int           `bson:"stars" json:"stars,omitempty"`
	Labels                     []BoardLabel  `bson:"labels" json:"labels,omitempty"`
	Members                    []BoardMember `bson:"members" json:"members,omitempty"`
	Color                      string        `bson:"color" json:"color,omitempty"`
	SubtasksDefaultBoardId     *string       `bson:"subtasksDefaultBoardId" json:"subtasksDefaultBoardId,omitempty"`
	SubtasksDefaultListId      *string       `bson:"subtasksDefaultListId" json:"subtasksDefaultListId,omitempty"`
	DateSettingsDefaultBoardId *string       `bson:"dateSettingsDefaultBoardId" json:"dateSettingsDefaultBoardId,omitempty"`
	DateSettingsDefaultListId  *string       `bson:"dateSettingsDefaultListId" json:"dateSettingsDefaultListId,omitempty"`
	AllowsSubtasks             bool          `bson:"allowsSubtasks" json:"allowsSubtasks,omitempty"`
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
	PresentParentTask          string        `bson:"presentParentTask" json:"presentParentTask,omitempty"`
	IsOvertime                 bool          `bson:"isOvertime" json:"isOvertime,omitempty"`
	Type                       string        `bson:"type" json:"type,omitempty"`
	Slug                       BoardSlug     `bson:"slug" json:"slug,omitempty"`
	Watchers                   []interface{} `bson:"watchers" json:"watchers,omitempty"`
	AllowsCardNumber           bool          `bson:"allowsCardNumber" json:"allowsCardNumber,omitempty"`
	AllowsShowLists            bool          `bson:"allowsShowLists" json:"allowsShowLists,omitempty"`
}

type BoardLabelID string
type BoardLabelName string
type BoardLabel struct {
	ID    BoardLabelID   `bson:"_id" json:"_id,omitempty"`
	Name  BoardLabelName `bson:"name" json:"name,omitempty"`
	Color string         `bson:"color" json:"color,omitempty"`
}

type BoardMember struct {
	UserID        UserID `bson:"userId" json:"userId,omitempty"`
	IsAdmin       bool   `bson:"isAdmin" json:"isAdmin,omitempty"`
	IsActive      bool   `bson:"isActive" json:"isActive,omitempty"`
	IsNoComments  bool   `bson:"isNoComments" json:"isNoComments,omitempty"`
	IsCommentOnly bool   `bson:"isCommentOnly" json:"isCommentOnly,omitempty"`
	IsWorker      bool   `bson:"isWorker" json:"isWorker,omitempty"`
}

type BoardID string
type BoardSlug string
type BoardTitle string

func (boardID BoardID) GetDocument(ctx context.Context, wekan *Wekan) (Board, error) {
	return wekan.GetBoardFromID(ctx, boardID)
}

func (boardID BoardID) Check(ctx context.Context, wekan *Wekan) error {
	_, err := wekan.GetBoardFromID(ctx, boardID)
	return err
}

// NewBoardLabel retourne un objet BoardLabel
func NewBoardLabel(name string, color string) BoardLabel {
	return BoardLabel{
		ID:    BoardLabelID(newId6()),
		Name:  BoardLabelName(name),
		Color: color,
	}
}

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

// GetBoardFromSlug retourne l'objet board à partir du champ .slug
func (wekan *Wekan) GetBoardFromSlug(ctx context.Context, slug BoardSlug) (Board, error) {
	var board Board
	err := wekan.db.Collection("boards").FindOne(ctx, bson.M{"slug": slug}).Decode(&board)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return Board{}, boardNotFoundWithSlug(slug)
		}
		return Board{}, UnexpectedMongoError{err}
	}
	return board, nil
}

// GetBoardFromTitle retourne l'objet board à partir du champ title
func (wekan *Wekan) GetBoardFromTitle(ctx context.Context, title BoardTitle) (Board, error) {
	var board Board
	err := wekan.db.Collection("boards").FindOne(ctx, bson.M{"title": title}).Decode(&board)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return Board{}, boardNotFoundWithTitle(title)
		}
		return Board{}, UnexpectedMongoError{err}
	}
	return board, nil
}

// GetBoardFromID retourne l'objet board à partir du champ ._id
func (wekan *Wekan) GetBoardFromID(ctx context.Context, id BoardID) (Board, error) {
	var board Board
	if err := wekan.db.Collection("boards").FindOne(ctx, bson.M{"_id": id}).Decode(&board); err != nil {
		if err == mongo.ErrNoDocuments {
			return Board{}, boardNotFoundWithId(id)
		}
		return Board{}, UnexpectedMongoError{err}
	}
	return board, nil
}

// GetMember teste si l'utilisateur fait partie de l'array members
func (board Board) GetMember(userID UserID) BoardMember {
	for _, boardMember := range board.Members {
		if boardMember.UserID == userID {
			return boardMember
		}
	}
	return BoardMember{}
}

// UserIsMember teste si l'utilisateur est membre de la board, activé ou non
func (board Board) UserIsMember(user User) bool {
	return board.GetMember(user.ID) != BoardMember{}
}

// UserIsActiveMember teste si l'utilisateur est activé sur la board, s'il est absent il est alors considéré comme inactif
func (board Board) UserIsActiveMember(user User) bool {
	return board.GetMember(user.ID).IsActive
}

// AddMemberToBoard ajoute un objet BoardMember sur la board
func (wekan *Wekan) AddMemberToBoard(ctx context.Context, boardID BoardID, boardMember BoardMember) error {
	if err := wekan.AssertPrivileged(ctx); err != nil {
		return err
	}

	if _, err := wekan.GetUserFromID(ctx, boardMember.UserID); err != nil {
		return err
	}

	toInsertBoardMember := boardMember
	// l'utilisateur est activé par la méthode EnableBoardMember pour prendre en charge l'insertion de l'activity
	toInsertBoardMember.IsActive = false

	_, err := wekan.db.Collection("boards").UpdateOne(ctx,
		bson.M{"_id": boardID},
		bson.M{
			"$push": bson.M{
				"members": toInsertBoardMember,
			},
		})
	if err != nil {
		return UnexpectedMongoError{err}
	}

	if boardMember.IsActive {
		err = wekan.EnableBoardMember(ctx, boardID, boardMember.UserID)
		if err != nil {
			return err
		}
	}
	return nil
}

// EnableBoardMember active l'utilisateur dans la propriété `member` d'une board
func (wekan *Wekan) EnableBoardMember(ctx context.Context, boardID BoardID, userID UserID) error {
	if err := wekan.AssertPrivileged(ctx); err != nil {
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
		return errors.Wrap(err, "erreur pendant l'insertion d'une activité")
	}
	return nil
}

// DisableBoardMember desactive l'utilisateur dans la propriété `member` d'une board
func (wekan *Wekan) DisableBoardMember(ctx context.Context, boardID BoardID, userID UserID) error {
	if err := wekan.AssertPrivileged(ctx); err != nil {
		return err
	}
	if wekan.adminUserID == userID {
		return ForbiddenOperationError{
			ProtectedUserError{userID},
		}
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
func (wekan *Wekan) EnsureUserIsActiveBoardMember(ctx context.Context, boardID BoardID, userID UserID) (bool, error) {
	if err := wekan.AssertPrivileged(ctx); err != nil {
		return false, err
	}
	board, err := boardID.GetDocument(ctx, wekan)
	if err != nil {
		return false, errors.Wrapf(err, "erreur pendant la récupération d'un document board ??? %s", boardID)
	}
	user, err := userID.GetDocument(ctx, wekan)
	if err != nil {
		return false, errors.Wrapf(err, "erreur pendant la récupération d'un document user ??? %s", userID)
	}
	if board.UserIsActiveMember(user) {
		return false, nil // l'utilisateur est déjà membre actif pas d'action requise
	}
	if board.UserIsMember(user) {
		return true, wekan.EnableBoardMember(ctx, board.ID, user.ID)
	}
	return true, wekan.AddMemberToBoard(ctx, board.ID, BoardMember{user.ID, false, true, false, false, false})
}

// EnsureUserIsInactiveBoardMember fait en sorte de désactiver un utilisateur sur une board lorsqu'il est participant
func (wekan *Wekan) EnsureUserIsInactiveBoardMember(ctx context.Context, boardID BoardID, userID UserID) (bool, error) {
	if err := wekan.AssertPrivileged(ctx); err != nil {
		return false, err
	}
	board, err := boardID.GetDocument(ctx, wekan)
	if err != nil {
		return false, err
	}
	user, err := userID.GetDocument(ctx, wekan)
	if err != nil {
		return false, err
	}
	if board.UserIsActiveMember(user) {
		return true, wekan.DisableBoardMember(ctx, board.ID, user.ID)
	}
	return false, nil
}

func (wekan *Wekan) EnsureUserIsBoardAdmin(ctx context.Context, boardID BoardID, userID UserID) (bool, error) {
	if err := wekan.AssertPrivileged(ctx); err != nil {
		return false, err
	}

	_, err := wekan.EnsureUserIsActiveBoardMember(ctx, boardID, userID)
	if err != nil {
		return false, err
	}
	stats, err := wekan.db.Collection("boards").UpdateOne(ctx, bson.M{"_id": boardID},
		bson.M{
			"$set": bson.M{"members.$[member].isAdmin": true},
		},
		&options.UpdateOptions{
			ArrayFilters: &options.ArrayFilters{
				Filters: bson.A{bson.M{"member.userId": userID}}},
		},
	)
	if err != nil {
		return false, UnexpectedMongoError{err}
	}
	return stats.ModifiedCount != 0, nil
}

func BuildBoard(title string, slug string, boardType string) Board {
	board := Board{
		ID:         BoardID(newId()),
		Title:      BoardTitle(title),
		Permission: "private",
		Type:       boardType,
		Slug:       BoardSlug(slug),
		Archived:   false,
		CreatedAt:  toMongoTime(time.Now()),
		ModifiedAt: toMongoTime(time.Now()),
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
	if err := wekan.AssertPrivileged(ctx); err != nil {
		return err
	}

	_, err := wekan.insertActivity(ctx, newActivityCreateBoard(wekan.adminUserID, board.ID))
	if err != nil {
		return err
	}
	_, err = wekan.db.Collection("boards").InsertOne(ctx, board)
	if err != nil {
		return UnexpectedMongoError{err}
	}
	return nil
}

func (wekan *Wekan) InsertBoardLabel(ctx context.Context, board Board, boardLabel BoardLabel) error {
	if err := wekan.AssertPrivileged(ctx); err != nil {
		return err
	}

	if board.GetLabelByName(boardLabel.Name) != (BoardLabel{}) {
		return BoardLabelAlreadyExistsError{boardLabel, board}
	}
	stats, err := wekan.db.Collection("boards").UpdateOne(ctx,
		bson.M{
			"_id": board.ID,
		},
		bson.M{
			"$push": bson.M{
				"labels": boardLabel,
			},
		})
	if err != nil {
		return UnexpectedMongoError{err}
	}
	if stats.ModifiedCount != 1 {
		return boardNotFoundWithId(board.ID)
	}
	return err
}

// SelectBoardsFromMemberID retourne les boards où on trouve le memberID passé en paramètre
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

// SelectDomainBoards retourne les boards correspondant à la slugDomainRegexp
func (wekan *Wekan) SelectDomainBoards(ctx context.Context) ([]Board, error) {
	var boards []Board
	query := bson.M{
		"slug": primitive.Regex{Pattern: wekan.slugDomainRegexp, Options: "i"},
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

// HasLabelName est vrai lorsque la board dispose du labelName passé en paramètre
func (board Board) HasLabelName(name BoardLabelName) bool {
	for _, label := range board.Labels {
		if label.Name == name {
			return true
		}
	}
	return false
}

func (board Board) HasAnyLabelNames(names []BoardLabelName) bool {
	for _, name := range names {
		if board.HasLabelName(name) {
			return true
		}
	}
	return false
}

func (board Board) HasAllLabelNames(names []BoardLabelName) bool {
	for _, name := range names {
		if !board.HasLabelName(name) {
			return false
		}
	}
	return true
}
