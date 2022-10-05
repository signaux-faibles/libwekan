package libwekan

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type CardID string

type Poker struct {
	Question             bool     `bson:"question"`
	One                  []string `bson:"one"`
	Two                  []string `bson:"two"`
	Three                []string `bson:"three"`
	Five                 []string `bson:"five"`
	Eight                []string `bson:"eight"`
	Thirteen             []string `bson:"thirteen"`
	Twenty               []string `bson:"twenty"`
	Forty                []string `bson:"forty"`
	OneHundred           []string `bson:"oneHundred"`
	Unsure               []string `bson:"unsure"`
	End                  *bool    `bson:"end"`
	AllowNonBoardMembers bool     `bson:"allowNonBoardMembers"`
}
type Vote struct {
	Question             string   `bson:"question"`
	Positive             []string `bson:"positive"`
	Negative             []string `bson:"negative"`
	End                  *bool    `bson:"end"`
	Public               bool     `bson:"public"`
	AllowNonBoardMembers bool     `bson:"allowNonBoardMembers"`
}

type CustomField struct {
	ID    string `json:"_id" bson:"_id"`
	Value string `json:"value" bson:"value"`
}

type Card struct {
	ID               CardID         `bson:"_id"`
	Title            string         `bson:"title"`
	Members          []UserID       `bson:"members"`
	LabelIDs         []BoardLabelID `bson:"labelIds"`
	CustomFields     []CustomField  `bson:"customFields"`
	ListID           ListID         `bson:"listId"`
	BoardID          BoardID        `bson:"boardId"`
	Sort             float64        `bson:"sort"`
	SwimlaneID       SwimlaneID     `bson:"swimlaneId"`
	Type             string         `bson:"type"`
	Archived         bool           `bson:"archived"`
	ParentID         CardID         `bson:"parentId,omitempty"`
	CoverID          string         `bson:"coverId"`
	CreatedAt        time.Time      `bson:"createdAt"`
	ModifiedAt       time.Time      `bson:"modifiedAt"`
	DateLastActivity time.Time      `bson:"dateLastActivity"`
	Description      string         `bson:"description"`
	RequestedBy      UserID         `bson:"requestedBy"`
	AssignedBy       UserID         `bson:"assignedBy"`
	Assignees        []UserID       `bson:"assignees"`
	SpentTime        int            `bson:"spentTime"`
	IsOverTime       bool           `bson:"isOvertime"`
	UserID           UserID         `bson:"userId"`
	SubtaskSort      int            `bson:"subtaskSort"`
	LinkedID         CardID         `bson:"linkedId"`
	Vote             Vote           `bson:"vote"`
	Poker            Poker          `bson:"poker"`
	TargetIDGantt    []string       `bson:"targetId_gantt"`
	LinkTypeGantt    []string       `bson:"linkType_gantt"`
	LinkIDGantt      []string       `bson:"linkId_gantt"`
	StartAt          time.Time      `bson:"startAt"`
}

func BuildCard(boardID BoardID, listID ListID, swimlaneID SwimlaneID, title string, description string, userID UserID) Card {
	return Card{
		ID:               CardID(newId()),
		Title:            title,
		ListID:           listID,
		BoardID:          boardID,
		SwimlaneID:       swimlaneID,
		Members:          []UserID{},
		LabelIDs:         []BoardLabelID{},
		Type:             "card",
		CreatedAt:        toMongoTime(time.Now()),
		ModifiedAt:       toMongoTime(time.Now()),
		DateLastActivity: toMongoTime(time.Now()),
		Description:      description,
		UserID:           userID,
		TargetIDGantt:    []string{},
		LinkTypeGantt:    []string{},
		LinkIDGantt:      []string{},
		StartAt:          toMongoTime(time.Now()),
	}
}

func (cardID CardID) Check(ctx context.Context, wekan *Wekan) error {
	_, err := wekan.GetCardFromID(ctx, cardID)
	return err
}

func (card *Card) AddMember(memberID UserID) {
	if !(contains(card.Members, memberID)) {
		card.Members = append(card.Members, memberID)
	}
}

func (wekan *Wekan) SelectCardsFromQuery(ctx context.Context, query bson.M) ([]Card, error) {
	var cards []Card
	cur, err := wekan.db.Collection("cards").Find(ctx, query)
	if err != nil {
		return nil, UnexpectedMongoError{err}
	}
	err = cur.All(ctx, &cards)
	if err != nil {
		return nil, UnexpectedMongoError{err}
	}
	return cards, nil
}

func (wekan *Wekan) SelectCardsFromUserID(ctx context.Context, userID UserID) ([]Card, error) {
	return wekan.SelectCardsFromQuery(ctx, bson.M{"userId": userID})
}

func (wekan *Wekan) SelectCardsFromMemberID(ctx context.Context, userID UserID) ([]Card, error) {
	return wekan.SelectCardsFromQuery(ctx, bson.M{"members": userID})
}

func (wekan *Wekan) SelectCardsFromBoardID(ctx context.Context, boardID BoardID) ([]Card, error) {
	return wekan.SelectCardsFromQuery(ctx, bson.M{"boardId": boardID})
}

func (wekan *Wekan) SelectCardsFromSwimlaneID(ctx context.Context, swimlaneID SwimlaneID) ([]Card, error) {
	return wekan.SelectCardsFromQuery(ctx, bson.M{"swimlaneId": swimlaneID})
}

func (wekan *Wekan) SelectCardsFromListID(ctx context.Context, listID ListID) ([]Card, error) {
	return wekan.SelectCardsFromQuery(ctx, bson.M{"listId": listID})
}

func (wekan *Wekan) GetCardFromID(ctx context.Context, cardID CardID) (Card, error) {
	cards, err := wekan.SelectCardsFromQuery(ctx, bson.M{"_id": cardID})
	if err != nil {
		return Card{}, UnexpectedMongoError{err}
	}
	if len(cards) == 0 {
		return Card{}, CardNotFoundError{cardID}
	}
	if len(cards) > 1 {
		return Card{}, UnexpectedMongoError{errors.New("erreur fatale, cette requête ne peut retourner qu'un objet")}
	}
	return cards[0], nil
}

func (wekan *Wekan) InsertCard(ctx context.Context, card Card) error {
	if err := wekan.CheckDocuments(
		ctx,
		card.BoardID,
		card.ListID,
		card.SwimlaneID,
	); err != nil {
		return err
	}
	if _, err := wekan.db.Collection("cards").InsertOne(ctx, card); err != nil {
		return UnexpectedMongoError{err}
	}
	return nil
}
